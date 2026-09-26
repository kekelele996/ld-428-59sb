export interface ApiEnvelope<T> {
  code: number;
  message: string;
  data: T;
}

// ApiError 携带后端返回的业务错误码与中文提示。
export class ApiError extends Error {
  status: number;
  code: number;

  constructor(status: number, code: number, message: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
  }
}

function readToken(): string | null {
  try {
    const raw = localStorage.getItem('artvault-auth');
    if (!raw) return null;
    return (JSON.parse(raw) as { token?: string }).token ?? null;
  } catch {
    return null;
  }
}

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const token = readToken();
  const headers: Record<string, string> = { 'Content-Type': 'application/json', ...(init?.headers as Record<string, string> | undefined) };
  if (token) headers.Authorization = `Bearer ${token}`;

  const response = await fetch(path, { ...init, headers });

  let payload: ApiEnvelope<T> | T | null = null;
  const text = await response.text();
  if (text) {
    payload = JSON.parse(text) as ApiEnvelope<T> | T;
  }

  if (!response.ok) {
    const envelope = payload as ApiEnvelope<T> | null;
    throw new ApiError(response.status, envelope?.code ?? response.status, envelope?.message ?? `Request failed: ${response.status}`);
  }

  if (payload && typeof payload === 'object' && 'data' in payload) {
    return (payload as ApiEnvelope<T>).data;
  }
  return payload as T;
}
