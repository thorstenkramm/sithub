import { apiRequest } from './client';
import type { SingleResponse } from './types';

export interface VersionAttributes {
  version: string;
}

export function fetchVersion() {
  return apiRequest<SingleResponse<VersionAttributes>>('/api/v1/version');
}
