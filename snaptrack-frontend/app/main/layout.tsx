'use client';

import { useState } from 'react';
import Header from '../components/Header';
import Sidebar from '../components/Sidebar';
import { AuthProvider } from '../context/AuthContext';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(true); // Sidebar open by default

  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen);

  return (
    <html lang="en">
      <body className="flex flex-col min-h-screen bg-gray-50">
        <AuthProvider>
          {/* Header */}
          <Header toggleSidebar={toggleSidebar} />

          <div className="flex flex-1 mt-16">
            {/* Sidebar */}
            <Sidebar isOpen={isSidebarOpen} toggleSidebar={toggleSidebar} />

            {/* Main Content */}
            <main
              className={`flex-1 transition-all duration-300 ${
                isSidebarOpen ? 'ml-52' : 'ml-16'
              } p-6`}
            >
              {children}
            </main>
          </div>
        </AuthProvider>
      </body>
    </html>
  );
}