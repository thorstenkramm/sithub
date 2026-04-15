import { apiRequest } from './client';
import type { CollectionResponse } from './types';
import { buildWeekQuery } from './weekQuery';

export interface MatrixDayMeta {
  date: string;
  weekday: string;
}

export interface MatrixCell {
  date: string;
  availability: 'free' | 'occupied';
  booker_name?: string;
  booker_user_id?: string;
  booked_by_me: boolean;
  booking_id?: string;
}

export interface MatrixItem {
  item_id: string;
  item_name: string;
  equipment: string[];
  warning?: string;
  reserved?: boolean;
  cells: MatrixCell[];
}

export interface ItemGroupMatrixAttributes {
  item_group_id: string;
  item_group_name: string;
  days: MatrixDayMeta[];
  items: MatrixItem[];
}

export function fetchWeeklyMatrix(areaId: string, week?: string, days?: number) {
  const qs = buildWeekQuery(week, days);
  return apiRequest<CollectionResponse<ItemGroupMatrixAttributes>>(
    `/api/v1/areas/${areaId}/item-groups/matrix${qs ? `?${qs}` : ''}`
  );
}
