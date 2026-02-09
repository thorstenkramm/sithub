import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface ItemAttributes {
  name: string;
  equipment: string[];
  availability: 'available' | 'occupied';
  warning?: string;
  // Admin-only fields (present when user is admin and item is occupied)
  booking_id?: string;
  booker_name?: string;
}

export function fetchItems(itemGroupId: string, date?: string) {
  const params = date ? `?date=${encodeURIComponent(date)}` : '';
  return apiRequest<CollectionResponse<ItemAttributes>>(`/api/v1/item-groups/${itemGroupId}/items${params}`);
}
