import { motion } from 'framer-motion';
import { Dialog } from '@headlessui/react';
import { Backup, deleteBackup } from '@/app/lib/api';

interface DeleteModelProps {
  backup: Backup;
  onClose: () => void;
  onSuccess: () => void;
  onError: () => void;
  token: string;
}

const DeleteModel = ({ backup, onClose, onSuccess, onError, token }: DeleteModelProps) => {
  const handleDelete = async () => {
    const response = await deleteBackup(backup.id, token);
    if (response.success) {
      onSuccess(); // Let parent handle success toast
    } else {
      onError(response.message);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center"
    >
      <motion.div
        initial={{ scale: 0.8, y: 50 }}
        animate={{ scale: 1, y: 0 }}
        className="bg-white p-6 rounded-lg w-full max-w-md"
      >
        <Dialog.Title className="text-xl font-semibold text-gray-800">Delete Backup</Dialog.Title>
        <p className="mt-2 text-sm text-gray-600">
          Are you sure you want to delete the backup for <strong>{backup.app}</strong>?
        </p>
        <div className="flex justify-end space-x-2 mt-4">
          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-200 rounded-lg hover:bg-gray-300"
          >
            Cancel
          </button>
          <button
            onClick={handleDelete}
            className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600"
          >
            Delete
          </button>
        </div>
      </motion.div>
    </motion.div>
  );
};

export default DeleteModel;