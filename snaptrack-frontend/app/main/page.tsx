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

export default function Dashboard() {
  const { metrics } = useSocket();

  if (!metrics)
    return <p className="text-center mt-20 text-gray-500">Loading metrics...</p>;

  const formatBytesToGB = (bytes: number) =>
    (bytes / 1024 ** 3).toFixed(2) + " GB";

  const cpuPercent = metrics.cpuPercent ?? 0;
  const ramPercent = metrics.ramPercent ?? 0;
  const diskPercent = metrics.diskPercent ?? 0;
  const ramTotalBytes = metrics.ramTotalBytes ?? 1;
  const diskTotalBytes = metrics.diskTotal ?? 1;
  const netInBytes = metrics.netInBytes ?? 0;
  const netOutBytes = metrics.netOutBytes ?? 0;
  const uptimeSeconds = metrics.uptimeSeconds ?? 0;

  const ramUsed = (ramPercent / 100) * ramTotalBytes;
  const diskUsed = (diskPercent / 100) * diskTotalBytes;
  const networkTotal = netInBytes + netOutBytes;
  const networkUsedGB = networkTotal / 1024 ** 3;
  const networkPercent = Math.min((networkUsedGB / 10) * 100, 100);
  const uptimePercent = Math.min((uptimeSeconds / (24 * 3600)) * 100, 100);
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
          percent={cpuPercent}
          color="#ff6b6b"
          usageText={`${cpuPercent.toFixed(1)}% used`}
        />
        <CircularProgress
          label="RAM Usage"
          percent={ramPercent}
          color="#4caf50"
          usageText={`${formatBytesToGB(ramUsed)} / ${formatBytesToGB(ramTotalBytes)}`}
        />
        <CircularProgress
          label="Disk Usage"
          percent={diskPercent}
          color="#1e90ff"
          usageText={`${formatBytesToGB(diskUsed)} / ${formatBytesToGB(diskTotalBytes)}`}
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
          usageText={formatUptime(uptimeSeconds)}
        />

      </div>
    </div>
  );
}
