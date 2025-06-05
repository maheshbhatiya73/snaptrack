"use client";

import { useSocket } from "@/app/context/SocketContext";
import React from "react";
import { motion, AnimatePresence } from "framer-motion";
import { FaMicrochip, FaMemory, FaSpinner, FaSyncAlt, FaHdd } from "react-icons/fa";

export default function ProcessListPage() {
  const { runningProcesses } = useSocket();

  const rowVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { opacity: 1, y: 0, transition: { duration: 0.3 } },
    exit: { opacity: 0, y: -20, transition: { duration: 0.2 } },
  };

  function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
}


  return (
    <div>
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="max-w-10xl mx-auto bg-white rounded-2xl shadow-xl p-6"
      >
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-4xl font-extrabold text-sky-900 tracking-tight">
            Running Processes
          </h1>
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            className="flex items-center gap-2 bg-sky-600 text-white px-5 py-2.5 rounded-full shadow-md hover:bg-sky-700 transition-all duration-300"
            onClick={() => window.location.reload()}
          >
            <FaSyncAlt className="text-lg" />
            <span className="font-medium">Refresh</span>
          </motion.button>
        </div>

        {!runningProcesses ? (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="flex flex-col items-center justify-center h-64"
          >
            <FaSpinner className="animate-spin text-sky-600 text-5xl mb-4" />
            <p className="text-sky-800 text-lg font-medium">Loading process data...</p>
          </motion.div>
        ) : (
          <div className="overflow-x-auto rounded-xl shadow-sm border border-sky-100">
            <table className="min-w-full table-auto bg-white">
              <thead>
                <tr className="bg-sky-600 text-white">
                  <th className="p-4 text-left font-semibold text-sm uppercase tracking-wide whitespace-nowrap">
                    PID
                  </th>

                  <th className="p-4 text-left font-semibold text-sm uppercase tracking-wide whitespace-nowrap">
                    <div className="inline-flex items-center gap-2 whitespace-nowrap">
                      <FaMicrochip className="text-lg" /> <span>Name</span>
                    </div>
                  </th>

                  <th className="p-4 text-left font-semibold text-sm uppercase tracking-wide whitespace-nowrap">
                    CPU %
                  </th>

                  <th className="p-4 text-left font-semibold text-sm uppercase tracking-wide whitespace-nowrap">
                    <div className="inline-flex items-center gap-2 whitespace-nowrap">
                      <FaMemory className="text-lg" /> <span>Memory %</span>
                    </div>
                  </th>

                  <th className="p-4 text-left font-semibold text-sm uppercase tracking-wide whitespace-nowrap">
                    <div className="inline-flex items-center gap-2 whitespace-nowrap">
                      <FaHdd className="text-lg" /> <span>Disk Read</span>
                    </div>
                  </th>

                  <th className="p-4 text-left font-semibold text-sm uppercase tracking-wide whitespace-nowrap">
                    <div className="inline-flex items-center gap-2 whitespace-nowrap">
                      <FaHdd className="text-lg" /> <span>Disk Write</span>
                    </div>
                  </th>
                </tr>
              </thead>

              <tbody>
                <AnimatePresence>
                  {runningProcesses.map((proc) => (
                    <motion.tr
                      key={proc.pid}
                      variants={rowVariants}
                      initial="hidden"
                      animate="visible"
                      exit="exit"
                      className="border-b border-sky-100 hover:bg-sky-50 transition-colors duration-200"
                    >
                      <td className="p-4 text-sky-900 font-medium">{proc.pid}</td>
                      <td className="p-4 text-sky-900 font-medium">{proc.name}</td>
                      <td className="p-4 text-sky-900">{proc.cpu_percent.toFixed(2)}</td>
                      <td className="p-4 text-sky-900">{proc.mem_percent.toFixed(2)}</td>
                      <td className="p-4 text-sky-900">{formatBytes(proc.read_bytes)}</td>
                      <td className="p-4 text-sky-900">{formatBytes(proc.write_bytes)}</td>
                    </motion.tr>
                  ))}
                </AnimatePresence>
              </tbody>
            </table>
          </div>
        )}
      </motion.div>
    </div>
  );
}
