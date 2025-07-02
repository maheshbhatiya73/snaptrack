"use client";
import { useSocket } from '@/app/context/SocketContext';
import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, ArrowUp, ArrowDown } from 'lucide-react';

const ProcessList = () => {
  const { runningProcesses } = useSocket();
  const [processes, setProcesses] = useState<any>([]);
  const [sortConfig, setSortConfig] = useState({ key: 'cpu_percent', direction: 'desc' });

  useEffect(() => {
    if (runningProcesses) {
      setProcesses(runningProcesses);
    }
  }, [runningProcesses]);

  const sortData = (key:any) => {
    let direction = 'desc';
    if (sortConfig.key === key && sortConfig.direction === 'desc') {
      direction = 'asc';
    }

    const sorted = [...processes].sort((a, b) => {
      if (a[key] < b[key]) return direction === 'asc' ? -1 : 1;
      if (a[key] > b[key]) return direction === 'asc' ? 1 : -1;
      return 0;
    });

    setProcesses(sorted);
    setSortConfig({ key, direction });
  };

  const getSortIcon = (key:any) => {
    if (sortConfig.key !== key) return <ArrowUpDown className="ml-2 h-4 w-4" />;
    return sortConfig.direction === 'asc' ? (
      <ArrowUp className="ml-2 h-4 w-4" />
    ) : (
      <ArrowDown className="ml-2 h-4 w-4" />
    );
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.5 }}
      className="min-h-screen p-6 font-sans"
    >
      <div className="max-w-8xl mx-auto">
        <motion.h1
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.5, delay: 0.2 }}
          className="text-3xl font-bold mb-8 text-green-400"
        >
          System Processes
        </motion.h1>
        <div className="bg-gray-800 rounded-xl shadow-2xl overflow-hidden border border-gray-700">
          <Table>
            <TableHeader>
              <TableRow className="bg-gray-900 hover:bg-gray-900">
                {[
                  { key: 'pid', label: 'PID' },
                  { key: 'name', label: 'Name' },
                  { key: 'cpu_percent', label: 'CPU %' },
                  { key: 'mem_percent', label: 'Memory %' },
                  { key: 'status', label: 'Status' },
                  { key: 'read_bytes', label: 'Read Bytes' },
                  { key: 'write_bytes', label: 'Write Bytes' },
                ].map((column) => (
                  <TableHead key={column.key}>
                    <Button
                      variant="ghost"
                      onClick={() => sortData(column.key)}
                      className="text-green-400 hover:text-green-300 hover:bg-gray-700 transition-colors duration-200"
                    >
                      {column.label}
                      {getSortIcon(column.key)}
                    </Button>
                  </TableHead>
                ))}
              </TableRow>
            </TableHeader>
            <TableBody>
              <AnimatePresence>
                {processes.map((process:any, index:number) => (
                  <motion.tr
                    key={process.pid}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    transition={{ duration: 0.3, delay: index * 0.05 }}
                    className={`border-b border-gray-700 ${
                      index % 2 === 0 ? 'bg-gray-800' : 'bg-gray-850'
                    } hover:bg-gray-700 transition-colors duration-200`}
                  >
                    <TableCell className="text-green-400">{process.pid}</TableCell>
                    <TableCell className="text-green-400 font-medium">{process.name}</TableCell>
                    <TableCell className="text-green-400">{process.cpu_percent.toFixed(2)}</TableCell>
                    <TableCell className="text-green-400">{process.mem_percent.toFixed(2)}</TableCell>
                    <TableCell className="text-green-400 capitalize">{process.status}</TableCell>
                    <TableCell className="text-green-400">{process.read_bytes.toLocaleString()}</TableCell>
                    <TableCell className="text-green-400">{process.write_bytes.toLocaleString()}</TableCell>
                  </motion.tr>
                ))}
              </AnimatePresence>
            </TableBody>
          </Table>
        </div>
      </div>
    </motion.div>
  );
};

export default ProcessList;