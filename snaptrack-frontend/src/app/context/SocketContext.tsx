"use client";

import { useLinuxToast } from "@/lib/use-linux-toast";
import React, { createContext, useContext, useEffect, useRef, useState } from "react";

interface ServiceInfo {
  name: string;
  status: string;
  uptime: string;
  memory: string;
  version: string;
}

interface Metrics {
  diskTotal: number;
  uptimeSeconds: number;
  netOutBytes: number;
  netInBytes: number;
  ramUsedBytes: number;
  ramTotalBytes: number;
  diskUsed: number;
  diskUsedBytes: number;
  cpuPercent: number;
  ramPercent: number;
  diskPercent: number;
}

interface RunningProcess {
  pid: number;
  name: string;
  cpu_percent: number;
  mem_percent: number;
  status: string;
}

interface FirewallRule {
  id: string;
  protocol: string;
  port: string;
  source: string;
  destination: string;
  action: string;
}

interface RunningPort {
  protocol: string;
  port: number;
  process: string;
  pid: number;
}

interface SocketContextType {
  socket: WebSocket | null;
  metrics: Metrics | null;
  services: ServiceInfo[] | null;
  firewallRules: FirewallRule[] | null;
  runningPorts: RunningPort[] | null;
  logs: { [service: string]: string[] };
  runningProcesses: RunningProcess[] | null;
  sendAction: (
    type: "start" | "stop" | "restart" | "logs" | "stop_port" | "add_port" | "add_rule",
    data: any
  ) => void;
}

const SocketContext = createContext<SocketContextType>({
  socket: null,
  metrics: null,
  services: null,
  firewallRules: null,
  runningPorts: null,
  logs: {},
  sendAction: () => {},
  runningProcesses: null,
});

export const useSocket = () => useContext(SocketContext);

export const SocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const socketRef = useRef<WebSocket | null>(null);
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [services, setServices] = useState<ServiceInfo[] | null>(null);
  const [firewallRules, setFirewallRules] = useState<FirewallRule[] | null>(null);
  const [runningPorts, setRunningPorts] = useState<RunningPort[] | null>(null);
  const [logs, setLogs] = useState<{ [service: string]: string[] }>({});
  const [runningProcesses, setRunningProcesses] = useState<RunningProcess[] | null>(null);
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const { success, error } = useLinuxToast();
  
  const sendAction = (
    type: any,
    data: any
  ) => {
    const ws = socketRef.current;
    if (ws && ws.readyState === WebSocket.OPEN) {
      let actionType = type;
      if (["start", "stop", "restart", "logs"].includes(type)) {
        actionType = `services_${type}`;
      } else if (type === "add_rule") {
        actionType = `firewall_${type}`;
      }
      const message = JSON.stringify({ type: actionType, data });
      ws.send(message);
    } else {
      error("WebSocket not connected");
    }
  };

  useEffect(() => {
    if (socketRef.current) return; 

    const ws = new WebSocket("ws://localhost:8000/ws");
    socketRef.current = ws;
    setSocket(ws);

    let hasConnected = false;

    ws.onopen = () => {
      if (!hasConnected) {
        success("✅ Connected to server");
        hasConnected = true;
      }
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log("🔵 WS Message:", data);

        switch (data.type) {
          case "services":
            setServices(data.services);
            break;
          case "metrics":
            setMetrics(data.stats);
            break;
          case "firewalls":
            setFirewallRules(data.rules);
            break;
          case "ports":
            setRunningPorts(data.ports);
            break;
          case "process_list":
            setRunningProcesses(data.processes);
            break;
          case "log":
            setLogs((prev) => ({
              ...prev,
              [data.service]: data.log.split("\n").filter((line: string) => line.trim()),
            }));
            break;
          case "action_response":
            break;
          default:
            console.warn("⚠️ Unknown message type:", data.type);
        }
      } catch (err) {
        console.error("❗ Failed to parse WS message", err);
        error("Error processing server message");
      }
    };

    ws.onclose = () => {
      console.log("❌ WebSocket disconnected");
    };

    ws.onerror = (err) => {
      console.log("❗ WebSocket error:", err);
    };

    return () => {
      if (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING) {
        console.log("🧹 Closing WebSocket");
        ws.close();
      }
      socketRef.current = null;
    };
  }, []);

  return (
    <SocketContext.Provider
      value={{
        socket: socketRef.current,
        metrics,
        services,
        firewallRules,
        runningPorts,
        logs,
        sendAction,
        runningProcesses,
      }}
    >
      {children}
    </SocketContext.Provider>
  );
};
