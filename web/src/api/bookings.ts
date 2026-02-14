import { apiRequest, parseErrorDetail, ApiError } from './client';
import type { CollectionResponse, SingleResponse } from './types';

export interface BookingAttributes {
  item_id: string;
  user_id: string;
  booking_date: string;
  created_at: string;
  note: string;
}

export interface MyBookingAttributes {
  item_id: string;
  item_name: string;
  item_group_id: string;
  item_group_name: string;
  area_id: string;
  area_name: string;
  booking_date: string;
  created_at: string;
  booked_by_user_id: string;
  booked_by_user_name: string;
  booked_for_me: boolean;
  is_guest?: boolean;
  guest_name?: string;
  guest_email?: string;
  note: string;
}

export interface CreateBookingPayload {
  data: {
    type: 'bookings';
    attributes: {
      item_id: string;
      booking_date?: string;
      booking_dates?: string[];
      for_user_id?: string;
      for_user_name?: string;
      is_guest?: boolean;
      guest_email?: string;
    };
  };
}

export interface MultiDayBookingResult {
  created: Array<{
    type: string;
    id: string;
    attributes: BookingAttributes;
  }>;
  conflicts?: string[];
}

export interface BookOnBehalfOptions {
  forUserId: string;
  forUserName?: string;
}

export interface GuestBookingOptions {
  guestName: string;
  guestEmail?: string;
}

function buildBookingAttributes(params: {
  itemId: string;
  bookingDate?: string;
  bookingDates?: string[];
  onBehalf?: BookOnBehalfOptions;
  guest?: GuestBookingOptions;
}): CreateBookingPayload['data']['attributes'] {
  const attrs: CreateBookingPayload['data']['attributes'] = {
    item_id: params.itemId
  };

  if (params.bookingDate) {
    attrs.booking_date = params.bookingDate;
  }
  if (params.bookingDates) {
    attrs.booking_dates = params.bookingDates;
  }

  if (params.onBehalf) {
    attrs.for_user_id = params.onBehalf.forUserId;
    if (params.onBehalf.forUserName) {
      attrs.for_user_name = params.onBehalf.forUserName;
    }
  }

  if (params.guest) {
    attrs.is_guest = true;
    attrs.for_user_name = params.guest.guestName;
    if (params.guest.guestEmail) {
      attrs.guest_email = params.guest.guestEmail;
    }
  }

  return attrs;
}

export function createBooking(
  itemId: string,
  bookingDate: string,
  onBehalf?: BookOnBehalfOptions,
  guest?: GuestBookingOptions
) {
  const payload: CreateBookingPayload = {
    data: {
      type: 'bookings',
      attributes: buildBookingAttributes({
        itemId,
        bookingDate,
        onBehalf,
        guest
      })
    }
  };

  return apiRequest<SingleResponse<BookingAttributes>>('/api/v1/bookings', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export function createMultiDayBooking(
  itemId: string,
  bookingDates: string[],
  onBehalf?: BookOnBehalfOptions,
  guest?: GuestBookingOptions
) {
  const payload: CreateBookingPayload = {
    data: {
      type: 'bookings',
      attributes: buildBookingAttributes({
        itemId,
        bookingDates,
        onBehalf,
        guest
      })
    }
  };

  return apiRequest<MultiDayBookingResult>('/api/v1/bookings', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export function fetchMyBookings() {
  return apiRequest<CollectionResponse<MyBookingAttributes>>('/api/v1/bookings');
}

export interface BookingHistoryParams {
  from?: string;
  to?: string;
}

export function fetchBookingHistory(params?: BookingHistoryParams) {
  const queryParams = new URLSearchParams();
  if (params?.from) queryParams.set('from', params.from);
  if (params?.to) queryParams.set('to', params.to);

  const queryString = queryParams.toString();
  const url = queryString ? `/api/v1/bookings/history?${queryString}` : '/api/v1/bookings/history';

  return apiRequest<CollectionResponse<MyBookingAttributes>>(url);
}

export function updateBookingNote(bookingId: string, note: string) {
  return apiRequest<SingleResponse<BookingAttributes>>(`/api/v1/bookings/${bookingId}`, {
    method: 'PATCH',
    body: JSON.stringify({
      data: {
        type: 'bookings',
        id: bookingId,
        attributes: { note }
      }
    })
  });
}

export async function cancelBooking(bookingId: string): Promise<void> {
  const response = await fetch(`/api/v1/bookings/${bookingId}`, {
    method: 'DELETE',
    headers: {
      Accept: 'application/vnd.api+json'
    }
  });

  if (!response.ok) {
    const detail = await parseErrorDetail(response);
    throw new ApiError(`Request failed: ${response.status}`, response.status, detail);
  }
}
