import { useRouter } from 'vue-router';
import { ApiError } from '../api/client';

/**
 * Composable for handling 401/403 API errors with navigation.
 * Returns a handler that redirects on auth errors and returns true if handled.
 */
export function useAuthErrorHandler() {
  const router = useRouter();

  const handleAuthError = async (err: unknown): Promise<boolean> => {
    if (err instanceof ApiError && err.status === 401) {
      window.location.href = '/login';
      return true;
    }
    if (err instanceof ApiError && err.status === 403) {
      await router.push('/access-denied');
      return true;
    }
    return false;
  };

  return { handleAuthError };
}
