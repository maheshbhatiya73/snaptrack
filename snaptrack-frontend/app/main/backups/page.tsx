"use client";
import { useState, useEffect, Fragment } from 'react';
import { FaPlus, FaEdit, FaTrash, FaCheckCircle, FaTimesCircle, FaExclamationCircle, FaHourglassHalf, FaQuestionCircle, FaEye } from 'react-icons/fa';
import { Dialog, Transition } from '@headlessui/react';
import { Backup, getBackups } from '@/app/lib/api';

type BackupStatus = 'success' | 'failed' | 'pending' | 'started' | 'scheduled' | 'unknown' | string;

type BackupType = Backup & {
    status?: BackupStatus;
    logs?: {
        id: string;
        status: BackupStatus;
        message: string;
        createdAt: string;
    }[];
};

import { Tooltip } from "react-tooltip";
import CreateModel from './CreateModel';
import UpdateModel from './UpdateModel';
import DeleteModel from './DeleteModel';
import { useAuth } from '@/app/context/AuthContext';
import { useToast } from '@/app/components/Toast';

const Home = () => {
    const [backups, setBackups] = useState<BackupType[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [limit] = useState(10);
    const [isCreateOpen, setIsCreateOpen] = useState(false);
    const [isUpdateOpen, setIsUpdateOpen] = useState(false);
    const [isDeleteOpen, setIsDeleteOpen] = useState(false);
    const [isLogsOpen, setIsLogsOpen] = useState(false);
    const [selectedBackup, setSelectedBackup] = useState<BackupType | null>(null);
    const [selectedBackupForLogs, setSelectedBackupForLogs] = useState<BackupType | null>(null);
    const { token } = useAuth();
    const { addToast } = useToast();

    useEffect(() => {
        const fetchBackups = async () => {
            if (!token) {
                return;
            }
            try {
                const response = await getBackups(token, page, limit);
                if (response.success) {
                    setBackups(response.data as Backup[]);
                    setTotal(response.total || 0);
                } else {
                    addToast('Failed to fetch backups', 'error');
                }
            } catch (error) {
                addToast('Error fetching backups', 'error');
            }
        };
        fetchBackups();
    }, [page, limit, token]);

    const formatCountdown = (nextRun?: string) => {
        if (!nextRun || nextRun === '0001-01-01T00:00:00Z') return 'N/A';
        const next = new Date(nextRun);
        const now = new Date();
        const diff = next.getTime() - now.getTime();
        if (diff <= 0) return 'Now';
        const hours = Math.floor(diff / 1000 / 60 / 60);
        const minutes = Math.floor((diff / 1000 / 60) % 60);
        if (hours > 0) return `in ${hours} hour${hours > 1 ? 's' : ''} ${minutes} minute${minutes !== 1 ? 's' : ''}`;
        return `in ${minutes} minute${minutes !== 1 ? 's' : ''}`;
    };

    const formatTimestamp = (timestamp: string) => {
        return new Date(timestamp).toLocaleString('en-US', {
            year: 'numeric',
            month: 'short',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
        });
    };

    const statusIcons = {
        success: <FaCheckCircle className="inline mr-1 text-green-600" />,
        failed: <FaTimesCircle className="inline mr-1 text-red-600" />,
        pending: <FaExclamationCircle className="inline mr-1 text-yellow-600" />,
        started: <FaHourglassHalf className="inline mr-1 text-blue-600" />,
        scheduled: <FaHourglassHalf className="inline mr-1 text-blue-400" />,
        unknown: <FaQuestionCircle className="inline mr-1 text-gray-600" />,
    };

    const formatStatus = (status?: string) => {
        if (!status) return 'Unknown';
        return status.charAt(0).toUpperCase() + status.slice(1);
    };

    return (
        <div className="container mx-auto p-6">
            <div
                className="flex justify-between items-center mb-6"
            >
                <h1 className="text-3xl font-bold text-gray-800">Backup Management</h1>
                <button
                    onClick={() => setIsCreateOpen(true)}
                    className="flex items-center px-4 py-2 bg-sky-400 text-white rounded-lg hover:bg-opacity-90 transition"
                >
                    <FaPlus className="mr-2" /> New Backup
                </button>
            </div>
            <div className="overflow-x-auto bg-white rounded-xl shadow-lg">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            {['App', 'Type', 'Size', 'FileType', 'Status', 'Next Run', 'Actions'].map((header) => (
                                <th
                                    key={header}
                                    className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider select-none"
                                >
                                    {header}
                                </th>
                            ))}
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                            {backups.map((backup, idx) => (
                                <tr
                                    key={idx}
                                    className={`cursor-pointer ${idx % 2 === 0 ? 'bg-white' : 'bg-gray-50'} hover:bg-sky-50`}
                                >
                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 rounded-l-lg">
                                        {backup.app}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700 capitalize">{backup.type}</td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{backup.size || 'N/A'}</td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{backup.fileType || 'N/A'}</td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span
                                            className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold
                                                ${backup.status === 'success' ? 'bg-green-100 text-green-800' :
                                                    backup.status === 'failed' ? 'bg-red-100 text-red-800' :
                                                    backup.status === 'started' ? 'bg-blue-100 text-blue-800' :
                                                    backup.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                                                    'bg-gray-100 text-gray-800'}`}
                                        >
                                            {statusIcons[backup.status || 'unknown'] || null}
                                            {formatStatus(backup.status)}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                                        {formatCountdown(backup.nextRun)}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm flex space-x-4 rounded-r-lg">
                                        <button
                                            onClick={() => {
                                                setSelectedBackup(backup);
                                                setIsUpdateOpen(true);
                                            }}
                                            className="text-sky-500 hover:text-sky-700 transition"
                                            data-tooltip-id="edit-tooltip"
                                            data-tooltip-content="Edit Backup"
                                        >
                                            <FaEdit size={18} />
                                        </button>
                                        <button
                                            onClick={() => {
                                                setSelectedBackup(backup);
                                                setIsDeleteOpen(true);
                                            }}
                                            className="text-red-500 hover:text-red-700 transition"
                                            data-tooltip-id="delete-tooltip"
                                            data-tooltip-content="Delete Backup"
                                        >
                                            <FaTrash size={18} />
                                        </button>
                                        <button
                                            onClick={() => {
                                                setSelectedBackupForLogs(backup);
                                                setIsLogsOpen(true);
                                            }}
                                            className="text-gray-500 hover:text-gray-700 transition"
                                            data-tooltip-id="logs-tooltip"
                                            data-tooltip-content="View Logs"
                                        >
                                            <FaEye size={18} />
                                        </button>
                                        <Tooltip id="edit-tooltip" place="top" />
                                        <Tooltip id="delete-tooltip" place="top" />
                                        <Tooltip id="logs-tooltip" place="top" />
                                    </td>
                                </tr>
                            ))}
                    </tbody>
                </table>
            </div>
        <div className="flex justify-between items-center mt-4">
            <p className="text-sm text-gray-600">
                Showing {backups.length} of {total} backups
            </p>
            <div className="space-x-2">
                <button
                onClick={() => setPage((p) => Math.max(p - 1, 1))}
                disabled={page === 1}
                className="px-4 py-2 bg-gray-200 rounded-lg disabled:opacity-50 hover:bg-gray-300 transition"
                >
                Previous
                </button>
                <button
                onClick={() => setPage((p) => p + 1)}
                disabled={page * limit >= total}
                className="px-4 py-2 bg-gray-200 rounded-lg disabled:opacity-50 hover:bg-gray-300 transition"
                >
                Next
                </button>
            </div>
        </div>
            <Transition show={isCreateOpen} as={Fragment}>
                <Dialog open={isCreateOpen} onClose={() => setIsCreateOpen(false)} className="relative z-50">
                    <CreateModel
                        onClose={() => setIsCreateOpen(false)}
                        onSuccess={(backup: BackupType) => {
                            setIsCreateOpen(false);
                            setBackups((prev) => [backup, ...prev]);
                            addToast(`Backup ${backup.app} created successfully`, 'success');
                        }}
                        onError={(error: string) => {
                            addToast(error || 'Failed to create backup', 'error');
                        }}
                        token={token}
                    />
                </Dialog>
            </Transition>
            <Transition show={isUpdateOpen} as={Fragment}>
                <Dialog open={isUpdateOpen} onClose={() => setIsUpdateOpen(false)} className="relative z-50">
                    {selectedBackup && (
                        <UpdateModel
                            backup={selectedBackup}
                            onClose={() => setIsUpdateOpen(false)}
                            onSuccess={(updatedBackup: BackupType) => {
                                setIsUpdateOpen(false);
                                setBackups((prev) =>
                                    prev.map((b) => (b.id === updatedBackup.id ? updatedBackup : b))
                                );
                                addToast(`Backup ${updatedBackup.app} updated successfully`, 'success');
                            }}
                            onError={(error: string) => {
                                setIsUpdateOpen(false);
                                addToast(error || 'Failed to update backup', 'error');
                            }}
                            token={token}
                        />
                    )}
                </Dialog>
            </Transition>
            <Transition show={isDeleteOpen} as={Fragment}>
                <Dialog open={isDeleteOpen} onClose={() => setIsDeleteOpen(false)} className="relative z-50">
                    {selectedBackup && (
                        <DeleteModel
                            backup={selectedBackup}
                            onClose={() => setIsDeleteOpen(false)}
                            onSuccess={() => {
                                setIsDeleteOpen(false);
                                setBackups((prev) => prev.filter((b) => b.id !== selectedBackup.id));
                                addToast(`Backup ${selectedBackup.app} deleted successfully`, 'success');
                            }}
                            onError={(error: string) => {
                                setIsDeleteOpen(false);
                                addToast(error || 'Failed to delete backup', 'error');
                            }}
                            token={token}
                        />
                    )}
                </Dialog>
            </Transition>
            <Transition show={isLogsOpen} as={Fragment}>
                <Dialog open={isLogsOpen} onClose={() => setIsLogsOpen(false)} className="relative z-50">
                    <Transition.Child
                        as={Fragment}
                        enter="ease-out duration-300"
                        enterFrom="opacity-0"
                        enterTo="opacity-100"
                        leave="ease-in duration-200"
                        leaveFrom="opacity-100"
                        leaveTo="opacity-0"
                    >
                        <div className="fixed inset-0 bg-black bg-opacity-25" />
                    </Transition.Child>
                    <div className="fixed inset-0 overflow-y-auto">
                        <div className="flex min-h-full items-center justify-center p-4 text-center">
                            <Transition.Child
                                as={Fragment}
                                enter="ease-out duration-300"
                                enterFrom="opacity-0 scale-95"
                                enterTo="opacity-100 scale-100"
                                leave="ease-in duration-200"
                                leaveFrom="opacity-100 scale-100"
                                leaveTo="opacity-0 scale-95"
                            >
                                <Dialog.Panel className="w-full max-w-4xl transform overflow-hidden rounded-2xl bg-white p-6 text-left align-middle shadow-xl transition-all">
                                    <Dialog.Title as="h3" className="text-lg font-medium leading-6 text-gray-900">
                                        Logs for Backup: {selectedBackupForLogs?.app || 'N/A'}
                                    </Dialog.Title>
                                    <div className="mt-4">
                                        {selectedBackupForLogs?.logs && selectedBackupForLogs.logs.length > 0 ? (
                                            <div className="overflow-x-auto">
                                                <table className="min-w-full divide-y divide-gray-200">
                                                    <thead className="bg-gray-50">
                                                        <tr>
                                                            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
                                                            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                                                            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Message</th>
                                                            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created At</th>
                                                        </tr>
                                                    </thead>
                                                    <tbody className="bg-white divide-y divide-gray-200">
                                                        {selectedBackupForLogs.logs.map((log) => (
                                                            <tr key={log.id}>
                                                                <td className="px-4 py-2 whitespace-nowrap text-sm text-gray-500">{log.id}</td>
                                                                <td className="px-4 py-2 whitespace-nowrap">
                                                                    <span
                                                                        className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium
                                                                            ${log.status === 'success' ? 'bg-green-100 text-green-800' :
                                                                                log.status === 'failed' ? 'bg-red-100 text-red-800' :
                                                                                log.status === 'started' ? 'bg-blue-100 text-blue-800' :
                                                                                log.status === 'scheduled' ? 'bg-blue-100 text-blue-800' :
                                                                                'bg-gray-100 text-gray-800'}`}
                                                                    >
                                                                        {statusIcons[log.status as keyof typeof statusIcons] || null}
                                                                        {formatStatus(log.status)}
                                                                    </span>
                                                                </td>
                                                                <td className="px-4 py-2 text-sm text-gray-500">{log.message}</td>
                                                                <td className="px-4 py-2 whitespace-nowrap text-sm text-gray-500">
                                                                    {formatTimestamp(log.createdAt)}
                                                                </td>
                                                            </tr>
                                                        ))}
                                                    </tbody>
                                                </table>
                                            </div>
                                        ) : (
                                            <p className="text-sm text-gray-500">No logs available for this backup.</p>
                                        )}
                                    </div>
                                    <div className="mt-6">
                                        <button
                                            type="button"
                                            className="inline-flex justify-center px-4 py-2 text-sm font-medium text-white bg-sky-400 border border-transparent rounded-md hover:bg-sky-500 focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-sky-500"
                                            onClick={() => setIsLogsOpen(false)}
                                        >
                                            Close
                                        </button>
                                    </div>
                                </Dialog.Panel>
                            </Transition.Child>
                        </div>
                    </div>
                </Dialog>
            </Transition>
        </div>
    );
};

export default Home;
