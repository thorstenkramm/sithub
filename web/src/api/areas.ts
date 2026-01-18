import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface AreaAttributes {
  name: string;
  description?: string;
  floor_plan?: string;
}

export function fetchAreas() {
  return apiRequest<CollectionResponse<AreaAttributes>>('/api/v1/areas');
}
