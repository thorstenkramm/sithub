import { apiRequest } from './client';
import type { SingleResponse } from './types';

export interface SettingsAttributes {
  weeks_in_advanced: number;
}

export function fetchSettings() {
  return apiRequest<SingleResponse<SettingsAttributes>>('/api/v1/settings');
}
