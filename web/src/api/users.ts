import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface UserListAttributes {
  display_name: string;
  email: string;
  is_admin: boolean;
  auth_source: string;
  role: string;
}

export function fetchUsers() {
  return apiRequest<CollectionResponse<UserListAttributes>>('/api/v1/users');
}
