"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { login } from "@/app/lib/api";
import { useAuth } from "@/app/context/AuthContext";

export default function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const router = useRouter();
  const { setToken } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const response = await login(username, password);

      if (response.success && response.token) {
        await setToken(response.token); // Wait for token verification
        router.push("/main"); // Explicit redirect
      } else {
        setError(response.message || "Authentication failed");
      }
    } catch (err) {
      setError("Failed to connect to server");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100 transition-colors duration-300">
      <div className="max-w-md w-full bg-white rounded-2xl shadow-2xl p-8 transform hover:scale-105 transition-transform duration-300">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Root Login</h1>
        </div>
        <p className="text-gray-600 mb-8">Access your server with admin credentials</p>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="username" className="block text-sm font-medium text-gray-700">
              Username
            </label>
            <div className="relative mt-1">
              <input
                id="username"
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
                placeholder="root"
                required
              />
              <span className="absolute inset-y-0 right-3 flex items-center text-gray-400">👤</span>
            </div>
          </div>
          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700">
              Password
            </label>
            <div className="relative mt-1">
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
                placeholder="••••••••"
                required
              />
              <span className="absolute inset-y-0 right-3 flex items-center text-gray-400">🔒</span>
            </div>
          </div>
          {error && (
            <p className="text-red-500 text-sm animate-pulse" role="alert">
              {error}
            </p>
          )}
          <button
            type="submit"
            disabled={loading}
            className={`w-full py-3 px-4 rounded-lg text-white font-semibold ${
              loading ? "bg-gray-500 cursor-not-allowed" : "bg-blue-600 hover:bg-blue-700"
            } transition duration-300 flex items-center justify-center gap-2`}
          >
            {loading ? (
              <>
                <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
                Logging in...
              </>
            ) : (
              <>
                <span>Login</span>
                <span>🚀</span>
              </>
            )}
          </button>
        </form>
        <p className="mt-4 text-center text-sm text-gray-600">
          Forgot password? <a href="/forgot" className="text-blue-500 hover:underline">Reset</a>
        </p>
      </div>
    </div>
  );
}