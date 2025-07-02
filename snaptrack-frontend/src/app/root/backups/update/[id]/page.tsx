'use client';
import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { motion, AnimatePresence } from 'framer-motion';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Folder, Calendar, Send, FileArchive, Clock, AlertCircle } from 'lucide-react';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { useAuth } from '@/app/store/useAuth';
import { getBackupById, updateBackup } from '@/lib/api';
import { useLinuxToast } from '@/lib/use-linux-toast';

export default function UpdateBackupPage() {
  const { token } = useAuth();
  const { id } = useParams();
  const router = useRouter();
  const [formData, setFormData] = useState({
    app: '',
    type: 'manual',
    sourcePath: '',
    destinationPath: '',
    fileType: 'zip',
    scheduleKind: 'one-time',
    scheduleDate: new Date(),
  });
  const [errors, setErrors] = useState({
    app: '',
    sourcePath: '',
    destinationPath: '',
    scheduleDate: '',
  });
  const [isLoading, setIsLoading] = useState(true);
  const [fetchError, setFetchError] = useState(null);
  const { success, error } = useLinuxToast();


  // Fetch backup data by ID
  useEffect(() => {
    async function fetchBackup() {
      try {
        const response = await getBackupById(id);
        const backup = response;
        setFormData({
          app: backup.app,
          type: backup.type,
          sourcePath: backup.sourcePath,
          destinationPath: backup.destinationPath,
          fileType: backup.fileType,
          scheduleKind: backup.schedule.kind,
          scheduleDate: new Date(backup.schedule.date),
        });
        setIsLoading(false);
      } catch (err:any) {
        setFetchError(err.message || 'Failed to fetch backup');
        setIsLoading(false);
      }
    }
    if (id && token) {
      fetchBackup();
    }
  }, [id, token]);

  // Handle input changes
  const handleInputChange = (e:any) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    setErrors((prev) => ({ ...prev, [name]: '' }));
  };

  // Handle select changes
  const handleSelectChange = (name:any, value:any) => {
    setFormData((prev) => ({ ...prev, [name]: value }));
    setErrors((prev) => ({ ...prev, [name]: '' }));
  };

  // Handle date change
  const handleDateChange = (date:any) => {
    setFormData((prev) => ({ ...prev, scheduleDate: date }));
    setErrors((prev) => ({ ...prev, scheduleDate: '' }));
  };

  // Manual validation
  const validateForm = () => {
    const newErrors = {
      app: '',
      sourcePath: '',
      destinationPath: '',
      scheduleDate: '',
    };
    let isValid = true;

    if (!formData.app) {
      newErrors.app = 'App name is required';
      isValid = false;
    } else if (formData.app.length > 50) {
      newErrors.app = 'App name must be 50 characters or less';
      isValid = false;
    }

    if (!formData.sourcePath) {
      newErrors.sourcePath = 'Source path is required';
      isValid = false;
    } else if (!/^\/[a-zA-Z0-9\/_-]+$/.test(formData.sourcePath)) {
      newErrors.sourcePath = 'Invalid path format';
      isValid = false;
    }

    if (!formData.destinationPath) {
      newErrors.destinationPath = 'Destination path is required';
      isValid = false;
    } else if (!/^\/[a-zA-Z0-9\/_-]+$/.test(formData.destinationPath)) {
      newErrors.destinationPath = 'Invalid path format';
      isValid = false;
    }

    if (formData.scheduleKind === 'one-time' && formData.scheduleDate < new Date()) {
      newErrors.scheduleDate = 'Schedule date must be in the future';
      isValid = false;
    }

    setErrors(newErrors);
    return isValid;
  };

  // Handle form submission
 const onSubmit = async (e:any) => {
  e.preventDefault();

  if (validateForm()) {
    try {
      const backupData = {
        app: formData.app,
        type: formData.type,
        sourcePath: formData.sourcePath,
        destinationPath: formData.destinationPath,
        fileType: formData.fileType,
        schedule: {
          kind: formData.scheduleKind,
          date: formData.scheduleDate.toISOString(),
        },
      };

      await updateBackup(id, backupData);
      success('Backup updated successfully');
      router.push('/root/backups');
    } catch (err: any) {
      error(err.message || 'Failed to update backup');
      setErrors((prev) => ({
        ...prev,
        general: err.message || 'Failed to update backup',
      }));
    }
  }
};

  if (isLoading) {
    return (
      <div className="min-h-screen text-white p-4 sm:p-8">
        <p>Loading backup...</p>
      </div>
    );
  }

  if (fetchError) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black text-white p-4 sm:p-8">
        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="text-red-400 flex items-center gap-2"
        >
          <AlertCircle className="w-5 h-5" />
          {fetchError}
        </motion.p>
      </div>
    );
  }

  return (
    <div className="p-4 sm:p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Card className="bg-gray-800/80 border-gray-700/50 shadow-xl hover:shadow-green-500/30 transition-shadow duration-300 backdrop-blur-sm">
          <CardHeader>
            <CardTitle className="flex items-center gap-3 text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-green-400 to-cyan-300">
              <FileArchive className="w-8 h-8 text-green-400" />
              Update Backup
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
                  onValueChange={(value) => handleSelectChange('type', value)}
                >
                  <SelectTrigger className="bg-gray-900/50 border-gray-700 text-white focus:ring-2 focus:ring-green-500">
                    <SelectValue placeholder="Select type" />
                  </SelectTrigger>
                  <SelectContent className="bg-gray-900 border-gray-700 text-white">
                    <SelectItem value="manual">Manual</SelectItem>
                    <SelectItem value="automatic">Automatic</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Source Path */}
              <div className="space-y-2">
                <Label htmlFor="sourcePath" className="text-green-400 font-medium flex items-center gap-2">
                  <Folder className="w-5 h-5" /> Source Path
                </Label>
                <Input
                  id="sourcePath"
                  name="sourcePath"
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
                  <FileArchive className="w-5 h-5" /> File Type
                </Label>
                <Select
                  value={formData.fileType}
                  onValueChange={(value) => handleSelectChange('fileType', value)}
                >
                  <SelectTrigger className="bg-gray-900/50 border-gray-700 text-white focus:ring-2 focus:ring-green-500">
                    <SelectValue placeholder="Select file type" />
                  </SelectTrigger>
                  <SelectContent className="bg-gray-900 border-gray-700 text-white">
                    <SelectItem value="zip">ZIP</SelectItem>
                    <SelectItem value="tar.gz">TAR.GZ</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Schedule */}
              <div className="space-y-2">
                <Label htmlFor="scheduleKind" className="text-green-400 font-medium flex items-center gap-2">
                  <Calendar className="w-5 h-5" /> Schedule
                </Label>
                <Select
                  value={formData.scheduleKind}
                  onValueChange={(value) => handleSelectChange('scheduleKind', value)}
                >
                  <SelectTrigger className="bg-gray-900/50 border-gray-700 text-white focus:ring-2 focus:ring-green-500">
                    <SelectValue placeholder="Select schedule" />
                  </SelectTrigger>
                  <SelectContent className="bg-gray-900 border-gray-700 text-white">
                    <SelectItem value="one-time">One-Time</SelectItem>
                    <SelectItem value="hourly">Hourly</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Schedule Date */}
              {formData.scheduleKind === 'one-time' && (
                <div className="space-y-2 md:col-span-2">
                  <Label htmlFor="scheduleDate" className="text-green-400 font-medium flex items-center gap-2">
                    <Clock className="w-5 h-5" /> Schedule Date/Time
                  </Label>
                  <DatePicker
                    selected={formData.scheduleDate}
                    onChange={handleDateChange}
                    showTimeSelect
                    dateFormat="Pp"
                    minDate={new Date()}
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
                </div>
              )}

              {/* Submit Button */}
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
                    Update Backup
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
