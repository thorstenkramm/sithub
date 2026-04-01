import { apiRequest } from './client';
import type { CollectionResponse, SingleResponse } from './types';

export interface FloorPlanPositionAttributes {
  floor_plan: string;
  item_id: string;
  label?: string;
  x: number;
  y: number;
  width: number;
  height: number;
  border_width: number;
  created_at: string;
  updated_at: string;
}

export function fetchFloorPlanPositions(floorPlan: string) {
  return apiRequest<CollectionResponse<FloorPlanPositionAttributes>>(
    `/api/v1/floor-plan-positions?floor_plan=${encodeURIComponent(floorPlan)}`
  );
}

export function createFloorPlanPosition(data: {
  floor_plan: string;
  item_id: string;
  label?: string;
  x: number;
  y: number;
  width: number;
  height: number;
  border_width?: number;
}) {
  return apiRequest<SingleResponse<FloorPlanPositionAttributes>>(
    '/api/v1/floor-plan-positions',
    {
      method: 'POST',
      body: JSON.stringify({
        data: {
          type: 'floor-plan-positions',
          attributes: data
        }
      })
    }
  );
}

export function updateFloorPlanPosition(
  id: string,
  data: Partial<{ label: string; x: number; y: number; width: number; height: number; border_width: number }>
) {
  return apiRequest<SingleResponse<FloorPlanPositionAttributes>>(
    `/api/v1/floor-plan-positions/${encodeURIComponent(id)}`,
    {
      method: 'PUT',
      body: JSON.stringify({
        data: {
          type: 'floor-plan-positions',
          attributes: data
        }
      })
    }
  );
}

export function deleteFloorPlanPosition(id: string) {
  return apiRequest<void>(
    `/api/v1/floor-plan-positions/${encodeURIComponent(id)}`,
    { method: 'DELETE' }
  );
}
