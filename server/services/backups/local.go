package backups

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"snaptrack/db"
	"time"
)

// countingReader wraps an io.Reader and calls onRead callback on every read
type countingReader struct {
	r      io.Reader
	onRead func(int64)
}

func (cr countingReader) Read(p []byte) (int, error) {
	n, err := cr.r.Read(p)
	if n > 0 {
		cr.onRead(int64(n))
	}
	return n, err
}

// runLocalBackup handles local backups with progress updates
func (bs *BackupService) runLocalBackup(backup db.Backup, progress *db.BackupProgress) (int64, string, error) {
	fmt.Println("getting data")

	// Step 1: Calculate total bytes
	progress.Message = "Scanning files..."
	bs.BroadcastProgress(progress)

	totalBytes, err := bs.calculateTotalBytes(backup.Source)
	if err != nil {
		bs.updateProgress(progress, 0, "failed", fmt.Sprintf("Failed to calculate total bytes: %v", err))
		return 0, "", err
	}
	progress.TotalBytes = &totalBytes
	progress.Progress = 0
	progress.BytesProcessed = 0
	bs.BroadcastProgress(progress)

	// Step 2: Run backup based on file type
	switch backup.FileType {
	case "tar":
		return bs.createTarBackupWithProgress(backup.Source, backup.Destination, progress)
	case "zip":
		return bs.createZipBackupWithProgress(backup.Source, backup.Destination, progress)
	case "raw":
		return bs.createRawBackupWithProgress(backup.Source, backup.Destination, progress)
	default:
		return 0, "", fmt.Errorf("unsupported file type: %s", backup.FileType)
	}
}

// ------------------- TAR BACKUP -------------------
func (bs *BackupService) createTarBackupWithProgress(source, destinationDir string, progress *db.BackupProgress) (int64, string, error) {
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return 0, "", fmt.Errorf("failed to create destination directory: %v", err)
	}

	tarFilePath := filepath.Join(destinationDir, "backup.tar")
	f, err := os.Create(tarFilePath)
	if err != nil {
		return 0, "", err
	}
	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	start := time.Now()
	lastUpdate := time.Now()
	delta := int64(0)

	onRead := func(n int64) {
		progress.BytesProcessed += n
		delta += n
		now := time.Now()
		if delta >= 1<<20 || now.Sub(lastUpdate) >= time.Second {
			elapsed := now.Sub(start).Seconds()
			var speed int64
			if elapsed > 0 {
				speed = int64(float64(progress.BytesProcessed) / elapsed)
			}
			var eta int64
			if speed > 0 && progress.TotalBytes != nil {
				eta = (*progress.TotalBytes - progress.BytesProcessed) / speed
			}
			progress.SpeedBPS = &speed
			progress.ETASeconds = &eta
			if progress.TotalBytes != nil {
				progress.Progress = int(progress.BytesProcessed * 100 / *progress.TotalBytes)
			}
			progress.UpdatedAt = now
			db.DB.Save(progress)
			bs.BroadcastProgress(progress)
			lastUpdate = now
			delta = 0
		}
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = rel

		if err = tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			progress.CurrentFile = &rel
			bs.BroadcastProgress(progress)
			in, err := os.Open(path)
			if err != nil {
				return err
			}
			defer in.Close()
			_, err = io.Copy(tw, countingReader{r: in, onRead: onRead})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return 0, "", err
	}

	progress.Progress = 100
	progress.Message = "Calculating checksum..."
	bs.BroadcastProgress(progress)

	info, err := os.Stat(tarFilePath)
	if err != nil {
		return 0, "", err
	}
	checksum, err := bs.calculateFileChecksum(tarFilePath)
	if err != nil {
		return 0, "", err
	}

	progress.Message = "Backup completed successfully"
	bs.BroadcastProgress(progress)

	return info.Size(), checksum, nil
}

