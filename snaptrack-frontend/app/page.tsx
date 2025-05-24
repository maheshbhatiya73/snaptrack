"use client";

import React, { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/app/context/AuthContext";

export default function HomePage() {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (isLoading) return;

    console.log("HomePage - isAuthenticated:", isAuthenticated);
    if (isAuthenticated) {
      console.log("Redirecting to /main from /");
      router.replace("/main");
    } else {
      console.log("Redirecting to /auth/login from /");
      router.replace("/auth/login");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return <div className="min-h-screen flex items-center justify-center bg-gray-100">Loading...</div>;
  }

  // Render nothing while redirecting
  return null;
}