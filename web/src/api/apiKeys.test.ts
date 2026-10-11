import { createApiKey, deleteApiKey, fetchApiKeys } from './apiKeys';

describe('API key requests', () => {
  const fetchMock = vi.fn();
  beforeEach(() => {
    vi.stubGlobal('fetch', fetchMock);
    fetchMock.mockReset();
  });
  afterEach(() => vi.unstubAllGlobals());

  it('loads the key collection', async () => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ data: [] })));
    await expect(fetchApiKeys()).resolves.toEqual({ data: [] });
    expect(fetchMock).toHaveBeenCalledWith('/api/v1/api-keys', expect.any(Object));
  });

  it('posts a JSON:API create request and returns the one-time key', async () => {
    const response = { data: { id: 'id', attributes: { key: 'secret' } } };
    fetchMock.mockResolvedValue(new Response(JSON.stringify(response), { status: 201 }));
    const attributes = { name: 'Export', scope: 'read-bookings' as const, expires_in_months: 12 as const };
    await expect(createApiKey(attributes)).resolves.toEqual(response);
    expect(fetchMock).toHaveBeenCalledWith('/api/v1/api-keys', expect.objectContaining({
      method: 'POST', body: JSON.stringify({ data: { type: 'api-keys', attributes } })
    }));
  });

  it('revokes by encoded ID without parsing an empty 204 body', async () => {
    const response = new Response(null, { status: 204 });
    const json = vi.spyOn(response, 'json');
    fetchMock.mockResolvedValue(response);
    await expect(deleteApiKey('key/id')).resolves.toBeUndefined();
    expect(json).not.toHaveBeenCalled();
    expect(fetchMock).toHaveBeenCalledWith('/api/v1/api-keys/key%2Fid', expect.objectContaining({ method: 'DELETE' }));
  });
});
