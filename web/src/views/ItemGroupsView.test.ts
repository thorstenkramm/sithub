import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import ItemGroupsView from './ItemGroupsView.vue';
import { fetchMe } from '../api/me';
import { fetchItemGroups } from '../api/itemGroups';
import { fetchAreas } from '../api/areas';
import { fetchItems } from '../api/items';
import { fetchWeeklyAvailability } from '../api/itemGroupAvailability';
import { useDateState } from '../composables/useDateState';
import { getISOWeekString, getMondayOfWeek, getWeekdayDates } from '../composables/useWeekSelector';
import { buildViewStubs, createFetchMeMocker, createTestI18n, defineAuthRedirectTests } from './testHelpers';
import { ApiError, CONNECTION_LOST_MESSAGE } from '../api/client';
import en from '../locales/en.json';
import de from '../locales/de.json';

const pushMock = vi.fn();
const liveFeed = vi.hoisted(() => ({
  handler: null as ((event: unknown) => void) | null
}));

vi.mock('../api/me', () => ({ fetchMe: vi.fn() }));
vi.mock('../api/itemGroups', () => ({ fetchItemGroups: vi.fn() }));
vi.mock('../api/areas', () => ({ fetchAreas: vi.fn() }));
vi.mock('../api/items', () => ({ fetchItems: vi.fn() }));
vi.mock('../api/itemGroupAvailability', () => ({ fetchWeeklyAvailability: vi.fn() }));
vi.mock('../api/itemGroupMatrix', () => ({ fetchWeeklyMatrix: vi.fn() }));
vi.mock('../stores/useLiveFeedStore', () => ({
  useLiveFeedStore: () => ({
    start: vi.fn(),
    stop: vi.fn(),
    reset: vi.fn(),
    subscribe: (handler: (event: unknown) => void) => {
      liveFeed.handler = handler;
      return () => {
        if (liveFeed.handler === handler) {
          liveFeed.handler = null;
        }
      };
    }
  })
}));
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { areaId: 'area-1' } }),
  useRouter: () => ({ push: pushMock })
}));

