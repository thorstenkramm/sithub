import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { createBooking, fetchMyBookings } from './bookings';

const mockFetch = vi.fn();

describe('createBooking', () => {
  beforeEach(() => {
    global.fetch = mockFetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('sends POST request with correct JSON:API payload', async () => {
    const mockResponse = {
      data: {
        type: 'bookings',
        id: 'booking-123',
        attributes: {
          desk_id: 'desk-1',
          user_id: 'user-1',
          booking_date: '2026-01-20',
          created_at: '2026-01-19T10:00:00Z'
        }
      }
    };

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockResponse)
    });

    const result = await createBooking('desk-1', '2026-01-20');

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/bookings', {
      method: 'POST',
      body: JSON.stringify({
        data: {
          type: 'bookings',
          attributes: {
            desk_id: 'desk-1',
            booking_date: '2026-01-20'
          }
        }
      }),
      headers: expect.any(Headers)
    });

    expect(result.data.id).toBe('booking-123');
    expect(result.data.attributes.desk_id).toBe('desk-1');
  });

  // Issue 5: Improve error test to verify status code is propagated
  it('throws ApiError with correct status on conflict response', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 409
    });

    await expect(createBooking('desk-1', '2026-01-20')).rejects.toMatchObject({
      status: 409
    });
  });

  it('throws ApiError with correct status on bad request response', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 400
    });

    await expect(createBooking('desk-1', '2026-01-20')).rejects.toMatchObject({
      status: 400
    });
  });
});

describe('fetchMyBookings', () => {
  beforeEach(() => {
    global.fetch = mockFetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('sends GET request to /api/v1/bookings', async () => {
    const mockResponse = {
      data: [
        {
          id: 'booking-1',
          type: 'bookings',
          attributes: {
            desk_id: 'desk-1',
            desk_name: 'Desk 1',
            room_id: 'room-1',
            room_name: 'Room 101',
            area_id: 'area-1',
            area_name: 'Main Office',
            booking_date: '2026-01-20',
            created_at: '2026-01-19T10:00:00Z'
          }
        }
      ]
    };

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockResponse)
    });

    const result = await fetchMyBookings();

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/bookings', {
      headers: expect.any(Headers)
    });

    expect(result.data).toHaveLength(1);
    expect(result.data[0].attributes.desk_name).toBe('Desk 1');
    expect(result.data[0].attributes.room_name).toBe('Room 101');
    expect(result.data[0].attributes.area_name).toBe('Main Office');
  });

  it('throws ApiError on error response', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 500
    });

    await expect(fetchMyBookings()).rejects.toMatchObject({
      status: 500
    });
  });
});
