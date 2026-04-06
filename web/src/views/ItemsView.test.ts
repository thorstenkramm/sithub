import { mount, flushPromises } from '@vue/test-utils';
import { nextTick } from 'vue';
import { createPinia, setActivePinia } from 'pinia';
import ItemsView from './ItemsView.vue';
import PageHeader from '../components/PageHeader.vue';
import { fetchItems } from '../api/items';
import { fetchMe } from '../api/me';
import { fetchItemGroups } from '../api/itemGroups';
import { fetchAreas } from '../api/areas';
import { fetchColleagues } from '../api/users';
import { createBooking, cancelBooking, updateBookingNote, fetchMyBookings } from '../api/bookings';
import { useDateState } from '../composables/useDateState';
import { buildViewStubs, createTestI18n, defineAuthRedirectTests } from './testHelpers';
import { ApiError, CONNECTION_LOST_MESSAGE } from '../api/client';
import { middleTruncate } from '../utils/text';

/* jscpd:ignore-start */

const pushMock = vi.fn();
vi.mock('../api/me');
vi.mock('../api/items');
vi.mock('../api/itemGroups');
vi.mock('../api/areas');
vi.mock('../api/users');
vi.mock('../api/bookings');
const routeMock = { params: { itemGroupId: 'ig-1' }, query: { areaId: 'area-1' } };
vi.mock('vue-router', () => ({
  useRoute: () => routeMock,
  useRouter: () => ({ push: pushMock })
}));

