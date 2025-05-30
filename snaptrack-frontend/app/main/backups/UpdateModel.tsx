import { useState } from 'react';
import { motion } from 'framer-motion';
import { Dialog, Listbox, Transition } from '@headlessui/react';
import { FaAngleDown, FaTimes } from 'react-icons/fa';
import { Backup, Schedule, updateBackup } from '@/app/lib/api';

interface UpdateModelProps {
  backup: Backup;
  onClose: () => void;
  onSuccess: () => void;
  onError: () => void;
  token: string;
}

const backupTypes = ['manual', 'full', 'incremental'];
const fileTypes = ['zip', 'tar', 'tar.gz'];
const scheduleKinds = ['one-time', 'hourly', 'daily', 'weekly', 'monthly'];

const UpdateModel = ({ backup, onClose, onSuccess, onError, token }: UpdateModelProps) => {
  const [form, setForm] = useState<Partial<Backup>>(backup);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const response = await updateBackup(backup.id, form, token);
    if (response.success && response.data) {
      onSuccess(); 
    } else {
      onError(); 
    }
  };


  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50"
    >
      <motion.div
        initial={{ scale: 0.95, y: 20 }}
        animate={{ scale: 1, y: 0 }}
        className="bg-white w-full max-w-3xl mx-4 rounded-2xl shadow-2xl p-8 relative"
      >
        <Dialog.Title className="text-2xl font-bold text-gray-800 flex justify-between items-center mb-6">
          ✏️ Update Backup
          <button onClick={onClose}>
            <FaTimes className="text-gray-400 hover:text-gray-600" />
          </button>
        </Dialog.Title>

        <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* App Name */}
          <div className="col-span-2">
            <label className="text-sm font-medium text-gray-700">App Name</label>
            <input
              type="text"
              value={form.app}
              onChange={(e) => setForm({ ...form, app: e.target.value })}
              className="w-full mt-1 px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-sky-500"
              required
            />
          </div>

          {/* Backup Type Dropdown */}
          <Listbox
            value={form.type}
            onChange={(value) => setForm({ ...form, type: value })}
          >
            <div className="space-y-1">
              <Listbox.Label className="text-sm font-medium text-gray-700">Backup Type</Listbox.Label>
              <div className="relative">
                <Listbox.Button className="w-full cursor-pointer rounded-xl border border-gray-300 bg-white py-2 pl-4 pr-10 text-left shadow-sm focus:outline-none focus:ring-2 focus:ring-sky-500">
                  <span className="block capitalize">{form.type}</span>
                  <span className="absolute inset-y-0 right-0 flex items-center pr-3">
                    <FaAngleDown className="h-4 w-4 text-gray-400" />
                  </span>
                </Listbox.Button>
                <Transition
                  enter="transition duration-100 ease-out"
                  enterFrom="transform scale-95 opacity-0"
                  enterTo="transform scale-100 opacity-100"
                  leave="transition duration-75 ease-in"
                  leaveFrom="transform scale-100 opacity-100"
                  leaveTo="transform scale-95 opacity-0"
                >
                  <Listbox.Options className="absolute z-10 mt-1 w-full rounded-xl bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                    {backupTypes.map((type) => (
                      <Listbox.Option
                        key={type}
                        value={type}
                        className={({ active }) =>
                          `cursor-pointer select-none px-4 py-2 capitalize ${active ? 'bg-sky-100 text-sky-700' : 'text-gray-800'
                          }`
                        }
                      >
                        {type}
                      </Listbox.Option>
                    ))}
                  </Listbox.Options>
                </Transition>
              </div>
            </div>
          </Listbox>

          {/* File Type Dropdown */}
          <Listbox
            value={form.fileType}
            onChange={(value) => setForm({ ...form, fileType: value })}
          >
            <div className="space-y-1">
              <Listbox.Label className="text-sm font-medium text-gray-700">File Type</Listbox.Label>
              <div className="relative">
                <Listbox.Button className="w-full cursor-pointer rounded-xl border border-gray-300 bg-white py-2 pl-4 pr-10 text-left shadow-sm focus:outline-none focus:ring-2 focus:ring-sky-500">
                  <span className="block">{form.fileType}</span>
                  <span className="absolute inset-y-0 right-0 flex items-center pr-3">
                    <FaAngleDown className="h-4 w-4 text-gray-400" />
                  </span>
                </Listbox.Button>
                <Transition
                  enter="transition duration-100 ease-out"
                  enterFrom="transform scale-95 opacity-0"
                  enterTo="transform scale-100 opacity-100"
                  leave="transition duration-75 ease-in"
                  leaveFrom="transform scale-100 opacity-100"
                  leaveTo="transform scale-95 opacity-0"
                >
                  <Listbox.Options className="absolute z-10 mt-1 w-full rounded-xl bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                    {fileTypes.map((type) => (
                      <Listbox.Option
                        key={type}
                        value={type}
                        className={({ active }) =>
                          `cursor-pointer select-none px-4 py-2 ${active ? 'bg-sky-100 text-sky-700' : 'text-gray-800'
                          }`
                        }
                      >
                        {type}
                      </Listbox.Option>
                    ))}
                  </Listbox.Options>
                </Transition>
              </div>
            </div>
          </Listbox>

          {/* Source Path */}
          <div className="col-span-2">
            <label className="text-sm font-medium text-gray-700">Source Path</label>
            <input
              type="text"
              value={form.sourcePath}
              onChange={(e) => setForm({ ...form, sourcePath: e.target.value })}
              className="w-full mt-1 px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-sky-500"
              required
            />
          </div>

          {/* Destination Path */}
          <div className="col-span-2">
            <label className="text-sm font-medium text-gray-700">Destination Path</label>
            <input
              type="text"
              value={form.destinationPath}
              onChange={(e) => setForm({ ...form, destinationPath: e.target.value })}
              className="w-full mt-1 px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-sky-500"
              required
            />
          </div>

          {/* Schedule Kind Dropdown */}
          <Listbox
            value={form.schedule?.kind}
            onChange={(value) =>
              setForm({
                ...form,
                schedule: { ...form.schedule, kind: value },
              })
            }
          >
            <div className="col-span-2 space-y-1">
              <Listbox.Label className="text-sm font-medium text-gray-700">Schedule</Listbox.Label>
              <div className="relative">
                <Listbox.Button className="w-full cursor-pointer rounded-xl border border-gray-300 bg-white py-2 pl-4 pr-10 text-left shadow-sm focus:outline-none focus:ring-2 focus:ring-sky-500">
                  <span className="block capitalize">{form.schedule?.kind}</span>
                  <span className="absolute inset-y-0 right-0 flex items-center pr-3">
                    <FaAngleDown className="h-4 w-4 text-gray-400" />
                  </span>
                </Listbox.Button>
                <Transition
                  enter="transition duration-100 ease-out"
                  enterFrom="transform scale-95 opacity-0"
                  enterTo="transform scale-100 opacity-100"
                  leave="transition duration-75 ease-in"
                  leaveFrom="transform scale-100 opacity-100"
                  leaveTo="transform scale-95 opacity-0"
                >
                  <Listbox.Options className="absolute z-10 mt-1 w-full rounded-xl bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                    {scheduleKinds.map((kind) => (
                      <Listbox.Option
                        key={kind}
                        value={kind}
                        className={({ active }) =>
                          `cursor-pointer select-none px-4 py-2 capitalize ${active ? 'bg-sky-100 text-sky-700' : 'text-gray-800'
                          }`
                        }
                      >
                        {kind}
                      </Listbox.Option>
                    ))}
                  </Listbox.Options>
                </Transition>
              </div>
            </div>
          </Listbox>

          {/* Buttons */}
          <div className="col-span-2 flex justify-end gap-4 mt-6">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-sky-500 text-white rounded-xl hover:bg-sky-600 transition-all"
            >
              Update Backup
            </button>
          </div>
        </form>
      </motion.div>
    </motion.div>
  );
};

export default UpdateModel;
