import { apiRequest } from './client';
import type { SingleResponse } from './types';

export interface BookingAttributes {
  desk_id: string;
  user_id: string;
  booking_date: string;
  created_at: string;
}

export interface CreateBookingPayload {
  data: {
    type: 'bookings';
    attributes: {
      desk_id: string;
      booking_date: string;
    };
  };
}

export function createBooking(deskId: string, bookingDate: string) {
  const payload: CreateBookingPayload = {
    data: {
      type: 'bookings',
      attributes: {
        desk_id: deskId,
        booking_date: bookingDate
      }
    }
  };

  return apiRequest<SingleResponse<BookingAttributes>>('/api/v1/bookings', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}
