// lib/use-linux-toast.ts
'use client';

import { toast } from 'sonner';

export function useLinuxToast() {
  return {
    success: (msg: string, desc?: string) =>
      toast.success(msg, {
        description: desc,
      }),
    error: (msg: string, desc?: string) =>
      toast.error(msg, {
        description: desc,
      }),
    info: (msg: string, desc?: string) =>
      toast(msg, {
        description: desc,
      }),
  };
}
