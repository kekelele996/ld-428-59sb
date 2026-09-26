import { apiPaths } from '../constants/apiPaths';
import type { Reservation, ReservationCreateInput, ReservationCreateResult, SessionAvailability } from '../types/reservation';
import { reservations as mockReservations } from '../utils/mockData';
import { request } from '../utils/request';

// createReservation 提交预约；失败原因（场次已满、未开放等）需要展示给观众，不吞错误。
export async function createReservation(exhibitionId: string, input: ReservationCreateInput): Promise<ReservationCreateResult> {
  return request<ReservationCreateResult>(`${apiPaths.exhibitions}/${exhibitionId}/reservations`, {
    method: 'POST',
    body: JSON.stringify(input),
  });
}

export async function fetchSessionAvailability(exhibitionId: string, visitDate: string): Promise<SessionAvailability[]> {
  try {
    return await request<SessionAvailability[]>(`${apiPaths.exhibitions}/${exhibitionId}/reservations/availability?date=${visitDate}`);
  } catch {
    return [];
  }
}

export async function fetchExhibitionReservations(exhibitionId: string): Promise<Reservation[]> {
  try {
    return await request<Reservation[]>(`${apiPaths.exhibitions}/${exhibitionId}/reservations`);
  } catch {
    return mockReservations.filter((item) => item.exhibitionId === exhibitionId);
  }
}

export async function cancelReservation(id: string): Promise<Reservation> {
  return request<Reservation>(`${apiPaths.reservations}/${id}/cancel`, { method: 'PATCH' });
}
