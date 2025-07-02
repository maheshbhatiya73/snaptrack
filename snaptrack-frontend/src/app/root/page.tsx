'use client';
import { useSocket } from '../context/SocketContext';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Cpu, HardDrive, MemoryStick, Clock3, Download, Upload } from 'lucide-react';
import { CircularProgressbar, buildStyles } from 'react-circular-progressbar';
import 'react-circular-progressbar/dist/styles.css';
import { motion } from 'framer-motion';
import { useEffect, useState } from 'react';


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

    const netInSpeed = netInHistory.length >= 2 ? netInHistory[netInHistory.length - 1] - netInHistory[netInHistory.length - 2] : 0;
    const netOutSpeed = netOutHistory.length >= 2 ? netOutHistory[netOutHistory.length - 1] - netOutHistory[netOutHistory.length - 2] : 0;

    const netInPeak = Math.max(...netInHistory);
    const netOutPeak = Math.max(...netOutHistory);

    const netInTrend = netInSpeed > 0 ? '↑' : netInSpeed < 0 ? '↓' : '→';
    const netOutTrend = netOutSpeed > 0 ? '↑' : netOutSpeed < 0 ? '↓' : '→';

    useEffect(() => {
        if (!metrics) return;

        setNetInHistory((prev) => [...prev.slice(-9), metrics.netInBytes]);
        setNetOutHistory((prev) => [...prev.slice(-9), metrics.netOutBytes]);
    }, [metrics]);

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
                            <MemoryStick className="w-5 h-5 text-green-400" />
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
                <Card className="bg-gray-800/50  border-gray-700/50 shadow-lg hover:shadow-green-500/30 transition-all duration-300 min-w-[200px]">
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
                <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700/50 shadow-lg hover:shadow-green-500/30 transition-all duration-300 min-w-[200px]">
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
        </div>
    );
}