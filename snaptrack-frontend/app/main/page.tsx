'use client';

import { motion } from 'framer-motion';
import { useAuth } from '../context/AuthContext';

export default function MainPage() {
  const { setToken } = useAuth();

  const handleLogout = () => {
    setToken(null);
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className=""
    >
      <div className="bg-white rounded-2xl shadow-xl p-8 mx-auto">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Server Dashboard</h1>
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            onClick={handleLogout}
            className="py-2 px-4 gradient-bg text-white rounded-lg hover:shadow-lg transition duration-300"
          >
            Logout
          </motion.button>
        </div>
        <p className="text-gray-600 mb-6">Welcome, admin! Manage your server resources below.</p>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.2 }}
            className="p-6 bg-gray-100 rounded-xl shadow-sm hover:shadow-md transition-shadow"
          >
            <h2 className="text-xl font-semibold text-gray-800">System Status</h2>
            <p className="text-gray-600">All systems operational</p>
          </motion.div>
          <motion.div
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.3 }}
            className="p-6 bg-gray-100 rounded-xl shadow-sm hover:shadow-md transition-shadow"
          >
            <h2 className="text-xl font-semibold text-gray-800">Recent Activity</h2>
            <p className="text-gray-600">No recent activity</p>
          </motion.div>
        </div>
      </div>
    </motion.div>
  );
}