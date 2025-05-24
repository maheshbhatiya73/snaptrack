"use client";

import React, { createContext, useContext, useEffect, useState } from "react";

interface SocketContextType {
  socket: WebSocket | null;
}

const SocketContext = createContext<SocketContextType>({ socket: null });

export const useSocket = () => useContext(SocketContext);

export const SocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);

  useEffect(() => {
    console.log("🔧 Initializing WebSocket...");

    const socket = new WebSocket("ws://localhost:8000/ws");

    socket.onopen = () => {
      console.log("✅ Connected");
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
    <SocketContext.Provider value={{ socket }}>
      {children}
    </SocketContext.Provider>
  );
};