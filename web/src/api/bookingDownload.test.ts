import { downloadBookingsCsv } from './bookings';
import { apiDownload } from './client';

vi.mock('./client', () => ({ apiDownload: vi.fn(), apiRequest: vi.fn(), throwApiError: vi.fn() }));

it.each(['office', ''])('downloads all rows for area %s and uses a fallback filename', async (area) => {
  const blob = new Blob(['csv']);
  vi.mocked(apiDownload).mockResolvedValue({ blob, filename: null });
  const result = await downloadBookingsCsv({ from: '2026-10-01', to: '2026-10-08', area });
  expect(result).toEqual({ blob, filename: 'sithub-bookings-2026-10-01_2026-10-08.csv' });
  const url = new URL(vi.mocked(apiDownload).mock.lastCall![0] as string, 'https://example.com');
  const filters = ['filter[from]', 'filter[to]'];
  expect([...url.searchParams.keys()]).toEqual(area ? [...filters, 'filter[area]'] : filters);
  expect(url.searchParams.get('filter[area]')).toBe(area || null);
});
