import { ref } from 'vue';
import { apiErrorReason } from '../api/client';

export function useApi() {
  const loading = ref(false);
  const error = ref<string | null>(null);

  async function run<T>(fn: () => Promise<T>): Promise<T> {
    loading.value = true;
    error.value = null;
    try {
      return await fn();
    } catch (err) {
      // Backend detail (untranslated) or HTTP status, see apiErrorReason.
      error.value = err instanceof Error ? apiErrorReason(err) : 'Unknown error';
      throw err;
    } finally {
      loading.value = false;
    }
  }

  return {
    loading,
    error,
    run
  };
}
