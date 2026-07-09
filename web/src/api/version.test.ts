import { fetchVersion } from './version';

describe('fetchVersion', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('fetches and returns the running version', async () => {
    const mockResponse = {
      data: {
        id: 'version',
        type: 'version',
        attributes: { version: '1.2.3' }
      }
    };

    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockResponse)
    } as Response);

    const result = await fetchVersion();
    expect(result.data.attributes.version).toBe('1.2.3');
    expect(fetch).toHaveBeenCalledWith(
      '/api/v1/version',
      expect.objectContaining({ headers: expect.any(Headers) })
    );
  });
});
