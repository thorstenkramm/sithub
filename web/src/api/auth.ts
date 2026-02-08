import { apiRequest } from './client';
import type { SingleResponse } from './types';
import type { UserAttributes } from './me';

export function loginLocal(email: string, password: string) {
  return apiRequest<SingleResponse<UserAttributes>>('/api/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password })
  });
}

export async function logout(): Promise<void> {
  try {
    await fetch('/api/v1/auth/logout', { method: 'POST' });
  } catch {
    // Silently ignore network errors — the user is logging out regardless.
  }
}
