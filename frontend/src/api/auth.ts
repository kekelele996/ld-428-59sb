import { apiPaths } from '../constants/apiPaths';
import { request } from '../utils/request';

export type Role = 'Admin' | 'Curator' | 'Artist' | 'Viewer';

export interface AuthUser {
  id: string;
  username: string;
  name: string;
  role: Role;
}

interface LoginResponse {
  token: string;
  user: AuthUser;
}

export async function login(username: string, password: string): Promise<LoginResponse> {
  return request<LoginResponse>(apiPaths.authLogin, {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  });
}
