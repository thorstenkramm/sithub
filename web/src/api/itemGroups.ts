import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface ItemGroupAttributes {
  name: string;
  description?: string;
  floor_plan?: string;
}

export function fetchItemGroups(areaId: string) {
  return apiRequest<CollectionResponse<ItemGroupAttributes>>(`/api/v1/areas/${areaId}/item-groups`);
}
