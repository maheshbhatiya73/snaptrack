'use client';
import { useState, ChangeEvent, FormEvent, JSX } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Folder, Calendar, Send, FileArchive, Clock, AlertCircle } from 'lucide-react';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { createBackup } from '@/lib/api';
import { useLinuxToast } from '@/lib/use-linux-toast';
import { useRouter } from 'next/navigation';

// Type definitions
type BackupType = 'manual' | 'full' | 'incremental';
type FileType = 'zip' | 'tar' | 'tar.gz';
type ScheduleKind = 'one-time' | 'hourly' | 'daily' | 'weekly' | 'monthly';

interface Schedule {
  kind: ScheduleKind;
  date: string;
}

interface FormData {
  app: string;
  type: BackupType;
  sourcePath: string;
  destinationPath: string;
  fileType: FileType;
  schedule: Schedule;
  runNow?: boolean;
}

interface FormErrors {
  app: string;
  sourcePath: string;
  destinationPath: string;
  scheduleDate: string;
}

type FormFieldNames = keyof FormErrors;

export default function CreateBackupPage(): JSX.Element {
  // Form state with proper typing
  const [formData, setFormData] = useState<FormData>({
    app: '',
    type: 'manual',
    sourcePath: '',
    destinationPath: '',
    fileType: 'zip',
    schedule: {
      kind: 'one-time',
      date: new Date().toISOString(),
    },
    runNow: false,
  });

  // Error state with proper typing
  const [errors, setErrors] = useState<FormErrors>({
    app: '',
    sourcePath: '',
    destinationPath: '',
    scheduleDate: '',
  });
  const { success, error } = useLinuxToast();
  const router = useRouter();

  // Handle input changes with proper typing
  const handleInputChange = (e: ChangeEvent<HTMLInputElement>): void => {
    const { name, value } = e.target;
    const fieldName = name as FormFieldNames;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    // Clear error when user starts typing
    setErrors((prev) => ({ ...prev, [fieldName]: '' }));
  };

  // Handle select changes with proper typing
  const handleSelectChange = <T extends keyof FormData>(name: T, value: FormData[T]): void => {
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));

    // Clear error if the field has an error state
    if (name in errors) {
      setErrors((prev) => ({ ...prev, [name as FormFieldNames]: '' }));
    }
  };

  const handleScheduleKindChange = (kind: ScheduleKind) => {
    setFormData((prev) => ({
      ...prev,
      schedule: {
        ...prev.schedule,
        kind,
      },
    }));
  };

  // Handle date change with proper typing
  const handleDateChange = (date: Date | null) => {
    if (date) {
      setFormData((prev) => ({
        ...prev,
        schedule: {
          ...prev.schedule,
          date: date.toISOString(),
        },
      }));
      setErrors((prev) => ({ ...prev, scheduleDate: '' }));
    }
  };

  const handleRunNowChange = (e: ChangeEvent<HTMLInputElement>) => {
    const checked = e.target.checked;
    setFormData((prev) => ({
      ...prev,
      runNow: checked,
    }));
  };

  // Manual validation with proper return type
  const validateForm = (): boolean => {
    const newErrors: FormErrors = {
      app: '',
      sourcePath: '',
      destinationPath: '',
      scheduleDate: '',
    };
    let isValid = true;

    if (!formData.app) {
      newErrors.app = 'App name is required';
      isValid = false;
    }

    if (!formData.sourcePath) {
      newErrors.sourcePath = 'Source path is required';
      isValid = false;
    }

    if (!formData.destinationPath) {
      newErrors.destinationPath = 'Destination path is required';
      isValid = false;
    }

    if (
      formData.schedule.kind === 'one-time' &&
      !formData.runNow &&
      new Date(formData.schedule.date) < new Date()
    ) {
      newErrors.scheduleDate = 'Schedule date must be in the future';
      isValid = false;
    }

    setErrors(newErrors);
    return isValid;
  };

  

  // Handle form submission with proper typing
  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (validateForm()) {
      console.log('Form submitted:', formData);
      const initialFormData: FormData = {
        app: '',
        type: 'manual',
        sourcePath: '',
        destinationPath: '',
        fileType: 'zip',
        schedule: {
          kind: 'one-time',
          date: new Date().toISOString(),
        },
        runNow: false,
      };

      const initialErrors: FormErrors = {
        app: '',
        sourcePath: '',
        destinationPath: '',
        scheduleDate: '',
      };

      try {
        const response = await createBackup(formData);
        success('Backup created successfully');
        setFormData(initialFormData);
        setErrors(initialErrors);
        router.push('/root/backups');
      } catch (err: any) {
        error(err.message || 'Unexpected error occurred');
      }
    }
  };

  return (
    <div className="text-white p-4 sm:p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Card className="bg-gray-800/80 border-gray-700/50 shadow-xl hover:shadow-green-500/30 transition-shadow duration-300 backdrop-blur-sm">
          <CardHeader>
            <CardTitle className="flex items-center gap-3 text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-green-400 to-cyan-300">
              <FileArchive className="w-8 h-8 text-green-400" />
              Create New Backup
            </CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={onSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* App Name */}
              <div className="space-y-2">
                <Label htmlFor="app" className="text-green-400 font-medium flex items-center gap-2">
                  <Folder className="w-5 h-5" /> App Name
                </Label>
                <Input
                  id="app"
                  name="app"
                  type="text"
                  value={formData.app}
                  onChange={handleInputChange}
                  placeholder="e.g., My App Backup"
                  className="bg-gray-900/50 border-gray-700 text-white placeholder-gray-500 focus:ring-2 focus:ring-green-500 transition-all duration-200"
                />
                <AnimatePresence>
                  {errors.app && (
                    <motion.p
                      initial={{ opacity: 0, y: -10 }}
                      animate={{ opacity: 1, y: 0 }}
                      exit={{ opacity: 0, y: -10 }}
                      className="text-red-400 text-sm flex items-center gap-1"
                    >
                      <AlertCircle className="w-4 h-4" /> {errors.app}
                    </motion.p>
                  )}
                </AnimatePresence>
              </div>

              {/* Backup Type */}
              <div className="space-y-2">
                <Label htmlFor="type" className="text-green-400 font-medium flex items-center gap-2">
                  <FileArchive className="w-5 h-5" /> Backup Type
                </Label>
                <Select
                  value={formData.type}
                  onValueChange={(value: BackupType) => handleSelectChange('type', value)}
                >
                  <SelectTrigger className="bg-gray-900/50 border-gray-700 text-white focus:ring-2 focus:ring-green-500">
                    <SelectValue placeholder="Select type" />
                  </SelectTrigger>
                  <SelectContent className="bg-gray-900 border-gray-700 text-white">
                    <SelectItem value="manual">Manual</SelectItem>
                    <SelectItem value="full">Full</SelectItem>
                    <SelectItem value="incremental">Incremental</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Source Path */}
              {/* Source Path */}
              <div className="space-y-2">
                <Label htmlFor="sourcePath" className="text-green-400 font-medium flex items-center gap-2">
                  <Folder className="w-5 h-5" /> Source Path
                </Label>
                <Input
                  id="sourcePath"
                  name="sourcePath"
                  type="text"
                  value={formData.sourcePath}
                  onChange={handleInputChange}
                  placeholder="/home/user/data"
                  className="bg-gray-900/50 border-gray-700 text-white placeholder-gray-500 focus:ring-2 focus:ring-green-500"
                />
                <AnimatePresence>
                  {errors.sourcePath && (
                    <motion.p
                      initial={{ opacity: 0, y: -10 }}
                      animate={{ opacity: 1, y: 0 }}
                      exit={{ opacity: 0, y: -10 }}
                      className="text-red-400 text-sm flex items-center gap-1"
                    >
                      <AlertCircle className="w-4 h-4" /> {errors.sourcePath}
                    </motion.p>
                  )}
                </AnimatePresence>
              </div>

              {/* Destination Path */}
              <div className="space-y-2">
                <Label htmlFor="destinationPath" className="text-green-400 font-medium flex items-center gap-2">
                  <Folder className="w-5 h-5" /> Destination Path
                </Label>
                <Input
                  id="destinationPath"
                  name="destinationPath"
                  type="text"
                  value={formData.destinationPath}
                  onChange={handleInputChange}
                  placeholder="/home/user/backups"
                  className="bg-gray-900/50 border-gray-700 text-white placeholder-gray-500 focus:ring-2 focus:ring-green-500"
                />
                <AnimatePresence>
                  {errors.destinationPath && (
                    <motion.p
                      initial={{ opacity: 0, y: -10 }}
                      animate={{ opacity: 1, y: 0 }}
                      exit={{ opacity: 0, y: -10 }}
                      className="text-red-400 text-sm flex items-center gap-1"
                    >
                      <AlertCircle className="w-4 h-4" /> {errors.destinationPath}
                    </motion.p>
                  )}
                </AnimatePresence>
              </div>


              {/* File Type */}
              <div className="space-y-2">
                <Label htmlFor="fileType" className="text-green-400 font-medium flex items-center gap-2">
                  <FileArchive className="w-5 h-5" /> File ShakespeareType
                </Label>
                <Select
                  value={formData.fileType}
                  onValueChange={(value: FileType) => handleSelectChange('fileType', value)}
                >
                  <SelectTrigger className="bg-gray-900/50 border-gray-700 text-white focus:ring-2 focus:ring-green-500">
                    <SelectValue placeholder="Select file type" />
                  </SelectTrigger>
                  <SelectContent className="bg-gray-900 border-gray-700 text-white">
                    <SelectItem value="zip">ZIP</SelectItem>
                    <SelectItem value="tar">TAR</SelectItem>
                    <SelectItem value="tar.gz">TAR.GZ</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Schedule */}
              <div className="space-y-2">
                <Label htmlFor="schedule.kind" className="text-green-400 font-medium flex items-center gap-2">
                  <Calendar className="w-5 h-5" /> Schedule
                </Label>
                <Select
                  value={formData.schedule.kind}
                  onValueChange={handleScheduleKindChange}
                >
                  <SelectTrigger className="bg-gray-900/50 border-gray-700 text-white focus:ring-2 focus:ring-green-500">
                    <SelectValue placeholder="Select schedule" />
                  </SelectTrigger>
                  <SelectContent className="bg-gray-900 border-gray-700 textunionfs">
                    <SelectItem value="one-time">One-Time</SelectItem>
                    <SelectItem value="hourly">Hourly</SelectItem>
                    <SelectItem value="daily">Daily</SelectItem>
                    <SelectItem value="weekly">Weekly</SelectItem>
                    <SelectItem value="monthly">Monthly</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Schedule Date */}
              {formData.schedule.kind === 'one-time' && (
                <div className="space-y-2 md:col-span-2">
                  {/* Schedule Date Picker */}
                  <Label htmlFor="scheduleDate" className="text-green-400 font-medium flex items-center gap-2">
                    <Clock className="w-5 h-5" /> Schedule Date/Time
                  </Label>
                  <DatePicker
                    selected={new Date(formData.schedule.date)}
                    onChange={handleDateChange}
                    showTimeSelect
                    dateFormat="Pp"
                    minDate={new Date()}
                    disabled={formData.runNow === true}
                    className="w-full bg-gray-900/50 border-gray-700 text-white rounded-md p-2 focus:ring-2 focus:ring-green-500 transition-all duration-200"
                  />
                  <AnimatePresence>
                    {errors.scheduleDate && (
                      <motion.p
                        initial={{ opacity: 0, y: -10 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -10 }}
                        className="text-red-400 text-sm flex items-center gap-1"
                      >
                        <AlertCircle className="w-4 h-4" /> {errors.scheduleDate}
                      </motion.p>
                    )}
                  </AnimatePresence>

                  {/* Run Now Checkbox */}
                  <div className="flex items-center space-x-2 pt-2">
                    <input
                      type="checkbox"
                      id="runNow"
                      checked={formData.runNow || false}
                      onChange={handleRunNowChange}
                      className="w-4 h-4 text-green-600 focus:ring-green-500 border-gray-300 rounded"
                    />
                    <Label htmlFor="runNow" className="text-sm text-gray-300">
                      Run backup immediately instead of scheduled time
                    </Label>
                  </div>
                </div>
              )}

              <div className="md:col-span-2 w-40">
                <motion.div
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  className="mt-4"
                >
                  <Button
                    type="submit"
                    className="w-full bg-gradient-to-r from-green-500 to-cyan-500 hover:from-green-600 hover:to-cyan-600 text-white font-semibold py-3 rounded-lg transition-all duration-200"
                  >
                    <Send className="w-5 h-5 mr-2" />
                    Create Backup
                  </Button>
                </motion.div>
              </div>
            </form>
          </CardContent>
        </Card>
      </motion.div>
    </div>
  );
}