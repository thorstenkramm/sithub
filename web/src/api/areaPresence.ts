import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface PresenceAttributes {
  user_id: string;
  user_name: string;
  item_id: string;
  item_name: string;
  item_group_id: string;
  item_group_name: string;
  note: string;
}

export async function fetchAreaPresence(
  areaId: string,
  date: string
): Promise<CollectionResponse<PresenceAttributes>> {
  return apiRequest<CollectionResponse<PresenceAttributes>>(
    `/api/v1/areas/${areaId}/presence?date=${date}`
  );
}
