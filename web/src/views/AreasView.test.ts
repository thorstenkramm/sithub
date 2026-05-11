import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import AreasView from './AreasView.vue';
import { fetchAreas } from '../api/areas';
import { fetchMe } from '../api/me';
import { buildViewStubs, createFetchMeMocker, createTestI18n, defineAuthRedirectTests } from './testHelpers';
import { ApiError, CONNECTION_LOST_MESSAGE } from '../api/client';
import { __resetLegacyPurgeForTests } from '../composables/useFavorites';

const pushMock = vi.fn();

vi.mock('../api/me', () => ({ fetchMe: vi.fn() }));
vi.mock('../api/areas', () => ({ fetchAreas: vi.fn() }));
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));

describe('AreasView', () => {
  const stubs = buildViewStubs([
    'v-card-item',
    'v-card-subtitle',
    'v-card-actions',
    'v-avatar',
    'v-icon',
    'router-link'
  ]);
  const fetchMeMock = fetchMe as unknown as ReturnType<typeof vi.fn>;
  const mockFetchMeBase = createFetchMeMocker(fetchMeMock);
  const mockFetchMe = (isAdmin: boolean) => mockFetchMeBase('Ada Lovelace', isAdmin);

  const mockFetchAreas = (count: number) => {
    const fetchAreasMock = fetchAreas as unknown as ReturnType<typeof vi.fn>;
    fetchAreasMock.mockResolvedValue({
      data: Array.from({ length: count }, (_, index) => ({
        id: `area-${index + 1}`,
        type: 'areas',
        attributes: {
          name: `Area ${index + 1}`,
          description: index === 0 ? 'Main area' : undefined,
          floor_plan: index === 0 ? 'area.svg' : undefined
        }
      }))
    });
  };

  const mountView = () =>
    mount(AreasView, {
      global: {
        stubs,
        plugins: [createPinia(), createTestI18n()]
      }
    });

  beforeEach(() => {
    setActivePinia(createPinia());
    pushMock.mockReset();
    localStorage.clear();
    __resetLegacyPurgeForTests();
    mockFetchAreas(0);
  });

  it('shows page header with title', async () => {
    mockFetchMe(false);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Areas');
  });

  it('shows admin controls for admins only', async () => {
    mockFetchMe(true);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Cancel booking (admin)');
  });

  it('hides admin controls for non-admin users', async () => {
    mockFetchMe(false);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).not.toContain('Cancel booking (admin)');
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);

  it('shows an empty state when no areas exist', async () => {
    mockFetchMe(false);
    mockFetchAreas(0);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('No areas available');
  });

  it('renders the areas list when data exists', async () => {
    mockFetchMe(false);
    mockFetchAreas(2);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Area 1');
    expect(wrapper.text()).toContain('Area 2');
  });

  it('shows a connection lost error when user loading fails', async () => {
    fetchMeMock.mockRejectedValue(new ApiError(CONNECTION_LOST_MESSAGE, 0));
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain(CONNECTION_LOST_MESSAGE);
  });

  it('shows select button label on area tiles', async () => {
    mockFetchMe(false);
    mockFetchAreas(1);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Select');
  });

  it('shows the Favorites virtual tile when at least one favorite exists', async () => {
    localStorage.setItem('sithub_favorite_items', JSON.stringify([{
      areaId: 'area-1',
      itemId: 'desk-1',
      itemName: 'Desk 1',
      itemGroupId: 'ig-1',
      itemGroupName: 'Room 1'
    }]));
    mockFetchMe(false);
    mockFetchAreas(1);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.find('[data-cy="favorites-tile"]').exists()).toBe(true);
  });

  it('shows the Favorites virtual tile even when no real areas exist', async () => {
    localStorage.setItem('sithub_favorite_items', JSON.stringify([{
      areaId: 'area-1',
      itemId: 'desk-1',
      itemName: 'Desk 1',
      itemGroupId: 'ig-1',
      itemGroupName: 'Room 1'
    }]));
    mockFetchMe(false);
    mockFetchAreas(0);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.find('[data-cy="favorites-tile"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="areas-empty"]').exists()).toBe(false);
  });

  it('hides the Favorites virtual tile when no favorites exist', async () => {
    localStorage.removeItem('sithub_favorite_items');
    mockFetchMe(false);
    mockFetchAreas(1);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.find('[data-cy="favorites-tile"]').exists()).toBe(false);
  });

});
