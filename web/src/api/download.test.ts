import { apiDownload, ApiError, CONNECTION_LOST_MESSAGE } from './client';

describe('apiDownload', () => {
  afterEach(() => vi.unstubAllGlobals());

  it.each(['attachment; filename="bookings.csv"', 'attachment; filename=bookings.csv', null])(
    'reads a blob and filename from %s', async (header) => {
      const blob = new Blob(['csv']);
      const fetchMock = vi.fn().mockResolvedValue({ ok: true, headers: new Headers(header ? { 'Content-Disposition': header } : {}), blob: async () => blob });
      vi.stubGlobal('fetch', fetchMock);
      expect(await apiDownload('/export')).toEqual({ blob, filename: header ? 'bookings.csv' : null });
      expect(fetchMock.mock.calls[0][1].headers.get('Accept')).toBe('text/csv');
    });

  it('preserves backend errors', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, status: 403, json: async () => ({ errors: [{ detail: 'forbidden' }] }) }));
    await expect(apiDownload('/export')).rejects.toMatchObject({ status: 403, detail: 'forbidden' });
  });

  it('reports connection failures', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('offline')));
    await expect(apiDownload('/export')).rejects.toBeInstanceOf(ApiError);
    await expect(apiDownload('/export')).rejects.toMatchObject({ status: 0 });
  });

  it('reports a connection failure when the response body is interrupted', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true, headers: new Headers(), blob: async () => { throw new TypeError('terminated'); }
    }));
    await expect(apiDownload('/export')).rejects.toMatchObject({ status: 0, message: CONNECTION_LOST_MESSAGE });
  });
});
