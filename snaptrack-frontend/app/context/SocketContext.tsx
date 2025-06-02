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

interface LogMessage {
  type: string;
  service: string;
  log: string;
}

interface SocketContextType {
  socket: WebSocket | null;
  metrics: Metrics | null;
  services: ServiceInfo[] | null;
  logs: { [service: string]: string[] };
  sendAction: (type: "start" | "stop" | "restart" | "logs", service: string) => void;
}

const SocketContext = createContext<SocketContextType>({
  socket: null,
  metrics: null,
  services: null,
  logs: {},
  sendAction: () => {},
});

export const useSocket = () => useContext(SocketContext);

export const SocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [services, setServices] = useState<ServiceInfo[] | null>(null);
  const [logs, setLogs] = useState<{ [service: string]: string[] }>({});

  const sendAction = (type: "start" | "stop" | "restart" | "logs", service: string) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      const action = { type, service };
      socket.send(JSON.stringify(action));
      console.log(`Sent action: ${JSON.stringify(action)}`);
    } else {
      console.error("WebSocket is not open");
      toast.error("WebSocket connection not established");
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

        if (data.type === "services") {
          setServices(data.services);
        } else if (data.type === "metrics") {
          setMetrics(data.stats);
        } else if (data.type === "log") {
          setLogs((prev) => ({
            ...prev,
            [data.service]: data.log.split("\n").filter((line: string) => line.trim()),
          }));
        } else if (data.type === "action_response") {
          toast[data.success ? "success" : "error"](data.message, {
            position: "top-right",
            autoClose: 3000,
          });
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
    <SocketContext.Provider value={{ socket, metrics, services, logs, sendAction }}>
      {children}
    </SocketContext.Provider>
  );
};