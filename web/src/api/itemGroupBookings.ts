import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface ItemGroupBookingAttributes {
  item_id: string;
  item_name: string;
  user_id: string;
  user_name: string;
  booking_date: string;
  is_guest?: boolean;
  note: string;
}

export function fetchItemGroupBookings(itemGroupId: string, date?: string) {
  const params = date ? `?date=${encodeURIComponent(date)}` : '';
  return apiRequest<CollectionResponse<ItemGroupBookingAttributes>>(
    `/api/v1/item-groups/${itemGroupId}/bookings${params}`
  );
}
