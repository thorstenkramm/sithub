import {
  ApiError,
  apiRequest,
  formatApiErrorMessage,
  isConnectionError,
  throwApiError,
  CONNECTION_LOST_MESSAGE
} from './client';

const JSON_API = 'application/vnd.api+json';

const setFetchMock = (response: {
  ok: boolean;
  status?: number;
  json?: () => Promise<unknown>;
}) => {
  const fetchMock = vi.fn().mockResolvedValue({
    ok: response.ok,
    status: response.status ?? 200,
    json: response.json ?? (async () => ({ ok: true }))
  });

  globalThis.fetch = fetchMock as unknown as typeof fetch;
  return fetchMock;
};

describe('apiRequest', () => {
  const originalFetch = globalThis.fetch;

  afterEach(() => {
    globalThis.fetch = originalFetch;
  });

  it('sets JSON:API headers by default', async () => {
    const fetchMock = setFetchMock({ ok: true });

    await apiRequest('/api/v1/ping');

    const [, init] = fetchMock.mock.calls[0];
    const headers = init?.headers as Headers;

    expect(headers.get('Accept')).toBe(JSON_API);
    expect(headers.get('Content-Type')).toBe(JSON_API);
  });

  it('preserves explicit content type', async () => {
    const fetchMock = setFetchMock({ ok: true });

    await apiRequest('/api/v1/ping', {
      headers: {
        'Content-Type': 'application/custom'
      }
    });

    const [, init] = fetchMock.mock.calls[0];
    const headers = init?.headers as Headers;

    expect(headers.get('Content-Type')).toBe('application/custom');
  });

  it('throws ApiError on non-OK responses', async () => {
    setFetchMock({
      ok: false,
      status: 500,
      json: async () => ({})
    });

    await expect(apiRequest('/api/v1/ping')).rejects.toBeInstanceOf(ApiError);
  });

  describe('error detail parsing', () => {
    const expectApiError = async (
      url: string,
      expectedStatus: number,
      expectedDetail: string | null
    ) => {
      await expect(apiRequest(url)).rejects.toSatisfy((err: unknown) => {
        expect(err).toBeInstanceOf(ApiError);
        const apiErr = err as ApiError;
        expect(apiErr.status).toBe(expectedStatus);
        expect(apiErr.detail).toBe(expectedDetail);
        return true;
      });
    };

    it('includes error detail from JSON:API error response', async () => {
      setFetchMock({
        ok: false,
        status: 409,
        json: async () => ({
          errors: [{ status: '409', title: 'Conflict', detail: 'Desk is already booked', code: 'conflict' }]
        })
      });

      await expectApiError('/api/v1/bookings', 409, 'Desk is already booked');
    });

    it('sets detail to null when error response has no detail', async () => {
      setFetchMock({ ok: false, status: 400, json: async () => ({}) });

      await expectApiError('/api/v1/bookings', 400, null);
    });

    it('handles JSON parse errors gracefully', async () => {
      setFetchMock({
        ok: false,
        status: 500,
        json: async () => {
          throw new Error('Invalid JSON');
        }
      });

      await expectApiError('/api/v1/ping', 500, null);
    });

    it('keeps valid details when other errors are malformed', async () => {
      setFetchMock({
        ok: false,
        status: 422,
        json: async () => ({ errors: [null, { detail: 'valid detail' }, { detail: 42 }, { detail: {} }, { detail: '' }] })
      });
      await expectApiError('/api/v1/bookings', 422, 'valid detail');
    });

    it('joins multiple error details with a newline', async () => {
      setFetchMock({
        ok: false,
        status: 422,
        json: async () => ({
          errors: [{ detail: 'first problem' }, { title: 'no detail' }, { detail: 'second problem' }]
        })
      });

      await expectApiError('/api/v1/bookings', 422, 'first problem\nsecond problem');
    });
  });

  describe('connection error handling', () => {
    it('throws ApiError with status 0 when fetch fails', async () => {
      globalThis.fetch = vi.fn().mockRejectedValue(new TypeError('Failed to fetch')) as unknown as typeof fetch;

      await expect(apiRequest('/api/v1/ping')).rejects.toSatisfy((err: unknown) => {
        expect(err).toBeInstanceOf(ApiError);
        const apiErr = err as ApiError;
        expect(apiErr.status).toBe(0);
        expect(apiErr.message).toBe(CONNECTION_LOST_MESSAGE);
        return true;
      });
    });

    it('isConnectionError returns true for status 0 ApiError', () => {
      const err = new ApiError(CONNECTION_LOST_MESSAGE, 0);
      expect(isConnectionError(err)).toBe(true);
    });

    it('isConnectionError returns false for other ApiErrors', () => {
      const err = new ApiError('Not Found', 404);
      expect(isConnectionError(err)).toBe(false);
    });

    it('isConnectionError returns false for non-ApiError', () => {
      expect(isConnectionError(new Error('something'))).toBe(false);
      expect(isConnectionError(null)).toBe(false);
    });
  });
});

describe('throwApiError', () => {
  it('throws ApiError with status and backend detail', async () => {
    const response = {
      status: 400,
      json: async () => ({ errors: [{ detail: 'booking_date is too far in the future' }] })
    } as unknown as Response;

    await expect(throwApiError(response)).rejects.toMatchObject({
      message: 'Request failed: 400',
      status: 400,
      detail: 'booking_date is too far in the future'
    });
  });

  it('throws ApiError with null detail when the body is empty', async () => {
    const response = {
      status: 502,
      json: async () => {
        throw new SyntaxError('Unexpected end of JSON input');
      }
    } as unknown as Response;

    await expect(throwApiError(response)).rejects.toSatisfy((err: unknown) => {
      expect(err).toBeInstanceOf(ApiError);
      expect((err as ApiError).status).toBe(502);
      expect((err as ApiError).detail).toBeNull();
      return true;
    });
  });
});

describe('formatApiErrorMessage', () => {
  const prefix = 'Unable to book item. API error:';

  it('appends the backend detail verbatim on a new line', () => {
    const err = new ApiError('Request failed: 400', 400, 'booking_date is too far in the future: 2026-08-24');

    expect(formatApiErrorMessage(prefix, err)).toBe(
      `${prefix}\nbooking_date is too far in the future: 2026-08-24`
    );
  });

  it('falls back to the HTTP status when the detail is missing', () => {
    expect(formatApiErrorMessage(prefix, new ApiError('Request failed: 502', 502))).toBe(`${prefix}\nHTTP 502`);
  });

  it('falls back to the HTTP status when the detail is empty', () => {
    expect(formatApiErrorMessage(prefix, new ApiError('Request failed: 500', 500, ''))).toBe(`${prefix}\nHTTP 500`);
  });

  it('shows the connection lost message for network errors', () => {
    const err = new ApiError(CONNECTION_LOST_MESSAGE, 0);

    expect(formatApiErrorMessage(prefix, err)).toBe(`${prefix}\n${CONNECTION_LOST_MESSAGE}`);
  });

  it('uses the message of a plain Error', () => {
    expect(formatApiErrorMessage(prefix, new Error('boom'))).toBe(`${prefix}\nboom`);
  });

  it('stringifies non-Error values', () => {
    expect(formatApiErrorMessage(prefix, 'weird')).toBe(`${prefix}\nweird`);
  });

  it('puts exactly one newline between prefix and reason', () => {
    const message = formatApiErrorMessage(prefix, new ApiError('x', 403, 'reserved'));

    expect(message.split('\n')).toEqual([prefix, 'reserved']);
  });
});
