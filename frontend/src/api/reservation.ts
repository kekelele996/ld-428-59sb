import { apiPaths } from '../constants/apiPaths';
import type { Reservation, ReservationCreateInput } from '../types/reservation';
import { request } from '../utils/request';

export async function createReservation(input: ReservationCreateInput): Promise<Reservation> {
  return request<Reservation>(apiPaths.reservations, {
    method: 'POST',
    body: JSON.stringify(input),
  });
}

export async function fetchReservations(exhibitionId: string): Promise<Reservation[]> {
  return request<Reservation[]>(apiPaths.exhibitionReservations(exhibitionId));
}

export async function cancelReservation(reservationId: string): Promise<Reservation> {
  return request<Reservation>(apiPaths.cancelReservation(reservationId), { method: 'POST' });
}
