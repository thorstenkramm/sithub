import { fetchMe, changePassword } from './me';
import { apiRequest } from './client';

vi.mock('./client', () => ({
  apiRequest: vi.fn()
}));

const apiRequestMock = apiRequest as unknown as ReturnType<typeof vi.fn>;

describe('fetchMe', () => {
  it('calls apiRequest with /api/v1/me', async () => {
    apiRequestMock.mockResolvedValue({ data: { attributes: { display_name: 'Test User' } } });

    await fetchMe();

    expect(apiRequestMock).toHaveBeenCalledWith('/api/v1/me');
  });

  it('returns user attributes from the response', async () => {
    const mockResponse = {
      data: {
        id: '1',
        type: 'users',
        attributes: {
          display_name: 'Alice',
          email: 'alice@example.com',
          is_admin: true,
          auth_source: 'internal'
        }
      }
    };
    apiRequestMock.mockResolvedValue(mockResponse);

    const result = await fetchMe();

    expect(result.data.attributes.display_name).toBe('Alice');
    expect(result.data.attributes.is_admin).toBe(true);
  });
});

describe('changePassword', () => {
  it('sends PATCH /api/v1/me with current and new password', async () => {
    apiRequestMock.mockResolvedValue({ data: { attributes: { display_name: 'Test User' } } });

    await changePassword('oldpass', 'newpass12345678');

    expect(apiRequestMock).toHaveBeenCalledWith('/api/v1/me', {
      method: 'PATCH',
      body: JSON.stringify({
        data: {
          attributes: {
            current_password: 'oldpass',
            new_password: 'newpass12345678'
          }
        }
      })
    });
  });

  it('propagates API errors on failure', async () => {
    apiRequestMock.mockRejectedValue(new Error('Bad Request'));

    await expect(changePassword('old', 'short')).rejects.toThrow('Bad Request');
  });
});
