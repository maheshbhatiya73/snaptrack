'use client';

import { motion, AnimatePresence } from 'framer-motion';
import { FaHome, FaChartBar, FaCog, FaSignOutAlt } from 'react-icons/fa';

const navItems = [
  { icon: <FaHome />, label: 'Home', path: '/' },
  { icon: <FaChartBar />, label: 'Analytics', path: '/analytics' },
  { icon: <FaCog />, label: 'Settings', path: '/settings' },
  { icon: <FaSignOutAlt />, label: 'Logout', path: '/logout' },
];

export default function Sidebar({ isOpen, toggleSidebar }: { isOpen: boolean; toggleSidebar: () => void }) {
  return (
    <AnimatePresence>
      <motion.aside
        initial={{ width: isOpen ? 200 : 80 }}
        animate={{ width: isOpen ? 200 : 80 }}
        transition={{ duration: 0.3 }}
        className="fixed top-16 left-0  h-[calc(100vh-4rem)] bg-white shadow-lg z-10"
      >
        <nav className="flex flex-col p-4  space-y-4">
          {navItems.map((item, index) => (
            <motion.div
              key={index}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              className="relative group flex items-center space-x-2 text-gray-700 hover:bg-blue-50 p-2 rounded-lg cursor-pointer"
            >
              <span className="text-xl">{item.icon}</span>
              {isOpen && <span className="font-medium">{item.label}</span>}
              {!isOpen && (
                <span className="tooltip">{item.label}</span>
              )}
            </motion.div>
          ))}
        </nav>
      </motion.aside>
    </AnimatePresence>
  );
}