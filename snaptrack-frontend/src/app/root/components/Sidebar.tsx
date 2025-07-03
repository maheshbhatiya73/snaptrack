'use client';

import { Home, Server, Shield, Settings, LogOut, HardDrive, Activity, Cpu, ServerCog, icons } from 'lucide-react';
import Link from 'next/link';

const navItems = [
  { label: 'Dashboard', href: '/root', icon: Home },
  { label: 'Backups', href: '/root/backups', icon: HardDrive },
  { label: 'Monitor', href: '/root/monitor', icon: Activity },
  { label: 'Process', href: '/root/process', icon: Cpu },
  { label: 'Services', href: '/root/services', icon: ServerCog },
  { label: "firewall", href: "/root/firewall", icon: Shield},
  { label: 'Settings', href: '/root/settings', icon: Settings },
];

export default function Sidebar() {
  return (
    <aside className="w-64 bg-zinc-900 text-green-400 h-screen p-4 border-r border-zinc-800 font-mono hidden md:block">
      <div className="text-lg font-bold mb-6 text-white">SnapTrack</div>
      <nav className="space-y-2">
        {navItems.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className="flex items-center gap-3 px-3 py-2 rounded hover:bg-zinc-800 transition-colors"
          >
            <item.icon className="w-5 h-5" />
            {item.label}
          </Link>
        ))}
        <button
          className="mt-6 flex items-center gap-3 px-3 py-2 text-red-500 hover:bg-red-900/20 w-full rounded"
          onClick={() => {
            localStorage.removeItem('auth_token');
            window.location.href = '/auth/';
          }}
        >
          <LogOut className="w-5 h-5" />
          Logout
        </button>
      </nav>
    </aside>
  );
}
