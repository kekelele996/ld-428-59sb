import { apiPaths } from '../constants/apiPaths';
import type { Session, SessionCreateInput } from '../types/reservation';
import { request } from '../utils/request';

export async function fetchSessions(exhibitionId: string): Promise<Session[]> {
  return request<Session[]>(apiPaths.sessions(exhibitionId));
}

export async function createSession(exhibitionId: string, input: SessionCreateInput): Promise<Session> {
  return request<Session>(apiPaths.sessions(exhibitionId), {
    method: 'POST',
    body: JSON.stringify(input),
  });
}