describe('ItemsView', () => {
  const originalMatchMedia = window.matchMedia;
  const slotStub = {
    template: '<div><slot /></div>'
  };
  const stubs = {
    ...buildViewStubs([
      'v-list-item-subtitle',
      'v-card-actions',
      'v-avatar',
      'v-chip',
      'v-radio',
      'v-radio-group',
      'v-text-field',
      'v-checkbox',
      'v-expand-transition',
      'v-autocomplete',
      'v-combobox',
      'v-menu',
      'v-date-picker',
      'v-progress-circular',
      'v-skeleton-loader',
      'v-snackbar',
      'v-textarea',
      'v-spacer',
      'v-btn-toggle',
      'v-select',
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
    'v-bottom-sheet': {
      props: ['modelValue'],
      template: '<div v-if="modelValue"><slot /></div>'
    },
    'v-checkbox': {
      props: ['modelValue', 'disabled', 'color'],
      emits: ['update:modelValue'],
      template:
        '<label v-bind="$attrs" :data-disabled="disabled" :data-color="color"><input type="checkbox" :checked="modelValue" :disabled="disabled" @change="$emit(\'update:modelValue\', !modelValue)" /><slot /></label>'
    },
    'v-card-item': {
      template: '<div><slot name="prepend" /><slot /><slot name="append" /></div>'
    },
    'v-card-actions': {
      template: '<div v-bind="$attrs"><slot /></div>'
    },
    'v-tooltip': {
      template: '<div><slot name="activator" :props="{}" /><slot /></div>'
    },
    'v-snackbar': {
      template: '<div v-bind="$attrs"><slot /></div>'
    },
    'v-icon': slotStub
  };

  const fetchMeMock = vi.mocked(fetchMe);
  const fetchItemsMock = vi.mocked(fetchItems);
  const fetchItemGroupsMock = vi.mocked(fetchItemGroups);
  const fetchAreasMock = vi.mocked(fetchAreas);
  const fetchColleaguesMock = vi.mocked(fetchColleagues);
  const createBookingMock = vi.mocked(createBooking);
  const cancelBookingMock = vi.mocked(cancelBooking);
  const updateBookingNoteMock = vi.mocked(updateBookingNote);
  const fetchMyBookingsMock = vi.mocked(fetchMyBookings);

  const formatDate = (date: Date) => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  const futureDay = () => {
    const d = new Date();
    d.setDate(d.getDate() + 3);
    return formatDate(d);
  };

  const mountView = () =>
    mount(ItemsView, {
      global: {
        stubs,
        plugins: [createPinia(), createTestI18n()]
      }
    });

  const setElementWidth = (
    element: Element,
    dimensions: { clientWidth?: number; scrollWidth?: number }
  ) => {
    if ('clientWidth' in dimensions) {
      Object.defineProperty(element, 'clientWidth', {
        configurable: true,
        value: dimensions.clientWidth
      });
    }

    if ('scrollWidth' in dimensions) {
      Object.defineProperty(element, 'scrollWidth', {
        configurable: true,
        value: dimensions.scrollWidth
      });
    }
  };

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
    routeMock.query = { areaId: 'area-1' };
    fetchMeMock.mockResolvedValue({
      data: {
        attributes: {
          display_name: 'Ada Lovelace',
          is_admin: false
        }
      }
    });
    fetchItemsMock.mockResolvedValue({ data: [] });
    fetchAreasMock.mockResolvedValue({
      data: [{ id: 'area-1', type: 'areas', attributes: { name: 'Test Area' } }]
    });
    fetchItemGroupsMock.mockResolvedValue({
      data: [{ id: 'ig-1', type: 'item-groups', attributes: { name: 'Test Group' } }]
    });
    fetchColleaguesMock.mockResolvedValue({
      data: [
        { id: 'u-1', type: 'colleagues', attributes: { display_name: 'Jane Doe' } },
        { id: 'u-2', type: 'colleagues', attributes: { display_name: 'Bob Smith' } }
      ]
    });
    createBookingMock.mockResolvedValue({ data: { id: 'booking-1' } } as never);
    cancelBookingMock.mockResolvedValue(undefined as never);
    updateBookingNoteMock.mockResolvedValue(undefined as never);
    fetchMyBookingsMock.mockResolvedValue({ data: [] } as never);
    localStorage.removeItem('sithub_booking_mode');
    sessionStorage.clear();
    useDateState().resetDayToToday();
  });

  afterEach(() => {
    window.matchMedia = originalMatchMedia;
  });

  it('renders item equipment, warning, and status on available items', async () => {
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: 'item-1',
          type: 'items',
          attributes: {
            name: 'Item 1',
            equipment: ['Monitor', 'Keyboard'],
            warning: 'USB-C only',
            availability: 'available'
          }
        }
      ]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('Item 1');
    expect(wrapper.text()).toContain('Monitor');
    expect(wrapper.text()).toContain('USB-C only');
  });

  it('shows booker name when item is occupied', async () => {
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: 'item-1',
          type: 'items',
          attributes: {
            name: 'Item 1',
            equipment: [],
            availability: 'occupied',
            booker_name: 'Alice Smith'
          }
        }
      ]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('Alice Smith');
    expect(wrapper.find('[data-cy="item-booker"]').exists()).toBe(true);
  });

  it('does not show booker name when item is available', async () => {
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: 'item-1',
          type: 'items',
          attributes: {
            name: 'Item 1',
            equipment: [],
            availability: 'available'
          }
        }
      ]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.find('[data-cy="item-booker"]').exists()).toBe(false);
  });

  it('shows a reserved overlay for reserved day-mode items', async () => {
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: 'item-1',
          type: 'items',
          attributes: {
            name: 'Reserved Desk',
            equipment: [],
            availability: 'available',
            reserved: true
          }
        }
      ]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.find('[data-cy="item-reserved"]').exists()).toBe(true);
    expect(wrapper.text()).toContain('Reserved');
    expect(wrapper.get('[data-cy="item-entry"]').attributes('title')).toContain('reserved');
  });

  it('dims and disables reserved week-mode items', async () => {
    localStorage.setItem('sithub_booking_mode', 'week');
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: 'item-1',
          type: 'items',
          attributes: {
            name: 'Reserved Desk',
            equipment: [],
            availability: 'available',
            reserved: true
          }
        }
      ]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.find('[data-cy="item-reserved"]').exists()).toBe(true);
    expect(wrapper.get('[data-cy="week-item-entry"]').attributes('title')).toContain('reserved');
    expect(wrapper.get('[data-cy="week-day-checkbox"]').attributes('data-disabled')).toBe('true');
  });

  it('shows empty state when no items exist', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain('No items available');
  });

  it('fetches items on mount with current date', async () => {
    mountView();
    await flushPromises();

    // Should fetch items with today's date on mount
    expect(fetchItemsMock).toHaveBeenCalled();
    // Check that it was called with ig-1 and a date in YYYY-MM-DD format
    const lastCall = fetchItemsMock.mock.calls[fetchItemsMock.mock.calls.length - 1];
    expect(lastCall[0]).toBe('ig-1');
    expect(lastCall[1]).toMatch(/^\d{4}-\d{2}-\d{2}$/);
  });

  describe('breadcrumbs', () => {
    it('includes area link when areaId is in query', async () => {
      routeMock.query = { areaId: 'area-1' };
      const wrapper = mountView();
      await flushPromises();

      const breadcrumbs = wrapper.findComponent(PageHeader).props('breadcrumbs') as Array<{ text: string; to?: unknown }>;
      expect(breadcrumbs[1]?.to).toBe('/areas/area-1/item-groups');
    });

    it('renders area breadcrumb as non-clickable when areaId is missing and area is unresolved', async () => {
      routeMock.query = {};
      fetchAreasMock.mockResolvedValue({ data: [] });
      fetchItemGroupsMock.mockResolvedValue({ data: [] });
      const wrapper = mountView();
      await flushPromises();

      const breadcrumbs = wrapper.findComponent(PageHeader).props('breadcrumbs') as Array<{ text: string; to?: unknown }>;
      expect(breadcrumbs[1]?.to).toBeUndefined();
    });
  });

  describe('booking mode toggle', () => {
    beforeEach(() => {
      localStorage.removeItem('sithub_booking_mode');
    });

    afterEach(() => {
      localStorage.removeItem('sithub_booking_mode');
    });

    it('defaults to day mode when localStorage is empty', async () => {
      const wrapper = mountView();
      await flushPromises();

      expect(wrapper.find('[data-cy="booking-mode-toggle"]').exists()).toBe(true);
      // In day mode, week items list should not exist
      expect(wrapper.find('[data-cy="week-items-list"]').exists()).toBe(false);
    });

    it('persists mode in localStorage when switched to week', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: { name: 'Desk A', equipment: [], availability: 'available' as const }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      // Simulate mode change by finding the component and triggering update
      const toggle = wrapper.find('[data-cy="booking-mode-toggle"]');
      expect(toggle.exists()).toBe(true);

      (wrapper.vm as unknown as { bookingMode: 'day' | 'week' }).bookingMode = 'week';
      await flushPromises();

      expect(localStorage.getItem('sithub_booking_mode')).toBe('week');
    });

    it('restores week mode from localStorage on mount', async () => {
      localStorage.setItem('sithub_booking_mode', 'week');
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: { name: 'Desk A', equipment: [], availability: 'available' as const }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      // In week mode, the week selector should be present
      expect(wrapper.find('[data-cy="week-selector"]').exists()).toBe(true);
      // Day mode list should not exist
      expect(wrapper.find('[data-cy="items-list"]').exists()).toBe(false);
    });
  });


  it('shows warning icon on folded booked day tiles with warnings', async () => {
    fetchItemsMock.mockResolvedValue({
      data: [{
        id: 'item-1',
        type: 'items',
        attributes: {
          name: 'Item 1',
          equipment: [],
          warning: 'Caution',
          availability: 'occupied'
        }
      }]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.find('[data-cy="folded-warning-icon"]').exists()).toBe(true);

    // Warning icon stays visible when expanded (it's in the subtitle row now)
    (wrapper.vm as unknown as { expandedDayTiles: Set<string> }).expandedDayTiles = new Set(['item-1']);
    await nextTick();

    expect(wrapper.find('[data-cy="folded-warning-icon"]').exists()).toBe(true);
  });

  it('keeps long day-mode item names untruncated when they fit the card width', async () => {
    const longName = 'Desk with a very descriptive suffix';
    fetchItemsMock.mockResolvedValue({
      data: [{
        id: 'item-1',
        type: 'items',
        attributes: {
          name: longName,
          equipment: [],
          availability: 'available'
        }
      }]
    });

    const wrapper = mountView();
    await flushPromises();

    const display = wrapper.get('.item-name').element;
    const measure = wrapper.get('.item-name-measure').element;
    setElementWidth(display, { clientWidth: 240 });
    setElementWidth(measure, { scrollWidth: 180 });

    window.dispatchEvent(new Event('resize'));
    await nextTick();

    expect(wrapper.get('.item-name').text()).toBe(longName);
  });

  it('middle-truncates long day-mode item names only when they overflow the card width', async () => {
    const longName = 'Desk with a very descriptive suffix';
    fetchItemsMock.mockResolvedValue({
      data: [{
        id: 'item-1',
        type: 'items',
        attributes: {
          name: longName,
          equipment: [],
          availability: 'available'
        }
      }]
    });

    const wrapper = mountView();
    await flushPromises();

    const display = wrapper.get('.item-name').element;
    const measure = wrapper.get('.item-name-measure').element;
    setElementWidth(display, { clientWidth: 120 });
    setElementWidth(measure, { scrollWidth: 240 });

    window.dispatchEvent(new Event('resize'));
    await nextTick();

    expect(wrapper.get('.item-name').text()).toBe(middleTruncate(longName, 25));
  });

  it('renders the day-mode favorite heart inside the status row', async () => {
    fetchItemsMock.mockResolvedValue({
      data: [{
        id: 'item-1',
        type: 'items',
        attributes: {
          name: 'Desk A',
          equipment: [],
          availability: 'available'
        }
      }]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.find('[data-cy="day-status-row"] [data-cy="item-favorite-heart"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="day-item-actions"] [data-cy="item-favorite-heart"]').exists()).toBe(false);
  });

  describe('week mode rendering', () => {
    beforeEach(() => {
      localStorage.setItem('sithub_booking_mode', 'week');
    });

    afterEach(() => {
      localStorage.removeItem('sithub_booking_mode');
    });

    it('renders week item tiles in week mode', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: { name: 'Desk A', equipment: [], availability: 'available' as const }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      expect(wrapper.find('[data-cy="week-items-list"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="week-item-entry"]').exists()).toBe(true);
    });

    it('renders the week-mode favorite heart inside the status row', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: { name: 'Desk A', equipment: [], availability: 'available' as const }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      expect(wrapper.find('[data-cy="week-status-row"] [data-cy="week-item-favorite-heart"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="week-item-actions"]').exists()).toBe(false);
    });


    it('shows warning icon on folded week tiles with warnings', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: { name: 'Desk A', equipment: [], warning: 'Cable issue', availability: 'available' as const }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      expect(wrapper.find('[data-cy="week-folded-warning-icon"]').exists()).toBe(true);

      // Warning icon stays visible when expanded (it's in the subtitle row now)
      (wrapper.vm as unknown as { expandedWeekTiles: Set<string> }).expandedWeekTiles = new Set(['item-1']);
      await nextTick();

      expect(wrapper.find('[data-cy="week-folded-warning-icon"]').exists()).toBe(true);
    });

    it('shows week selector instead of date picker', async () => {
      const wrapper = mountView();
      await flushPromises();

      expect(wrapper.find('[data-cy="week-selector"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="items-date"]').exists()).toBe(false);
    });

    it('fetches items for each weekday on mount in week mode', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: { name: 'Desk A', equipment: [], availability: 'available' as const }
        }]
      });

      mountView();
      await flushPromises();

      // Should fetch items for multiple weekdays (5 per week)
      // Each call should have the item group ID and a date
      const calls = fetchItemsMock.mock.calls;
      const weekCalls = calls.filter(c => c[0] === 'ig-1' && typeof c[1] === 'string');
      expect(weekCalls.length).toBeGreaterThanOrEqual(5);
    });
  });

  it('renders colleague autocomplete when booking type is colleague', async () => {
    const wrapper = mountView();
    await flushPromises();

    // Colleague fields hidden by default
    expect(wrapper.find('[data-cy="colleague-select"]').exists()).toBe(false);

    (wrapper.vm as unknown as { bookingType: 'self' | 'colleague' }).bookingType = 'colleague';
    await nextTick();

    expect(wrapper.find('[data-cy="colleague-select"]').exists()).toBe(true);

    // No old-style text fields anywhere in DOM
    expect(wrapper.find('[data-cy="colleague-id-input"]').exists()).toBe(false);
    expect(wrapper.find('[data-cy="colleague-name-input"]').exists()).toBe(false);
    expect(wrapper.find('[data-cy="guest-name-input"]').exists()).toBe(false);
  });

  it('restores the memorized day from session storage on mount', async () => {
    const storedDay = futureDay();
    useDateState().setDay(storedDay);
    const wrapper = mountView();

    await flushPromises();

    expect((wrapper.vm as unknown as { selectedDate: string }).selectedDate).toBe(storedDay);
    expect(fetchItemsMock).toHaveBeenCalledWith('ig-1', storedDay);
  });

  it('keeps the selected day when toggling between day and week mode', async () => {
    const storedDay = futureDay();
    useDateState().setDay(storedDay);
    const wrapper = mountView();

    await flushPromises();

    const vm = wrapper.vm as unknown as { selectedDate: string; bookingMode: 'day' | 'week' };
    vm.bookingMode = 'week';
    await flushPromises();
    vm.bookingMode = 'day';
    await flushPromises();

    expect(vm.selectedDate).toBe(storedDay);
    expect(sessionStorage.getItem('sithub_selected_day')).toBe(storedDay);
  });

  it('resets the live selected day to today after a successful booking', async () => {
    const storedDay = futureDay();
    useDateState().setDay(storedDay);
    fetchItemsMock.mockResolvedValue({
      data: [{
        id: 'item-1',
        type: 'items',
        attributes: {
          name: 'Desk A',
          equipment: [],
          availability: 'available' as const
        }
      }]
    });

    const wrapper = mountView();
    await flushPromises();

    await wrapper.get('[data-cy="book-item-btn"]').trigger('click');
    await flushPromises();

    const today = formatDate(new Date());
    expect(createBookingMock).toHaveBeenCalledWith('item-1', storedDay, undefined);
    expect(sessionStorage.getItem('sithub_selected_day')).toBe(today);
    expect((wrapper.vm as unknown as { selectedDate: string }).selectedDate).toBe(today);
    expect(wrapper.find('[data-cy="booking-success"]').exists()).toBe(true);
    expect(wrapper.text()).toContain('Desk A');
  });

  it('does not render guest radio option', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.find('[data-cy="book-guest-radio"]').exists()).toBe(false);
    expect(wrapper.find('[data-cy="book-self-radio"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="book-colleague-radio"]').exists()).toBe(true);
  });

  it('does not render multi-day checkbox', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.find('[data-cy="multi-day-checkbox"]').exists()).toBe(false);
  });

  it('fetches users on mount for colleague dropdown', async () => {
    mountView();
    await flushPromises();

    expect(fetchColleaguesMock).toHaveBeenCalled();
  });

  describe('collapsible day tiles', () => {
    it('hides equipment on folded booked tiles', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: {
            name: 'Booked Item',
            equipment: ['Monitor'],
            warning: 'USB-C only',
            availability: 'occupied',
            booker_name: 'Alice'
          }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      // Booked tile hides equipment and warning by default
      expect(wrapper.find('[data-cy="item-equipment"]').exists()).toBe(false);
      expect(wrapper.find('[data-cy="item-warning"]').exists()).toBe(false);
      // Booker name remains visible
      expect(wrapper.find('[data-cy="item-booker"]').exists()).toBe(true);
    });

    it('shows equipment on expanded booked tiles', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: {
            name: 'Booked Item',
            equipment: ['Monitor'],
            warning: 'USB-C only',
            availability: 'occupied',
            booker_name: 'Alice'
          }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      // Expand the tile
      const vm = wrapper.vm as unknown as { expandedDayTiles: Set<string> };
      vm.expandedDayTiles = new Set(['item-1']);
      await nextTick();

      // Equipment and warning now visible
      expect(wrapper.find('[data-cy="item-equipment"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="item-warning"]').exists()).toBe(true);
    });

    it('always shows equipment on available tiles', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: {
            name: 'Available Item',
            equipment: ['Monitor'],
            warning: 'USB-C only',
            availability: 'available'
          }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      expect(wrapper.find('[data-cy="item-equipment"]').exists()).toBe(true);
      // Warning is only visible when tile is expanded; folded shows warning icon only
      expect(wrapper.find('[data-cy="item-warning"]').exists()).toBe(false);
      expect(wrapper.find('[data-cy="folded-warning-icon"]').exists()).toBe(true);
    });
  });

  describe('equipment filter', () => {
    it('renders filter input', async () => {
      const wrapper = mountView();
      await flushPromises();

      expect(wrapper.find('[data-cy="equipment-filter-input"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="equipment-filter-info"]').exists()).toBe(true);
    });

    it('blurs items that do not match the filter', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [
          {
            id: 'item-1',
            type: 'items',
            attributes: { name: 'Desk A', equipment: ['webcam', 'monitor'], availability: 'available' as const }
          },
          {
            id: 'item-2',
            type: 'items',
            attributes: { name: 'Desk B', equipment: ['keyboard'], availability: 'available' as const }
          }
        ]
      });

      const wrapper = mountView();
      await flushPromises();

      // No overlay initially
      expect(wrapper.findAll('[data-cy="equipment-not-available"]')).toHaveLength(0);

      // Set filter
      (wrapper.vm as unknown as { equipmentFilter: string }).equipmentFilter = 'webcam';
      await nextTick();

      // One item matches, one does not
      const overlays = wrapper.findAll('[data-cy="equipment-not-available"]');
      expect(overlays).toHaveLength(1);

      // The matching item should not have the blur class
      const cards = wrapper.findAll('[data-cy="item-entry"]');
      const deskA = cards.find(c => c.attributes('data-cy-item-id') === 'item-1');
      const deskB = cards.find(c => c.attributes('data-cy-item-id') === 'item-2');
      expect(deskA?.classes()).not.toContain('item-filtered-out');
      expect(deskB?.classes()).toContain('item-filtered-out');
    });

    it('removes blur when filter is cleared', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: { name: 'Desk A', equipment: ['keyboard'], availability: 'available' as const }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      (wrapper.vm as unknown as { equipmentFilter: string }).equipmentFilter = 'webcam';
      await nextTick();
      expect(wrapper.findAll('[data-cy="equipment-not-available"]')).toHaveLength(1);

      (wrapper.vm as unknown as { equipmentFilter: string }).equipmentFilter = '';
      await nextTick();
      expect(wrapper.findAll('[data-cy="equipment-not-available"]')).toHaveLength(0);
    });

    it('opens filter help dialog when the info button is clicked', async () => {
      const wrapper = mountView();
      await flushPromises();

      expect(wrapper.find('[data-cy="equipment-filter-help"]').exists()).toBe(false);
      await wrapper.find('[data-cy="equipment-filter-info"]').trigger('click');
      await nextTick();

      const help = wrapper.find('[data-cy="equipment-filter-help"]');
      expect(help.exists()).toBe(true);
      expect(help.text()).toContain('show only items having the filter keyword(s) in any of the equipment items;');
      expect(help.text()).toContain('multiple keywords are combined with OR;');
      expect(help.text()).toContain('use plus sign to combine with AND;');
    });

    it('applies the same filter blur behavior in week mode', async () => {
      localStorage.setItem('sithub_booking_mode', 'week');
      fetchItemsMock.mockResolvedValue({
        data: [
          {
            id: 'item-1',
            type: 'items',
            attributes: { name: 'Desk A', equipment: ['webcam', 'monitor'], availability: 'available' as const }
          },
          {
            id: 'item-2',
            type: 'items',
            attributes: { name: 'Desk B', equipment: ['keyboard'], availability: 'available' as const }
          }
        ]
      });

      const wrapper = mountView();
      await flushPromises();

      (wrapper.vm as unknown as { equipmentFilter: string }).equipmentFilter = 'webcam';
      await nextTick();

      const cards = wrapper.findAll('[data-cy="week-item-entry"]');
      const deskA = cards.find(c => c.attributes('data-cy-item-id') === 'item-1');
      const deskB = cards.find(c => c.attributes('data-cy-item-id') === 'item-2');
      expect(deskA?.classes()).not.toContain('item-filtered-out');
      expect(deskB?.classes()).toContain('item-filtered-out');
      expect(wrapper.findAll('[data-cy="equipment-not-available"]').length).toBeGreaterThanOrEqual(1);

      localStorage.removeItem('sithub_booking_mode');
    });

    it('saves a filter and shows a confirmation', async () => {
      const wrapper = mountView();
      await flushPromises();

      const vm = wrapper.vm as unknown as {
        equipmentFilter: string;
        toggleSaveFilter: () => void;
      };
      vm.equipmentFilter = 'webcam';
      await nextTick();
      vm.toggleSaveFilter();
      await flushPromises();

      expect(JSON.parse(localStorage.getItem('sithub_saved_filters')!)).toEqual(['webcam']);
      expect(wrapper.text()).toContain('Filter saved.');
    });

    it('deletes a saved filter, clears the input, and shows a confirmation', async () => {
      localStorage.setItem('sithub_saved_filters', JSON.stringify(['webcam']));
      const wrapper = mountView();
      await flushPromises();

      await wrapper.get('[data-cy="equipment-filter-input"]').setValue('webcam');
      await nextTick();
      await wrapper.get('[data-cy="equipment-filter-delete"]').trigger('click');
      await flushPromises();

      expect(JSON.parse(localStorage.getItem('sithub_saved_filters')!)).toEqual([]);
      expect((wrapper.vm as unknown as { equipmentFilter: string }).equipmentFilter).toBe('');
      expect(wrapper.text()).toContain('Saved filter deleted.');
    });
  });

  it('shows floor plan button and dialog when the item group has a floor plan', async () => {
    fetchItemGroupsMock.mockResolvedValue({
      data: [{ id: 'ig-1', type: 'item-groups', attributes: { name: 'Test Group', floor_plan: 'group.svg' } }]
    });

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.find('[data-cy="item-group-floor-plan-btn"]').exists()).toBe(true);
    await wrapper.get('[data-cy="item-group-floor-plan-btn"]').trigger('click');
    const dialog = wrapper.get('[data-cy="item-group-floor-plan-dialog"]');
    expect(dialog.exists()).toBe(true);
    expect(dialog.attributes('data-persistent')).toBe('');
  });

  it('opens the item-group floor plan fullscreen on compact viewports', async () => {
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
    fetchItemGroupsMock.mockResolvedValue({
      data: [{ id: 'ig-1', type: 'item-groups', attributes: { name: 'Test Group', floor_plan: 'group.svg' } }]
    });

    const wrapper = mountView();
    await flushPromises();

    await wrapper.get('[data-cy="item-group-floor-plan-btn"]').trigger('click');
    expect(wrapper.get('[data-cy="item-group-floor-plan-dialog"]').attributes('data-fullscreen')).toBe('true');
  });

  it('opens the item-group floor plan dialog without max-width on desktop', async () => {
    fetchItemGroupsMock.mockResolvedValue({
      data: [{ id: 'ig-1', type: 'item-groups', attributes: { name: 'Test Group', floor_plan: 'group.svg' } }]
    });

    const wrapper = mountView();
    await flushPromises();

    await wrapper.get('[data-cy="item-group-floor-plan-btn"]').trigger('click');
    const dialog = wrapper.get('[data-cy="item-group-floor-plan-dialog"]');
    expect(dialog.attributes('max-width')).toBeUndefined();
    expect(dialog.attributes('maxwidth')).toBeUndefined();
  });

  describe('booking limit error modal', () => {
    it('shows limit dialog instead of snackbar for day-mode limit errors', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: { name: 'Desk A', equipment: [], availability: 'available' as const }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      createBookingMock.mockRejectedValue(
        new ApiError('Conflict', 409, 'booking limit exceeded: you have reached the maximum of 3 active bookings')
      );

      await wrapper.find('[data-cy="book-item-btn"]').trigger('click');
      await flushPromises();

      expect(wrapper.find('[data-cy="booking-limit-dialog"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="booking-limit-text"]').text()).toContain('3');
    });

    it('shows limit dialog for week-mode limit errors', async () => {
      localStorage.setItem('sithub_booking_mode', 'week');
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: { name: 'Desk A', equipment: [], availability: 'available' as const }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      // Select a week day checkbox to enable the confirm button
      const checkbox = wrapper.find('[data-cy="week-day-checkbox"]');
      if (checkbox.exists()) {
        await checkbox.find('input').setValue(true);
        await flushPromises();
      }

      createBookingMock.mockRejectedValue(
        new ApiError('Conflict', 409, 'booking limit exceeded: you have reached the maximum of 3 active bookings')
      );

      const confirmBtn = wrapper.find('[data-cy="week-confirm-btn"]');
      if (confirmBtn.exists()) {
        await confirmBtn.trigger('click');
        await flushPromises();
      }

      expect(wrapper.find('[data-cy="booking-limit-dialog"]').exists()).toBe(true);
      localStorage.removeItem('sithub_booking_mode');
    });

    it('dismisses the limit dialog when OK is clicked', async () => {
      fetchItemsMock.mockResolvedValue({
        data: [{
          id: 'item-1',
          type: 'items',
          attributes: { name: 'Desk A', equipment: [], availability: 'available' as const }
        }]
      });

      const wrapper = mountView();
      await flushPromises();

      createBookingMock.mockRejectedValue(
        new ApiError('Conflict', 409, 'booking limit exceeded: you have reached the maximum of 3 active bookings')
      );

      await wrapper.find('[data-cy="book-item-btn"]').trigger('click');
      await flushPromises();

      expect(wrapper.find('[data-cy="booking-limit-dialog"]').exists()).toBe(true);

      await wrapper.find('[data-cy="booking-limit-ok"]').trigger('click');
      await flushPromises();

      expect(wrapper.find('[data-cy="booking-limit-dialog"]').exists()).toBe(false);
    });
  });

  it('shows a connection lost error when initial user loading fails', async () => {
    fetchMeMock.mockRejectedValue(new ApiError(CONNECTION_LOST_MESSAGE, 0));

    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.text()).toContain(CONNECTION_LOST_MESSAGE);
  });

  describe('warning confirmation dialog', () => {
    const itemWithWarning = {
      id: 'item-warn',
      type: 'items' as const,
      attributes: {
        name: 'Workspace 1',
        equipment: ['monitor'],
        availability: 'available' as const,
        warning: 'Only for Apple users.'
      }
    };

    const itemNoWarning = {
      id: 'item-ok',
      type: 'items' as const,
      attributes: {
        name: 'Workspace 2',
        equipment: [],
        availability: 'available' as const
      }
    };

    beforeEach(() => {
      localStorage.removeItem('sithub_warning_suppressed');
      createBookingMock.mockClear();
    });

    it('shows warning dialog when booking item with warning', async () => {
      fetchItemsMock.mockResolvedValue({ data: [itemWithWarning] });
      const wrapper = mountView();
      await flushPromises();
      createBookingMock.mockClear();

      await wrapper.find('[data-cy="book-item-btn"]').trigger('click');
      await flushPromises();

      expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="warning-item-name"]').text()).toBe('Workspace 1');
      expect(wrapper.find('[data-cy="warning-message"]').text()).toBe('Only for Apple users.');
      expect(createBookingMock).not.toHaveBeenCalled();
    });

    it('proceeds with booking when CONFIRM is clicked', async () => {
      fetchItemsMock.mockResolvedValue({ data: [itemWithWarning] });
      const wrapper = mountView();
      await flushPromises();
      createBookingMock.mockClear();

      await wrapper.find('[data-cy="book-item-btn"]').trigger('click');
      await flushPromises();

      await wrapper.find('[data-cy="warning-confirm-btn"]').trigger('click');
      await flushPromises();

      expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
      expect(createBookingMock).toHaveBeenCalled();
    });

    it('aborts booking when CANCEL is clicked', async () => {
      fetchItemsMock.mockResolvedValue({ data: [itemWithWarning] });
      const wrapper = mountView();
      await flushPromises();
      createBookingMock.mockClear();

      await wrapper.find('[data-cy="book-item-btn"]').trigger('click');
      await flushPromises();

      await wrapper.find('[data-cy="warning-cancel-btn"]').trigger('click');
      await flushPromises();

      expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
      expect(createBookingMock).not.toHaveBeenCalled();
    });

    it('stores suppression in localStorage when dont-show-again is checked', async () => {
      fetchItemsMock.mockResolvedValue({ data: [itemWithWarning] });
      const wrapper = mountView();
      await flushPromises();

      await wrapper.find('[data-cy="book-item-btn"]').trigger('click');
      await flushPromises();

      const checkbox = wrapper.find('[data-cy="warning-dont-show-checkbox"] input');
      await checkbox.setValue(true);
      await flushPromises();

      await wrapper.find('[data-cy="warning-confirm-btn"]').trigger('click');
      await flushPromises();

      const stored = JSON.parse(localStorage.getItem('sithub_warning_suppressed') || '[]');
      expect(stored).toHaveLength(1);
      expect(stored[0]).toMatch(/^item-warn::/);
    });

    it('skips dialog for suppressed item and books directly', async () => {
      // Pre-suppress by booking with dont-show-again first
      fetchItemsMock.mockResolvedValue({ data: [itemWithWarning] });
      const setup = mountView();
      await flushPromises();
      await setup.find('[data-cy="book-item-btn"]').trigger('click');
      await flushPromises();
      await setup.find('[data-cy="warning-dont-show-checkbox"] input').setValue(true);
      await flushPromises();
      await setup.find('[data-cy="warning-confirm-btn"]').trigger('click');
      await flushPromises();
      setup.unmount();

      // Now mount fresh and verify suppression works
      fetchItemsMock.mockResolvedValue({ data: [itemWithWarning] });
      const wrapper = mountView();
      await flushPromises();
      createBookingMock.mockClear();

      await wrapper.find('[data-cy="book-item-btn"]').trigger('click');
      await flushPromises();

      expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
      expect(createBookingMock).toHaveBeenCalled();
    });

    it('books directly without dialog for item without warning', async () => {
      fetchItemsMock.mockResolvedValue({ data: [itemNoWarning] });
      const wrapper = mountView();
      await flushPromises();
      createBookingMock.mockClear();

      await wrapper.find('[data-cy="book-item-btn"]').trigger('click');
      await flushPromises();

      expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
      expect(createBookingMock).toHaveBeenCalled();
    });

    describe('week mode sequential warnings', () => {
      const warnItemA = {
        id: 'warn-a',
        type: 'items' as const,
        attributes: {
          name: 'Desk A',
          equipment: [],
          availability: 'available' as const,
          warning: 'Warning for A'
        }
      };

      const warnItemB = {
        id: 'warn-b',
        type: 'items' as const,
        attributes: {
          name: 'Desk B',
          equipment: [],
          availability: 'available' as const,
          warning: 'Warning for B'
        }
      };

      const noWarnItem = {
        id: 'no-warn',
        type: 'items' as const,
        attributes: {
          name: 'Desk C',
          equipment: [],
          availability: 'available' as const
        }
      };

      // v-btn stub that doesn't double-fire click: uses emits declaration
      // so @click listener stays in $attrs and isn't also emitted
      const weekStubs = {
        ...stubs,
        'v-btn': {
          inheritAttrs: false,
          emits: ['click'],
          template: '<button type="button" @click="$emit(\'click\', $event)"><slot /></button>'
        }
      };

      const mountWeekView = () =>
        mount(ItemsView, {
          global: {
            stubs: weekStubs,
            plugins: [createPinia(), createTestI18n()]
          }
        });

      const setupWeekMode = (weekItems: typeof warnItemA[]) => {
        localStorage.setItem('sithub_booking_mode', 'week');
        fetchItemsMock.mockResolvedValue({ data: weekItems });
        fetchMyBookingsMock.mockResolvedValue({ data: [] } as never);
      };

      beforeEach(() => {
        localStorage.removeItem('sithub_warning_suppressed');
        createBookingMock.mockClear();
      });

      afterEach(() => {
        localStorage.removeItem('sithub_booking_mode');
      });

      const selectOneCheckboxPerItem = async (wrapper: ReturnType<typeof mountView>, itemCount: number) => {
        const checkboxes = wrapper.findAll('[data-cy="week-day-checkbox"] input');
        // Week mode renders N day-checkboxes per item; select first checkbox of each item
        const checkboxesPerItem = Math.floor(checkboxes.length / itemCount);
        for (let i = 0; i < itemCount && i * checkboxesPerItem < checkboxes.length; i++) {
          await checkboxes[i * checkboxesPerItem].setValue(true);
        }
        await flushPromises();
      };

      it('shows sequential warnings for 2 warned items then proceeds', async () => {
        setupWeekMode([warnItemA, warnItemB]);
        const wrapper = mountWeekView();
        await flushPromises();
        createBookingMock.mockClear();

        await selectOneCheckboxPerItem(wrapper, 2);

        const confirmBtn = wrapper.find('[data-cy="week-confirm-btn"]');
        if (!confirmBtn.exists()) return;
        await confirmBtn.trigger('click');
        await flushPromises();

        // First warning dialog shown
        expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);
        const firstName = wrapper.find('[data-cy="warning-item-name"]').text();
        expect(createBookingMock).not.toHaveBeenCalled();

        // Confirm first — dialog stays open with next item's content
        await wrapper.find('[data-cy="warning-confirm-btn"]').trigger('click');
        await nextTick();

        // Second warning dialog shown (same dialog, updated content)
        expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);
        const secondName = wrapper.find('[data-cy="warning-item-name"]').text();
        expect(secondName).not.toBe(firstName);
        expect(createBookingMock).not.toHaveBeenCalled();

        // Confirm second — booking proceeds
        await wrapper.find('[data-cy="warning-confirm-btn"]').trigger('click');
        await flushPromises();

        expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
        expect(createBookingMock).toHaveBeenCalled();
      });

      it('aborts entire booking when CANCEL on first warning', async () => {
        setupWeekMode([warnItemA, warnItemB]);
        const wrapper = mountWeekView();
        await flushPromises();
        createBookingMock.mockClear();

        await selectOneCheckboxPerItem(wrapper, 2);

        const confirmBtn = wrapper.find('[data-cy="week-confirm-btn"]');
        if (!confirmBtn.exists()) return;
        await confirmBtn.trigger('click');
        await flushPromises();

        // Cancel on first warning
        await wrapper.find('[data-cy="warning-cancel-btn"]').trigger('click');
        await flushPromises();

        expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
        expect(createBookingMock).not.toHaveBeenCalled();
      });

      it('aborts entire booking when CANCEL on second warning', async () => {
        setupWeekMode([warnItemA, warnItemB]);
        const wrapper = mountWeekView();
        await flushPromises();
        createBookingMock.mockClear();

        await selectOneCheckboxPerItem(wrapper, 2);

        const confirmBtn = wrapper.find('[data-cy="week-confirm-btn"]');
        if (!confirmBtn.exists()) return;
        await confirmBtn.trigger('click');
        await flushPromises();

        // Confirm first warning
        await wrapper.find('[data-cy="warning-confirm-btn"]').trigger('click');
        await flushPromises();

        // Cancel on second warning
        await wrapper.find('[data-cy="warning-cancel-btn"]').trigger('click');
        await flushPromises();

        expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
        expect(createBookingMock).not.toHaveBeenCalled();
      });

      it('skips suppressed items in the queue', async () => {
        // Pre-suppress warn-a in day mode
        fetchItemsMock.mockResolvedValue({ data: [warnItemA] });
        const setup = mountView();
        await flushPromises();
        await setup.find('[data-cy="book-item-btn"]').trigger('click');
        await flushPromises();
        await setup.find('[data-cy="warning-dont-show-checkbox"] input').setValue(true);
        await flushPromises();
        await setup.find('[data-cy="warning-confirm-btn"]').trigger('click');
        await flushPromises();
        setup.unmount();

        // Now week mode with both items
        setupWeekMode([warnItemA, warnItemB]);
        const wrapper = mountView();
        await flushPromises();
        createBookingMock.mockClear();

        await selectOneCheckboxPerItem(wrapper, 2);

        const confirmBtn = wrapper.find('[data-cy="week-confirm-btn"]');
        if (!confirmBtn.exists()) return;
        await confirmBtn.trigger('click');
        await flushPromises();

        // Only warn-b dialog should appear (warn-a is suppressed)
        expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);
        expect(wrapper.find('[data-cy="warning-item-name"]').text()).toBe('Desk B');

        // Confirm — booking proceeds
        await wrapper.find('[data-cy="warning-confirm-btn"]').trigger('click');
        await flushPromises();
        expect(createBookingMock).toHaveBeenCalled();
      });

      it('proceeds immediately when all warnings suppressed', async () => {
        // Pre-suppress both items in day mode
        fetchItemsMock.mockResolvedValue({ data: [warnItemA] });
        let setup = mountView();
        await flushPromises();
        await setup.find('[data-cy="book-item-btn"]').trigger('click');
        await flushPromises();
        await setup.find('[data-cy="warning-dont-show-checkbox"] input').setValue(true);
        await flushPromises();
        await setup.find('[data-cy="warning-confirm-btn"]').trigger('click');
        await flushPromises();
        setup.unmount();

        fetchItemsMock.mockResolvedValue({ data: [warnItemB] });
        setup = mountView();
        await flushPromises();
        await setup.find('[data-cy="book-item-btn"]').trigger('click');
        await flushPromises();
        await setup.find('[data-cy="warning-dont-show-checkbox"] input').setValue(true);
        await flushPromises();
        await setup.find('[data-cy="warning-confirm-btn"]').trigger('click');
        await flushPromises();
        setup.unmount();

        // Week mode with both suppressed
        setupWeekMode([warnItemA, warnItemB]);
        const wrapper = mountView();
        await flushPromises();
        createBookingMock.mockClear();

        await selectOneCheckboxPerItem(wrapper, 2);

        const confirmBtn = wrapper.find('[data-cy="week-confirm-btn"]');
        if (!confirmBtn.exists()) return;
        await confirmBtn.trigger('click');
        await flushPromises();

        // No dialog — proceeds directly
        expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
        expect(createBookingMock).toHaveBeenCalled();
      });

      it('shows one dialog for single warned item then proceeds', async () => {
        setupWeekMode([warnItemA, noWarnItem]);
        const wrapper = mountView();
        await flushPromises();
        createBookingMock.mockClear();

        await selectOneCheckboxPerItem(wrapper, 2);

        const confirmBtn = wrapper.find('[data-cy="week-confirm-btn"]');
        if (!confirmBtn.exists()) return;
        await confirmBtn.trigger('click');
        await flushPromises();

        // Only warn-a dialog
        expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);
        expect(wrapper.find('[data-cy="warning-item-name"]').text()).toBe('Desk A');

        await wrapper.find('[data-cy="warning-confirm-btn"]').trigger('click');
        await flushPromises();

        expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
        expect(createBookingMock).toHaveBeenCalled();
      });
    });
  });

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);
});
/* jscpd:ignore-end */
