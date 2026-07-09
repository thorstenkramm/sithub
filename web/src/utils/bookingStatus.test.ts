import { deriveBookingStatus } from './bookingStatus';
import type { MyBookingAttributes } from '../api/bookings';

function makeAttrs(overrides: Partial<MyBookingAttributes> = {}): MyBookingAttributes {
  return {
    item_id: 'desk-1',
    item_name: 'Desk 1',
    item_group_id: 'ig-1',
    item_group_name: 'Group 1',
    area_id: 'area-1',
    area_name: 'Area 1',
    booking_date: '2026-07-10',
    created_at: '2026-07-01T00:00:00Z',
    booked_by_user_id: '',
    booked_by_user_name: '',
    booked_for_me: false,
    note: '',
    ...overrides
  };
}

describe('deriveBookingStatus', () => {
  it('returns "guest" for a guest booking', () => {
    expect(deriveBookingStatus(makeAttrs({ is_guest: true }))).toBe('guest');
  });

  it('prioritizes guest over booked-for-me and on-behalf', () => {
    const attrs = makeAttrs({
      is_guest: true,
      booked_for_me: true,
      booked_by_user_id: 'colleague-1'
    });
    expect(deriveBookingStatus(attrs)).toBe('guest');
  });

  it('returns "booked-for-me" when a colleague booked for the user', () => {
    const attrs = makeAttrs({ booked_for_me: true, booked_by_user_id: 'colleague-1' });
    expect(deriveBookingStatus(attrs)).toBe('booked-for-me');
  });

  it('returns "on-behalf" when the user booked for someone else', () => {
    const attrs = makeAttrs({ booked_by_user_id: 'colleague-1', booked_for_me: false });
    expect(deriveBookingStatus(attrs)).toBe('on-behalf');
  });

  it('returns null for a plain self-booking', () => {
    expect(deriveBookingStatus(makeAttrs())).toBeNull();
  });
});
