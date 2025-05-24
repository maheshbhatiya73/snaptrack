"use client";

import React from "react";
import { useAuth } from "@/app/context/AuthContext";

export default function MainPage() {
  const { setToken } = useAuth();

  const handleLogout = () => {
    setToken(null);
  };

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-4xl mx-auto bg-white rounded-2xl shadow-2xl p-8">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Server Dashboard</h1>
          <button
            onClick={handleLogout}
            className="py-2 px-4 bg-red-600 text-white rounded-lg hover:bg-red-700 transition duration-300"
          >
            Logout
          </button>
        </div>
        <p className="text-gray-600 mb-4">Welcome, admin! Manage your server resources below.</p>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="p-4 bg-gray-50 rounded-lg">
            <h2 className="text-xl font-semibold text-gray-800">System Status</h2>
            <p className="text-gray-600">All systems operational</p>
          </div>
          <div className="p-4 bg-gray-50 rounded-lg">
            <h2 className="text-xl font-semibold text-gray-800">Recent Activity</h2>
            <p className="text-gray-600">No recent activity</p>
          </div>
        </div>
      </div>
    </div>
  );
}