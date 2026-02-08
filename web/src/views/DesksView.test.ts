import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import DesksView from './DesksView.vue';
import { fetchDesks } from '../api/desks';
import { fetchMe } from '../api/me';
import { fetchRooms } from '../api/rooms';
import { fetchAreas } from '../api/areas';
import { buildViewStubs, defineAuthRedirectTests } from './testHelpers';

/* jscpd:ignore-start */

const pushMock = vi.fn();
vi.mock('../api/me');
vi.mock('../api/desks');
vi.mock('../api/rooms');
vi.mock('../api/areas');
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { roomId: 'room-1' } }),
  useRouter: () => ({ push: pushMock })
}));

describe('DesksView', () => {
  const stubs = buildViewStubs([
    'v-list-item-subtitle',
    'v-card-item',
    'v-card-actions',
    'v-avatar',
    'v-icon',
    'v-chip',
    'v-radio',
    'v-radio-group',
    'v-text-field',
    'v-checkbox',
    'v-expand-transition',
    'v-menu',
    'v-date-picker',
    'v-skeleton-loader',
    'router-link'
  ]);

  const fetchMeMock = vi.mocked(fetchMe);
  const fetchDesksMock = vi.mocked(fetchDesks);
  const fetchRoomsMock = vi.mocked(fetchRooms);
  const fetchAreasMock = vi.mocked(fetchAreas);

  const mountView = () =>
    mount(DesksView, {
      global: {
        stubs,
        plugins: [createPinia()]
      }
    });

  beforeEach(() => {
    setActivePinia(createPinia());
    pushMock.mockReset();
    fetchMeMock.mockResolvedValue({
      data: {
        attributes: {
          display_name: 'Ada Lovelace',
          is_admin: false
        }
      }
    });
    fetchDesksMock.mockResolvedValue({ data: [] });
    fetchAreasMock.mockResolvedValue({
      data: [{ id: 'area-1', type: 'areas', attributes: { name: 'Test Area' } }]
    });
    fetchRoomsMock.mockResolvedValue({
      data: [{ id: 'room-1', type: 'rooms', attributes: { name: 'Test Room' } }]
    });
  });

  it('renders desk equipment, warning, and status', async () => {
    fetchDesksMock.mockResolvedValue({
      data: [
        {
          id: 'desk-1',
          type: 'desks',
          attributes: {
            name: 'Desk 1',
            equipment: ['Monitor', 'Keyboard'],
            warning: 'USB-C only',
            availability: 'occupied'
          }
        }
      ]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('Desk 1');
    expect(wrapper.text()).toContain('Monitor');
    expect(wrapper.text()).toContain('USB-C only');
  });

  it('shows empty state when no desks exist', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('No desks available');
  });

  it('fetches desks on mount with current date', async () => {
    mountView();
    await flushPromises();

    // Should fetch desks with today's date on mount
    expect(fetchDesksMock).toHaveBeenCalled();
    // Check that it was called with room-1 and a date in YYYY-MM-DD format
    const lastCall = fetchDesksMock.mock.calls[fetchDesksMock.mock.calls.length - 1];
    expect(lastCall[0]).toBe('room-1');
    expect(lastCall[1]).toMatch(/^\d{4}-\d{2}-\d{2}$/);
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);
});
/* jscpd:ignore-end */
