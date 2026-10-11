import { createRouter, createWebHistory } from 'vue-router';
import { fetchMe } from '../api/me';

const mockBeforeEach = vi.fn();

vi.mock('vue-router', () => ({
  createRouter: vi.fn(() => ({ name: 'router', beforeEach: mockBeforeEach })),
  createWebHistory: vi.fn(() => ({ name: 'history' }))
}));

const authStoreMock = { isAuthenticated: false, isAdmin: false, isDataExporter: false, setUser: vi.fn() };

vi.mock('../stores/useAuthStore', () => ({
  useAuthStore: vi.fn(() => authStoreMock)
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
            path: '/areas/:areaId/item-groups',
            name: 'item-groups'
          }),
          expect.objectContaining({
            path: '/areas/:areaId/presence',
            name: 'area-presence'
          }),
          expect.objectContaining({
            path: '/item-groups/:itemGroupId/items',
            name: 'items'
          }),
          expect.objectContaining({
            path: '/item-groups/:itemGroupId/bookings',
            name: 'item-group-bookings'
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
            path: '/admin/users',
            name: 'admin-users',
            meta: { requiresAdmin: true }
          }),
          expect.objectContaining({
            path: '/admin/api-keys',
            name: 'api-keys',
            meta: { requiresAdmin: true }
          }),
          expect.objectContaining({
            path: '/data-export',
            name: 'data-export',
            meta: { requiresDataExporter: true }
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

describe('router guard for /admin/users', () => {
  type Guard = (to: { name: string; meta: Record<string, unknown> }) => Promise<unknown>;
  const adminUsersRoute = { name: 'admin-users', meta: { requiresAdmin: true } };
  const fetchMeMock = fetchMe as unknown as ReturnType<typeof vi.fn>;

  async function loadGuard(): Promise<Guard> {
    vi.resetModules();
    await import('./index');
    return mockBeforeEach.mock.lastCall?.[0] as Guard;
  }

  beforeEach(() => {
    authStoreMock.isAuthenticated = true;
    authStoreMock.isAdmin = false;
  });

  it('should redirect non-admins to access denied', async () => {
    const guard = await loadGuard();

    await expect(guard(adminUsersRoute)).resolves.toEqual({ name: 'access-denied' });
  });

  it('should let admins reach the users page', async () => {
    authStoreMock.isAdmin = true;
    const guard = await loadGuard();

    await expect(guard(adminUsersRoute)).resolves.toBe(true);
  });

  it('should redirect to login when the session cannot be loaded', async () => {
    authStoreMock.isAuthenticated = false;
    fetchMeMock.mockRejectedValueOnce(new Error('401'));
    const guard = await loadGuard();

    await expect(guard(adminUsersRoute)).resolves.toEqual({ name: 'login' });
  });
});

describe('router guard for /data-export', () => {
  type Guard = (to: { name: string; meta: Record<string, unknown> }) => Promise<unknown>;
  const dataExportRoute = { name: 'data-export', meta: { requiresDataExporter: true } };

  async function loadGuard(): Promise<Guard> {
    vi.resetModules();
    await import('./index');
    return mockBeforeEach.mock.lastCall?.[0] as Guard;
  }

  beforeEach(() => {
    authStoreMock.isAuthenticated = true;
    authStoreMock.isAdmin = false;
    authStoreMock.isDataExporter = false;
  });

  it.each([
    { name: 'a regular user', isAdmin: false },
    { name: 'an administrator without the flag', isAdmin: true }
  ])('should redirect $name to access denied', async ({ isAdmin }) => {
    authStoreMock.isAdmin = isAdmin;
    const guard = await loadGuard();

    await expect(guard(dataExportRoute)).resolves.toEqual({ name: 'access-denied' });
  });

  it('should let a data exporter reach the page', async () => {
    authStoreMock.isDataExporter = true;
    const guard = await loadGuard();

    await expect(guard(dataExportRoute)).resolves.toBe(true);
  });
});
