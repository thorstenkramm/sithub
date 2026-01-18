import { mount, flushPromises } from '@vue/test-utils';
import AreasView from './AreasView.vue';
import { ApiError } from '../api/client';
import { fetchAreas } from '../api/areas';
import { fetchMe } from '../api/me';

const pushMock = vi.fn();

vi.mock('../api/me', () => ({
  fetchMe: vi.fn()
}));

vi.mock('../api/areas', () => ({
  fetchAreas: vi.fn()
}));

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: pushMock
  })
}));

describe('AreasView', () => {
  const slotStub = {
    template: '<div><slot /></div>'
  };

  const stubs = {
    'v-container': slotStub,
    'v-row': slotStub,
    'v-col': slotStub,
    'v-card': slotStub,
    'v-card-title': slotStub,
    'v-card-text': slotStub,
    'v-btn': slotStub,
    'v-list': slotStub,
    'v-list-item': slotStub,
    'v-list-item-title': slotStub,
    'v-progress-linear': slotStub,
    'v-alert': slotStub
  };

  const mockFetchMe = (isAdmin: boolean) => {
    const fetchMeMock = fetchMe as unknown as ReturnType<typeof vi.fn>;
    fetchMeMock.mockResolvedValue({
      data: {
        attributes: {
          display_name: 'Ada Lovelace',
          is_admin: isAdmin
        }
      }
    });
  };

  const mockFetchAreas = (count: number) => {
    const fetchAreasMock = fetchAreas as unknown as ReturnType<typeof vi.fn>;
    fetchAreasMock.mockResolvedValue({
      data: Array.from({ length: count }, (_, index) => ({
        id: `area-${index + 1}`,
        type: 'areas',
        attributes: {
          name: `Area ${index + 1}`,
          sort_order: index,
          created_at: '2026-01-18T00:00:00Z',
          updated_at: '2026-01-18T00:00:00Z'
        }
      }))
    });
  };

  const mountView = () =>
    mount(AreasView, {
      global: {
        stubs
      }
    });

  beforeEach(() => {
    pushMock.mockReset();
    mockFetchAreas(0);
  });

  it('shows the signed-in user name', async () => {
    mockFetchMe(false);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Signed in as Ada Lovelace');
  });

  it('shows admin controls for admins only', async () => {
    mockFetchMe(true);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Cancel booking (admin)');
  });

  it('hides admin controls for non-admin users', async () => {
    mockFetchMe(false);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).not.toContain('Cancel booking (admin)');
  });

  it('redirects to login on 401', async () => {
    const fetchMeMock = fetchMe as unknown as ReturnType<typeof vi.fn>;
    fetchMeMock.mockRejectedValue(new ApiError('Unauthorized', 401));

    const originalLocation = window.location;
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: { href: 'http://localhost/' }
    });

    mount(AreasView, {
      global: {
        stubs: {
          'v-container': slotStub,
          'v-row': slotStub,
          'v-col': slotStub,
          'v-card': slotStub,
          'v-card-title': slotStub,
          'v-card-text': slotStub,
          'v-btn': slotStub,
          'v-list': slotStub,
          'v-list-item': slotStub,
          'v-list-item-title': slotStub,
          'v-progress-linear': slotStub,
          'v-alert': slotStub
        }
      }
    });

    await flushPromises();

    expect(window.location.href).toBe('/oauth/login');
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: originalLocation
    });
  });

  it('redirects to access denied on 403', async () => {
    const fetchMeMock = fetchMe as unknown as ReturnType<typeof vi.fn>;
    fetchMeMock.mockRejectedValue(new ApiError('Forbidden', 403));

    mountView();

    await flushPromises();

    expect(pushMock).toHaveBeenCalledWith('/access-denied');
  });

  it('shows an empty state when no areas exist', async () => {
    mockFetchMe(false);
    mockFetchAreas(0);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('No areas available.');
  });

  it('renders the areas list when data exists', async () => {
    mockFetchMe(false);
    mockFetchAreas(2);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Area 1');
    expect(wrapper.text()).toContain('Area 2');
  });
});
