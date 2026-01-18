import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface DeskAttributes {
  name: string;
  equipment: string[];
  availability: 'available' | 'occupied';
  warning?: string;
}

export function fetchDesks(roomId: string, date?: string) {
  const params = date ? `?date=${encodeURIComponent(date)}` : '';
  return apiRequest<CollectionResponse<DeskAttributes>>(`/api/v1/rooms/${roomId}/desks${params}`);
}
