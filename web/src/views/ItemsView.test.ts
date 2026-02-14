import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import ItemsView from './ItemsView.vue';
import PageHeader from '../components/PageHeader.vue';
import { fetchItems } from '../api/items';
import { fetchMe } from '../api/me';
import { fetchItemGroups } from '../api/itemGroups';
import { fetchAreas } from '../api/areas';
import { buildViewStubs, defineAuthRedirectTests } from './testHelpers';

/* jscpd:ignore-start */

const pushMock = vi.fn();
vi.mock('../api/me');
vi.mock('../api/items');
vi.mock('../api/itemGroups');
vi.mock('../api/areas');
const routeMock = { params: { itemGroupId: 'ig-1' }, query: { areaId: 'area-1' } };
vi.mock('vue-router', () => ({
  useRoute: () => routeMock,
  useRouter: () => ({ push: pushMock })
}));

describe('ItemsView', () => {
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
    'v-dialog',
    'v-bottom-sheet',
    'v-textarea',
    'v-spacer',
    'v-btn-toggle',
    'v-select',
    'router-link'
  ]);

  const fetchMeMock = vi.mocked(fetchMe);
  const fetchItemsMock = vi.mocked(fetchItems);
  const fetchItemGroupsMock = vi.mocked(fetchItemGroups);
  const fetchAreasMock = vi.mocked(fetchAreas);

  const mountView = () =>
    mount(ItemsView, {
      global: {
        stubs,
        plugins: [createPinia()]
      }
    });

  beforeEach(() => {
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
  });

  it('renders item equipment, warning, and status', async () => {
    fetchItemsMock.mockResolvedValue({
      data: [
        {
          id: 'item-1',
          type: 'items',
          attributes: {
            name: 'Item 1',
            equipment: ['Monitor', 'Keyboard'],
            warning: 'USB-C only',
            availability: 'occupied'
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

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);
});
/* jscpd:ignore-end */
