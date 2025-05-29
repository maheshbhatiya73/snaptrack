"use client"
import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { FaPlus, FaEdit, FaTrash, FaCheckCircle } from 'react-icons/fa';
import { Dialog, Transition } from '@headlessui/react';
import { Backup, getBackups } from '@/app/lib/api';
import CreateModel from './CreateModel';
import UpdateModel from './UpdateModel';
import DeleteModel from './DeleteModel';
import { useAuth } from '@/app/context/AuthContext';

const Home = () => {
    const [backups, setBackups] = useState<Backup[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [limit] = useState(10);
    const [isCreateOpen, setIsCreateOpen] = useState(false);
    const [isUpdateOpen, setIsUpdateOpen] = useState(false);
    const [isDeleteOpen, setIsDeleteOpen] = useState(false);
    const [selectedBackup, setSelectedBackup] = useState<Backup | null>(null);
    const { token } = useAuth()

    useEffect(() => {
        const fetchBackups = async () => {
            const response = await getBackups(token, page, limit);
            if (response.success) {
                setBackups(response.data as Backup[]);
                setTotal(response.total || 0);
            }
        };
        fetchBackups();
    }, [page, token]);

    const statusIcons = {
        success: <FaCheckCircle className="inline mr-1 text-green-600" />,
        failed: <FaTimesCircle className="inline mr-1 text-red-600" />,
        pending: <FaExclamationTriangle className="inline mr-1 text-yellow-600" />,
    };

    return (
        <div className="container mx-auto p-6">
            <motion.div
                initial={{ opacity: 0, y: -20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5 }}
                className="flex justify-between items-center mb-6"
            >
                <h1 className="text-3xl font-bold text-gray-800">Backup Management</h1>
                <button
                    onClick={() => setIsCreateOpen(true)}
                    className="flex items-center px-4 py-2 bg-sky-400 text-white rounded-lg hover:bg-opacity-90 transition"
                >
                    <FaPlus className="mr-2" /> New Backup
                </button>
            </motion.div>

            {/* Table */}
            <div className="overflow-x-auto bg-white rounded-xl shadow-lg">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            {['App', 'Type', 'Size', 'Date', 'Status', 'Actions'].map((header) => (
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
                        <AnimatePresence>
                            {backups.map((backup, idx) => (
                                <motion.tr
                                    key={backup.id}
                                    initial={{ opacity: 0, y: 10 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    exit={{ opacity: 0, y: 10 }}
                                    transition={{ duration: 0.2 }}
                                    className={`cursor-pointer ${idx % 2 === 0 ? 'bg-white' : 'bg-gray-50'
                                        } hover:bg-sky-50`}
                                >
                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 rounded-l-lg">
                                        {backup.app}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{backup.type}</td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{backup.size}</td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                                        {new Date(backup.date).toLocaleDateString()}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span
                                            className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold
                  ${backup.status === 'success'
                                                    ? 'bg-green-100 text-green-800'
                                                    : backup.status === 'failed'
                                                        ? 'bg-red-100 text-red-800'
                                                        : 'bg-yellow-100 text-yellow-800'
                                                }
                `}
                                        >
                                            {statusIcons[backup.status] || null}
                                            {backup.status.charAt(0).toUpperCase() + backup.status.slice(1)}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm flex space-x-4 rounded-r-lg">
                                        <button
                                            onClick={() => {
                                                setSelectedBackup(backup);
                                                setIsUpdateOpen(true);
                                            }}
                                            className="text-sky-500 hover:text-sky-700 transition"
                                            data-tip="Edit Backup"
                                        >
                                            <FaEdit size={18} />
                                        </button>
                                        <button
                                            onClick={() => {
                                                setSelectedBackup(backup);
                                                setIsDeleteOpen(true);
                                            }}
                                            className="text-red-500 hover:text-red-700 transition"
                                            data-tip="Delete Backup"
                                        >
                                            <FaTrash size={18} />
                                        </button>
                                        <ReactTooltip place="top" effect="solid" />
                                    </td>
                                </motion.tr>
                            ))}
                        </AnimatePresence>
                    </tbody>
                </table>
            </div>

            {/* Pagination */}
            <div className="flex justify-between items-center mt-4">
                <p className="text-sm text-gray-600">
                    Showing {backups.length} of {total} backups
                </p>
                <div className="space-x-2">
                    <button
                        onClick={() => setPage((p) => Math.max(p - 1, 1))}
                        disabled={page === 1}
                        className="px-4 py-2 bg-gray-200 rounded-lg disabled:opacity-50"
                    >
                        Previous
                    </button>
                    <button
                        onClick={() => setPage((p) => p + 1)}
                        disabled={page * limit >= total}
                        className="px-4 py-2 bg-gray-200 rounded-lg disabled:opacity-50"
                    >
                        Next
                    </button>
                </div>
            </div>

            {/* Modals */}
            <Transition show={isCreateOpen} as="div">
                <Dialog open={isCreateOpen} onClose={() => setIsCreateOpen(false)} className="relative z-50">
                    <CreateModel
                        onClose={() => setIsCreateOpen(false)}
                        onSuccess={() => {
                            setIsCreateOpen(false);
                            // Refresh backups
                        }}
                        token={token}
                    />
                </Dialog>
            </Transition>

            <Transition show={isUpdateOpen} as="div">
                <Dialog open={isUpdateOpen} onClose={() => setIsUpdateOpen(false)} className="relative z-50">
                    {selectedBackup && (
                        <UpdateModel
                            backup={selectedBackup}
                            onClose={() => setIsUpdateOpen(false)}
                            onSuccess={() => {
                                setIsUpdateOpen(false);
                                // Refresh backups
                            }}
                            token={token}
                        />
                    )}
                </Dialog>
            </Transition>

            <Transition show={isDeleteOpen} as="div">
                <Dialog open={isDeleteOpen} onClose={() => setIsDeleteOpen(false)} className="relative z-50">
                    {selectedBackup && (
                        <DeleteModel
                            backup={selectedBackup}
                            onClose={() => setIsDeleteOpen(false)}
                            onSuccess={() => {
                                setIsDeleteOpen(false);
                                // Refresh backups
                            }}
                            token={token}
                        />
                    )}
                </Dialog>
            </Transition>
        </div>
    );
};

export default Home;