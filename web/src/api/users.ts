import { ApiError, CONNECTION_LOST_MESSAGE, apiRequest, throwApiError } from './client';
import type { CollectionResponse, PaginatedCollectionResponse, SingleResponse } from './types';

export interface UserListAttributes {
  display_name: string;
  email: string;
  is_admin: boolean;
  is_data_exporter: boolean;
  auth_source: string;
  role: string;
  /** RFC 3339 timestamp of the last login, empty if the user never logged in. */
  last_login: string;
  deleted: boolean;
  created_at: string;
  updated_at: string;
}

export interface UserDeletionImpactAttributes {
  has_bookings: boolean;
  future_booking_count: number;
  mode: 'delete' | 'deactivate';
}

export interface PageRequest {
  limit: number;
  offset: number;
}

/** Fetches one page of users (admin only), ordered by display name. */
export function fetchUsers({ limit, offset }: PageRequest) {
  const params = new URLSearchParams({
    'page[limit]': String(limit),
    'page[offset]': String(offset)
  });
  return apiRequest<PaginatedCollectionResponse<UserListAttributes>>(`/api/v1/users?${params.toString()}`);
}

export interface UpdateUserAttributes {
  is_data_exporter?: boolean;
}

/** Updates the given attributes of a user (admin only); omitted attributes stay unchanged. */
export function updateUser(id: string, attributes: UpdateUserAttributes) {
  return apiRequest<SingleResponse<UserListAttributes>>(`/api/v1/users/${id}`, {
    method: 'PATCH',
    body: JSON.stringify({ data: { type: 'users', id, attributes } })
  });
}

export function fetchUserDeletionImpact(id: string) {
  return apiRequest<SingleResponse<UserDeletionImpactAttributes>>(`/api/v1/users/${id}/deletion-impact`);
}

export async function deleteUser(id: string) {
  let response: Response;
  try {
    response = await fetch(`/api/v1/users/${id}`, {
      method: 'DELETE',
      headers: { Accept: 'application/vnd.api+json' }
    });
  } catch {
    throw new ApiError(CONNECTION_LOST_MESSAGE, 0);
  }
  if (!response.ok) {
    await throwApiError(response);
  }
  if (response.status === 204) {
    return null;
  }
  return response.json() as Promise<SingleResponse<UserListAttributes>>;
}

export interface ColleagueAttributes {
  display_name: string;
}

export function fetchColleagues() {
  return apiRequest<CollectionResponse<ColleagueAttributes>>('/api/v1/colleagues');
}
