import { apiRequest } from './client';
import type { SingleResponse } from './types';

export interface UserAttributes {
  display_name: string;
  is_admin: boolean;
}

export function fetchMe() {
  return apiRequest<SingleResponse<UserAttributes>>('/api/v1/me');
}
