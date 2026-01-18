import { fetchMe } from './me';
import { apiRequest } from './client';

type ApiRequestType = typeof apiRequest;

vi.mock('./client', () => ({
  apiRequest: vi.fn()
}));

describe('fetchMe', () => {
  it('calls apiRequest with /api/v1/me', async () => {
    const apiRequestMock = apiRequest as unknown as ReturnType<typeof vi.fn>;
    apiRequestMock.mockResolvedValue({ data: { attributes: { display_name: 'Test User' } } });

    await fetchMe();

    expect(apiRequestMock).toHaveBeenCalledWith('/api/v1/me');
  });
});
