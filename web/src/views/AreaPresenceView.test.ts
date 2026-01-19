import { mount, flushPromises } from '@vue/test-utils';
import AreaPresenceView from './AreaPresenceView.vue';
import { fetchAreaPresence } from '../api/areaPresence';
import {
  buildViewStubs,
  expectLoginRedirect,
  expectAccessDeniedRedirect
} from './testHelpers';
import { ApiError } from '../api/client';

const pushMock = vi.fn();
vi.mock('../api/areaPresence');
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { areaId: 'area-1' } }),
  useRouter: () => ({ push: pushMock })
}));

describe('AreaPresenceView', () => {
  const stubs = buildViewStubs(['v-list-item-subtitle']);

  const fetchAreaPresenceMock = vi.mocked(fetchAreaPresence);

  const mountView = () =>
    mount(AreaPresenceView, {
      global: {
        stubs
      }
    });

  beforeEach(() => {
    pushMock.mockReset();
    fetchAreaPresenceMock.mockResolvedValue({ data: [] });
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
    expect(wrapper.text()).toContain('Room One - Desk 1');
  });

  it('shows empty state when no one is present', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('No one is scheduled for this date.');
  });

  it('fetches presence when the date changes', async () => {
    const wrapper = mountView();
    await flushPromises();

    await wrapper.get('[data-cy="presence-date"]').setValue('2026-01-20');
    await flushPromises();

    expect(fetchAreaPresenceMock).toHaveBeenLastCalledWith('area-1', '2026-01-20');
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
