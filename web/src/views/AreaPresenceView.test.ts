import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import AreaPresenceView from './AreaPresenceView.vue';
import { fetchAreaPresence } from '../api/areaPresence';
import { fetchAreas } from '../api/areas';
import {
  buildViewStubs,
  expectLoginRedirect,
  expectAccessDeniedRedirect
} from './testHelpers';
import { ApiError } from '../api/client';

const pushMock = vi.fn();
vi.mock('../api/areaPresence');
vi.mock('../api/areas');
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { areaId: 'area-1' } }),
  useRouter: () => ({ push: pushMock })
}));

describe('AreaPresenceView', () => {
  const stubs = buildViewStubs([
    'v-list-item-subtitle',
    'v-card-item',
    'v-card-subtitle',
    'v-avatar',
    'v-icon',
    'v-chip',
    'v-menu',
    'v-date-picker',
    'v-text-field',
    'v-skeleton-loader',
    'router-link'
  ]);

  const fetchAreaPresenceMock = vi.mocked(fetchAreaPresence);
  const fetchAreasMock = vi.mocked(fetchAreas);

  const mountView = () =>
    mount(AreaPresenceView, {
      global: {
        stubs,
        plugins: [createPinia()]
      }
    });

  beforeEach(() => {
    setActivePinia(createPinia());
    pushMock.mockReset();
    fetchAreaPresenceMock.mockResolvedValue({ data: [] });
    fetchAreasMock.mockResolvedValue({
      data: [{ id: 'area-1', type: 'areas', attributes: { name: 'Test Area' } }]
    });
  });

  it('renders presence list', async () => {
    fetchAreaPresenceMock.mockResolvedValue({
      data: [
        {
          id: 'booking-1',
          type: 'presence',
          attributes: {
            user_id: 'user-1',
            user_name: 'Alice Smith',
            desk_id: 'desk-1',
            desk_name: 'Desk 1',
            room_id: 'room-1',
            room_name: 'Room One'
          }
        }
      ]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('Alice Smith');
    expect(wrapper.text()).toContain('Room One');
    expect(wrapper.text()).toContain('Desk 1');
  });

  it('shows empty state when no one is present', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('No one scheduled');
  });

  it('fetches presence on mount', async () => {
    const wrapper = mountView();
    await flushPromises();

    // Should fetch with today's date on mount
    expect(fetchAreaPresenceMock).toHaveBeenCalled();
    const lastCall = fetchAreaPresenceMock.mock.calls[fetchAreaPresenceMock.mock.calls.length - 1];
    expect(lastCall[0]).toBe('area-1');
    expect(lastCall[1]).toMatch(/^\d{4}-\d{2}-\d{2}$/);
  });

  it('shows error for area not found', async () => {
    fetchAreaPresenceMock.mockRejectedValue(new ApiError('Not found', 404));

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('Area not found.');
  });

  it('redirects to login on 401', async () => {
    fetchAreaPresenceMock.mockRejectedValue(new ApiError('Unauthorized', 401));
    await expectLoginRedirect(mountView);
  });

  it('redirects to access denied on 403', async () => {
    fetchAreaPresenceMock.mockRejectedValue(new ApiError('Forbidden', 403));
    await expectAccessDeniedRedirect(mountView, pushMock);
  });
});
