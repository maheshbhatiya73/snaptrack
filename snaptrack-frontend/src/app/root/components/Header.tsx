'use client';

export default function Header() {
  return (
    <header className="h-14 border-b border-zinc-800 bg-zinc-900 px-4 flex items-center justify-between text-green-400 font-mono">
      <div className="text-sm">Linux Agent Monitor</div>
      <div className="text-xs text-zinc-400">v1.0.0</div>
    </header>
  );
}
