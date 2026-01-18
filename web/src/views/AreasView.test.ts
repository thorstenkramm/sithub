import { mount, flushPromises } from '@vue/test-utils';
import AreasView from './AreasView.vue';
import { ApiError } from '../api/client';
import { fetchMe } from '../api/me';

vi.mock('../api/me', () => ({
  fetchMe: vi.fn()
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
    'v-btn': slotStub
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

  const mountView = () =>
    mount(AreasView, {
      global: {
        stubs
      }
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
          'v-btn': slotStub
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
});
