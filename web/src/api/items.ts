import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface ItemAttributes {
  name: string;
  equipment: string[];
  availability: 'available' | 'occupied';
  warning?: string;
  booker_name?: string; // present when item is occupied
  booking_id?: string; // admin-only, present when item is occupied
  note?: string; // present when item is occupied and has a note
}

export function fetchItems(itemGroupId: string, date?: string) {
  const params = date ? `?date=${encodeURIComponent(date)}` : '';
  return apiRequest<CollectionResponse<ItemAttributes>>(`/api/v1/item-groups/${itemGroupId}/items${params}`);
}
