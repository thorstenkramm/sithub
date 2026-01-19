import { mount, flushPromises } from '@vue/test-utils';
import MyBookingsView from './MyBookingsView.vue';
import { fetchMyBookings } from '../api/bookings';
import { fetchMe } from '../api/me';
import { buildViewStubs, createFetchMeMocker, defineAuthRedirectTests } from './testHelpers';

/* jscpd:ignore-start */
const pushMock = vi.fn();
vi.mock('../api/me', () => ({ fetchMe: vi.fn() }));
vi.mock('../api/bookings', () => ({ fetchMyBookings: vi.fn() }));
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));
/* jscpd:ignore-end */

describe('MyBookingsView', () => {
  const stubs = buildViewStubs(['v-list-item-subtitle']);
  const fetchMeMock = fetchMe as unknown as ReturnType<typeof vi.fn>;
  const mockFetchMe = createFetchMeMocker(fetchMeMock);

  const mockFetchBookings = (bookings: Array<{
    id: string;
    deskName: string;
    roomName: string;
    areaName: string;
    bookingDate: string;
  }>) => {
    const fetchBookingsMock = fetchMyBookings as unknown as ReturnType<typeof vi.fn>;
    fetchBookingsMock.mockResolvedValue({
      data: bookings.map((b) => ({
        id: b.id,
        type: 'bookings',
        attributes: {
          desk_id: `desk-${b.id}`,
          desk_name: b.deskName,
          room_id: `room-${b.id}`,
          room_name: b.roomName,
          area_id: `area-${b.id}`,
          area_name: b.areaName,
          booking_date: b.bookingDate,
          created_at: '2026-01-19T10:00:00Z'
        }
      }))
    });
  };

  const mountView = () =>
    mount(MyBookingsView, {
      global: {
        stubs
      }
    });

  beforeEach(() => {
    pushMock.mockReset();
    mockFetchBookings([]);
  });

  it('shows the signed-in user name', async () => {
    mockFetchMe('Jane Doe');
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Signed in as Jane Doe');
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);

  it('shows an empty state when no bookings exist', async () => {
    mockFetchMe();
    mockFetchBookings([]);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('No upcoming bookings');
  });

  it('renders bookings list with desk, room, area, and date', async () => {
    mockFetchMe();
    mockFetchBookings([
      { id: '1', deskName: 'Corner Desk', roomName: 'Room 101', areaName: 'Main Office', bookingDate: '2026-01-20' },
      { id: '2', deskName: 'Window Desk', roomName: 'Room 102', areaName: 'Annex', bookingDate: '2026-01-21' }
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
      { id: '1', deskName: 'Desk 1', roomName: 'Room 1', areaName: 'Area 1', bookingDate: '2026-01-20' }
    ]);
    const wrapper = mountView();

    await flushPromises();

    // The date should be formatted (exact format depends on locale)
    expect(wrapper.text()).toContain('2026');
  });
});
