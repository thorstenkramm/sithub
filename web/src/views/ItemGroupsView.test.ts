import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import ItemGroupsView from './ItemGroupsView.vue';
import { fetchMe } from '../api/me';
import { fetchItemGroups } from '../api/itemGroups';
import { fetchAreas } from '../api/areas';
import { fetchWeeklyAvailability } from '../api/itemGroupAvailability';
import { buildViewStubs, createFetchMeMocker, defineAuthRedirectTests } from './testHelpers';

const pushMock = vi.fn();

vi.mock('../api/me', () => ({ fetchMe: vi.fn() }));
vi.mock('../api/itemGroups', () => ({ fetchItemGroups: vi.fn() }));
vi.mock('../api/areas', () => ({ fetchAreas: vi.fn() }));
vi.mock('../api/itemGroupAvailability', () => ({ fetchWeeklyAvailability: vi.fn() }));
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
    'v-skeleton-loader',
    'v-select',
    'router-link'
  ]);
  const fetchMeMock = fetchMe as unknown as ReturnType<typeof vi.fn>;
  const fetchAvailabilityMock = fetchWeeklyAvailability as unknown as ReturnType<typeof vi.fn>;
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

  const mockAvailability = (data: unknown[] = []) => {
    fetchAvailabilityMock.mockResolvedValue({ data });
  };

  const makeDays = (availabilities: number[]) => {
    const baseDate = new Date(2026, 2, 16); // 2026-03-16 Monday
    const weekdays = ['MO', 'TU', 'WE', 'TH', 'FR'];
    return availabilities.map((available, i) => {
      const d = new Date(baseDate);
      d.setDate(baseDate.getDate() + i);
      const dateStr = d.toISOString().slice(0, 10);
      return { date: dateStr, weekday: weekdays[i], total: 2, available };
    });
  };

  const mockAvailabilityForIG1 = (availabilities: number[]) => {
    mockAvailability([{
      id: 'ig-1',
      type: 'item-group-availability',
      attributes: {
        item_group_id: 'ig-1',
        item_group_name: 'Item Group 1',
        days: makeDays(availabilities)
      }
    }]);
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
    fetchAvailabilityMock.mockReset();
    mockFetchMe();
    mockFetchAreas();
    mockFetchItemGroups(0);
    mockAvailability();
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

  it('navigates to items with areaId in query', async () => {
    mockFetchItemGroups(1);
    const wrapper = mountView();
    await flushPromises();

    const card = wrapper.find('[data-cy="item-group-item"]');
    await card.trigger('click');

    expect(pushMock).toHaveBeenCalledWith({
      name: 'items',
      params: { itemGroupId: 'ig-1' },
      query: { areaId: 'area-1' }
    });
  });

  it('renders the week selector card', async () => {
    mockFetchItemGroups(1);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.find('[data-cy="week-selector-card"]').exists()).toBe(true);
  });

  it('uses locale-aware date formatting for week labels', async () => {
    mockFetchItemGroups(1);
    const formatterSpy = vi.spyOn(Intl, 'DateTimeFormat');

    const wrapper = mountView();
    await flushPromises();

    expect(formatterSpy).toHaveBeenCalledWith(undefined, {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    });

    formatterSpy.mockRestore();

    expect(wrapper.find('[data-cy="week-selector"]').exists()).toBe(true);
  });

  it('fetches availability on mount', async () => {
    mockFetchItemGroups(1);
    mockAvailability();
    mountView();

    await flushPromises();

    expect(fetchAvailabilityMock).toHaveBeenCalledWith('area-1', expect.any(String));
  });

  it('shows availability indicators when data is returned', async () => {
    mockFetchItemGroups(1);
    mockAvailabilityForIG1([1, 0, 2, 2, 2]);
    const wrapper = mountView();

    await flushPromises();

    const indicators = wrapper.find('[data-cy="availability-indicators"]');
    expect(indicators.exists()).toBe(true);
    expect(indicators.text()).toContain('MO');
    expect(indicators.text()).toContain('TU');
    expect(indicators.text()).toContain('FR');
  });

  it('shows available dot for days with availability', async () => {
    mockFetchItemGroups(1);
    mockAvailabilityForIG1([2, 0, 1, 0, 2]);
    const wrapper = mountView();

    await flushPromises();

    const dots = wrapper.findAll('.indicator-dot');
    expect(dots.length).toBe(5);
    expect(dots[0].classes()).toContain('dot-available');
    expect(dots[1].classes()).toContain('dot-booked');
    expect(dots[2].classes()).toContain('dot-available');
    expect(dots[3].classes()).toContain('dot-booked');
    expect(dots[4].classes()).toContain('dot-available');
  });

  it('handles availability fetch failure gracefully', async () => {
    mockFetchItemGroups(1);
    fetchAvailabilityMock.mockRejectedValue(new Error('Network error'));
    const wrapper = mountView();

    await flushPromises();

    // Should still render item groups without availability indicators
    expect(wrapper.text()).toContain('Item Group 1');
    expect(wrapper.find('[data-cy="availability-indicators"]').exists()).toBe(false);
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);
});
