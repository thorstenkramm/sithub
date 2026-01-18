import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface AreaAttributes {
  name: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export function fetchAreas() {
  return apiRequest<CollectionResponse<AreaAttributes>>('/api/v1/areas');
}