// ------------------- ZIP BACKUP -------------------
func (bs *BackupService) createZipBackupWithProgress(source, destinationDir string, progress *db.BackupProgress) (int64, string, error) {
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return 0, "", fmt.Errorf("failed to create destination directory: %v", err)
	}

	zipFilePath := filepath.Join(destinationDir, "backup.zip")
	f, err := os.Create(zipFilePath)
	if err != nil {
		return 0, "", err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	start := time.Now()
	lastUpdate := time.Now()
	delta := int64(0)

	onRead := func(n int64) {
		progress.BytesProcessed += n
		delta += n
		now := time.Now()
		if delta >= 1<<20 || now.Sub(lastUpdate) >= time.Second {
			elapsed := now.Sub(start).Seconds()
			var speed int64
			if elapsed > 0 {
				speed = int64(float64(progress.BytesProcessed) / elapsed)
			}
			var eta int64
			if speed > 0 && progress.TotalBytes != nil {
				eta = (*progress.TotalBytes - progress.BytesProcessed) / speed
			}
			progress.SpeedBPS = &speed
			progress.ETASeconds = &eta
			if progress.TotalBytes != nil {
				progress.Progress = int(progress.BytesProcessed * 100 / *progress.TotalBytes)
			}
			progress.UpdatedAt = now
			db.DB.Save(progress)
			bs.BroadcastProgress(progress)
			lastUpdate = now
			delta = 0
		}
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		header := &zip.FileHeader{
			Name:     rel,
			Method:   zip.Deflate,
			Modified: info.ModTime(),
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		progress.CurrentFile = &rel
		bs.BroadcastProgress(progress)

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(writer, countingReader{r: in, onRead: onRead})
		return err
	})

	if err != nil {
		return 0, "", err
	}

	progress.Progress = 100
	progress.Message = "Calculating checksum..."
	bs.BroadcastProgress(progress)

	info, err := os.Stat(zipFilePath)
	if err != nil {
		return 0, "", err
	}
	checksum, err := bs.calculateFileChecksum(zipFilePath)
	if err != nil {
		return 0, "", err
	}

	progress.Message = "Backup completed successfully"
	bs.BroadcastProgress(progress)
	return info.Size(), checksum, nil
}

// ------------------- RAW BACKUP -------------------
func (bs *BackupService) createRawBackupWithProgress(source, destination string, progress *db.BackupProgress) (int64, string, error) {
	if err := os.MkdirAll(destination, 0755); err != nil {
		return 0, "", fmt.Errorf("failed to create destination directory: %v", err)
	}

	start := time.Now()
	lastUpdate := time.Now()
	delta := int64(0)

	onRead := func(n int64) {
		progress.BytesProcessed += n
		delta += n
		now := time.Now()
		if delta >= 1<<20 || now.Sub(lastUpdate) >= time.Second {
			elapsed := now.Sub(start).Seconds()
			var speed int64
			if elapsed > 0 {
				speed = int64(float64(progress.BytesProcessed) / elapsed)
			}
			var eta int64
			if speed > 0 && progress.TotalBytes != nil {
				eta = (*progress.TotalBytes - progress.BytesProcessed) / speed
			}
			progress.SpeedBPS = &speed
			progress.ETASeconds = &eta
			if progress.TotalBytes != nil {
				progress.Progress = int(progress.BytesProcessed * 100 / *progress.TotalBytes)
			}
			progress.UpdatedAt = now
			db.DB.Save(progress)
			bs.BroadcastProgress(progress)
			lastUpdate = now
			delta = 0
		}
	}

	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destination, rel)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		progress.CurrentFile = &rel
		bs.BroadcastProgress(progress)

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, countingReader{r: in, onRead: onRead})
		return err
	})

	if err != nil {
		return 0, "", err
	}

	progress.Progress = 100
	progress.Message = "Calculating checksum..."
	bs.BroadcastProgress(progress)

	checksum, err := bs.calculateDirectoryChecksum(destination)
	if err != nil {
		return 0, "", fmt.Errorf("failed to calculate checksum: %v", err)
	}

	progress.Message = "Backup completed successfully"
	bs.BroadcastProgress(progress)
	return *progress.TotalBytes, checksum, nil
}

// ------------------- CHECKSUM HELPERS -------------------
func (bs *BackupService) calculateFileChecksum(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func (bs *BackupService) calculateDirectoryChecksum(dirPath string) (string, error) {
	hash := sha256.New()
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			hash.Write(data)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
