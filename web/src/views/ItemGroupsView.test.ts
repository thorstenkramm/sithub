import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import ItemGroupsView from './ItemGroupsView.vue';
import { fetchMe } from '../api/me';
import { fetchItemGroups } from '../api/itemGroups';
import { fetchAreas } from '../api/areas';
import { buildViewStubs, createFetchMeMocker, defineAuthRedirectTests } from './testHelpers';

const pushMock = vi.fn();

vi.mock('../api/me', () => ({ fetchMe: vi.fn() }));
vi.mock('../api/itemGroups', () => ({ fetchItemGroups: vi.fn() }));
vi.mock('../api/areas', () => ({ fetchAreas: vi.fn() }));
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { areaId: 'area-1' } }),
  useRouter: () => ({ push: pushMock })
}));

describe('ItemGroupsView', () => {
  const stubs = buildViewStubs([
    'v-card-item',
    'v-card-subtitle',
    'v-card-actions',
    'v-avatar',
    'v-icon',
    'router-link'
  ]);
  const fetchMeMock = fetchMe as unknown as ReturnType<typeof vi.fn>;
  const mockFetchMe = () => createFetchMeMocker(fetchMeMock)('Ada Lovelace');

  const mockFetchAreas = () => {
    const fetchAreasMock = fetchAreas as unknown as ReturnType<typeof vi.fn>;
    fetchAreasMock.mockResolvedValue({
      data: [{ id: 'area-1', type: 'areas', attributes: { name: 'Test Area' } }]
    });
  };

  const mockFetchItemGroups = (count: number) => {
    const fetchItemGroupsMock = fetchItemGroups as unknown as ReturnType<typeof vi.fn>;
    fetchItemGroupsMock.mockResolvedValue({
      data: Array.from({ length: count }, (_, index) => ({
        id: `ig-${index + 1}`,
        type: 'item-groups',
        attributes: {
          name: `Item Group ${index + 1}`,
          description: index === 0 ? 'Main group' : undefined
        }
      }))
    });
  };

  const mountView = () =>
    mount(ItemGroupsView, {
      global: {
        stubs,
        plugins: [createPinia()]
      }
    });

  beforeEach(() => {
    setActivePinia(createPinia());
    pushMock.mockReset();
    mockFetchMe();
    mockFetchAreas();
    mockFetchItemGroups(0);
  });

  it('shows page header with title', async () => {
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Item Groups');
  });

  it('shows an empty state when no item groups exist', async () => {
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('No item groups available');
  });

  it('renders the item group list when data exists', async () => {
    mockFetchItemGroups(2);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Item Group 1');
    expect(wrapper.text()).toContain('Item Group 2');
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);
});
