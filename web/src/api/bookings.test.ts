import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { createBooking } from './bookings';

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
