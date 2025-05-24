'use client';

import { FaBars } from 'react-icons/fa';
import { FaUserCircle } from 'react-icons/fa';
import { motion } from 'framer-motion';

export default function Header({ toggleSidebar }: { toggleSidebar: () => void }) {
  return (
    <header className="fixed top-0 left-0 w-full bg-white shadow-md z-20">
      <div className="flex items-center justify-between p-4  mx-auto">
        {/* Left: Sidebar toggle + Logo */}
        <div className="flex items-center space-x-4">
         
          <span className="text-2xl ml-2 font-semibold text-blue-600 tracking-tight">SnapTrack</span>
           <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={toggleSidebar}
            className="text-2xl text-gray-700"
          >
            <FaBars />
          </motion.button>
        </div>

        {/* Right: Profile Icon */}
        <motion.div
          whileHover={{ scale: 1.1 }}
          className="text-gray-700 cursor-pointer text-3xl"
        >
          <FaUserCircle />
        </motion.div>
      </div>
    </header>
  );
}
