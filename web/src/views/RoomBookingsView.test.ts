import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import RoomBookingsView from './RoomBookingsView.vue';
import { fetchRoomBookings } from '../api/roomBookings';
import { fetchAreas } from '../api/areas';
import { fetchRooms } from '../api/rooms';
import {
  buildViewStubs,
  expectLoginRedirect,
  expectAccessDeniedRedirect
} from './testHelpers';
import { ApiError } from '../api/client';

/* jscpd:ignore-start */

const pushMock = vi.fn();
vi.mock('../api/roomBookings');
vi.mock('../api/areas');
vi.mock('../api/rooms');
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { roomId: 'room-1' } }),
  useRouter: () => ({ push: pushMock })
}));

describe('RoomBookingsView', () => {
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

  const fetchRoomBookingsMock = vi.mocked(fetchRoomBookings);
  const fetchAreasMock = vi.mocked(fetchAreas);
  const fetchRoomsMock = vi.mocked(fetchRooms);

  const mountView = () =>
    mount(RoomBookingsView, {
      global: {
        stubs,
        plugins: [createPinia()]
      }
    });

  beforeEach(() => {
    setActivePinia(createPinia());
    pushMock.mockReset();
    fetchRoomBookingsMock.mockResolvedValue({ data: [] });
    fetchAreasMock.mockResolvedValue({
      data: [{ id: 'area-1', type: 'areas', attributes: { name: 'Test Area' } }]
    });
    fetchRoomsMock.mockResolvedValue({
      data: [{ id: 'room-1', type: 'rooms', attributes: { name: 'Test Room' } }]
    });
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

    expect(wrapper.text()).toContain('No bookings');
  });

  it('fetches bookings on mount', async () => {
    const wrapper = mountView();
    await flushPromises();

    // Should fetch with today's date on mount
    expect(fetchRoomBookingsMock).toHaveBeenCalled();
    const lastCall = fetchRoomBookingsMock.mock.calls[fetchRoomBookingsMock.mock.calls.length - 1];
    expect(lastCall[0]).toBe('room-1');
    expect(lastCall[1]).toMatch(/^\d{4}-\d{2}-\d{2}$/);
  });

  it('shows error for room not found', async () => {
    fetchRoomBookingsMock.mockRejectedValue(new ApiError('Not found', 404));

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('Room not found.');
  });

  it('redirects to login on 401', async () => {
    fetchRoomBookingsMock.mockRejectedValue(new ApiError('Unauthorized', 401));
    await expectLoginRedirect(mountView);
  });

  it('redirects to access denied on 403', async () => {
    fetchRoomBookingsMock.mockRejectedValue(new ApiError('Forbidden', 403));
    await expectAccessDeniedRedirect(mountView, pushMock);
  });
});
/* jscpd:ignore-end */
