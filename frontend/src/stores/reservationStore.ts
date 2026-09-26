import { create } from 'zustand';

import { cancelReservation, createReservation, fetchExhibitionReservations, fetchSessionAvailability } from '../api/reservation';
import type { Reservation, ReservationCreateInput, ReservationCreateResult, SessionAvailability } from '../types/reservation';

interface ReservationState {
  reservations: Reservation[];
  availability: SessionAvailability[];
  submitting: boolean;
  loadReservations: (exhibitionId: string) => Promise<void>;
  loadAvailability: (exhibitionId: string, visitDate: string) => Promise<void>;
  submitReservation: (exhibitionId: string, input: ReservationCreateInput) => Promise<ReservationCreateResult>;
  cancel: (id: string) => Promise<void>;
}

export const useReservationStore = create<ReservationState>((set) => ({
  reservations: [],
  availability: [],
  submitting: false,
  loadReservations: async (exhibitionId) => set({ reservations: await fetchExhibitionReservations(exhibitionId) }),
  loadAvailability: async (exhibitionId, visitDate) => set({ availability: await fetchSessionAvailability(exhibitionId, visitDate) }),
  submitReservation: async (exhibitionId, input) => {
    set({ submitting: true });
    try {
      return await createReservation(exhibitionId, input);
    } finally {
      set({ submitting: false });
    }
  },
  cancel: async (id) => {
    const updated = await cancelReservation(id);
    set((state) => ({ reservations: state.reservations.map((item) => (item.id === id ? updated : item)) }));
  },
}));
