import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import {
  cancelBooking,
  createBooking,
  fetchBookingExport,
  fetchMyBookings,
  updateBookingNote
} from './bookings';

const mockFetch = vi.fn();

const expectCreateBookingPayload = (attributes: Record<string, unknown>) => {
  expect(mockFetch).toHaveBeenCalledWith('/api/v1/bookings', {
    method: 'POST',
    body: JSON.stringify({
      data: {
        type: 'bookings',
        attributes
      }
    }),
    headers: expect.any(Headers)
  });
};


const mockCreateBookingResponse = (bookingId: string) => ({
  data: {
    type: 'bookings',
    id: bookingId,
    attributes: {
      item_id: 'item-1',
      user_id: 'colleague-1',
      booking_date: '2026-01-20',
      created_at: '2026-01-19T10:00:00Z'
    }
  }
});

const mockFetchSuccess = (response: unknown) => {
  mockFetch.mockResolvedValueOnce({
    ok: true,
    json: () => Promise.resolve(response)
  });
};
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
    const mockResponse = mockCreateBookingResponse('booking-123');
    mockFetchSuccess(mockResponse);

    await createBooking('item-1', '2026-01-20', {
      forUserId: 'colleague-1',
      forUserName: 'Jane Doe'
    });

    expectCreateBookingPayload({
      item_id: 'item-1',
      booking_date: '2026-01-20',
      for_user_id: 'colleague-1',
      for_user_name: 'Jane Doe'
    });
  });

  it('omits for_user_name when booking on behalf without a display name', async () => {
    const mockResponse = mockCreateBookingResponse('booking-124');
    mockFetchSuccess(mockResponse);

    await createBooking('item-1', '2026-01-20', {
      forUserId: 'colleague-1'
    });

    expectCreateBookingPayload({
      item_id: 'item-1',
      booking_date: '2026-01-20',
      for_user_id: 'colleague-1'
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

describe('fetchBookingExport', () => {
  beforeEach(() => {
    global.fetch = mockFetch;
    mockFetch.mockReset();
    mockFetch.mockResolvedValue({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ data: [], meta: { total: 0, limit: 20, offset: 0 }, links: {} })
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  const requestedUrl = () => new URL(String(mockFetch.mock.calls[0][0]), 'http://localhost');

  it('requests the range with paging and the area', async () => {
    await fetchBookingExport({ from: '2026-10-04', to: '2026-10-08', area: 'office', limit: 20, offset: 40 });

    const url = requestedUrl();
    expect(url.pathname).toBe('/api/v1/bookings/export');
    expect(Object.fromEntries(url.searchParams)).toEqual({
      'filter[from]': '2026-10-04',
      'filter[to]': '2026-10-08',
      'filter[area]': 'office',
      'page[limit]': '20',
      'page[offset]': '40'
    });
  });

  it('omits the area parameter for all areas', async () => {
    await fetchBookingExport({ from: '2026-10-04', to: '2026-10-08', area: '', limit: 20, offset: 0 });

    expect(requestedUrl().searchParams.has('filter[area]')).toBe(false);
  });

  it('rejects with the backend detail on a 400', async () => {
    mockFetch.mockResolvedValue({
      ok: false,
      status: 400,
      json: () => Promise.resolve({ errors: [{ detail: '"filter[from]" must not be after "filter[to]"' }] })
    });

    await expect(
      fetchBookingExport({ from: '2026-10-09', to: '2026-10-08', area: '', limit: 20, offset: 0 })
    ).rejects.toMatchObject({ status: 400, detail: '"filter[from]" must not be after "filter[to]"' });
  });
});
