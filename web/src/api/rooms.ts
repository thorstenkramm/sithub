import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface RoomAttributes {
  name: string;
  description?: string;
  floor_plan?: string;
}

export function fetchRooms(areaId: string) {
  return apiRequest<CollectionResponse<RoomAttributes>>(`/api/v1/areas/${areaId}/rooms`);
}
