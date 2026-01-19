import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import DesksView from './DesksView.vue';
import { fetchDesks } from '../api/desks';
import { fetchMe } from '../api/me';
import { buildViewStubs, defineAuthRedirectTests } from './testHelpers';

const pushMock = vi.fn();
vi.mock('../api/me');
vi.mock('../api/desks');
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { roomId: 'room-1' } }),
  useRouter: () => ({ push: pushMock })
}));

describe('DesksView', () => {
  const stubs = buildViewStubs(['v-list-item-subtitle']);

  const fetchMeMock = vi.mocked(fetchMe);
  const fetchDesksMock = vi.mocked(fetchDesks);

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
    expect(wrapper.text()).toContain('Status: Occupied');
  });

  it('shows empty state when no desks exist', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('No desks available.');
  });

  it('fetches desks when the date changes', async () => {
    const wrapper = mountView();
    await flushPromises();

    await wrapper.get('[data-cy=\"desks-date\"]').setValue('2026-01-20');
    await flushPromises();

    expect(fetchDesksMock).toHaveBeenLastCalledWith('room-1', '2026-01-20');
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);
});
