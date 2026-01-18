import { useApi } from './useApi';

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
});
