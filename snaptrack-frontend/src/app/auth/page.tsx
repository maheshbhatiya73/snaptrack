'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Terminal, KeyRound } from 'lucide-react';
import { motion } from 'framer-motion';
import { login } from '@/lib/api';
import { useEffect, useState } from 'react';
import { useLinuxToast } from '@/lib/use-linux-toast';
import { useAuth } from '../store/useAuth';
import { useRouter } from 'next/navigation';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const { success, error } = useLinuxToast();
  const { login: setLoginToken, initAuth, isAuthenticated } = useAuth();
  const router = useRouter()


  useEffect(() => {
    initAuth();
    if (isAuthenticated) {
      router.replace('/root');
    }
  }, [isAuthenticated]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const result = await login(username, password);
      if (result.success) {
        setLoginToken(result.token);
        success(`Welcome ${username}`, result.message);
        router.push('/root');
      } else {
        error("Authentication failed", result.message);
      }
    } catch (err: any) {
      error(err.message);
    }
  };

  return (
    <div className="min-h-screen bg-black text-green-400 flex items-center justify-center px-4">
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="w-full max-w-md"
      >
        <Card className="bg-zinc-900 border-zinc-800 shadow-xl rounded-xl">
          <CardHeader className="text-center space-y-2">
            <motion.div
              className="flex justify-center"
              initial={{ scale: 0.8 }}
              animate={{ scale: 1 }}
              transition={{ duration: 0.3, delay: 0.1 }}
            >
              <Terminal className="w-10 h-10 text-green-400" />
            </motion.div>
            <h1 className="text-xl text-white font-mono tracking-wide">StackRoost Login</h1>
            <p className="text-sm text-zinc-400 font-mono">
              Welcome to your Linux agent dashboard
            </p>
          </CardHeader>

          <CardContent>
            <form className="space-y-4" onSubmit={handleLogin}>
              <div className="relative">
                <Input
                  type="text"
                  placeholder="Username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="bg-zinc-800 text-green-400 placeholder-zinc-500 font-mono pl-10"
                />
                <span className="absolute left-3 top-2.5 text-green-500">
                  <Terminal className="w-4 h-4" />
                </span>
              </div>

              <div className="relative">
                <Input
                  type="password"
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="bg-zinc-800 text-green-400 placeholder-zinc-500 font-mono pl-10"
                />
                <span className="absolute left-3 top-2.5 text-green-500">
                  <KeyRound className="w-4 h-4" />
                </span>
              </div>

              <motion.div
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <Button type="submit" className="w-full font-mono bg-green-600 hover:bg-green-500 text-black">
                  Authenticate
                </Button>
              </motion.div>
            </form>
          </CardContent>
        </Card>

        <motion.div
          className="mt-4 text-center text-xs text-zinc-500 font-mono"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.8 }}
        >
          Secure access. Linux agent v1.0.0
        </motion.div>
      </motion.div>
    </div>
  );
}
