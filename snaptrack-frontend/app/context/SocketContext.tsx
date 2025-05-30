"use client";

import React, { createContext, useContext, useEffect, useState } from "react";

interface Metrics {
  diskTotal: number;
  uptimeSeconds: number;
  netOutBytes: any;
  netInBytes: any;
  ramUsedBytes: number;
  ramTotalBytes: number;
  diskTotalBytes: number;
  diskUsedBytes: number;
  cpuPercent: number;
  ramPercent: number;
  diskPercent: number;
}

interface SocketContextType {
  socket: WebSocket | null;
  metrics: Metrics | null;
}

const SocketContext = createContext<SocketContextType>({
  socket: null,
  metrics: null,
});

export const useSocket = () => useContext(SocketContext);

export const SocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [metrics, setMetrics] = useState<Metrics | null>(null);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8000/ws");

    socket.onopen = () => {
      console.log("✅ Connected");
    };

    socket.onmessage = (event) => {
      try {
        const data: Metrics = JSON.parse(event.data);
        console.log(data)
        setMetrics(data);
      } catch (e) {
        console.error("Failed to parse message", e);
      }
    };

    socket.onclose = () => {
      console.log("❌ Disconnected");
    };

    socket.onerror = (err) => {
      console.error("❗ Error:", err);
    };

    setSocket(socket);

    return () => {
      socket.close();
    };
  }, []);

  return (
    <SocketContext.Provider value={{ socket, metrics }}>
      {children}
    </SocketContext.Provider>
  );
};
