import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import RoomsView from './RoomsView.vue';
import { fetchMe } from '../api/me';
import { fetchRooms } from '../api/rooms';
import { fetchAreas } from '../api/areas';
import { buildViewStubs, createFetchMeMocker, defineAuthRedirectTests } from './testHelpers';

const pushMock = vi.fn();

vi.mock('../api/me', () => ({ fetchMe: vi.fn() }));
vi.mock('../api/rooms', () => ({ fetchRooms: vi.fn() }));
vi.mock('../api/areas', () => ({ fetchAreas: vi.fn() }));
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { areaId: 'area-1' } }),
  useRouter: () => ({ push: pushMock })
}));

describe('RoomsView', () => {
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

  const mockFetchRooms = (count: number) => {
    const fetchRoomsMock = fetchRooms as unknown as ReturnType<typeof vi.fn>;
    fetchRoomsMock.mockResolvedValue({
      data: Array.from({ length: count }, (_, index) => ({
        id: `room-${index + 1}`,
        type: 'rooms',
        attributes: {
          name: `Room ${index + 1}`,
          description: index === 0 ? 'Main room' : undefined
        }
      }))
    });
  };

  const mountView = () =>
    mount(RoomsView, {
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
    mockFetchRooms(0);
  });

  it('shows page header with title', async () => {
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Rooms');
  });

  it('shows an empty state when no rooms exist', async () => {
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('No rooms available');
  });

  it('renders the room list when data exists', async () => {
    mockFetchRooms(2);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Room 1');
    expect(wrapper.text()).toContain('Room 2');
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);
});
