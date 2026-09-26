import { create } from 'zustand';

import { login as loginApi, type AuthUser } from '../api/auth';

const STORAGE_KEY = 'artvault-auth';

interface AuthState {
  token: string | null;
  user: AuthUser | null;
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<AuthUser>;
  logout: () => void;
}

interface PersistedAuth {
  token: string;
  user: AuthUser;
}

function readPersisted(): PersistedAuth | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? (JSON.parse(raw) as PersistedAuth) : null;
  } catch {
    return null;
  }
}

const persisted = readPersisted();

export const useAuthStore = create<AuthState>((set) => ({
  token: persisted?.token ?? null,
  user: persisted?.user ?? null,
  isAuthenticated: Boolean(persisted),
  login: async (username, password) => {
    const { token, user } = await loginApi(username, password);
    localStorage.setItem(STORAGE_KEY, JSON.stringify({ token, user }));
    set({ token, user, isAuthenticated: true });
    return user;
  },
  logout: () => {
    localStorage.removeItem(STORAGE_KEY);
    set({ token: null, user: null, isAuthenticated: false });
  },
}));
