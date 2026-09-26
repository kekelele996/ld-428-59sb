// Session 展览预约场次。
export interface Session {
  id: string;
  exhibitionId: string;
  date: string; // YYYY-MM-DD
  startTime: string; // HH:MM
  endTime: string;
  capacity: number;
  bookedCount: number;
  createdAt: string;
  updatedAt: string;
}

export const ReservationStatus = {
  Booked: 'Booked',
  Canceled: 'Canceled',
} as const;

export type ReservationStatus = (typeof ReservationStatus)[keyof typeof ReservationStatus];

// Reservation 预约记录。
export interface Reservation {
  id: string;
  exhibitionId: string;
  sessionId: string;
  phone: string;
  partySize: number;
  status: ReservationStatus;
  createdAt: string;
  canceledAt?: string;
}

export interface ReservationCreateInput {
  sessionId: string;
  phone: string;
  partySize: number;
}

export interface SessionCreateInput {
  date: string;
  startTime: string;
  endTime: string;
  capacity: number;
}
