import { mount, flushPromises } from '@vue/test-utils';
import AreasView from './AreasView.vue';
import { fetchAreas } from '../api/areas';
import { fetchMe } from '../api/me';
import { buildViewStubs, createFetchMeMocker, defineAuthRedirectTests } from './testHelpers';

const pushMock = vi.fn();

vi.mock('../api/me', () => ({ fetchMe: vi.fn() }));
vi.mock('../api/areas', () => ({ fetchAreas: vi.fn() }));
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));

describe('AreasView', () => {
  const stubs = buildViewStubs();
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
          floor_plan: index === 0 ? 'floor_plans/area.svg' : undefined
        }
      }))
    });
  };

  const mountView = () =>
    mount(AreasView, {
      global: {
        stubs
      }
    });

  beforeEach(() => {
    pushMock.mockReset();
    mockFetchAreas(0);
  });

  it('shows the signed-in user name', async () => {
    mockFetchMe(false);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Signed in as Ada Lovelace');
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

    expect(wrapper.text()).toContain('No areas available.');
  });

  it('renders the areas list when data exists', async () => {
    mockFetchMe(false);
    mockFetchAreas(2);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Area 1');
    expect(wrapper.text()).toContain('Area 2');
  });
});
