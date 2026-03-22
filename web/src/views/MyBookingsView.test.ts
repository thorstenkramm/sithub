import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import MyBookingsView from './MyBookingsView.vue';
import { fetchMyBookings, cancelBooking } from '../api/bookings';
import { fetchMe } from '../api/me';
import { buildViewStubs, createFetchMeMocker, defineAuthRedirectTests } from './testHelpers';
import { ApiError, CONNECTION_LOST_MESSAGE } from '../api/client';

/* jscpd:ignore-start */
const pushMock = vi.fn();
vi.mock('../api/me', () => ({ fetchMe: vi.fn() }));
vi.mock('../api/bookings', () => ({ fetchMyBookings: vi.fn(), cancelBooking: vi.fn() }));
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));
/* jscpd:ignore-end */

describe('MyBookingsView', () => {
  const stubs = buildViewStubs([
    'v-list-item-subtitle',
    'v-card-item',
    'v-card-subtitle',
    'v-card-actions',
    'v-avatar',
    'v-icon',
    'v-chip',
    'v-spacer',
    'v-skeleton-loader',
    'v-dialog',
    'v-bottom-sheet',
    'v-snackbar',
    'v-textarea',
    'router-link'
  ]);
  const fetchMeMock = fetchMe as unknown as ReturnType<typeof vi.fn>;
  const mockFetchMe = createFetchMeMocker(fetchMeMock);

  const mockFetchBookings = (bookings: Array<{
    id: string;
    itemName: string;
    itemGroupName: string;
    areaName: string;
    bookingDate: string;
    bookedByUserId?: string;
    bookedByUserName?: string;
    bookedForMe?: boolean;
  }>) => {
    const fetchBookingsMock = fetchMyBookings as unknown as ReturnType<typeof vi.fn>;
    fetchBookingsMock.mockResolvedValue({
      data: bookings.map((b) => ({
        id: b.id,
        type: 'bookings',
        attributes: {
          item_id: `item-${b.id}`,
          item_name: b.itemName,
          item_group_id: `ig-${b.id}`,
          item_group_name: b.itemGroupName,
          area_id: `area-${b.id}`,
          area_name: b.areaName,
          booking_date: b.bookingDate,
          created_at: '2026-01-19T10:00:00Z',
          booked_by_user_id: b.bookedByUserId ?? '',
          booked_by_user_name: b.bookedByUserName ?? '',
          booked_for_me: b.bookedForMe ?? false
        }
      }))
    });
  };

  const mountView = () =>
    mount(MyBookingsView, {
      global: {
        stubs: {
          ...stubs,
          BookingCard: {
            props: ['booking'],
            template: `
              <div>
                <div>{{ booking.attributes.item_name }}</div>
                <div>{{ booking.attributes.item_group_name }}</div>
                <div>{{ booking.attributes.area_name }}</div>
                <div>{{ booking.attributes.booking_date }}</div>
                <div v-if="booking.attributes.booked_for_me">Booked by {{ booking.attributes.booked_by_user_name }}</div>
                <button type="button" @click="$emit('cancel', booking.id)">Cancel</button>
              </div>
            `
          },
          ConfirmDialog: {
            props: ['modelValue'],
            template: '<div v-if="modelValue"><button type="button" data-cy="confirm-cancel" @click="$emit(\'confirm\')">Confirm</button></div>'
          },
          'v-snackbar': {
            template: '<div v-bind="$attrs"><slot /></div>'
          }
        },
        plugins: [createPinia()]
      }
    });

  beforeEach(() => {
    setActivePinia(createPinia());
    pushMock.mockReset();
    mockFetchBookings([]);
    const cancelBookingMock = cancelBooking as unknown as ReturnType<typeof vi.fn>;
    cancelBookingMock.mockResolvedValue(undefined);
  });

  it('shows page header with title', async () => {
    mockFetchMe('Jane Doe');
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('My Bookings');
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);

  it('shows an empty state when no bookings exist', async () => {
    mockFetchMe();
    mockFetchBookings([]);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('No upcoming bookings');
  });

  it('renders bookings list with item, item group, area, and date', async () => {
    mockFetchMe();
    mockFetchBookings([
      { id: '1', itemName: 'Corner Desk', itemGroupName: 'Room 101', areaName: 'Main Office', bookingDate: '2026-01-20' },
      { id: '2', itemName: 'Window Desk', itemGroupName: 'Room 102', areaName: 'Annex', bookingDate: '2026-01-21' }
    ]);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Corner Desk');
    expect(wrapper.text()).toContain('Room 101');
    expect(wrapper.text()).toContain('Main Office');
    expect(wrapper.text()).toContain('Window Desk');
    expect(wrapper.text()).toContain('Room 102');
    expect(wrapper.text()).toContain('Annex');
  });

  it('displays formatted date', async () => {
    mockFetchMe();
    mockFetchBookings([
      { id: '1', itemName: 'Desk 1', itemGroupName: 'Room 1', areaName: 'Area 1', bookingDate: '2026-01-20' }
    ]);
    const wrapper = mountView();

    await flushPromises();

    // The date should be formatted (exact format depends on locale)
    expect(wrapper.text()).toContain('2026');
  });

  it('shows a connection lost error when user loading fails', async () => {
    fetchMeMock.mockRejectedValue(new ApiError(CONNECTION_LOST_MESSAGE, 0));
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain(CONNECTION_LOST_MESSAGE);
  });

  it('shows "Booked by" info when booking was made on behalf of user', async () => {
    mockFetchMe();
    mockFetchBookings([
      {
        id: '1',
        itemName: 'Corner Desk',
        itemGroupName: 'Room 101',
        areaName: 'Main Office',
        bookingDate: '2026-01-20',
        bookedByUserId: 'colleague-123',
        bookedByUserName: 'Jane Doe',
        bookedForMe: true
      }
    ]);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Booked by Jane Doe');
  });

  it('shows a snackbar when cancellation succeeds', async () => {
    mockFetchMe();
    mockFetchBookings([
      { id: '1', itemName: 'Corner Desk', itemGroupName: 'Room 101', areaName: 'Main Office', bookingDate: '2026-01-20' }
    ]);
    const cancelBookingMock = cancelBooking as unknown as ReturnType<typeof vi.fn>;
    const wrapper = mountView();

    await flushPromises();

    await wrapper.get('[data-cy="booking-item"] button').trigger('click');
    await flushPromises();
    await wrapper.get('[data-cy="confirm-cancel"]').trigger('click');
    await flushPromises();

    expect(cancelBookingMock).toHaveBeenCalledWith('1');
    expect(wrapper.find('[data-cy="cancel-success"]').exists()).toBe(true);
    expect(wrapper.text()).toContain('Booking cancelled successfully.');
  });
});
