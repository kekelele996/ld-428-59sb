import { create } from 'zustand';

import { createSession, fetchSessions } from '../api/session';
import type { Session, SessionCreateInput } from '../types/reservation';

interface SessionState {
  sessions: Session[];
  loadSessions: (exhibitionId: string) => Promise<void>;
  addSession: (exhibitionId: string, input: SessionCreateInput) => Promise<Session>;
}

export const useSessionStore = create<SessionState>((set) => ({
  sessions: [],
  loadSessions: async (exhibitionId) => set({ sessions: await fetchSessions(exhibitionId) }),
  addSession: async (exhibitionId, input) => {
    const created = await createSession(exhibitionId, input);
    set((state) => ({ sessions: [...state.sessions, created] }));
    return created;
  },
}));
