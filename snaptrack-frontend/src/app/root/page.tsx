'use client';
import { useSocket } from '../context/SocketContext';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Cpu, HardDrive, MemoryStick, Clock3, Download, Upload, Database, DatabaseBackup, WindArrowDown, Clock, Terminal, CheckCircle, AlertCircle } from 'lucide-react';
import { CircularProgressbar, buildStyles } from 'react-circular-progressbar';
import 'react-circular-progressbar/dist/styles.css';
import { motion } from 'framer-motion';
import { useEffect, useState } from 'react';
import { getAllBackups } from '@/lib/api';
import { useLinuxToast } from '@/lib/use-linux-toast';
import clsx from 'clsx';
import { format } from "date-fns";


function formatBytes(bytes: number) {
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    if (bytes === 0) return '0 B';
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`;
}

function formatUptime(seconds: number) {
    const days = Math.floor(seconds / (3600 * 24));
    const hours = Math.floor((seconds % (3600 * 24)) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${days}d ${hours}h ${minutes}m`;
}

function formatSpeed(bytesPerSec: number) {
    return `${formatBytes(bytesPerSec)}/s`;
}

export default function MonitorPage() {
    const { metrics } = useSocket();
    const [netInHistory, setNetInHistory] = useState<number[]>([]);
    const [netOutHistory, setNetOutHistory] = useState<number[]>([]);
    const [backups, setBackups] = useState<any[]>([]);
    const [isLoadingBackups, setIsLoadingBackups] = useState<boolean>(true);
    const { success, error } = useLinuxToast();

    const netInSpeed = netInHistory.length >= 2 ? netInHistory[netInHistory.length - 1] - netInHistory[netInHistory.length - 2] : 0;
    const netOutSpeed = netOutHistory.length >= 2 ? netOutHistory[netOutHistory.length - 1] - netOutHistory[netOutHistory.length - 2] : 0;

    const netInPeak = Math.max(...netInHistory, 0);
    const netOutPeak = Math.max(...netOutHistory, 0);

    const netInTrend = netInSpeed > 0 ? '↑' : netInSpeed < 0 ? '↓' : '→';
    const netOutTrend = netOutSpeed > 0 ? '↑' : netOutSpeed < 0 ? '↓' : '→';

    useEffect(() => {
        if (!metrics) return;

        setNetInHistory((prev) => [...prev.slice(-9), metrics.netInBytes]);
        setNetOutHistory((prev) => [...prev.slice(-9), metrics.netOutBytes]);
    }, [metrics]);

    useEffect(() => {
        async function fetchBackups() {
            try {
                setIsLoadingBackups(true);
                const response = await getAllBackups();
                const backupsArray = Array.isArray(response) ? response : (response || []);
                setBackups(backupsArray);
            } catch (err: any) {
                error(err.message || 'Failed to fetch backups');
            } finally {
                setIsLoadingBackups(false);
            }
        }
        fetchBackups();
    }, []);

    if (!metrics) {
        return (
            <div className="flex items-center justify-center h-full text-zinc-400 font-mono text-lg">
                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.5 }}
                >
                    Waiting for metrics...
                </motion.div>
            </div>
        );
    }

    const {
        cpuPercent,
        ramPercent,
        ramTotalBytes,
        ramUsedBytes,
        diskTotal,
        diskUsed,
        diskPercent,
        uptimeSeconds,
        netInBytes,
        netOutBytes,
    } = metrics;

    // Compute backup counters
    const totalBackups = backups.length;
    const pendingBackups = backups.filter(backup => backup.status === 'pending').length;
    const successfulBackups = backups.filter(backup => backup.status === 'success').length;
    const failedBackups = backups.filter(backup =>
        backup.logs && backup.logs.some((log: any) => log.status === 'failed')
    ).length;

    return (
        <div className="w-full h-full text-green-400 p-4 sm:p-6 overflow-x-auto">
            <motion.div
                className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 min-w-max"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5 }}
            >
                <Card className="bg-gray-800/50 border-gray-700/50 shadow-lg hover:shadow-green-500/20 transition-shadow duration-300 min-w-[200px]">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2 text-white text-lg font-semibold">
                            <Cpu className="w-5 h-5 text-green-400" />
                            CPU Usage
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="flex justify-center">
                        <div className="w-24 h-24 sm:w-32 sm:h-32">
                            <CircularProgressbar
                                value={cpuPercent}
                                text={`${cpuPercent.toFixed(0)}%`}
                                styles={buildStyles({
                                    textColor: '#ffffff',
                                    pathColor: '#10b981',
                                    trailColor: '#1e293b',
                                    textSize: '20px',
                                })}
                            />
                        </div>
                    </CardContent>
                </Card>
                <Card className="bg-gray-800/50 border-gray-700/50 shadow-lg hover:shadow-green-500/20 transition-shadow duration-300 min-w-[200px]">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2 text-white text-lg font-semibold">
                            <MemoryStick className="w-5" />
                            RAM Usage
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="flex flex-col items-center space-y-2">
                        <div className="w-24 h-24 sm:w-32 sm:h-32">
                            <CircularProgressbar
                                value={ramPercent}
                                text={`${ramPercent.toFixed(0)}%`}
                                styles={buildStyles({
                                    textColor: '#ffffff',
                                    pathColor: '#10b981',
                                    trailColor: '#1e293b',
                                    textSize: '20px',
                                })}
                            />
                        </div>
                        <div className="text-sm text-zinc-400 text-center">
                            {formatBytes(ramUsedBytes)} / {formatBytes(ramTotalBytes)}
                        </div>
                    </CardContent>
                </Card>
                <Card className="bg-gray-800/50 border-gray-700/50 shadow-lg hover:shadow-green-500/20 transition-shadow duration-300 min-w-[200px]">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2 text-white text-lg font-semibold">
                            <HardDrive className="w-5 h-5 text-green-400" />
                            Disk Usage
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="flex flex-col items-center space-y-2">
                        <div className="w-24 h-24 sm:w-32 sm:h-32">
                            <CircularProgressbar
                                value={diskPercent}
                                text={`${diskPercent.toFixed(0)}%`}
                                styles={buildStyles({
                                    textColor: '#ffffff',
                                    pathColor: '#10b981',
                                    trailColor: '#1e293b',
                                    textSize: '20px',
                                })}
                            />
                        </div>
                        <div className="text-sm text-zinc-400 text-center">
                            {formatBytes(diskUsed)} / {formatBytes(diskTotal)}
                        </div>
                    </CardContent>
                </Card>
                <Card className="bg-gray-800/50 border-gray-700/50 shadow-lg hover:shadow-green-500/30 transition-all duration-300 min-w-[200px]">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2 text-white text-lg font-semibold">
                            <Download className="w-5 h-5 text-green-400" />
                            Network In
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-3">
                        <div className="text-lg text-white font-mono">
                            {formatBytes(netInBytes)} <span className="text-green-400">{netInTrend}</span>
                        </div>
                        <div className="text-sm text-zinc-400">
                            Speed: {formatSpeed(Math.abs(netInSpeed))}
                        </div>
                        <div className="text-sm text-zinc-400">
                            Peak: {formatBytes(netInPeak)}
                        </div>
                    </CardContent>
                </Card>
                <Card className="bg-gray-800/50 border-gray-700/50 shadow-lg hover:shadow-green-500/30 transition-all duration-300 min-w-[200px]">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2 text-white text-lg font-semibold">
                            <Upload className="w-5 h-5 text-green-400" />
                            Network Out
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-3">
                        <div className="text-lg text-white font-mono">
                            {formatBytes(netOutBytes)} <span className="text-green-400">{netOutTrend}</span>
                        </div>
                        <div className="text-sm text-zinc-400">
                            Speed: {formatSpeed(Math.abs(netOutSpeed))}
                        </div>
                        <div className="text-sm text-zinc-400">
                            Peak: {formatBytes(netOutPeak)}
                        </div>
                    </CardContent>
                </Card>
                <Card className="bg-gray-800/50 border-gray-700/50 shadow-lg hover:shadow-green-500/20 transition-shadow duration-300 min-w-[200px]">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2 text-white text-lg font-semibold">
                            <Clock3 className="w-5 h-5 text-green-400" />
                            Uptime
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="flex flex-col items-center space-y-2">
                        <div className="w-24 h-24 sm:w-32 sm:h-32">
                            <CircularProgressbar
                                value={(uptimeSeconds % 86400) / 86400 * 100}
                                text={formatUptime(uptimeSeconds)}
                                styles={buildStyles({
                                    textColor: '#ffffff',
                                    pathColor: '#10b981',
                                    trailColor: '#1e293b',
                                    textSize: '14px',
                                })}
                            />
                        </div>
                    </CardContent>
                </Card>

            </motion.div>
            <div className="space-y-4 mt-10">
                <h2 className="text-xl font-bold text-white">Backups</h2>

                <div className="flex flex-wrap gap-4">
                    {/* Successful */}
                    <Card className="bg-gray-800/50 border-gray-700/50 shadow hover:shadow-green-500/20 transition-shadow duration-300 min-w-[200px]">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-white text-lg font-semibold">
                                <DatabaseBackup className="w-5 h-5 text-green-400" />
                                Successful
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            {isLoadingBackups ? (
                                <div className="text-sm text-zinc-400 text-center">Loading...</div>
                            ) : (
                                <div className="text-2xl text-green-400 font-bold">{successfulBackups}</div>
                            )}
                        </CardContent>
                    </Card>

                    {/* Failed */}
                    <Card className="bg-gray-800/50 border-gray-700/50 shadow hover:shadow-red-500/20 transition-shadow duration-300 min-w-[200px]">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-white text-lg font-semibold">
                                <WindArrowDown className="w-5 h-5 text-red-400" />
                                Failed
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            {isLoadingBackups ? (
                                <div className="text-sm text-zinc-400 text-center">Loading...</div>
                            ) : (
                                <div className="text-2xl text-red-400 font-bold">{failedBackups}</div>
                            )}
                        </CardContent>
                    </Card>

                    {/* Pending */}
                    <Card className="bg-gray-800/50 border-gray-700/50 shadow hover:shadow-yellow-500/20 transition-shadow duration-300 min-w-[200px]">
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2 text-white text-lg font-semibold">
                                <Clock className="w-5 h-5 text-yellow-400" />
                                Pending
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            {isLoadingBackups ? (
                                <div className="text-sm text-zinc-400 text-center">Loading...</div>
                            ) : (
                                <div className="text-2xl text-yellow-400 font-bold">{pendingBackups}</div>
                            )}
                        </CardContent>
                    </Card>
                </div>
            </div>
            <section className="space-y-6 mt-10">
                <h2 className="text-xl font-bold text-white">Backup Logs</h2>

                {backups.map((backup) => (
                    <div key={backup.id} className="border-l-2 border-green-500 pl-4 space-y-2">
                        {/* Backup Title */}
                        <div className="flex items-center gap-2 text-white font-semibold">
                            <Terminal className="w-4 h-4 text-green-400" />
                            <span>Backup:</span>
                            <span className="text-green-300">{backup.app}</span>
                            <span className="text-sm text-zinc-400 ml-2">({backup.type})</span>
                            <span
                                className={clsx(
                                    "ml-auto text-xs font-medium capitalize",
                                    backup.status === "pending" && "text-yellow-400",
                                    backup.status === "failed" && "text-red-400",
                                    backup.status === "success" && "text-green-400"
                                )}
                            >
                                {backup.status}
                            </span>
                        </div>

                        {/* Logs */}
                        <div className="ml-6 space-y-1">
                            {backup.logs.length === 0 ? (
                                <div className="text-sm text-zinc-500">No logs available.</div>
                            ) : (
                                backup.logs.map((log: any) => (
                                    <div key={log.id} className="flex items-start gap-2 text-sm text-zinc-300">
                                        {/* Status Icon */}
                                        {log.status === "failed" ? (
                                            <AlertCircle className="w-4 h-4 text-red-400 mt-[2px]" />
                                        ) : log.status === "success" ? (
                                            <CheckCircle className="w-4 h-4 text-green-400 mt-[2px]" />
                                        ) : (
                                            <Clock className="w-4 h-4 text-yellow-400 mt-[2px]" />
                                        )}

                                        {/* Message + Timestamp */}
                                        <div className="flex flex-col">
                                            <span>{log.message}</span>
                                            <span className="text-xs text-zinc-500">
                                                {format(new Date(log.createdAt), "yyyy-MM-dd HH:mm:ss")}
                                            </span>
                                        </div>
                                    </div>
                                ))
                            )}
                        </div>
                    </div>
                ))}
            </section>
        </div>
    );
}