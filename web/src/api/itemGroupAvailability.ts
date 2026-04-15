import { apiRequest } from './client';
import type { CollectionResponse } from './types';
import { buildWeekQuery } from './weekQuery';

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

export function fetchWeeklyAvailability(areaId: string, week?: string, days?: number) {
  const qs = buildWeekQuery(week, days);
  return apiRequest<CollectionResponse<ItemGroupAvailabilityAttributes>>(
    `/api/v1/areas/${areaId}/item-groups/availability${qs ? `?${qs}` : ''}`
  );
}
