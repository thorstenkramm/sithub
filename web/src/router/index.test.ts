import { createRouter, createWebHistory } from 'vue-router';

vi.mock('vue-router', () => ({
  createRouter: vi.fn(() => ({ name: 'router' })),
  createWebHistory: vi.fn(() => ({ name: 'history' }))
}));

describe('router', () => {
  it('creates router with history and routes', async () => {
    const module = await import('./index');

    expect(createWebHistory).toHaveBeenCalled();
    expect(createRouter).toHaveBeenCalledWith(
      expect.objectContaining({
        history: { name: 'history' },
        routes: [
          expect.objectContaining({
            path: '/',
            name: 'areas'
          })
        ]
      })
    );

    expect(module.default).toEqual({ name: 'router' });
  });
});
