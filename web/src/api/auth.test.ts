import { loginLocal, logout } from './auth';
import { apiRequest } from './client';

vi.mock('./client', () => ({
  apiRequest: vi.fn()
}));

const apiRequestMock = apiRequest as unknown as ReturnType<typeof vi.fn>;

describe('loginLocal', () => {
  it('sends POST /api/v1/auth/login with email and password', async () => {
    apiRequestMock.mockResolvedValue({
      data: { id: '1', type: 'users', attributes: { display_name: 'Test' } }
    });

    await loginLocal('test@example.com', 'password123');

    expect(apiRequestMock).toHaveBeenCalledWith('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email: 'test@example.com', password: 'password123' })
    });
  });

  it('returns the parsed user response', async () => {
    const mockResponse = {
      data: { id: '42', type: 'users', attributes: { display_name: 'Alice', email: 'alice@example.com' } }
    };
    apiRequestMock.mockResolvedValue(mockResponse);

    const result = await loginLocal('alice@example.com', 'securepass');

    expect(result).toEqual(mockResponse);
  });

  it('propagates API errors', async () => {
    apiRequestMock.mockRejectedValue(new Error('Unauthorized'));

    await expect(loginLocal('bad@example.com', 'wrong')).rejects.toThrow('Unauthorized');
  });
});

describe('logout', () => {
  it('sends POST /api/v1/auth/logout', async () => {
    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response());

    await logout();

    expect(fetchSpy).toHaveBeenCalledWith('/api/v1/auth/logout', { method: 'POST' });
    fetchSpy.mockRestore();
  });

  it('does not throw on network failure', async () => {
    const fetchSpy = vi.spyOn(globalThis, 'fetch').mockRejectedValue(new TypeError('Network error'));

    await expect(logout()).resolves.toBeUndefined();

    fetchSpy.mockRestore();
  });
});
