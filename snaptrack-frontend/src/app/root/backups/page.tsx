'use client';

import { useState, useEffect, useMemo, JSX } from 'react';
import { useRouter } from 'next/navigation';
import { motion } from 'framer-motion';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';
import { Eye, ArrowUpDown, AlertCircle, CheckCircle2, Clock, FileArchive, Pencil, Trash2 } from 'lucide-react';
import { useAuth } from '@/app/store/useAuth';
import { deleteBackup, getAllBackups } from '@/lib/api';
import { useLinuxToast } from '@/lib/use-linux-toast';

interface BackupLog {
  id: string;
  status: 'completed' | 'pending' | 'scheduled' | 'started';
  message: string;
  createdAt: string;
}

interface Backup {
  id: string;
  app: string;
  type: string;
  size: string;
  status: 'completed' | 'pending' | 'scheduled' | 'started';
  sourcePath: string;
  destinationPath: string;
  fileType: string;
  schedule: {
    kind: string;
  };
  createdAt: string;
  logs: BackupLog[];
}

interface SortConfig {
  key: keyof Backup;
  direction: 'asc' | 'desc';
}

export default function BackupTable() {
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [sortConfig, setSortConfig] = useState<SortConfig>({ key: 'createdAt', direction: 'desc' });
  const [selectedBackup, setSelectedBackup] = useState<Backup | null>(null);
  const [isLogsModalOpen, setIsLogsModalOpen] = useState<boolean>(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState<boolean>(false);
  const [backupToDelete, setBackupToDelete] = useState<Backup | null>(null);
  const [backups, setBackups] = useState<any>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const router = useRouter();
  const { success, error } = useLinuxToast();

  useEffect(() => {
    async function fetchBackups() {
      try {
        setIsLoading(true);
        const response = await getAllBackups();
        const backupsArray = Array.isArray(response) ? response : (response || []);
        setBackups(backupsArray);
      } catch (err: any) {
        error(err.message || 'Failed to fetch backups');
      } finally {
        setIsLoading(false);
      }
    }
    fetchBackups();
  }, []);

  const filteredData = useMemo(() => {
    let result = [...backups];
    if (searchTerm) {
      result = result.filter((backup) =>
        backup.app.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }
    result.sort((a: Backup, b: Backup) => {
      const key = sortConfig.key;
      const direction = sortConfig.direction === 'asc' ? 1 : -1;
      if (key === 'createdAt') {
        return direction * (new Date(a[key]).getTime() - new Date(b[key]).getTime());
      }
      return direction * (a[key] as string).localeCompare(b[key] as string);
    });
    return result;
  }, [searchTerm, sortConfig, backups]);

  const handleSort = (key: keyof Backup) => {
    setSortConfig((prev) => ({
      key,
      direction: prev.key === key && prev.direction === 'asc' ? 'desc' : 'asc',
    }));
  };

  const openLogsModal = (backup: Backup) => {
    setSelectedBackup(backup);
    setIsLogsModalOpen(true);
  };

  const openDeleteModal = (backup: Backup) => {
    setBackupToDelete(backup);
    setIsDeleteModalOpen(true);
  };

  function getRelativeTimeLabel(dateString: string): string {
    const now = new Date();
    const target = new Date(dateString);
    const diffMs = target.getTime() - now.getTime();

    if (diffMs <= 0) return 'now';

    const diffMinutes = Math.floor(diffMs / (1000 * 60));
    const hours = Math.floor(diffMinutes / 60);
    const minutes = diffMinutes % 60;

    if (hours > 0) {
      return `in ${hours}h ${minutes}m`;
    }

    return `in ${minutes} min`;
  }

  const handleDelete = async () => {
    if (!backupToDelete) return;
    try {
      await deleteBackup(backupToDelete.id);
      setBackups((prev: any) => prev.filter((backup: any) => backup.id !== backupToDelete.id));
      setIsDeleteModalOpen(false);
      setBackupToDelete(null);
      success('Backup deleted successfully');
    } catch (err: any) {
      error(err.message || 'Failed to delete backup');
    }
  };

  const handleCreateBackup = () => {
    router.push('/root/backups/create');
  };

  const handleEditBackup = (id: string) => {
    router.push(`/root/backups/update/${id}`);
  };

  const statusIcons: Record<string, JSX.Element> = {
    completed: <CheckCircle2 className="w-5 h-5 text-green-500" aria-label="Completed" />,
    pending: <Clock className="w-5 h-5 text-yellow-500" aria-label="Pending" />,
    scheduled: <Clock className="w-5 h-5 text-blue-500" aria-label="Scheduled" />,
    started: <FileArchive className="w-5 h-5 text-blue-500" aria-label="Started" />,
  };

  return (
    <div className="mt-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="mx-auto"
      >
        <div className="flex justify-between items-center mb-6">
          <Input
            placeholder="Search by app name..."
            value={searchTerm}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchTerm(e.target.value)}
            className="max-w-md bg-gray-900/50 border-gray-700 text-white placeholder-gray-500 focus:ring-2 focus:ring-green-500"
            aria-label="Search backups"
          />
          <Button
            onClick={handleCreateBackup}
            className="bg-gradient-to-r from-green-500 to-cyan-500 hover:from-green-600 hover:to-cyan-600 text-white font-semibold py-2 px-4 rounded-lg"
            aria-label="Create new backup"
          >
            <FileArchive className="w-5 h-5 mr-2" />
            Create Backup
          </Button>
        </div>

        {isLoading ? (
          <p className="text-white">Loading backups...</p>
        ) : (
          <div className="overflow-x-auto">
            <Table className="bg-gray-800/80 border-gray-700/50 rounded-lg">
              <TableHeader>
                <TableRow className="border-gray-700">
                  <TableHead
                    className="text-green-400 font-medium sticky top-0 bg-gray-800/80 cursor-pointer"
                    onClick={() => handleSort('app')}
                  >
                    App Name
                    <ArrowUpDown className="inline ml-2 w-4 h-4" aria-hidden="true" />
                  </TableHead>
                  <TableHead className="text-green-400 font-medium sticky top-0 bg-gray-800/80">
                    Type
                  </TableHead>
                  <TableHead className="text-green-400 font-medium sticky top-0 bg-gray-800/80">
                    Size
                  </TableHead>
                  <TableHead
                    className="text-green-400 font-medium sticky top-0 bg-gray-800/80 cursor-pointer"
                    onClick={() => handleSort('status')}
                  >
                    Status
                    <ArrowUpDown className="inline ml-2 w-4 h-4" aria-hidden="true" />
                  </TableHead>
                  <TableHead className="text-green-400 font-medium sticky top-0 bg-gray-800/80">
                    Source Path
                  </TableHead>
                  <TableHead className="text-green-400 font-medium sticky top-0 bg-gray-800/80">
                    Destination Path
                  </TableHead>
                  <TableHead className="text-green-400 font-medium sticky top-0 bg-gray-800/80">
                    File Type
                  </TableHead>
                  <TableHead className="text-green-400 font-medium sticky top-0 bg-gray-800/80">
                    Schedule
                  </TableHead>
                  <TableHead
                    className="text-green-400 font-medium sticky top-0 bg-gray-800/80 cursor-pointer"
                    onClick={() => handleSort('createdAt')}
                  >
                    Created At
                    <ArrowUpDown className="inline ml-2 w-4 h-4" aria-hidden="true" />
                  </TableHead>
                  <TableHead className="text-green-400 font-medium sticky top-0 bg-gray-800/80">
                    Actions
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredData.map((backup: Backup) => (
                  <TableRow
                    key={backup.id}
                    className="border-gray-700 hover:bg-gray-700/50 transition-colors"
                  >
                    <TableCell className="text-white">{backup.app}</TableCell>
                    <TableCell className="text-white">{backup.type}</TableCell>
                    <TableCell className="text-white">{backup.size}</TableCell>
                    <TableCell className="text-white flex items-center gap-2">
                      {statusIcons[backup.status] || <AlertCircle className="w-5 h-5 text-red-500" aria-label="Unknown" />}
                      {backup.status}
                    </TableCell>
                    <TableCell className="text-white">{backup.sourcePath}</TableCell>
                    <TableCell className="text-white">{backup.destinationPath}</TableCell>
                    <TableCell className="text-white">{backup.fileType}</TableCell>
                    <TableCell className="text-white">{backup.schedule.kind}</TableCell>
                    <TableCell className="text-white">
                      {new Date(backup.createdAt).toLocaleString()}
                    </TableCell>
                    <TableCell className="text-white flex gap-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => openLogsModal(backup)}
                        aria-label={`View logs for ${backup.app}`}
                      >
                        <Eye className="w-5 h-5 text-green-400" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleEditBackup(backup.id)}
                        aria-label={`Edit backup ${backup.app}`}
                      >
                        <Pencil className="w-5 h-5 text-blue-400" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => openDeleteModal(backup)}
                        aria-label={`Delete backup ${backup.app}`}
                      >
                        <Trash2 className="w-5 h-5 text-red-400" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
        <Dialog open={isLogsModalOpen} onOpenChange={setIsLogsModalOpen}>
          <DialogContent className="bg-gray-800 border-gray-700 text-white max-w-2xl">
            <DialogHeader>
              <DialogTitle className="text-green-400 flex items-center gap-2">
                <FileArchive className="w-6 h-6" />
                Logs for {selectedBackup?.app}
              </DialogTitle>
            </DialogHeader>
            <div className="max-h-緩 overflow-y-auto">
              <Table className="bg-gray-900/50">
                <TableHeader>
                  <TableRow className="border-gray-700">
                    <TableHead className="text-green-400 font-medium">Status</TableHead>
                    <TableHead className="text-green-400 font-medium">Message</TableHead>
                    <TableHead className="text-green-400 font-medium">Timestamp</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {selectedBackup?.logs.map((log: BackupLog) => (
                    <TableRow key={log.id} className="border-gray-700">
                      <TableCell className="text-white flex items-center gap-2">
                        {statusIcons[log.status] || <AlertCircle className="w-5 h-5 text-red-500" aria-label="Unknown" />}
                        {log.status}
                      </TableCell>
                      <TableCell className="text-white">{log.message}</TableCell>
                      <TableCell className="text-white">
                        {new Date(log.createdAt).toLocaleString()}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </DialogContent>
        </Dialog>

        <Dialog open={isDeleteModalOpen} onOpenChange={setIsDeleteModalOpen}>
          <DialogContent className="bg-gray-800 border-gray-700 text-white">
            <DialogHeader>
              <DialogTitle className="text-red-400 flex items-center gap-2">
                <Trash2 className="w-6 h-6" />
                Confirm Delete
              </DialogTitle>
            </DialogHeader>
            <p>Are you sure you want to delete the backup "{backupToDelete?.app}"? This action cannot be undone.</p>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setIsDeleteModalOpen(false)}
                className="border-gray-600 text-black hover:bg-gray-700"
              >
                Cancel
              </Button>
              <Button
                variant="destructive"
                onClick={handleDelete}
                className="bg-red-600 hover:bg-red-700"
              >
                Delete
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </motion.div>
    </div>
  );
}