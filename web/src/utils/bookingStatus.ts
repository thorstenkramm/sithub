import type { MyBookingAttributes } from '../api/bookings';

/**
 * The derived booking status used by the My Bookings tiles and table.
 * Returns null for a plain self-booking that has no special relationship.
 */
export type BookingStatus = 'guest' | 'booked-for-me' | 'on-behalf' | null;

/**
 * Derives the display status of a booking from its attributes.
 * Priority: guest, then booked-for-me, then on-behalf; otherwise null (self).
 * Shared between MyBookingsView (table) and BookingCard (tiles) so the
 * derivation lives in one place.
 */
export function deriveBookingStatus(attrs: MyBookingAttributes): BookingStatus {
  if (attrs.is_guest) return 'guest';
  if (attrs.booked_for_me) return 'booked-for-me';
  if (attrs.booked_by_user_id && !attrs.booked_for_me) return 'on-behalf';
  return null;
}
