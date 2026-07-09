import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useAreaDayGuard } from '../useAreaDayGuard';
import { fetchMyBookings } from '../../api/bookings';

vi.mock('../../api/bookings', () => ({
  fetchMyBookings: vi.fn()
}));

const fetchMyBookingsMock = vi.mocked(fetchMyBookings);

// Mirrors the real API: a plain self-booking carries NO booked_by/booked_for_me/
// for_user_name/is_guest attributes at all (the backend omits them). Extra
// attributes for on-behalf/guest rows are supplied via `extra`.
function booking(
  id: string,
  areaId: string,
  date: string,
  itemId: string,
  itemName = itemId,
  extra: Record<string, unknown> = {}
) {
  return {
    id,
    type: 'bookings' as const,
    attributes: {
      item_id: itemId,
      item_name: itemName,
      item_group_id: 'ig',
      item_group_name: 'IG',
      area_id: areaId,
      area_name: 'Area',
      booking_date: date,
      created_at: '',
      note: '',
      ...extra
    }
  };
}

const ctx = {
  areaId: 'area-1',
  date: '2099-01-01',
  newItemId: 'desk-2',
  newItemName: 'Desk 2'
};

describe('useAreaDayGuard', () => {
  beforeEach(() => {
    fetchMyBookingsMock.mockReset();
  });

  it('proceeds without a dialog when no conflicting booking exists', async () => {
    fetchMyBookingsMock.mockResolvedValue({ data: [] } as never);
    const guard = useAreaDayGuard();

    const decision = await guard.guard(ctx);

    expect(decision).toEqual({ proceed: true });
    expect(guard.show.value).toBe(false);
  });

  it('ignores bookings in a different area or on a different day', async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        booking('b1', 'area-OTHER', '2099-01-01', 'desk-1'),
        booking('b2', 'area-1', '2099-01-02', 'desk-1')
      ]
    } as never);
    const guard = useAreaDayGuard();

    const decision = await guard.guard(ctx);

    expect(decision).toEqual({ proceed: true });
    expect(guard.show.value).toBe(false);
  });

  it('ignores on-behalf and guest bookings so only the user\'s OWN seat conflicts (D1)', async () => {
    // for_user_id/for_user_name: made on a colleague's behalf; is_guest: made
    // for a guest. Neither occupies the current user's seat, so neither must be
    // swapped. The first row intentionally omits for_user_name to match a
    // backend display-name lookup miss.
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        booking('b1', 'area-1', '2099-01-01', 'desk-1', 'Desk 1', {
          booked_by_user_id: 'me',
          for_user_id: 'u-1'
        }),
        booking('b2', 'area-1', '2099-01-01', 'desk-3', 'Desk 3', {
          is_guest: true,
          guest_name: 'Guest'
        })
      ]
    } as never);
    const guard = useAreaDayGuard();

    const decision = await guard.guard(ctx);

    expect(decision).toEqual({ proceed: true });
    expect(guard.show.value).toBe(false);
  });

  it('scopes the conflict to the colleague\'s seat when forUserId is set', async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        // The user's own booking — must NOT conflict with a colleague booking.
        booking('own-1', 'area-1', '2099-01-01', 'desk-1', 'Desk 1'),
        // A booking the user already made for colleague u-1.
        booking('col-1', 'area-1', '2099-01-01', 'desk-3', 'Desk 3', {
          booked_by_user_id: 'me',
          for_user_id: 'u-1',
          for_user_name: 'Jane Doe'
        })
      ]
    } as never);
    const guard = useAreaDayGuard();

    const decision = guard.guard({ ...ctx, forUserId: 'u-1', forUserName: 'Jane Doe' });
    await Promise.resolve();
    await Promise.resolve();

    expect(guard.show.value).toBe(true);
    expect(guard.existingItemName.value).toBe('Desk 3');
    expect(guard.conflictForUserName.value).toBe('Jane Doe');

    guard.confirm();
    await expect(decision).resolves.toEqual({ proceed: true, existingBookingId: 'col-1' });
  });

  it('proceeds for a colleague booking when only the user themselves has a conflict', async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [booking('own-1', 'area-1', '2099-01-01', 'desk-1', 'Desk 1')]
    } as never);
    const guard = useAreaDayGuard();

    const decision = await guard.guard({ ...ctx, forUserId: 'u-1', forUserName: 'Jane Doe' });

    expect(decision).toEqual({ proceed: true });
    expect(guard.show.value).toBe(false);
  });

  it('proceeds for a colleague booking when the conflict belongs to a DIFFERENT colleague', async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        booking('col-2', 'area-1', '2099-01-01', 'desk-3', 'Desk 3', {
          booked_by_user_id: 'me',
          for_user_id: 'u-2',
          for_user_name: 'Other Person'
        })
      ]
    } as never);
    const guard = useAreaDayGuard();

    const decision = await guard.guard({ ...ctx, forUserId: 'u-1', forUserName: 'Jane Doe' });

    expect(decision).toEqual({ proceed: true });
    expect(guard.show.value).toBe(false);
  });

  it('treats a booking a colleague made FOR the user as the user\'s own seat', async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [
        booking('b1', 'area-1', '2099-01-01', 'desk-1', 'Desk 1', {
          booked_by_user_id: 'colleague',
          booked_by_user_name: 'Colleague',
          booked_for_me: true
        })
      ]
    } as never);
    const guard = useAreaDayGuard();

    const decision = guard.guard(ctx);
    await Promise.resolve();
    await Promise.resolve();

    expect(guard.show.value).toBe(true);
    guard.cancel();
    await expect(decision).resolves.toEqual({ proceed: false });
  });

  it('no-ops when the area cannot be determined (empty areaId)', async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [booking('b1', 'area-1', '2099-01-01', 'desk-1')]
    } as never);
    const guard = useAreaDayGuard();

    const decision = await guard.guard({ ...ctx, areaId: '' });

    expect(decision).toEqual({ proceed: true });
    expect(guard.show.value).toBe(false);
  });

  it('opens the dialog on a conflict and returns the id WITHOUT cancelling (D2)', async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [booking('b1', 'area-1', '2099-01-01', 'desk-1', 'Desk 1')]
    } as never);
    const guard = useAreaDayGuard();

    const decision = guard.guard(ctx);
    // Give the fetch microtask a chance to resolve.
    await Promise.resolve();
    await Promise.resolve();

    expect(guard.show.value).toBe(true);
    expect(guard.existingItemName.value).toBe('Desk 1');
    expect(guard.existingDate.value).toBe('2099-01-01');
    expect(guard.newItemName.value).toBe('Desk 2');

    // Confirm resolves with the conflicting id for a caller-side create-then-
    // cancel; the guard itself never cancels.
    guard.confirm();
    await expect(decision).resolves.toEqual({ proceed: true, existingBookingId: 'b1' });
    expect(guard.show.value).toBe(false);
  });

  it('resolves proceed=false and makes no change when cancelled', async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [booking('b1', 'area-1', '2099-01-01', 'desk-1', 'Desk 1')]
    } as never);
    const guard = useAreaDayGuard();

    const decision = guard.guard(ctx);
    await Promise.resolve();
    await Promise.resolve();
    expect(guard.show.value).toBe(true);

    guard.cancel();
    await expect(decision).resolves.toEqual({ proceed: false });
    expect(guard.show.value).toBe(false);
  });

  it('uses a preloaded booking list without fetching', async () => {
    const guard = useAreaDayGuard();
    const preloaded = [booking('b1', 'area-1', '2099-01-01', 'desk-1', 'Desk 1')];

    const decision = guard.guard(ctx, preloaded as never);
    await Promise.resolve();

    expect(fetchMyBookingsMock).not.toHaveBeenCalled();
    expect(guard.show.value).toBe(true);

    guard.confirm();
    await expect(decision).resolves.toEqual({ proceed: true, existingBookingId: 'b1' });
  });
});
