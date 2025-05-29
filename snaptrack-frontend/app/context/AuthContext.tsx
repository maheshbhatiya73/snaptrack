"use client";

import React, { createContext, useContext, useEffect, useState, useCallback } from "react";
import { useRouter, usePathname } from "next/navigation";
import { verifyToken } from "@/app/lib/api";

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  token: string | null;
  setToken: (token: string | null) => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  isLoading: true,
  token: null,
  setToken: async () => {},
});

export const useAuth = () => useContext(AuthContext);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [token, setTokenState] = useState<string | null>(null);


  const router = useRouter();
  const pathname = usePathname();

 const setToken = useCallback(async (token: string | null) => {
  console.log("setToken called with token:", token ? "present" : "null");
  if (token) {
    const valid = await verifyToken(token);
    console.log("Token valid:", valid);
    if (valid) {
      localStorage.setItem("token", token);
      setTokenState(token);
      setIsAuthenticated(true);
    } else {
      localStorage.removeItem("token");
      setTokenState(null);
      setIsAuthenticated(false);
    }
  } else {
    localStorage.removeItem("token");
    setTokenState(null);
    setIsAuthenticated(false);
  }
}, []);

  useEffect(() => {
  const checkAuth = async () => {
    const token = localStorage.getItem("token");
    console.log("Checking auth, token:", token ? "present" : "null");
    if (token) {
      const valid = await verifyToken(token);
      console.log("Initial token valid:", valid);
      setIsAuthenticated(valid);
      setTokenState(valid ? token : null);
      if (!valid) {
        localStorage.removeItem("token");
      }
    } else {
      setTokenState(null);
    }
    setIsLoading(false);
  };

  checkAuth();
}, []);

  useEffect(() => {
    if (isLoading) return;

    console.log("Auth redirect check - isAuthenticated:", isAuthenticated, "pathname:", pathname);
    if (isAuthenticated && pathname === "/auth/login") {
      console.log("Redirecting to /main from /auth/login");
      router.push("/main");
    } else if (!isAuthenticated && pathname !== "/auth/login" && pathname !== "/") {
      console.log("Redirecting to /auth/login from", pathname);
      router.push("/auth/login");
    }
  }, [isAuthenticated, pathname, isLoading, router]);

  return (
   <AuthContext.Provider value={{ isAuthenticated, isLoading, token, setToken }}>
  {children}
</AuthContext.Provider>

  );
};