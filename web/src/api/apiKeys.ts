import { apiRequest } from './client';
import type { CollectionResponse, SingleResponse } from './types';

export type ApiKeyScope = 'read-bookings';
export type ApiKeyExpiryMonths = 3 | 12 | 24;

export interface ApiKeyAttributes {
  name: string;
  key_prefix: string;
  scope: ApiKeyScope;
  expires_at: string;
  created_at: string;
  created_by: string;
}

export interface CreatedApiKeyAttributes extends ApiKeyAttributes {
  key: string;
}

export function fetchApiKeys() {
  return apiRequest<CollectionResponse<ApiKeyAttributes>>('/api/v1/api-keys');
}

export function createApiKey(attributes: {
  name: string;
  scope: ApiKeyScope;
  expires_in_months: ApiKeyExpiryMonths;
}) {
  return apiRequest<SingleResponse<CreatedApiKeyAttributes>>('/api/v1/api-keys', {
    method: 'POST',
    body: JSON.stringify({ data: { type: 'api-keys', attributes } })
  });
}

export function deleteApiKey(id: string) {
  return apiRequest<void>(`/api/v1/api-keys/${encodeURIComponent(id)}`, { method: 'DELETE' });
}
