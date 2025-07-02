'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../store/useAuth';
import Sidebar from './components/Sidebar';
import Header from './components/Header';
import { SocketProvider } from '../context/SocketContext';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, initAuth } = useAuth();
  const router = useRouter();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    initAuth();
    setTimeout(() => {
      if (!useAuth.getState().isAuthenticated) {
        router.replace('/auth/');
      }
      setLoading(false);
    }, 200); 
  }, []);

  if (loading) return null;

  return (
    <SocketProvider>
    <div className="flex h-screen  bg-black">
         <Sidebar />
      <div className="flex flex-col flex-1">
        <Header />
        <main className="p-4 overflow-y-auto">{children}</main>
      </div>
    </div>
    </SocketProvider>
  );
}
