package backups

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"snaptrack/db"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func validateRemoteServer(host string, user *string, keyPath *string, port *int, transferType *string) error {
    if host == "" || user == nil || keyPath == nil || port == nil {
        return fmt.Errorf("missing SSH credentials or host")
    }

    key, err := os.ReadFile(*keyPath)
    if err != nil {
        return fmt.Errorf("failed to read SSH key: %v", err)
    }
    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        return fmt.Errorf("failed to parse SSH key: %v", err)
    }
    config := &ssh.ClientConfig{
        User:            *user,
        Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout:         5 * time.Second,
    }
    client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, *port), config)
    if err != nil {
        return fmt.Errorf("cannot connect to remote server: %v", err)
    }
    client.Close()
    return nil
}

func (bs *BackupService) getRsyncTotalSize(source string, server db.Server, dest string) (int64, error) {
    sshCmd := fmt.Sprintf("ssh -i %q -p %d -o StrictHostKeyChecking=no", *server.SSHKeyPath, *server.SSHPort)
    args := []string{
        "-az",
        "--dry-run",
        "--stats",
        "-e", sshCmd,
        source,
        fmt.Sprintf("%s@%s:%s", *server.SSHUser, server.Host, dest),
    }
    cmd := exec.Command("rsync", args...)
    out, err := cmd.CombinedOutput()
    if err != nil {
        return 0, fmt.Errorf("rsync dry-run failed: %v: %s", err, strings.TrimSpace(string(out)))
    }
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Total file size:") {
			sizeStr := strings.TrimSpace(strings.TrimPrefix(line, "Total file size:"))
			sizeStr = strings.TrimSuffix(sizeStr, " bytes")
			sizeStr = strings.ReplaceAll(sizeStr, ",", "")
			return strconv.ParseInt(sizeStr, 10, 64)
		}
	}
	return 0, fmt.Errorf("could not find total file size in rsync output")
}

func (bs *BackupService) runRsyncBackup(source string, server db.Server, dest string, progress *db.BackupProgress) error {
    // Run rsync and stream native output to both terminal and websocket clients
    sshCmd := fmt.Sprintf("ssh -i %q -p %d -o StrictHostKeyChecking=no", *server.SSHKeyPath, *server.SSHPort)
    args := []string{
        "-az",
        "--progress",
        "-e", sshCmd,
        source,
        fmt.Sprintf("%s@%s:%s", *server.SSHUser, server.Host, dest),
    }

    cmd := exec.Command("rsync", args...)
    stdoutPipe, _ := cmd.StdoutPipe()
    stderrPipe, _ := cmd.StderrPipe()

    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start rsync: %v", err)
    }

    // Helper to print and broadcast a line
    emit := func(line string) {
        fmt.Println(line)
        progress.Message = line
        progress.UpdatedAt = time.Now()
        db.DB.Save(progress)
        bs.BroadcastProgress(progress)
    }

    // Stream stdout lines
    go func() {
        scanner := bufio.NewScanner(stdoutPipe)
        for scanner.Scan() {
            emit(scanner.Text())
        }
    }()

    // Stream stderr lines
    go func() {
        scanner := bufio.NewScanner(stderrPipe)
        for scanner.Scan() {
            emit(scanner.Text())
        }
    }()

    if err := cmd.Wait(); err != nil {
        bs.updateProgress(progress, progress.Progress, "failed", fmt.Sprintf("rsync failed: %v", err))
        return err
    }

    bs.updateProgress(progress, 100, "completed", "Backup completed successfully")
    return nil
}
func (bs *BackupService) runRemoteBackup(backup db.Backup, server db.Server, progress *db.BackupProgress) error {
    // Choose transfer type (default: rsync)
    transfer := "rsync"
    if server.TransferType != nil && *server.TransferType != "" {
        transfer = *server.TransferType
    }

    if transfer != "rsync" {
        // For now, only rsync is fully supported with progress; stream a notice
        bs.updateProgress(progress, progress.Progress, "running", fmt.Sprintf("Transfer type '%s' selected; falling back to rsync for progress support", transfer))
    }

    if backup.FileType == "raw" {
		source := backup.Source + "/."
		dest := backup.Destination
		total, err := bs.getRsyncTotalSize(source, server, dest)
		if err != nil {
			return err
		}
		t := total
		progress.TotalBytes = &t
		progress.BytesProcessed = 0
        err = bs.runRsyncBackup(source, server, dest, progress)
		if err != nil {
			return err
		}
		checksum, err := bs.calculateDirectoryChecksum(backup.Source)
		if err != nil {
			return err
		}
		backup.SizeBytes = total
		backup.Checksum = &checksum
		db.DB.Save(&backup)
		return nil
	} else {
		tempDir, err := os.MkdirTemp("", "backup-*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)

		totalUncomp, err := computeTotalSize(backup.Source)
		if err != nil {
			return err
		}
		t := totalUncomp
		progress.TotalBytes = &t
		progress.BytesProcessed = 0
		bs.updateProgress(progress, 0, "running", "Archiving locally...")
		var archSize int64
		var checksum string
		var archFile string
		if backup.FileType == "tar" {
			archSize, checksum, err = bs.createTarBackupWithProgress(backup.Source, tempDir, progress)
			archFile = filepath.Join(tempDir, "backup.tar")
		} else if backup.FileType == "zip" {
			archSize, checksum, err = bs.createZipBackupWithProgress(backup.Source, tempDir, progress)
			archFile = filepath.Join(tempDir, "backup.zip")
		} else {
			return fmt.Errorf("unsupported file type: %s", backup.FileType)
		}
		if err != nil {
			return err
		}
		archProcessed := progress.BytesProcessed
		overallTotal := archProcessed + archSize
		ot := overallTotal
		progress.TotalBytes = &ot
		progress.Progress = int(archProcessed * 100 / overallTotal)
		progress.Message = "Transferring to remote..."
		progress.BytesProcessed = archProcessed
		db.DB.Save(progress)
		bs.BroadcastProgress(progress)
        err = bs.runRsyncBackup(archFile, server, backup.Destination, progress)
		if err != nil {
			return err
		}
		backup.SizeBytes = archSize
		backup.Checksum = &checksum
		db.DB.Save(&backup)
		return nil
	}
}
