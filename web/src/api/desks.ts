import { apiRequest } from './client';
import type { CollectionResponse } from './types';

export interface DeskAttributes {
  name: string;
  equipment: string[];
  warning?: string;
}

export function fetchDesks(roomId: string) {
  return apiRequest<CollectionResponse<DeskAttributes>>(`/api/v1/rooms/${roomId}/desks`);
}
