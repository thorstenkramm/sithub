import { apiClient } from './client';
import type { JsonApiCollectionResponse } from './types';

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
): Promise<JsonApiCollectionResponse<PresenceAttributes>> {
  return apiClient.get<JsonApiCollectionResponse<PresenceAttributes>>(
    `/api/v1/areas/${areaId}/presence?date=${date}`
  );
}
