import { ApiError, apiRequest } from './client';

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
});
