import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface RoomBookingAttributes {
  desk_id: string;
  desk_name: string;
  user_id: string;
  user_name: string;
  booking_date: string;
}

export function fetchRoomBookings(roomId: string, date?: string) {
  const params = date ? `?date=${encodeURIComponent(date)}` : '';
  return apiRequest<CollectionResponse<RoomBookingAttributes>>(
    `/api/v1/rooms/${roomId}/bookings${params}`
  );
}
