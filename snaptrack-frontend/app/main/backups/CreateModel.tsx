import { useState } from 'react';
import { motion } from 'framer-motion';
import { Dialog, Listbox } from '@headlessui/react';
import { Backup, createBackup, Schedule } from '@/app/lib/api';
import { IoCheckmarkCircleOutline } from 'react-icons/io5';
import { FaAngleDown } from 'react-icons/fa';

interface CreateModelProps {
  onClose: () => void;
  onSuccess: () => void;
  onError: () => void;
  token: string;
}

const types: Backup['type'][] = ['manual', 'full', 'incremental'];
const fileTypes: Backup['fileType'][] = ['zip', 'tar', 'tar.gz'];
const scheduleKinds: Schedule['kind'][] = ['one-time', 'hourly', 'daily', 'weekly', 'monthly'];

const ListboxSelect = <T extends string>({
  label,
  value,
  options,
  onChange,
}: {
  label: string;
  value: T;
  options: T[];
  onChange: (val: T) => void;
}) => (
  <div className="w-full">
    <label className="text-sm font-medium text-gray-700 mb-1 block">{label}</label>
    <Listbox value={value} onChange={onChange}>
      <div className="relative mt-1">
        <Listbox.Button className="w-full cursor-pointer rounded-xl border border-gray-300 bg-white py-2 pl-3 pr-10 text-left shadow-sm focus:outline-none focus:ring-2 focus:ring-sky-400 focus:border-sky-400 sm:text-sm">
          <span className="block truncate capitalize">{value}</span>
          <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
            <FaAngleDown className="h-5 w-5 text-gray-400" aria-hidden="true" />
          </span>
        </Listbox.Button>
        <Listbox.Options className="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-xl bg-white py-1 text-base shadow-lg ring-1 ring-black/10 focus:outline-none sm:text-sm">
          {options.map((option) => (
            <Listbox.Option
              key={option}
              value={option}
              className={({ active }) =>
                `relative cursor-pointer select-none py-2 pl-10 pr-4 ${
                  active ? 'bg-sky-100 text-sky-900' : 'text-gray-900'
                }`
              }
            >
              {({ selected }) => (
                <>
                  <span className={`block truncate capitalize ${selected ? 'font-medium' : 'font-normal'}`}>
                    {option}
                  </span>
                  {selected && (
                    <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-sky-600">
                      <IoCheckmarkCircleOutline className="h-5 w-5" aria-hidden="true" />
                    </span>
                  )}
                </>
              )}
            </Listbox.Option>
          ))}
        </Listbox.Options>
      </div>
    </Listbox>
  </div>
);

const CreateModel = ({ onClose, onSuccess, onError, token }: CreateModelProps) => {
  const [form, setForm] = useState<Partial<Backup>>({
    app: '',
    type: 'manual',
    sourcePath: '',
    destinationPath: '',
    fileType: 'tar.gz',
    schedule: { kind: 'one-time' },
  });

  const handleSubmit = async (e: React.FormEvent) => {
  e.preventDefault();
  const response = await createBackup(form, token);
  if (response.success && response.data) {
    onSuccess(response.data); // Pass created backup
  } else {
    onError(response.message); // Let parent handle error toast
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
        initial={{ scale: 0.9, y: 30 }}
        animate={{ scale: 1, y: 0 }}
        className="bg-white w-full max-w-2xl mx-4 rounded-2xl shadow-2xl p-8 relative"
      >
        <Dialog.Title className="text-2xl font-semibold text-gray-800 mb-6">🛡️ Create New Backup</Dialog.Title>
        <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="col-span-2">
            <label className="text-sm font-medium text-gray-700">App</label>
            <input
              type="text"
              value={form.app}
              onChange={(e) => setForm({ ...form, app: e.target.value })}
              className="w-full mt-1 px-3 py-2 border border-gray-300 rounded-xl  focus:outline-none focus:ring-2 focus:ring-sky-400 focus:border-sky-400"
              required
            />
          </div>

          <ListboxSelect
            label="Backup Type"
            value={form.type!}
            options={types}
            onChange={(val) => setForm({ ...form, type: val })}
          />

          <ListboxSelect
            label="File Format"
            value={form.fileType!}
            options={fileTypes}
            onChange={(val) => setForm({ ...form, fileType: val })}
          />

          <div className="col-span-2">
            <label className="text-sm font-medium text-gray-700">Source Path</label>
            <input
              type="text"
              value={form.sourcePath}
              onChange={(e) => setForm({ ...form, sourcePath: e.target.value })}
              className="w-full mt-1 px-3 py-2 border border-gray-300 rounded-xl  focus:outline-none focus:ring-2 focus:ring-sky-400 focus:border-sky-400"
              required
            />
          </div>

          <div className="col-span-2">
            <label className="text-sm font-medium text-gray-700">Destination Path</label>
            <input
              type="text"
              value={form.destinationPath}
              onChange={(e) => setForm({ ...form, destinationPath: e.target.value })}
              className="w-full mt-1 px-3 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-sky-400 focus:border-sky-400"

              required
            />
          </div>

          <ListboxSelect
            label="Schedule Type"
            value={form.schedule?.kind!}
            options={scheduleKinds}
            onChange={(val) => setForm({ ...form, schedule: { ...form.schedule, kind: val } })}
          />

          <div className="col-span-2 flex justify-end mt-6 space-x-3">
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
              Create Backup
            </button>
          </div>
        </form>
      </motion.div>
    </motion.div>
  );
};

export default CreateModel;
