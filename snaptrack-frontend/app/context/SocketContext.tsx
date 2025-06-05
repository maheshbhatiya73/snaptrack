"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { toast, ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

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
  diskTotalBytes: number;
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

interface LogMessage {
  type: string;
  service: string;
  log: string;
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
    data: string | { port?: number; pid?: number; protocol?: string; source?: string; destination?: string; action?: string }
  ) => void;
}

const SocketContext = createContext<SocketContextType>({
  socket: null,
  metrics: null,
  services: null,
  firewallRules: null,
  runningPorts: null,
  logs: {},
  sendAction: () => { },
  runningProcesses: null,
});

export const useSocket = () => useContext(SocketContext);

export const SocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [services, setServices] = useState<ServiceInfo[] | null>(null);
  const [firewallRules, setFirewallRules] = useState<FirewallRule[] | null>(null);
  const [runningPorts, setRunningPorts] = useState<RunningPort[] | null>(null);
  const [logs, setLogs] = useState<{ [service: string]: string[] }>({});
  const [runningProcesses, setRunningProcesses] = useState<RunningProcess[] | null>(null);


  const sendAction = (
    type: "start" | "stop" | "restart" | "logs" | "stop_port" | "add_port" | "add_rule",
    data: string | { port?: number; pid?: number; protocol?: string; source?: string; destination?: string; action?: string }
  ) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      let actionType = type;
      if (["start", "stop", "restart", "logs"].includes(type)) {
        actionType = `services_${type}`;
      } else if (type === "add_rule") {
        actionType = `firewall_${type}`;
      }
      const message = JSON.stringify({
        type: actionType,
        data, // Send data directly, let JSON.stringify handle object serialization
      });
      console.log("Sending WebSocket message:", message); // Debug log
      socket.send(message);
    } else {
      console.error("WebSocket not connected");
      toast.error("WebSocket not connected");
    }
  };

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8000/ws");

    socket.onopen = () => {
      console.log("✅ WebSocket connected");
      toast.success("Connected to server");
    };

    socket.onmessage = (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);
        console.log("Received:", data);

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
            toast[data.success ? "success" : "error"](data.message, {
              position: "top-right",
              autoClose: 2000,
            });
            break;
          default:
            console.warn("Unknown message type:", data.type);
        }
      } catch (e) {
        console.error("Failed to parse message:", e);
        toast.error("Error processing server message");
      }
    };

    socket.onclose = () => {
      console.log("❌ WebSocket disconnected");
      toast.error("Disconnected from server");
    };

    socket.onerror = (err: Event) => {
      console.error("❗ WebSocket error:", err);
      toast.error("WebSocket error occurred");
    };

    setSocket(socket);

    return () => {
      socket.close();
      console.log("WebSocket cleanup");
    };
  }, []);

  return (
    <SocketContext.Provider
      value={{ socket, metrics, services, firewallRules, runningPorts, logs, sendAction, runningProcesses }}
    >
      {children}
      <ToastContainer position="top-right" autoClose={2000} theme="colored" />
    </SocketContext.Provider>
  );
};