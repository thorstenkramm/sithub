import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import ItemsView from './ItemsView.vue';
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
vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { itemGroupId: 'ig-1' } }),
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

  defineAuthRedirectTests(fetchMeMock, () => mountView(), pushMock);
});
/* jscpd:ignore-end */
