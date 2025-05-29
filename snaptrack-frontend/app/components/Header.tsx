'use client';

import { motion } from 'framer-motion';
import { FaBars, FaUserCircle } from 'react-icons/fa';

export default function Header({ toggleSidebar }: { toggleSidebar: () => void }) {
  return (
    <header className="fixed top-0 left-0 w-full bg-white shadow z-20">
      <div className="flex items-center justify-between px-6 py-4">
        <div className="flex items-center space-x-4">
          <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={toggleSidebar}
            className="text-xl text-gray-700"
          >
            <FaBars />
          </motion.button>
          <span className="text-xl font-bold text-blue-600 tracking-tight">ServerDeck</span>
        </div>
        <motion.div
          whileHover={{ scale: 1.1 }}
          className="text-gray-700 cursor-pointer text-2xl"
        >
          <FaUserCircle />
        </motion.div>
      </div>
    </header>
  );
}
