import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import BookingHistoryView from './BookingHistoryView.vue';
import { fetchBookingHistory } from '../api/bookings';
import { fetchMe } from '../api/me';
import { buildViewStubs, defineAuthRedirectTests } from './testHelpers';
import { ApiError, CONNECTION_LOST_MESSAGE } from '../api/client';

const pushMock = vi.fn();

vi.mock('../api/bookings', () => ({ fetchBookingHistory: vi.fn() }));
vi.mock('../api/me', () => ({ fetchMe: vi.fn() }));
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));

describe('BookingHistoryView', () => {
  const stubs = buildViewStubs([
    'v-list-item-subtitle',
    'v-card-item',
    'v-card-subtitle',
    'v-card-actions',
    'v-avatar',
    'v-icon',
    'v-chip',
    'v-menu',
    'v-date-picker',
    'v-spacer',
    'v-skeleton-loader',
    'v-text-field',
    'router-link'
  ]);

  const fetchMeMock = fetchMe as unknown as ReturnType<typeof vi.fn>;
  const fetchBookingHistoryMock = fetchBookingHistory as unknown as ReturnType<typeof vi.fn>;

  const mountView = () =>
    mount(BookingHistoryView, {
      global: {
        stubs,
        plugins: [createPinia()]
      }
    });

  beforeEach(() => {
    setActivePinia(createPinia());
    pushMock.mockReset();
    fetchMeMock.mockResolvedValue({
      data: { attributes: { display_name: 'Ada Lovelace', is_admin: false } }
    });
    fetchBookingHistoryMock.mockResolvedValue({ data: [] });
  });

  it('shows the empty state when no bookings exist', async () => {
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('No bookings found');
  });

  it('shows a connection lost error when user loading fails', async () => {
    fetchMeMock.mockRejectedValue(new ApiError(CONNECTION_LOST_MESSAGE, 0));
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain(CONNECTION_LOST_MESSAGE);
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);
});
