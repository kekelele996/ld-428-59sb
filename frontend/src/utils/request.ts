export interface ApiEnvelope<T> {
  code: number;
  message: string;
  data: T;
}

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) },
    ...init,
  });

  if (!response.ok) {
    let message = `Request failed: ${response.status}`;
    try {
      const payload = (await response.json()) as Partial<ApiEnvelope<unknown>>;
      if (payload && typeof payload.message === 'string' && payload.message) {
        message = payload.message;
      }
    } catch {
      // 响应体非 JSON 时保留默认错误信息
    }
    throw new Error(message);
  }

  const payload = (await response.json()) as ApiEnvelope<T> | T;
  if (payload && typeof payload === 'object' && 'data' in payload) {
    return (payload as ApiEnvelope<T>).data;
  }
  return payload as T;
}
