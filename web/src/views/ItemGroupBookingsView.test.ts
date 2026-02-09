import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import ItemGroupBookingsView from './ItemGroupBookingsView.vue';
import { fetchItemGroupBookings } from '../api/itemGroupBookings';
import { fetchAreas } from '../api/areas';
import { fetchItemGroups } from '../api/itemGroups';
import {
  buildViewStubs,
  expectLoginRedirect,
  expectAccessDeniedRedirect
} from './testHelpers';
import { ApiError } from '../api/client';

/* jscpd:ignore-start */

const pushMock = vi.fn();
vi.mock('../api/itemGroupBookings');
vi.mock('../api/areas');
vi.mock('../api/itemGroups');
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { itemGroupId: 'ig-1' } }),
  useRouter: () => ({ push: pushMock })
}));

describe('ItemGroupBookingsView', () => {
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

  const fetchItemGroupBookingsMock = vi.mocked(fetchItemGroupBookings);
  const fetchAreasMock = vi.mocked(fetchAreas);
  const fetchItemGroupsMock = vi.mocked(fetchItemGroups);

  const mountView = () =>
    mount(ItemGroupBookingsView, {
      global: {
        stubs,
        plugins: [createPinia()]
      }
    });

  beforeEach(() => {
    setActivePinia(createPinia());
    pushMock.mockReset();
    fetchItemGroupBookingsMock.mockResolvedValue({ data: [] });
    fetchAreasMock.mockResolvedValue({
      data: [{ id: 'area-1', type: 'areas', attributes: { name: 'Test Area' } }]
    });
    fetchItemGroupsMock.mockResolvedValue({
      data: [{ id: 'ig-1', type: 'item-groups', attributes: { name: 'Test Group' } }]
    });
  });

  it('renders bookings list', async () => {
    fetchItemGroupBookingsMock.mockResolvedValue({
      data: [
        {
          id: 'booking-1',
          type: 'bookings',
          attributes: {
            item_id: 'item-1',
            item_name: 'Item 1',
            user_id: 'user-1',
            user_name: 'Alice Smith',
            booking_date: '2026-01-20'
          }
        }
      ]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('Item 1');
    expect(wrapper.text()).toContain('Alice Smith');
  });

  it('shows empty state when no bookings exist', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('No bookings');
  });

  it('fetches bookings on mount', async () => {
    mountView();
    await flushPromises();

    // Should fetch with today's date on mount
    expect(fetchItemGroupBookingsMock).toHaveBeenCalled();
    const lastCall = fetchItemGroupBookingsMock.mock.calls[fetchItemGroupBookingsMock.mock.calls.length - 1];
    expect(lastCall[0]).toBe('ig-1');
    expect(lastCall[1]).toMatch(/^\d{4}-\d{2}-\d{2}$/);
  });

  it('shows error for item group not found', async () => {
    fetchItemGroupBookingsMock.mockRejectedValue(new ApiError('Not found', 404));

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('Item group not found.');
  });

  it('redirects to login on 401', async () => {
    fetchItemGroupBookingsMock.mockRejectedValue(new ApiError('Unauthorized', 401));
    await expectLoginRedirect(mountView);
  });

  it('redirects to access denied on 403', async () => {
    fetchItemGroupBookingsMock.mockRejectedValue(new ApiError('Forbidden', 403));
    await expectAccessDeniedRedirect(mountView, pushMock);
  });
});
/* jscpd:ignore-end */
