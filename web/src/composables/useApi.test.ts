import { useApi } from './useApi';
import { ApiError } from '../api/client';

const flush = () => new Promise((resolve) => setTimeout(resolve, 0));

describe('useApi', () => {
  it('toggles loading and captures errors', async () => {
    const { loading, error, run } = useApi();

    const failing = run(async () => {
      throw new Error('boom');
    });

    expect(loading.value).toBe(true);
    await expect(failing).rejects.toThrow('boom');
    await flush();

    expect(loading.value).toBe(false);
    expect(error.value).toBe('boom');
  });

  it('clears error on success', async () => {
    const { loading, error, run } = useApi();

    await run(async () => 'ok');

    expect(loading.value).toBe(false);
    expect(error.value).toBeNull();
  });

  it('stores the backend detail of an ApiError', async () => {
    const { error, run } = useApi();

    await expect(run(async () => {
      throw new ApiError('Request failed: 500', 500, 'database is locked');
    })).rejects.toBeInstanceOf(ApiError);

    expect(error.value).toBe('database is locked');
  });

  it('stores the HTTP status when an ApiError has no detail', async () => {
    const { error, run } = useApi();

    await expect(run(async () => {
      throw new ApiError('Request failed: 502', 502);
    })).rejects.toBeInstanceOf(ApiError);

    expect(error.value).toBe('HTTP 502');
  });
});
