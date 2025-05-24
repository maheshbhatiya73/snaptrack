// app/ClientWrapper.tsx
"use client";

import React from "react";
import { SocketProvider } from "./context/SocketContext";

export default function ClientWrapper({ children }: { children: React.ReactNode }) {
  return <SocketProvider>{children}</SocketProvider>;
}
