import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useAuthErrorHandler } from '../useAuthErrorHandler';
import { ApiError } from '../../api/client';

const pushMock = vi.fn();

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: pushMock })
}));

describe('useAuthErrorHandler', () => {
  const originalLocation = window.location;

  beforeEach(() => {
    pushMock.mockReset();
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: { href: '' }
    });
  });

  afterEach(() => {
    Object.defineProperty(window, 'location', { configurable: true, value: originalLocation });
  });

  describe('handleAuthError', () => {
    it('should redirect to login on 401', async () => {
      const { handleAuthError } = useAuthErrorHandler();

      await expect(handleAuthError(new ApiError('x', 401))).resolves.toBe(true);
      expect(window.location.href).toBe('/login');
    });

    it('should redirect to access denied on 403', async () => {
      const { handleAuthError } = useAuthErrorHandler();

      await expect(handleAuthError(new ApiError('x', 403))).resolves.toBe(true);
      expect(pushMock).toHaveBeenCalledWith('/access-denied');
    });

    it('should not handle other errors', async () => {
      const { handleAuthError } = useAuthErrorHandler();

      await expect(handleAuthError(new ApiError('x', 400))).resolves.toBe(false);
      await expect(handleAuthError(new Error('boom'))).resolves.toBe(false);
      expect(pushMock).not.toHaveBeenCalled();
    });
  });

  describe('handleUnauthorized', () => {
    it('should redirect to login on 401', () => {
      const { handleUnauthorized } = useAuthErrorHandler();

      expect(handleUnauthorized(new ApiError('x', 401))).toBe(true);
      expect(window.location.href).toBe('/login');
    });

    it('should leave 403 to the caller so the backend message can be shown', () => {
      const { handleUnauthorized } = useAuthErrorHandler();

      expect(handleUnauthorized(new ApiError('x', 403, 'This item is reserved.'))).toBe(false);
      expect(pushMock).not.toHaveBeenCalled();
      expect(window.location.href).toBe('');
    });
  });
});
