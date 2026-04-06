import { getAvatarUrl, uploadAvatar, deleteAvatar } from './avatars';

describe('getAvatarUrl', () => {
  it('returns the avatar URL for a given user ID', () => {
    expect(getAvatarUrl('user-123')).toBe('/api/v1/avatars/user-123');
  });

  it('encodes special characters in the user ID', () => {
    expect(getAvatarUrl('user with spaces')).toBe(
      '/api/v1/avatars/user%20with%20spaces',
    );
  });

  it('encodes slash characters in the user ID', () => {
    expect(getAvatarUrl('user/name')).toBe('/api/v1/avatars/user%2Fname');
  });
});

describe('uploadAvatar', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('sends POST /api/v1/me/avatar with FormData', async () => {
    const mockFetch = vi.fn().mockResolvedValue({ ok: true });
    vi.stubGlobal('fetch', mockFetch);

    const file = new File(['pixels'], 'photo.png', { type: 'image/png' });
    await uploadAvatar(file);

    expect(mockFetch).toHaveBeenCalledTimes(1);
    const [url, options] = mockFetch.mock.calls[0];
    expect(url).toBe('/api/v1/me/avatar');
    expect(options.method).toBe('POST');
    expect(options.body).toBeInstanceOf(FormData);
    expect((options.body as FormData).get('avatar')).toBe(file);
  });

  it('throws on non-ok response', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({ ok: false, status: 413 }),
    );

    const file = new File(['pixels'], 'photo.png', { type: 'image/png' });
    await expect(uploadAvatar(file)).rejects.toThrow('Upload failed: 413');
  });
});

describe('deleteAvatar', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('sends DELETE /api/v1/me/avatar', async () => {
    const mockFetch = vi.fn().mockResolvedValue({ ok: true });
    vi.stubGlobal('fetch', mockFetch);

    await deleteAvatar();

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/me/avatar', {
      method: 'DELETE',
    });
  });

  it('throws on non-ok response', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({ ok: false, status: 500 }),
    );

    await expect(deleteAvatar()).rejects.toThrow('Delete failed: 500');
  });
});
