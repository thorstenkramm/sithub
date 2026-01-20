import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface PresenceAttributes {
  user_id: string;
  user_name: string;
  desk_id: string;
  desk_name: string;
  room_id: string;
  room_name: string;
}

export async function fetchAreaPresence(
  areaId: string,
  date: string
): Promise<CollectionResponse<PresenceAttributes>> {
  return apiRequest<CollectionResponse<PresenceAttributes>>(
    `/api/v1/areas/${areaId}/presence?date=${date}`
  );
}
