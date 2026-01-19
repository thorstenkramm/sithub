import { mount, flushPromises } from '@vue/test-utils';
import RoomBookingsView from './RoomBookingsView.vue';
import { fetchRoomBookings } from '../api/roomBookings';
import { buildViewStubs, defineAuthRedirectTests, mockWindowLocation } from './testHelpers';
import { ApiError } from '../api/client';

const pushMock = vi.fn();
vi.mock('../api/roomBookings');
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { roomId: 'room-1' } }),
  useRouter: () => ({ push: pushMock })
}));

describe('RoomBookingsView', () => {
  const stubs = buildViewStubs(['v-list-item-subtitle']);

  const fetchRoomBookingsMock = vi.mocked(fetchRoomBookings);

  const mountView = () =>
    mount(RoomBookingsView, {
      global: {
        stubs
      }
    });

  beforeEach(() => {
    pushMock.mockReset();
    fetchRoomBookingsMock.mockResolvedValue({ data: [] });
  });

  it('renders bookings list', async () => {
    fetchRoomBookingsMock.mockResolvedValue({
      data: [
        {
          id: 'booking-1',
          type: 'bookings',
          attributes: {
            desk_id: 'desk-1',
            desk_name: 'Desk 1',
            user_id: 'user-1',
            user_name: 'Alice Smith',
            booking_date: '2026-01-20'
          }
        }
      ]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('Desk 1');
    expect(wrapper.text()).toContain('Alice Smith');
  });

  it('shows empty state when no bookings exist', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('No bookings for this date.');
  });

  it('fetches bookings when the date changes', async () => {
    const wrapper = mountView();
    await flushPromises();

    await wrapper.get('[data-cy="bookings-date"]').setValue('2026-01-20');
    await flushPromises();

    expect(fetchRoomBookingsMock).toHaveBeenLastCalledWith('room-1', '2026-01-20');
  });

  it('shows error for room not found', async () => {
    fetchRoomBookingsMock.mockRejectedValue(new ApiError('Not found', 404));

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('Room not found.');
  });

  it('redirects to login on 401', async () => {
    fetchRoomBookingsMock.mockRejectedValue(new ApiError('Unauthorized', 401));
    const restore = mockWindowLocation();

    mountView();
    await flushPromises();

    expect(window.location.href).toBe('/oauth/login');
    restore();
  });

  it('redirects to access denied on 403', async () => {
    fetchRoomBookingsMock.mockRejectedValue(new ApiError('Forbidden', 403));

    mountView();
    await flushPromises();

    expect(pushMock).toHaveBeenCalledWith('/access-denied');
  });
});
