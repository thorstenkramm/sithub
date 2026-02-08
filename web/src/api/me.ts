import { apiRequest } from './client';
import type { SingleResponse } from './types';

export interface UserAttributes {
  display_name: string;
  email: string;
  is_admin: boolean;
  auth_source: string;
  role: string;
}

export function fetchMe() {
  return apiRequest<SingleResponse<UserAttributes>>('/api/v1/me');
}

export function changePassword(currentPassword: string, newPassword: string) {
  return apiRequest<SingleResponse<UserAttributes>>('/api/v1/me', {
    method: 'PATCH',
    body: JSON.stringify({
      data: {
        attributes: {
          current_password: currentPassword,
          new_password: newPassword
        }
      }
    })
  });
}
