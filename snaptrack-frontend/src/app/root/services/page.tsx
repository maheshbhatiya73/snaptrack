"use client";
import React, { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { ArrowUpDown, ArrowUp, ArrowDown } from 'lucide-react';
import { useSocket } from '@/app/context/SocketContext';

const ServicesList = () => {
    const [services, setServices] = useState<any>([]);
    const [sortConfig, setSortConfig] = useState({ key: 'name', direction: 'asc' });
    const { services: liveServices } = useSocket();
    useEffect(() => {
        if (liveServices) {
            setServices(liveServices);
        }
    }, [liveServices]);

    const sortData = (key:any) => {
        let direction = 'asc';
        if (sortConfig.key === key && sortConfig.direction === 'asc') {
            direction = 'desc';
        }

        const sorted = [...services].sort((a, b) => {
            if (a[key] < b[key]) return direction === 'asc' ? -1 : 1;
            if (a[key] > b[key]) return direction === 'asc' ? 1 : -1;
            return 0;
        });

        setServices(sorted);
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

    const getStatusColor = (status:any) => {
        switch (status.toLowerCase()) {
            case 'active':
                return 'text-green-400';
            case 'inactive':
                return 'text-gray-400';
            case 'failed':
                return 'text-red-400';
            default:
                return 'text-gray-400';
        }
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
                    System Services
                </motion.h1>
                <div className="bg-gray-800 rounded-xl shadow-2xl overflow-hidden border border-gray-700">
                    <Table>
                        <TableHeader>
                            <TableRow className="bg-gray-900 hover:bg-gray-900">
                                {[
                                    { key: 'name', label: 'Service Name' },
                                    { key: 'status', label: 'Status' },
                                    { key: 'uptime', label: 'Uptime' },
                                    { key: 'memory', label: 'Memory' },
                                    { key: 'version', label: 'Version' },
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
                                {services.map((service:any, index:any) => (
                                    <motion.tr
                                        key={service.name}
                                        initial={{ opacity: 0, y: 20 }}
                                        animate={{ opacity: 1, y: 0 }}
                                        exit={{ opacity: 0, y: -20 }}
                                        transition={{ duration: 0.3, delay: index * 0.05 }}
                                        className={`border-b border-gray-700 ${index % 2 === 0 ? 'bg-gray-800' : 'bg-gray-850'
                                            } hover:bg-gray-700 transition-colors duration-200`}
                                    >
                                        <TableCell className="text-green-400 font-medium">{service.name}</TableCell>
                                        <TableCell className={getStatusColor(service.status)}>{service.status}</TableCell>
                                        <TableCell className="text-green-400">{service.uptime}</TableCell>
                                        <TableCell className="text-green-400">{service.memory}</TableCell>
                                        <TableCell className="text-green-400">{service.version || 'N/A'}</TableCell>
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

export default ServicesList;