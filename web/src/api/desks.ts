import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface DeskAttributes {
  name: string;
  equipment: string[];
  availability: 'available' | 'occupied';
  warning?: string;
  // Admin-only fields (present when user is admin and desk is occupied)
  booking_id?: string;
  booker_name?: string;
}

export function fetchDesks(roomId: string, date?: string) {
  const params = date ? `?date=${encodeURIComponent(date)}` : '';
  return apiRequest<CollectionResponse<DeskAttributes>>(`/api/v1/rooms/${roomId}/desks${params}`);
}
