import { create } from 'zustand';

import { cancelReservation, createReservation, fetchReservations } from '../api/reservation';
import type { Reservation, ReservationCreateInput } from '../types/reservation';

interface ReservationState {
  reservations: Reservation[];
  loadReservations: (exhibitionId: string) => Promise<void>;
  addReservation: (input: ReservationCreateInput) => Promise<Reservation>;
  cancel: (reservationId: string) => Promise<void>;
}

export const useReservationStore = create<ReservationState>((set) => ({
  reservations: [],
  loadReservations: async (exhibitionId) => set({ reservations: await fetchReservations(exhibitionId) }),
  addReservation: async (input) => createReservation(input),
  cancel: async (reservationId) => {
    const updated = await cancelReservation(reservationId);
    set((state) => ({
      reservations: state.reservations.map((item) => (item.id === reservationId ? updated : item)),
    }));
  },
}));
