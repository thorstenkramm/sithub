import { ref } from 'vue';
import { fetchMyBookings, type MyBookingAttributes } from '../api/bookings';
import type { JsonApiResource } from '../api/types';

/** Context describing the booking the user is about to create. */
export interface AreaDayGuardContext {
  /** Area the new item belongs to; empty means "cannot determine" (guard no-ops). */
  areaId: string;
  /** Target booking date (ISO YYYY-MM-DD). */
  date: string;
  /** The item being booked; used to skip a same-item match. */
  newItemId: string;
  /** Human-readable name of the item being booked (for the dialog message). */
  newItemName: string;
  /**
   * When booking on a colleague's behalf: their user id and display name.
   * The guard then checks the COLLEAGUE's seat — bookings the current user
   * already made for that colleague — instead of the user's own (FR178).
   */
  forUserId?: string;
  forUserName?: string;
}

/**
 * Result of a guard decision.
 *
 * - `proceed` is `true` when the caller may create the new booking.
 * - `existingBookingId` is set only when the user confirmed a swap; it is the
 *   id of the conflicting existing booking the caller must cancel AFTER the new
 *   booking has been created successfully (create-then-cancel, story 36.9 D2).
 */
export interface AreaDayGuardDecision {
  proceed: boolean;
  existingBookingId?: string;
}

/**
 * useAreaDayGuard enforces "at most one booking per area and day" (FR178) on
 * the client — for the current user's own seat, and, when booking on a
 * colleague's behalf (ctx.forUserId), for that colleague's seat. Before a
 * create runs, `await guard(ctx)`:
 *
 * - resolves `{ proceed: true }` immediately when the target seat has no OTHER
 *   booking in the same area on the same day (or the area is undeterminable);
 * - otherwise opens the confirmation dialog and resolves once the user decides:
 *   confirming resolves `{ proceed: true, existingBookingId }` WITHOUT cancelling
 *   anything; cancelling resolves `{ proceed: false }`.
 *
 * The guard never cancels a booking itself. Call sites CREATE the new booking
 * first and only cancel `existingBookingId` on success, so a failed create can
 * never leave the user with neither booking (story 36.9 D2). The reactive shape
 * mirrors useWarningConfirmation so it binds to a ConfirmDialog identically.
 * Each hosting component instantiates its own copy. An optional pre-loaded
 * booking list avoids a redundant fetch for callers that already hold one.
 */
export function useAreaDayGuard() {
  const show = ref(false);
  const existingItemName = ref('');
  const existingDate = ref('');
  const newItemName = ref('');
  // Display name of the colleague the pending booking is for; empty when the
  // conflict concerns the user's own seat. Drives the dialog message variant.
  const conflictForUserName = ref('');
  const loading = ref(false);

  let existingBookingId = '';
  let resolveDecision: ((decision: AreaDayGuardDecision) => void) | null = null;

  // findConflict returns the existing booking that occupies the target seat in
  // the same area on the same day, if any. fetchMyBookings returns rows
  // matching `user_id = ? OR booked_by_user_id = ?`.
  //
  // - Booking for the user themselves (no ctx.forUserId): the target seat is
  //   the user's OWN — plain self-bookings (which carry neither `for_user_name`
  //   nor `is_guest`; the backend omits `booked_for_me` for them) and bookings
  //   a colleague made FOR the current user. Rows the user made on a
  //   colleague's or guest's behalf are excluded, so swapping never cancels a
  //   colleague's or guest's booking (story 36.9 D1).
  // - Booking on a colleague's behalf (ctx.forUserId set): the target seat is
  //   the COLLEAGUE's — rows whose `for_user_id` matches. Only bookings the
  //   current user made for that colleague are visible here; bookings the
  //   colleague made themselves cannot be seen client-side.
  function findConflict(
    bookings: JsonApiResource<MyBookingAttributes>[],
    ctx: AreaDayGuardContext,
  ): JsonApiResource<MyBookingAttributes> | undefined {
    if (!ctx.areaId) return undefined;
    const occupiesTargetSeat = ctx.forUserId
      ? (a: MyBookingAttributes) => a.for_user_id === ctx.forUserId
      : (a: MyBookingAttributes) => !a.for_user_id && !a.for_user_name && a.is_guest !== true;
    return bookings.find(
      (b) =>
        occupiesTargetSeat(b.attributes) &&
        b.attributes.area_id === ctx.areaId &&
        b.attributes.booking_date === ctx.date &&
        b.attributes.item_id !== ctx.newItemId,
    );
  }

  function reset() {
    show.value = false;
    existingItemName.value = '';
    existingDate.value = '';
    newItemName.value = '';
    conflictForUserName.value = '';
    loading.value = false;
    existingBookingId = '';
    resolveDecision = null;
  }

  async function guard(
    ctx: AreaDayGuardContext,
    preloaded?: JsonApiResource<MyBookingAttributes>[],
  ): Promise<AreaDayGuardDecision> {
    let bookings = preloaded;
    if (!bookings) {
      try {
        const resp = await fetchMyBookings();
        bookings = resp.data;
      } catch {
        // On a fetch failure, proceed rather than block; the backend still
        // enforces item-level conflicts.
        return { proceed: true };
      }
    }

    const conflict = findConflict(bookings, ctx);
    if (!conflict) {
      return { proceed: true };
    }

    existingBookingId = conflict.id;
    existingItemName.value = conflict.attributes.item_name;
    existingDate.value = conflict.attributes.booking_date;
    newItemName.value = ctx.newItemName;
    conflictForUserName.value = ctx.forUserId
      ? (ctx.forUserName || conflict.attributes.for_user_name || '')
      : '';
    show.value = true;
    return new Promise<AreaDayGuardDecision>((resolve) => {
      resolveDecision = resolve;
    });
  }

  // confirm resolves the pending decision with the conflicting booking id so
  // the caller can create-then-cancel. It does NOT cancel here (story 36.9 D2).
  function confirm(): void {
    const decide = resolveDecision;
    const bookingId = existingBookingId;
    reset();
    decide?.({ proceed: true, existingBookingId: bookingId });
  }

  function cancel() {
    const decide = resolveDecision;
    reset();
    decide?.({ proceed: false });
  }

  return {
    show,
    existingItemName,
    existingDate,
    newItemName,
    conflictForUserName,
    loading,
    guard,
    confirm,
    cancel,
  };
}
