import { ReservationStatus, VisitSession } from './enums';

export interface Reservation {
  id: string;
  exhibitionId: string;
  visitDate: string;
  session: VisitSession;
  visitorCount: number;
  phone: string;
  status: ReservationStatus;
  createdAt: string;
}

export interface SessionAvailability {
  session: VisitSession;
  capacity: number;
  reserved: number;
  remaining: number;
}

export interface ReservationCreateInput {
  visitDate: string;
  session: VisitSession;
  visitorCount: number;
  phone: string;
}

export interface ReservationCreateResult {
  reservation: Reservation;
  duplicated: boolean;
}
