import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import AreaPresenceView from './AreaPresenceView.vue';
import { fetchAreaPresence } from '../api/areaPresence';
import { fetchAreas } from '../api/areas';
import {
  buildViewStubs,
  createTestI18n,
  expectLoginRedirect,
  expectAccessDeniedRedirect
} from './testHelpers';
import { ApiError, CONNECTION_LOST_MESSAGE } from '../api/client';

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
    'v-dialog',
    'v-bottom-sheet',
    'v-card-actions',
    'v-spacer',
    'router-link'
  ]);
  // Render the #prepend slot (avatar) in addition to the default slot.
  stubs['v-list-item'] = {
    template: '<div><slot name="prepend" /><slot /></div>'
  };

  const fetchAreaPresenceMock = vi.mocked(fetchAreaPresence);
  const fetchAreasMock = vi.mocked(fetchAreas);

  const mountView = () =>
    mount(AreaPresenceView, {
      global: {
        stubs,
        plugins: [createPinia(), createTestI18n()]
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
            item_id: 'item-1',
            item_name: 'Desk 1',
            item_group_id: 'ig-1',
            item_group_name: 'Room One'
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

  it('renders the avatar image and no initials for a user with a photo', async () => {
    fetchAreaPresenceMock.mockResolvedValue({
      data: [
        {
          id: 'booking-1',
          type: 'presence',
          attributes: {
            user_id: 'user-1',
            user_name: 'Alice Smith',
            item_id: 'item-1',
            item_name: 'Desk 1',
            item_group_id: 'ig-1',
            item_group_name: 'Room One'
          }
        }
      ]
    });

    const wrapper = mountView();
    await flushPromises();

    const img = wrapper.find('img.presence-avatar-img');
    expect(img.exists()).toBe(true);
    expect(img.attributes('src')).toContain('user-1');
    expect(wrapper.find('.presence-avatar-initials').exists()).toBe(false);
  });

  it('falls back to centered initials and hides the image after a photo load error', async () => {
    fetchAreaPresenceMock.mockResolvedValue({
      data: [
        {
          id: 'booking-1',
          type: 'presence',
          attributes: {
            user_id: 'user-1',
            user_name: 'Alice Smith',
            item_id: 'item-1',
            item_name: 'Desk 1',
            item_group_id: 'ig-1',
            item_group_name: 'Room One'
          }
        }
      ]
    });

    const wrapper = mountView();
    await flushPromises();

    await wrapper.find('img.presence-avatar-img').trigger('error');
    await flushPromises();

    expect(wrapper.find('img.presence-avatar-img').exists()).toBe(false);
    const initials = wrapper.find('.presence-avatar-initials');
    expect(initials.exists()).toBe(true);
    expect(initials.text()).toBe('AS');
  });

  it('shows empty state when no one is present', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('No one scheduled');
  });

  it('fetches presence on mount', async () => {
    mountView();
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

  it('shows a connection lost error when presence loading fails', async () => {
    fetchAreaPresenceMock.mockRejectedValue(new ApiError(CONNECTION_LOST_MESSAGE, 0));

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain(CONNECTION_LOST_MESSAGE);
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
