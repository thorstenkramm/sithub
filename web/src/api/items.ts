import { apiRequest } from "./client";
import type { CollectionResponse } from "./types";

export interface ItemAttributes {
  name: string;
  equipment: string[];
  availability: "available" | "occupied";
  warning?: string;
  icon?: string;
  booker_name?: string; // present when item is occupied
  booker_user_id?: string; // present when item is occupied (non-guest)
  booked_by_me?: boolean; // present when item is occupied
  booking_id?: string; // admin-only, present when item is occupied
  note?: string; // present when item is occupied and has a note
  reserved?: boolean; // true when item is reserved for other users
}

export function fetchItems(itemGroupId: string, date?: string) {
  const params = date ? `?date=${encodeURIComponent(date)}` : "";
  return apiRequest<CollectionResponse<ItemAttributes>>(
    `/api/v1/item-groups/${itemGroupId}/items${params}`,
  );
}
