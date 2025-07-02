// store/useAuth.ts
import { create } from 'zustand';

type AuthState = {
  token: string | null;
  isAuthenticated: boolean;
  login: (token: string) => void;
  logout: () => void;
  initAuth: () => void;
};

export const useAuth = create<AuthState>((set) => ({
  token: null,
  isAuthenticated: false,

  login: (token: string) => {
    localStorage.setItem('auth_token', token);
    set({ token, isAuthenticated: true });
  },

  logout: () => {
    localStorage.removeItem('auth_token');
    set({ token: null, isAuthenticated: false });
  },

  initAuth: () => {
    const storedToken = localStorage.getItem('auth_token');
    if (storedToken) {
      set({ token: storedToken, isAuthenticated: true });
    }
  },
}));
