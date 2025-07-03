"use client";
import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, ArrowUp, ArrowDown } from 'lucide-react';
import { useSocket } from '@/app/context/SocketContext';

const PortsList = () => {
  const { runningPorts } = useSocket();
  const [ports, setPorts] = useState([]);
  const [sortConfig, setSortConfig] = useState({ key: 'port', direction: 'asc' });

  useEffect(() => {
    if (runningPorts && runningPorts) {
      setPorts(runningPorts);
    }
  }, [runningPorts]);

  const sortData = (key) => {
    let direction = 'asc';
    if (sortConfig.key === key && sortConfig.direction === 'asc') {
      direction = 'desc';
    }

    const sorted = [...ports].sort((a, b) => {
      if (a[key] < b[key]) return direction === 'asc' ? -1 : 1;
      if (a[key] > b[key]) return direction === 'asc' ? 1 : -1;
      return 0;
    });

    setPorts(sorted);
    setSortConfig({ key, direction });
  };

  const getSortIcon = (key) => {
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
      className=" p-6 font-sans"
    >
      <div className="max-w-8xl mx-auto">
        <motion.h1
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.5, delay: 0.2 }}
          className="text-3xl font-bold mb-8 text-green-400"
        >
          Running Ports
        </motion.h1>
        <div className="bg-gray-800 rounded-xl shadow-2xl overflow-hidden border border-gray-700">
          <Table>
            <TableHeader>
              <TableRow className="bg-gray-900 hover:bg-gray-900">
                {[
                  { key: 'protocol', label: 'Protocol' },
                  { key: 'port', label: 'Port' },
                  { key: 'process', label: 'Process' },
                  { key: 'pid', label: 'PID' },
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
                {ports.map((port, index) => (
                  <motion.tr
                    key={`${port.port}-${port.pid}`}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    transition={{ duration: 0.3, delay: index * 0.05 }}
                    className={`border-b border-gray-700 ${
                      index % 2 === 0 ? 'bg-gray-800' : 'bg-gray-850'
                    } hover:bg-gray-700 transition-colors duration-200`}
                  >
                    <TableCell className="text-green-400 capitalize">{port.protocol}</TableCell>
                    <TableCell className="text-green-400">{port.port}</TableCell>
                    <TableCell className="text-green-400">{port.process || 'N/A'}</TableCell>
                    <TableCell className="text schwer-green-400">{port.pid || 'N/A'}</TableCell>
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

export default PortsList;