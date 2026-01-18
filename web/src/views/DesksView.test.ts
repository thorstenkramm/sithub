import { mount, flushPromises } from '@vue/test-utils';
import DesksView from './DesksView.vue';
import { fetchDesks } from '../api/desks';
import { fetchMe } from '../api/me';
import { buildViewStubs, defineAuthRedirectTests } from './testHelpers';

const pushMock = vi.fn();

vi.mock('../api/me');
vi.mock('../api/desks');

vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { roomId: 'room-1' } }),
  useRouter: () => ({
    push: pushMock
  })
}));

describe('DesksView', () => {
  const stubs = buildViewStubs(['v-list-item-subtitle']);

  const fetchMeMock = vi.mocked(fetchMe);
  const fetchDesksMock = vi.mocked(fetchDesks);

  const mountView = () =>
    mount(DesksView, {
      global: {
        stubs
      }
    });

  beforeEach(() => {
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

  it('renders desk equipment and warning', async () => {
    fetchDesksMock.mockResolvedValue({
      data: [
        {
          id: 'desk-1',
          type: 'desks',
          attributes: {
            name: 'Desk 1',
            equipment: ['Monitor', 'Keyboard'],
            warning: 'USB-C only'
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

    expect(wrapper.text()).toContain('No desks available.');
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);
});
