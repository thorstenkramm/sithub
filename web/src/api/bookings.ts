import { apiRequest, parseErrorDetail, ApiError } from './client';
import type { CollectionResponse, SingleResponse } from './types';

export interface BookingAttributes {
  desk_id: string;
  user_id: string;
  booking_date: string;
  created_at: string;
}

export interface MyBookingAttributes {
  desk_id: string;
  desk_name: string;
  room_id: string;
  room_name: string;
  area_id: string;
  area_name: string;
  booking_date: string;
  created_at: string;
  booked_by_user_id: string;
  booked_by_user_name: string;
  booked_for_me: boolean;
}

export interface CreateBookingPayload {
  data: {
    type: 'bookings';
    attributes: {
      desk_id: string;
      booking_date: string;
      for_user_id?: string;
      for_user_name?: string;
    };
  };
}

export interface BookOnBehalfOptions {
  forUserId: string;
  forUserName: string;
}

export function createBooking(
  deskId: string,
  bookingDate: string,
  onBehalf?: BookOnBehalfOptions
) {
  const payload: CreateBookingPayload = {
    data: {
      type: 'bookings',
      attributes: {
        desk_id: deskId,
        booking_date: bookingDate,
        ...(onBehalf && {
          for_user_id: onBehalf.forUserId,
          for_user_name: onBehalf.forUserName
        })
      }
    }
  };

  return apiRequest<SingleResponse<BookingAttributes>>('/api/v1/bookings', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export function fetchMyBookings() {
  return apiRequest<CollectionResponse<MyBookingAttributes>>('/api/v1/bookings');
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
