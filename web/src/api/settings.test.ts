import { fetchSettings } from './settings';

describe('fetchSettings', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('fetches and returns settings', async () => {
    const mockResponse = {
      data: {
        id: 'settings',
        type: 'settings',
        attributes: { weeks_in_advanced: 5 }
      }
    };

    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockResponse)
    } as Response);

    const result = await fetchSettings();
    expect(result.data.attributes.weeks_in_advanced).toBe(5);
    expect(fetch).toHaveBeenCalledWith(
      '/api/v1/settings',
      expect.objectContaining({ headers: expect.any(Headers) })
    );
  });
});
