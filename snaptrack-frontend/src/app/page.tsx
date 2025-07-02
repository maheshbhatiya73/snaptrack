'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from './store/useAuth';

export default function HomeRedirect() {
  const router = useRouter();
  const { isAuthenticated, initAuth } = useAuth();

  useEffect(() => {
    initAuth();
    if (isAuthenticated) {
      router.replace('/root');
    } else {
      router.replace('/auth/');
    }
  }, [isAuthenticated]);

  return null;
}
