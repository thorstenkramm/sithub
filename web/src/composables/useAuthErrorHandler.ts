import { useRouter } from 'vue-router';
import { ApiError } from '../api/client';

/**
 * Composable for handling 401/403 API errors with navigation.
 *
 * - `handleAuthError` is meant for page loads: 401 redirects to the login, 403 to Access Denied.
 * - `handleUnauthorized` is meant for user actions (booking, cancelling, ...): only an expired
 *   session (401) redirects; a 403 is left to the caller so the backend message can be shown.
 *
 * Both return true if the error was handled.
 */
export function useAuthErrorHandler() {
  const router = useRouter();

  const handleUnauthorized = (err: unknown): boolean => {
    if (err instanceof ApiError && err.status === 401) {
      window.location.href = '/login';
      return true;
    }
    return false;
  };

  const handleAuthError = async (err: unknown): Promise<boolean> => {
    if (handleUnauthorized(err)) {
      return true;
    }
    if (err instanceof ApiError && err.status === 403) {
      await router.push('/access-denied');
      return true;
    }
    return false;
  };

  return { handleAuthError, handleUnauthorized };
}
