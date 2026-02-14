import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { cancelBooking, createBooking, fetchMyBookings, updateBookingNote } from './bookings';

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
          item_id: 'item-1',
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

    const result = await createBooking('item-1', '2026-01-20');

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/bookings', {
      method: 'POST',
      body: JSON.stringify({
        data: {
          type: 'bookings',
          attributes: {
            item_id: 'item-1',
            booking_date: '2026-01-20'
          }
        }
      }),
      headers: expect.any(Headers)
    });

    expect(result.data.id).toBe('booking-123');
    expect(result.data.attributes.item_id).toBe('item-1');
  });

  it('sends POST request with for_user_id and for_user_name when booking on behalf', async () => {
    const mockResponse = {
      data: {
        type: 'bookings',
        id: 'booking-123',
        attributes: {
          item_id: 'item-1',
          user_id: 'colleague-1',
          booking_date: '2026-01-20',
          created_at: '2026-01-19T10:00:00Z'
        }
      }
    };

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockResponse)
    });

    await createBooking('item-1', '2026-01-20', {
      forUserId: 'colleague-1',
      forUserName: 'Jane Doe'
    });

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/bookings', {
      method: 'POST',
      body: JSON.stringify({
        data: {
          type: 'bookings',
          attributes: {
            item_id: 'item-1',
            booking_date: '2026-01-20',
            for_user_id: 'colleague-1',
            for_user_name: 'Jane Doe'
          }
        }
      }),
      headers: expect.any(Headers)
    });
  });

  // Issue 5: Improve error test to verify status code is propagated
  it('throws ApiError with correct status on conflict response', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 409
    });

    await expect(createBooking('item-1', '2026-01-20')).rejects.toMatchObject({
      status: 409
    });
  });

  it('throws ApiError with correct status on bad request response', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 400
    });

    await expect(createBooking('item-1', '2026-01-20')).rejects.toMatchObject({
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
            item_id: 'item-1',
            item_name: 'Desk 1',
            item_group_id: 'ig-1',
            item_group_name: 'Room 101',
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
    expect(result.data[0].attributes.item_name).toBe('Desk 1');
    expect(result.data[0].attributes.item_group_name).toBe('Room 101');
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

describe('updateBookingNote', () => {
  beforeEach(() => {
    global.fetch = mockFetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('sends PATCH request with correct JSON:API payload', async () => {
    const mockResponse = {
      data: {
        type: 'bookings',
        id: 'booking-123',
        attributes: {
          item_id: 'item-1',
          user_id: 'user-1',
          booking_date: '2026-01-20',
          created_at: '2026-01-19T10:00:00Z',
          note: 'Arriving late'
        }
      }
    };

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockResponse)
    });

    const result = await updateBookingNote('booking-123', 'Arriving late');

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/bookings/booking-123', {
      method: 'PATCH',
      body: JSON.stringify({
        data: {
          type: 'bookings',
          id: 'booking-123',
          attributes: { note: 'Arriving late' }
        }
      }),
      headers: expect.any(Headers)
    });

    expect(result.data.attributes.note).toBe('Arriving late');
  });

  it('throws ApiError on error response', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 404
    });

    await expect(updateBookingNote('nonexistent', 'test')).rejects.toMatchObject({
      status: 404
    });
  });
});

describe('cancelBooking', () => {
  beforeEach(() => {
    global.fetch = mockFetch;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('sends DELETE request to /api/v1/bookings/:id', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true
    });

    await cancelBooking('booking-123');

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/bookings/booking-123', {
      method: 'DELETE',
      headers: { Accept: 'application/vnd.api+json' }
    });
  });

  it('throws ApiError on 404 response', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 404,
      json: () =>
        Promise.resolve({
          errors: [{ detail: 'Booking not found' }]
        })
    });

    await expect(cancelBooking('nonexistent')).rejects.toMatchObject({
      status: 404,
      detail: 'Booking not found'
    });
  });

  it('throws ApiError on server error', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 500,
      json: () => Promise.resolve({})
    });

    await expect(cancelBooking('booking-123')).rejects.toMatchObject({
      status: 500
    });
  });
});
