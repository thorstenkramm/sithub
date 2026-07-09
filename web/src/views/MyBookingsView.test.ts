import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { h } from 'vue';
import MyBookingsView from './MyBookingsView.vue';
import { fetchMyBookings, cancelBooking } from '../api/bookings';
import { fetchMe } from '../api/me';
import { buildViewStubs, createFetchMeMocker, createTestI18n, defineAuthRedirectTests } from './testHelpers';
import { ApiError, CONNECTION_LOST_MESSAGE } from '../api/client';

/* jscpd:ignore-start */
const pushMock = vi.fn();
vi.mock('../api/me', () => ({ fetchMe: vi.fn() }));
vi.mock('../api/bookings', () => ({ fetchMyBookings: vi.fn(), cancelBooking: vi.fn() }));
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));
/* jscpd:ignore-end */

// Drives the viewport default: matches=true simulates a narrow (mobile) viewport.
const setMatchMedia = (matches: boolean) => {
  window.matchMedia = vi.fn().mockImplementation((query: string) => ({
    matches,
    media: query,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    addListener: vi.fn(),
    removeListener: vi.fn(),
    dispatchEvent: vi.fn()
  })) as unknown as typeof window.matchMedia;
};

// Minimal v-data-table stub that renders each item's cells via the named body slots
// (item.status / item.onBehalf / item.actions) so we can assert rendered content.
const dataTableStub = {
  props: ['headers', 'items'],
  setup(props: { items: Array<Record<string, unknown>> }, { slots }: { slots: Record<string, unknown> }) {
    return () =>
      h(
        'table',
        { 'data-cy': 'bookings-table' },
        props.items.map((item) =>
          h('tr', { key: item.id as string }, [
            h('td', {}, [item.date as string]),
            h('td', {}, [item.itemName as string]),
            h('td', {}, [item.area as string]),
            h('td', {}, slots['item.status'] ? (slots['item.status'] as (a: unknown) => unknown)({ item }) : []),
            h('td', {}, slots['item.onBehalf'] ? (slots['item.onBehalf'] as (a: unknown) => unknown)({ item }) : []),
            h('td', {}, slots['item.actions'] ? (slots['item.actions'] as (a: unknown) => unknown)({ item }) : [])
          ])
        )
      );
  }
};

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
    'v-switch',
    'v-tooltip',
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
    forUserName?: string;
    guestName?: string;
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
          booked_for_me: b.bookedForMe ?? false,
          for_user_name: b.forUserName,
          guest_name: b.guestName
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
          StatusChip: {
            props: ['status'],
            template: '<span :data-cy-status="status">{{ status }}</span>'
          },
          'v-data-table': dataTableStub,
          ConfirmDialog: {
            props: ['modelValue'],
            template: '<div v-if="modelValue"><button type="button" data-cy="confirm-cancel" @click="$emit(\'confirm\')">Confirm</button></div>'
          },
          'v-snackbar': {
            template: '<div v-bind="$attrs"><slot /></div>'
          }
        },
        plugins: [createPinia(), createTestI18n()]
      }
    });

  const singleBooking = () => [
    { id: '1', itemName: 'Corner Desk', itemGroupName: 'Room 101', areaName: 'Main Office', bookingDate: '2026-01-20' }
  ];

  // Mounts the view with a single-booking fixture and waits for async setup to settle.
  const mountReady = async () => {
    mockFetchMe();
    mockFetchBookings(singleBooking());
    const wrapper = mountView();
    await flushPromises();
    return wrapper;
  };

  // Triggers a cancel from the given trigger selector, confirms the dialog, and asserts the
  // shared confirm + success-snackbar path fires for booking '1'.
  const confirmCancelFrom = async (
    wrapper: ReturnType<typeof mountView>,
    triggerSelector: string
  ) => {
    const cancelBookingMock = cancelBooking as unknown as ReturnType<typeof vi.fn>;
    await wrapper.get(triggerSelector).trigger('click');
    await flushPromises();
    await wrapper.get('[data-cy="confirm-cancel"]').trigger('click');
    await flushPromises();

    expect(cancelBookingMock).toHaveBeenCalledWith('1');
    expect(wrapper.find('[data-cy="cancel-success"]').exists()).toBe(true);
    expect(wrapper.text()).toContain('Booking cancelled successfully.');
  };

  beforeEach(() => {
    setActivePinia(createPinia());
    localStorage.clear();
    setMatchMedia(false); // desktop default unless overridden
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

  it('defaults to the table view on desktop with empty storage', async () => {
    const wrapper = await mountReady();

    expect(wrapper.find('[data-cy="bookings-table"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="bookings-list"]').exists()).toBe(false);
    expect(wrapper.text()).toContain('Corner Desk');
    expect(wrapper.text()).toContain('Room 101');
    expect(wrapper.text()).toContain('Main Office');
  });

  const expectTileView = (wrapper: ReturnType<typeof mountView>) => {
    expect(wrapper.find('[data-cy="bookings-list"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="bookings-table"]').exists()).toBe(false);
  };

  it('defaults to the tile view on a narrow viewport', async () => {
    setMatchMedia(true);
    expectTileView(await mountReady());
  });

  it('restores a stored tile preference on desktop, overriding the viewport default', async () => {
    localStorage.setItem('sithub_my_bookings_view', 'cards');
    expectTileView(await mountReady());
  });

  it('renders bookings tiles with item, item group, area, and date', async () => {
    localStorage.setItem('sithub_my_bookings_view', 'cards');
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

  it('shows "Booked by" info in the tile view when booking was made on behalf of user', async () => {
    localStorage.setItem('sithub_my_bookings_view', 'cards');
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

  it('renders a StatusChip and on-behalf name in the table view', async () => {
    mockFetchMe();
    mockFetchBookings([
      {
        id: '1',
        itemName: 'Corner Desk',
        itemGroupName: 'Room 101',
        areaName: 'Main Office',
        bookingDate: '2026-01-20',
        bookedByUserId: 'colleague-123',
        forUserName: 'John Smith'
      }
    ]);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.find('[data-cy-status="on-behalf"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="table-on-behalf-of"]').text()).toContain('On behalf of John Smith');
  });

  it('cancels a booking from a table row using the confirm + snackbar path', async () => {
    const wrapper = await mountReady();
    await confirmCancelFrom(wrapper, '[data-cy="bookings-table"] [data-cy="cancel-btn"]');
  });

  it('shows a snackbar when cancellation succeeds from a tile', async () => {
    localStorage.setItem('sithub_my_bookings_view', 'cards');
    const wrapper = await mountReady();
    await confirmCancelFrom(wrapper, '[data-cy="booking-item"] button');
  });
});
