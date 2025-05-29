'use client';

import { useState } from 'react';
import Header from '../components/Header';
import Sidebar from '../components/Sidebar';
import { AuthProvider } from '../context/AuthContext';
import { ToastProvider } from '../components/Toast';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(true); // Sidebar open by default

  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen);

  return (
    <html lang="en">
      <body className="flex flex-col w-full min-h-screen bg-gray-50">
        <AuthProvider>
          <ToastProvider>
          {/* Header */}
          <Header toggleSidebar={toggleSidebar} />

          <div className="flex flex-1 mt-16">
            {/* Sidebar */}
            <Sidebar isOpen={isSidebarOpen} toggleSidebar={toggleSidebar} />

            {/* Main Content */}
            <main
              className={`flex-1 w-full transition-all duration-300 ${isSidebarOpen ? 'ml-60' : 'ml-20'
                } p-6`}
            >
              {children}
            </main>
          </div>
          </ToastProvider>
        </AuthProvider>
      </body>
    </html>
  );
}