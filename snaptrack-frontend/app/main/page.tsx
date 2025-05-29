"use client";

import React from "react";
import { useSocket } from "../context/SocketContext";

const circleRadius = 60;
const circumference = 2 * Math.PI * circleRadius;

interface Metrics {
  cpuPercent: number;
  ramPercent: number;
  ramTotalBytes: number;
  diskPercent: number;
  diskTotalBytes: number;
  netInBytes: number;
  netOutBytes: number;
  uptimeSeconds: number;
}

function CircularProgress({
  label,
  percent,
  color,
  usageText,
}: {
  label: string;
  percent: number;
  color: string;
  usageText?: string;
}) {
  const strokeDashoffset = circumference - (percent / 100) * circumference;

  return (
    <div className="w-36 mx-5 text-center select-none font-sans">
      <svg width={140} height={140} className="mx-auto">
        <circle
          stroke="#eee"
          fill="transparent"
          strokeWidth={12}
          r={circleRadius}
          cx={70}
          cy={70}
        />
        <circle
          stroke={color}
          fill="transparent"
          strokeWidth={12}
          strokeDasharray={circumference}
          strokeDashoffset={strokeDashoffset}
          r={circleRadius}
          cx={70}
          cy={70}
          style={{ transition: "stroke-dashoffset 0.7s ease" }}
          strokeLinecap="round"
        />
        <text
          x="50%"
          y="45%"
          dominantBaseline="middle"
          textAnchor="middle"
          className="text-3xl font-semibold fill-gray-800"
        >
          {percent.toFixed(1)}%
        </text>
      </svg>
      <div className="mt-2 text-lg font-bold text-gray-600">{label}</div>
      {usageText && (
        <div className="mt-1 text-sm font-medium text-gray-400">{usageText}</div>
      )}
    </div>
  );
}

const backups = [
  { id: 1, date: "2025-05-25 14:00", size: "1.2GB", status: "Success" },
  { id: 2, date: "2025-05-24 14:00", size: "1.1GB", status: "Success" },
  { id: 3, date: "2025-05-23 14:00", size: "1.3GB", status: "Failed" },
];

function BackupList() {
  return (
    <div className="mt-10 w-11/12 max-w-3xl bg-white rounded-xl shadow-lg p-6 font-sans">
      <h2 className="mb-4 text-2xl font-semibold text-gray-700">Recent Backups</h2>
      <table className="w-full border-collapse">
        <thead>
          <tr className="border-b-2 border-gray-200">
            <th className="text-left px-3 py-2">Date</th>
            <th className="text-left px-3 py-2">Size</th>
            <th className="text-left px-3 py-2">Status</th>
          </tr>
        </thead>
        <tbody>
          {backups.map(({ id, date, size, status }) => (
            <tr key={id} className={status === "Failed" ? "bg-red-100" : ""}>
              <td className="px-3 py-2">{date}</td>
              <td className="px-3 py-2">{size}</td>
              <td
                className={`px-3 py-2 font-semibold ${
                  status === "Failed" ? "text-red-700" : "text-green-700"
                }`}
              >
                {status}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default function Dashboard() {
  const { metrics } = useSocket();

  if (!metrics)
    return <p className="text-center mt-20 text-gray-500">Loading metrics...</p>;

  const formatBytesToGB = (bytes: number) =>
    (bytes / 1024 ** 3).toFixed(2) + " GB";

  const totalRAM = metrics.ramTotalBytes;
  const totalDisk = metrics.diskTotalBytes;

  const ramUsed = (metrics.ramPercent / 100) * totalRAM;
  const diskUsed = (metrics.diskPercent / 100) * totalDisk;

  const networkTotal = metrics.netInBytes + metrics.netOutBytes;
  const networkUsedGB = networkTotal / 1024 ** 3;
  const networkPercent = Math.min((networkUsedGB / 10) * 100, 100);
  
  const uptimePercent = Math.min((metrics.uptimeSeconds / (24 * 3600)) * 100, 100);
  const formatUptime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}h ${minutes}m`;
  };

  return (
    <div className=" flex flex-col items-center">
      <div className="flex flex-wrap justify-center gap-10">
        <CircularProgress
          label="CPU Usage"
          percent={metrics.cpuPercent}
          color="#ff6b6b"
          usageText={`${metrics.cpuPercent.toFixed(1)}% used`}
        />
        <CircularProgress
          label="RAM Usage"
          percent={metrics.ramPercent}
          color="#4caf50"
          usageText={`${formatBytesToGB(ramUsed)} / ${formatBytesToGB(totalRAM)}`}
        />
        <CircularProgress
          label="Disk Usage"
          percent={metrics.diskPercent}
          color="#1e90ff"
          usageText={`${formatBytesToGB(diskUsed)} / ${formatBytesToGB(totalDisk)}`}
        />
        <CircularProgress
          label="Network I/O"
          percent={networkPercent}
          color="#ffa500"
          usageText={`${formatBytesToGB(networkTotal)} used`}
        />
        <CircularProgress
          label="Uptime"
          percent={uptimePercent}
          color="#6a5acd"
          usageText={formatUptime(metrics.uptimeSeconds)}
        />
      </div>

      <BackupList />
    </div>
  );
}
