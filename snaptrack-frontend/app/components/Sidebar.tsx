'use client';

import { useRouter } from 'next/navigation';
import { motion, AnimatePresence } from 'framer-motion';
import {
  FaServer, FaDatabase, FaCogs, FaChartLine, FaFileAlt,
  FaSlidersH, FaSignOutAlt, FaTachometerAlt
} from 'react-icons/fa';

const navItems = [
  { icon: <FaTachometerAlt />, label: 'Dashboard', path: '/' },
  { icon: <FaServer />, label: 'Deployments', path: '/main/deployments' },
  { icon: <FaDatabase />, label: 'Backups', path: '/main/backups' },
  { icon: <FaChartLine />, label: 'Monitoring', path: '/main/monitoring' },
  { icon: <FaFileAlt />, label: 'Logs', path: '/main/logs' },
  { icon: <FaSlidersH />, label: 'Settings', path: '/main/settings' },
];

export default function Sidebar({ isOpen, toggleSidebar }: { isOpen: boolean; toggleSidebar: () => void }) {
  const router = useRouter();

  const handleLogout = () => {
    router.push('/logout');
  };

  return (
    <AnimatePresence>
      <motion.aside
        initial={{ width: isOpen ? 240 : 72 }}
        animate={{ width: isOpen ? 240 : 72 }}
        transition={{ duration: 0.3 }}
        className="fixed top-16 left-0 h-[calc(100vh-4rem)] bg-white text-gray-800 border-r z-10 shadow-sm flex flex-col justify-between"
      >
        <nav className="flex flex-col space-y-2 p-4">
          {navItems.map((item, idx) => (
            <motion.div
              key={idx}
              onClick={() => router.push(item.path)}
              whileHover={{ scale: 1.03 }}
              whileTap={{ scale: 0.95 }}
              className="relative group flex items-center hover:bg-gray-100 p-3 rounded-xl cursor-pointer transition-colors"
            >
              <span className="text-xl text-gray-600">{item.icon}</span>
              {isOpen && <span className="ml-3 font-medium text-sm">{item.label}</span>}
              {!isOpen && (
                <span className="absolute left-full ml-2 whitespace-nowrap bg-black text-white text-xs px-2 py-1 rounded shadow-lg opacity-0 group-hover:opacity-100 transition-opacity">
                  {item.label}
                </span>
              )}
            </motion.div>
          ))}
        </nav>

        {/* Logout Button */}
        <div className="p-4">
          <motion.div
            onClick={handleLogout}
            whileHover={{ scale: 1.03 }}
            whileTap={{ scale: 0.95 }}
            className="relative group flex items-center hover:bg-gray-100 p-3 rounded-xl cursor-pointer transition-colors"
          >
            <span className="text-xl text-gray-600"><FaSignOutAlt /></span>
            {isOpen && <span className="ml-3 font-medium text-sm">Logout</span>}
            {!isOpen && (
              <span className="absolute left-full ml-2 whitespace-nowrap bg-black text-white text-xs px-2 py-1 rounded shadow-lg opacity-0 group-hover:opacity-100 transition-opacity">
                Logout
              </span>
            )}
          </motion.div>
        </div>
      </motion.aside>
    </AnimatePresence>
  );
}
