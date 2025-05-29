"use client";
import { createContext, useContext, useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { FaCheckCircle, FaTimesCircle, FaExclamationTriangle } from 'react-icons/fa';

// Define toast types
type ToastType = 'success' | 'error' | 'info';
interface Toast {
  id: string;
  message: string;
  type: ToastType;
}

// Toast context
interface ToastContextType {
  addToast: (message: string, type: ToastType) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

// Toast provider
export const ToastProvider = ({ children }: { children: React.ReactNode }) => {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const addToast = (message: string, type: ToastType) => {
    const id = Math.random().toString(36).substr(2, 9);
    setToasts((prev) => [...prev, { id, message, type }]);
  };

  const removeToast = (id: string) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id));
  };

  return (
    <ToastContext.Provider value={{ addToast }}>
      <div className="fixed top-4 right-4 z-[1000] space-y-2">
        <AnimatePresence>
          {toasts.map((toast) => (
            <Toast key={toast.id} {...toast} onClose={() => removeToast(toast.id)} />
          ))}
        </AnimatePresence>
      </div>
      {children}
    </ToastContext.Provider>
  );
};

// Hook to use toast
export const useToast = () => {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
};

// Toast component
interface ToastProps extends Toast {
  onClose: () => void;
}

const Toast = ({ message, type, onClose }: ToastProps) => {
  const [progress, setProgress] = useState(100);

  useEffect(() => {
    const duration = 5000; // Toast duration in ms
    const interval = 50;
    const decrement = 100 / (duration / interval);

    const timer = setInterval(() => {
      setProgress((prev) => {
        if (prev <= 0) {
          clearInterval(timer);
          onClose();
          return 0;
        }
        return prev - decrement;
      });
    }, interval);

    return () => clearInterval(timer);
  }, [onClose]);

  const icons = {
    success: <FaCheckCircle className="text-green-500" />,
    error: <FaTimesCircle className="text-red-500" />,
    info: <FaExclamationTriangle className="text-yellow-500" />,
  };

  const bgColors = {
    success: 'bg-green-50 border-green-200',
    error: 'bg-red-50 border-red-200',
    info: 'bg-yellow-50 border-yellow-200',
  };

  return (
    <motion.div
      initial={{ opacity: 0, x: 100 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: 100 }}
      transition={{ duration: 0.3 }}
      className={`flex items-center p-4 rounded-lg shadow-lg border ${bgColors[type]} max-w-sm w-full`}
    >
      <div className="mr-3">{icons[type]}</div>
      <div className="flex-1 text-sm text-gray-800">{message}</div>
      <button onClick={onClose} className="ml-3 text-gray-500 hover:text-gray-700">
        <FaTimesCircle />
      </button>
      <div className="absolute bottom-0 left-0 h-1 bg-opacity-50 rounded-b-lg" style={{ width: `${progress}%`, backgroundColor: type === 'success' ? '#22c55e' : type === 'error' ? '#ef4444' : '#eab308' }} />
    </motion.div>
  );
};