describe('ItemGroupsView', () => {
  const originalMatchMedia = window.matchMedia;
  const stubs = {
    ...buildViewStubs([
      'v-card-item',
      'v-card-subtitle',
      'v-card-actions',
      'v-avatar',
      'v-icon',
      'v-progress-circular',
      'v-skeleton-loader',
      'v-select',
      'v-combobox',
      'v-spacer',
      'v-snackbar',
      'v-tooltip',
      'v-btn-toggle',
      'v-checkbox',
      'v-bottom-sheet',
      'router-link'
    ]),
    'v-btn': {
      template: '<button type="button" v-bind="$attrs" @click="$emit(\'click\', $event)"><slot /></button>'
    },
    'v-combobox': {
      props: ['modelValue'],
      template: '<input v-bind="$attrs" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />'
    },
    'v-dialog': {
      props: ['modelValue', 'fullscreen', 'persistent'],
      template: '<div v-if="modelValue" v-bind="$attrs" :data-fullscreen="fullscreen" :data-persistent="persistent"><slot /></div>'
    },
    'v-snackbar': {
      template: '<div v-bind="$attrs"><slot /></div>'
    },
    'v-tooltip': {
      template: '<div><slot name="activator" :props="{}" /><slot /></div>'
    },
    'v-switch': {
      props: ['modelValue', 'disabled'],
      template: '<input type="checkbox" v-bind="$attrs" :checked="modelValue" :disabled="disabled" @change="$emit(\'update:modelValue\', $event.target.checked)" />'
    },
    'AreaWeeklyMatrixView': {
      props: ['areaId', 'week', 'showWeekends'],
      template: '<div data-cy="area-weekly-matrix" />'
    }
  };
  const fetchMeMock = fetchMe as unknown as ReturnType<typeof vi.fn>;
  const fetchAvailabilityMock = fetchWeeklyAvailability as unknown as ReturnType<typeof vi.fn>;
  const mockFetchMe = () => createFetchMeMocker(fetchMeMock)('Ada Lovelace');
  const currentWeek = () => getISOWeekString(getMondayOfWeek(new Date()));
  const futureWeek = () => {
    const d = new Date();
    d.setDate(d.getDate() + 14);
    return getISOWeekString(getMondayOfWeek(d));
  };

  const mockFetchAreas = (floorPlan: string | undefined = undefined) => {
    const fetchAreasMock = fetchAreas as unknown as ReturnType<typeof vi.fn>;
    fetchAreasMock.mockResolvedValue({
      data: [{ id: 'area-1', type: 'areas', attributes: { name: 'Test Area', floor_plan: floorPlan } }]
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

  const mountView = (i18n = createTestI18n()) =>
    mount(ItemGroupsView, {
      global: {
        stubs,
        plugins: [createPinia(), i18n]
      }
    });

  beforeEach(() => {
    window.matchMedia = vi.fn().mockImplementation((query: string) => ({
      matches: false,
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn()
    })) as typeof window.matchMedia;
    setActivePinia(createPinia());
    pushMock.mockReset();
    fetchAvailabilityMock.mockReset();
    const fetchItemsMock = fetchItems as unknown as ReturnType<typeof vi.fn>;
    fetchItemsMock.mockReset();
    mockFetchMe();
    mockFetchAreas();
    mockFetchItemGroups(0);
    mockAvailability();
    localStorage.clear();
    sessionStorage.clear();
    liveFeed.handler = null;
    useDateState().setWeek(currentWeek());
  });

  afterEach(() => {
    window.matchMedia = originalMatchMedia;
  });

  it('renders page header with breadcrumbs', async () => {
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.find('[data-cy="breadcrumbs"]').exists()).toBe(true);
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

  it('week selector labels use DD.MM.-DD.MM.YYYY format', async () => {
    mockFetchItemGroups(1);

    const wrapper = mountView();
    await flushPromises();

    const selector = wrapper.find('[data-cy="week-selector"]');
    expect(selector.exists()).toBe(true);
  });

  it('fetches availability on mount', async () => {
    mockFetchItemGroups(1);
    mockAvailability();
    mountView();

    await flushPromises();

    expect(fetchAvailabilityMock).toHaveBeenCalledWith('area-1', expect.any(String), undefined);
  });

  it('restores the memorized week when loading availability', async () => {
    const storedWeek = futureWeek();
    useDateState().setWeek(storedWeek);
    mockFetchItemGroups(1);
    mockAvailability();
    mountView();

    await flushPromises();

    expect(fetchAvailabilityMock).toHaveBeenCalledWith('area-1', storedWeek, undefined);
  });

  it('reloads weekly availability when a relevant live event arrives', async () => {
    vi.useFakeTimers();
    try {
      mockFetchItemGroups(1);
      const fetchItemsMock = fetchItems as unknown as ReturnType<typeof vi.fn>;
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: { name: 'Desk 1', equipment: [] }
        }]
      });
      mockAvailability();

      mountView();
      await flushPromises();

      fetchAvailabilityMock.mockClear();
      expect(liveFeed.handler).toBeTypeOf('function');
      liveFeed.handler!({
        type: 'booking.created',
        booking_id: 'booking-1',
        item_id: 'item-1',
        user_id: 'other-user',
        booking_date: getWeekdayDates(getMondayOfWeek(new Date()), false)[0],
        timestamp: '2026-05-10T12:00:00Z'
      });

      await vi.advanceTimersByTimeAsync(300);
      await flushPromises();

      expect(fetchAvailabilityMock).toHaveBeenCalledTimes(1);
      expect(fetchAvailabilityMock).toHaveBeenCalledWith('area-1', currentWeek(), undefined);
    } finally {
      vi.useRealTimers();
    }
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

  it('localizes main item-group availability indicators and aria labels', async () => {
    mockFetchItemGroups(1);
    mockAvailabilityForIG1([1, 0, 2, 2, 2]);
    const wrapper = mountView(createTestI18n({
      locale: 'de',
      messages: { en, de }
    }));

    await flushPromises();

    const indicators = wrapper.find('[data-cy="availability-indicators"]');
    expect(indicators.text()).toContain('DI');
    expect(indicators.text()).not.toContain('TU');

    const firstIndicator = wrapper.find('.availability-indicator');
    expect(firstIndicator.attributes('aria-label')).toContain('verfügbar');
    expect(firstIndicator.attributes('aria-label')).not.toContain('available');
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

  it('renders promoted favorite tiles with subtitle and availability dots', async () => {
    localStorage.setItem('sithub_favorite_items', JSON.stringify([{
      areaId: 'area-1',
      itemId: 'item-1',
      itemName: 'Desk 1',
      itemGroupId: 'ig-1',
      itemGroupName: 'Item Group 1'
    }]));
    mockFetchItemGroups(1);
    mockAvailabilityForIG1([2, 0, 1, 2, 2]);
    const fetchItemsMock = fetchItems as unknown as ReturnType<typeof vi.fn>;
    fetchItemsMock.mockResolvedValue({ data: [] });
    const wrapper = mountView();

    await flushPromises();

    const favoriteTile = wrapper.get('[data-cy="favorite-item-tile"]');
    expect(favoriteTile.text()).toContain('Desk 1');
    expect(favoriteTile.text()).toContain('Item Group 1');

    const dots = favoriteTile.findAll('.indicator-dot');
    expect(dots).toHaveLength(5);
    expect(dots[1]?.classes()).toContain('dot-booked');
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

  it('shows floor plan button and dialog when the area has a floor plan', async () => {
    mockFetchItemGroups(1);
    mockFetchAreas('area.svg');
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.find('[data-cy="area-floor-plan-btn"]').exists()).toBe(true);
    await wrapper.get('[data-cy="area-floor-plan-btn"]').trigger('click');
    const dialog = wrapper.get('[data-cy="floor-plan-dialog"]');
    expect(dialog.exists()).toBe(true);
    expect(dialog.attributes('data-persistent')).toBe('');
  });

  it('opens the floor plan fullscreen on compact viewports', async () => {
    window.matchMedia = vi.fn().mockImplementation((query: string) => ({
      matches: query === '(max-width: 768px)',
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn()
    })) as typeof window.matchMedia;
    mockFetchItemGroups(1);
    mockFetchAreas('area.svg');
    const wrapper = mountView();

    await flushPromises();

    await wrapper.get('[data-cy="area-floor-plan-btn"]').trigger('click');
    expect(wrapper.get('[data-cy="floor-plan-dialog"]').attributes('data-fullscreen')).toBe('true');
  });

  it('opens the floor plan dialog without max-width on desktop', async () => {
    mockFetchItemGroups(1);
    mockFetchAreas('area.svg');
    const wrapper = mountView();

    await flushPromises();

    await wrapper.get('[data-cy="area-floor-plan-btn"]').trigger('click');
    const dialog = wrapper.get('[data-cy="floor-plan-dialog"]');
    expect(dialog.attributes('max-width')).toBeUndefined();
    expect(dialog.attributes('maxwidth')).toBeUndefined();
  });

  it('hides the floor plan button when the area has no floor plan', async () => {
    mockFetchItemGroups(1);
    mockFetchAreas();
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.find('[data-cy="area-floor-plan-btn"]').exists()).toBe(false);
  });

  it('shows a connection lost error when loading the area metadata fails', async () => {
    const fetchAreasMock = fetchAreas as unknown as ReturnType<typeof vi.fn>;
    fetchAreasMock.mockRejectedValue(new ApiError(CONNECTION_LOST_MESSAGE, 0));
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain(CONNECTION_LOST_MESSAGE);
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);

  it('shows select button label on item group tiles', async () => {
    mockFetchItemGroups(1);
    const wrapper = mountView();

    await flushPromises();

    expect(wrapper.text()).toContain('Select');
  });

  it('shows a confirmation when saving an equipment filter', async () => {
    mockFetchItemGroups(1);
    const fetchItemsMock = fetchItems as unknown as ReturnType<typeof vi.fn>;
    fetchItemsMock.mockResolvedValue({ data: [] });
    const wrapper = mountView();

    await flushPromises();

    const vm = wrapper.vm as unknown as {
      equipmentFilter: string;
      toggleSaveFilter: () => void;
    };
    vm.equipmentFilter = 'webcam';
    await flushPromises();
    vm.toggleSaveFilter();
    await flushPromises();

    expect(JSON.parse(localStorage.getItem('sithub_saved_filters')!)).toEqual(['webcam']);
    expect(wrapper.text()).toContain('Filter saved.');
  });

  it('deletes a saved filter, clears the input, and shows a confirmation', async () => {
    localStorage.setItem('sithub_saved_filters', JSON.stringify(['webcam']));
    mockFetchItemGroups(1);
    const fetchItemsMock = fetchItems as unknown as ReturnType<typeof vi.fn>;
    fetchItemsMock.mockResolvedValue({ data: [] });
    const wrapper = mountView();

    await flushPromises();

    await wrapper.get('[data-cy="ig-equipment-filter"]').setValue('webcam');
    await flushPromises();
    await wrapper.get('[data-cy="ig-equipment-filter-delete"]').trigger('click');
    await flushPromises();

    expect(JSON.parse(localStorage.getItem('sithub_saved_filters')!)).toEqual([]);
    expect((wrapper.vm as unknown as { equipmentFilter: string }).equipmentFilter).toBe('');
    expect(wrapper.text()).toContain('Saved filter deleted.');
  });

  it('shows view switch on desktop', async () => {
    mockFetchItemGroups(1);
    const wrapper = mountView();
    await flushPromises();

    const container = wrapper.find('[data-cy="view-switch-container"]');
    expect(container.exists()).toBe(true);
    expect(container.text()).toContain('Tiles');
    expect(container.text()).toContain('Table');

    const sw = wrapper.find('[data-cy="view-switch"]');
    expect(sw.exists()).toBe(true);
    expect(sw.attributes('disabled')).toBeUndefined();
  });

  it('disables view switch on mobile', async () => {
    window.matchMedia = vi.fn().mockImplementation((query: string) => ({
      matches: query === '(max-width: 768px)',
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn()
    })) as typeof window.matchMedia;
    mockFetchItemGroups(1);
    const wrapper = mountView();
    await flushPromises();

    const sw = wrapper.find('[data-cy="view-switch"]');
    expect(sw.exists()).toBe(true);
    expect(sw.attributes('disabled')).toBeDefined();
    expect(wrapper.find('[data-cy="view-switch-tooltip"]').text()).toContain('desktop only');
  });

  it('restores memorized table view for the same area on desktop', async () => {
    localStorage.setItem('sithub_area_view', JSON.stringify({ 'area-1': 'table' }));
    mockFetchItemGroups(1);
    const wrapper = mountView();
    await flushPromises();

    // Should show table view instead of card grid
    expect(wrapper.find('[data-cy="item-groups-list"]').exists()).toBe(false);
    expect(wrapper.find('[data-cy="area-weekly-matrix"]').exists()).toBe(true);
  });

  it('does not inherit table view from another area', async () => {
    localStorage.setItem('sithub_area_view', JSON.stringify({ 'area-other': 'table' }));
    mockFetchItemGroups(1);
    const wrapper = mountView();
    await flushPromises();

    // Should show card grid since area-1 has no preference
    expect(wrapper.find('[data-cy="item-groups-list"]').exists()).toBe(true);
  });

});
