import { createRouter, createWebHistory } from 'vue-router';

const mockBeforeEach = vi.fn();

vi.mock('vue-router', () => ({
  createRouter: vi.fn(() => ({ name: 'router', beforeEach: mockBeforeEach })),
  createWebHistory: vi.fn(() => ({ name: 'history' }))
}));

vi.mock('../stores/useAuthStore', () => ({
  useAuthStore: vi.fn(() => ({ isAuthenticated: false }))
}));

vi.mock('../api/me', () => ({
  fetchMe: vi.fn()
}));

describe('router', () => {
  it('creates router with history and routes', async () => {
    const module = await import('./index');

    expect(createWebHistory).toHaveBeenCalled();
    expect(createRouter).toHaveBeenCalledWith(
      expect.objectContaining({
        history: { name: 'history' },
        routes: expect.arrayContaining([
          expect.objectContaining({
            path: '/login',
            name: 'login',
            meta: { public: true }
          }),
          expect.objectContaining({
            path: '/',
            name: 'areas'
          }),
          expect.objectContaining({
            path: '/areas/:areaId/rooms',
            name: 'rooms'
          }),
          expect.objectContaining({
            path: '/areas/:areaId/presence',
            name: 'area-presence'
          }),
          expect.objectContaining({
            path: '/rooms/:roomId/desks',
            name: 'desks'
          }),
          expect.objectContaining({
            path: '/rooms/:roomId/bookings',
            name: 'room-bookings'
          }),
          expect.objectContaining({
            path: '/my-bookings',
            name: 'my-bookings'
          }),
          expect.objectContaining({
            path: '/bookings/history',
            name: 'booking-history'
          }),
          expect.objectContaining({
            path: '/access-denied',
            name: 'access-denied',
            meta: { public: true }
          })
        ])
      })
    );

    expect(mockBeforeEach).toHaveBeenCalledWith(expect.any(Function));
    expect(module.default).toEqual({ name: 'router', beforeEach: mockBeforeEach });
  });
});
