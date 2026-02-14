import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface DayAvailability {
  date: string;
  weekday: string;
  total: number;
  available: number;
}

export interface ItemGroupAvailabilityAttributes {
  item_group_id: string;
  item_group_name: string;
  days: DayAvailability[];
}

export function fetchWeeklyAvailability(areaId: string, week?: string) {
  const params = week ? `?week=${encodeURIComponent(week)}` : '';
  return apiRequest<CollectionResponse<ItemGroupAvailabilityAttributes>>(
    `/api/v1/areas/${areaId}/item-groups/availability${params}`
  );
}
