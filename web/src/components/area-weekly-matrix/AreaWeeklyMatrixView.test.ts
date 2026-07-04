import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import AreaWeeklyMatrixView from './AreaWeeklyMatrixView.vue';
import { fetchWeeklyMatrix } from '../../api/itemGroupMatrix';
import { createTestI18n } from '../../__tests__/helpers/i18n';
import { useAuthStore } from '../../stores/useAuthStore';
import { __resetLegacyPurgeForTests } from '../../composables/useFavorites';

vi.mock('../../api/itemGroupMatrix', () => ({ fetchWeeklyMatrix: vi.fn() }));
const liveFeed = vi.hoisted(() => ({
  handler: null as ((event: unknown) => void) | null
}));
vi.mock('../../stores/useLiveFeedStore', () => ({
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

const fetchMatrixMock = fetchWeeklyMatrix as unknown as ReturnType<typeof vi.fn>;

function makeMatrixResponse(opts: {
  groups?: Array<{
    id: string;
    name: string;
    items: Array<{
      id: string;
      name: string;
      equipment?: string[];
      warning?: string;
      reserved?: boolean;
      cells: Array<{
        date: string;
        availability: 'free' | 'occupied';
        booker_name?: string;
        booker_user_id?: string;
        booked_by_me?: boolean;
        booking_id?: string;
      }>;
    }>;
  }>;
  days?: Array<{ date: string; weekday: string }>;
} = {}) {
  const days = opts.days ?? [
    { date: '2099-04-13', weekday: 'MO' },
    { date: '2099-04-14', weekday: 'TU' },
    { date: '2099-04-15', weekday: 'WE' },
    { date: '2099-04-16', weekday: 'TH' },
    { date: '2099-04-17', weekday: 'FR' }
  ];

  const groups = opts.groups ?? [
    {
      id: 'ig-1',
      name: 'Room 101',
      items: [
        {
          id: 'desk-1',
          name: 'Desk 1',
          equipment: ['Dock', 'Monitor'],
          warning: 'Near window',
          cells: days.map(d => ({
            date: d.date,
            availability: 'free' as const,
            booked_by_me: false
          }))
        },
        {
          id: 'desk-2',
          name: 'Desk 2',
          equipment: [],
          cells: days.map(d => ({
            date: d.date,
            availability: 'free' as const,
            booked_by_me: false
          }))
        }
      ]
    },
    {
      id: 'ig-2',
      name: 'Room 102',
      items: [
        {
          id: 'desk-3',
          name: 'Desk 3',
          equipment: [],
          cells: days.map(d => ({
            date: d.date,
            availability: 'free' as const,
            booked_by_me: false
          }))
        }
      ]
    }
  ];

  return {
    data: groups.map(g => ({
      id: g.id,
      type: 'item-group-weekly-matrix',
      attributes: {
        item_group_id: g.id,
        item_group_name: g.name,
        days,
        items: g.items.map(item => ({
          item_id: item.id,
          item_name: item.name,
          equipment: item.equipment ?? [],
          warning: item.warning,
          reserved: item.reserved ?? false,
          cells: item.cells.map(c => ({
            date: c.date,
            availability: c.availability,
            booker_name: c.booker_name ?? '',
            booker_user_id: c.booker_user_id ?? '',
            booked_by_me: c.booked_by_me ?? false,
            booking_id: c.booking_id ?? ''
          }))
        }))
      }
    }))
  };
}

const stubs = {
  'v-alert': { template: '<div v-bind="$attrs"><slot /></div>' },
  'v-btn': {
    template: '<button v-bind="$attrs" @click="$emit(\'click\', $event)"><slot /></button>',
    emits: ['click']
  },
  'v-icon': { template: '<span v-bind="$attrs"><slot /></span>' },
  'v-tooltip': { template: '<div><slot name="activator" :props="{}" /><slot /></div>' },
  'v-avatar': { template: '<div v-bind="$attrs"><slot /></div>' },
  'v-img': { template: '<img v-bind="$attrs" />' },
  'v-skeleton-loader': { template: '<div v-bind="$attrs" />' },
  'v-snackbar': { template: '<div v-bind="$attrs"><slot /></div>' },
  'v-menu': { template: '<div v-bind="$attrs"><slot /></div>' },
  MatrixBookingPopover: { template: '<div data-cy="matrix-booking-popover-stub" />' },
  MatrixCancelPopover: { template: '<div data-cy="matrix-cancel-popover-stub" />' },
  LoadingState: {
    props: ['type', 'count'],
    template: '<div v-bind="$attrs" />'
  }
};

function makeOccupiedResponse(cellOverrides: Record<string, unknown> = {}) {
  return makeMatrixResponse({
    groups: [{
      id: 'ig-1',
      name: 'Room',
      items: [{
        id: 'desk-1',
        name: 'Desk',
        cells: [{
          date: '2099-04-13',
          availability: 'occupied' as const,
          booker_name: 'Someone',
          booker_user_id: 'other-user',
          booked_by_me: false,
          ...cellOverrides
        }]
      }]
    }],
    days: [{ date: '2099-04-13', weekday: 'MO' }]
  });
}

function mountWithAuth(userId: string, opts: { isAdmin?: boolean } = {}) {
  const pinia = createPinia();
  setActivePinia(pinia);
  const authStore = useAuthStore();
  authStore.userId = userId;
  if (opts.isAdmin) authStore.isAdmin = true;

  return {
    pinia,
    mount: () => mount(AreaWeeklyMatrixView, {
      props: { areaId: 'area-1', week: '2026-W16', showWeekends: false },
      global: { stubs, plugins: [pinia, createTestI18n()] }
    })
  };
}

function mountMatrix(props: Partial<{
  areaId: string;
  week: string;
  showWeekends: boolean;
  parsedEquipmentFilter: { exact: string[]; keywords: string[] }[];
}> = {}) {
  return mount(AreaWeeklyMatrixView, {
    props: {
      areaId: 'area-1',
      week: '2026-W16',
      showWeekends: false,
      ...props
    },
    global: {
      stubs,
      plugins: [createPinia(), createTestI18n()]
    }
  });
}

function mountMatrixWithStubs(
  customStubs: Record<string, unknown>,
  props: Partial<{
    areaId: string;
    week: string;
    showWeekends: boolean;
  }> = {}
) {
  return mount(AreaWeeklyMatrixView, {
    props: {
      areaId: 'area-1',
      week: '2026-W16',
      showWeekends: false,
      ...props
    },
    global: {
      stubs: {
        ...stubs,
        ...customStubs
      },
      plugins: [createPinia(), createTestI18n()]
    }
  });
}

describe('AreaWeeklyMatrixView', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    fetchMatrixMock.mockReset();
    localStorage.clear();
    __resetLegacyPurgeForTests();
    liveFeed.handler = null;
  });

  it('renders matrix container after data loads', async () => {
    fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
    const wrapper = mountMatrix();
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-container"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="matrix-loading"]').exists()).toBe(false);
  });

  it('reloads the matrix silently when a relevant live event arrives', async () => {
    vi.useFakeTimers();
    let resolveLiveRefresh: ((value: ReturnType<typeof makeMatrixResponse>) => void) | undefined;
    try {
      fetchMatrixMock
        .mockResolvedValueOnce(makeMatrixResponse())
        .mockImplementationOnce(() => new Promise(resolve => {
          resolveLiveRefresh = resolve;
        }));
      const wrapper = mountMatrix();
      await flushPromises();

      fetchMatrixMock.mockClear();
      expect(liveFeed.handler).toBeTypeOf('function');
      liveFeed.handler!({
        type: 'booking.created',
        booking_id: 'booking-1',
        item_id: 'desk-1',
        user_id: 'other-user',
        booking_date: '2099-04-13',
        timestamp: '2026-05-10T12:00:00Z'
      });

      await vi.advanceTimersByTimeAsync(300);

      expect(fetchMatrixMock).toHaveBeenCalledTimes(1);
      expect(fetchMatrixMock).toHaveBeenCalledWith('area-1', '2026-W16', 5);
      expect(wrapper.find('[data-cy="matrix-loading"]').exists()).toBe(false);

      resolveLiveRefresh?.(makeOccupiedResponse({ booker_name: 'Alice Smith' }));
      await flushPromises();

      expect(wrapper.find('[data-cy="matrix-loading"]').exists()).toBe(false);
      expect(wrapper.find('[data-cy="matrix-cell-occupied"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="matrix-cell-tooltip"]').text()).toBe('Alice Smith');
    } finally {
      vi.useRealTimers();
    }
  });

  it('shows error state on fetch failure', async () => {
    fetchMatrixMock.mockRejectedValue(new Error('fail'));
    const wrapper = mountMatrix();
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-error"]').exists()).toBe(true);
  });

  it('renders rooms and desks in configured order', async () => {
    fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
    const wrapper = mountMatrix();
    await flushPromises();

    // Room headers in order
    expect(wrapper.find('[data-cy="matrix-room-ig-1"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="matrix-room-ig-2"]').exists()).toBe(true);

    // Desk rows
    expect(wrapper.find('[data-cy="matrix-row-desk-1"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="matrix-row-desk-2"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="matrix-row-desk-3"]').exists()).toBe(true);
  });

  it('renders sticky header with weekday columns', async () => {
    fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
    const wrapper = mountMatrix();
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-header-row"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="matrix-day-MO"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="matrix-day-FR"]').exists()).toBe(true);
    expect(wrapper.find('.sticky-header').exists()).toBe(true);
  });

  it('has sticky left column class on desk labels', async () => {
    fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
    const wrapper = mountMatrix();
    await flushPromises();

    expect(wrapper.find('.matrix-desk-name.sticky-col').exists()).toBe(true);
  });

  it('defaults all rooms to expanded', async () => {
    fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
    const wrapper = mountMatrix();
    await flushPromises();

    // All desk rows visible (not collapsed)
    expect(wrapper.find('[data-cy="matrix-row-desk-1"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="matrix-row-desk-3"]').exists()).toBe(true);
  });

  it('collapses a room and shows summary counts', async () => {
    const resp = makeMatrixResponse({
      groups: [{
        id: 'ig-1',
        name: 'Room 101',
        items: [
          {
            id: 'desk-1',
            name: 'Desk 1',
            cells: [
              { date: '2026-04-13', availability: 'occupied', booker_name: 'Ada' },
              { date: '2026-04-14', availability: 'free' }
            ]
          },
          {
            id: 'desk-2',
            name: 'Desk 2',
            cells: [
              { date: '2026-04-13', availability: 'free' },
              { date: '2026-04-14', availability: 'occupied', booker_name: 'Bob' }
            ]
          }
        ]
      }],
      days: [
        { date: '2026-04-13', weekday: 'MO' },
        { date: '2026-04-14', weekday: 'TU' }
      ]
    });
    fetchMatrixMock.mockResolvedValue(resp);
    const wrapper = mountMatrix();
    await flushPromises();

    // Collapse ig-1
    await wrapper.find('[data-cy="matrix-room-toggle-ig-1"]').trigger('click');
    await flushPromises();

    // Desk rows should be hidden
    expect(wrapper.find('[data-cy="matrix-row-desk-1"]').exists()).toBe(false);

    // Summary counts visible
    const moSummary = wrapper.find('[data-cy="matrix-room-summary-ig-1-MO"]');
    expect(moSummary.exists()).toBe(true);
    expect(moSummary.text()).toBe('1/2');

    const tuSummary = wrapper.find('[data-cy="matrix-room-summary-ig-1-TU"]');
    expect(tuSummary.text()).toBe('1/2');
  });

  it('persists collapse state across remounts', async () => {
    fetchMatrixMock.mockResolvedValue(makeMatrixResponse());

    const wrapper1 = mountMatrix();
    await flushPromises();

    // Collapse ig-1
    await wrapper1.find('[data-cy="matrix-room-toggle-ig-1"]').trigger('click');
    wrapper1.unmount();

    // Remount
    const wrapper2 = mountMatrix();
    await flushPromises();

    // ig-1 desk rows should still be hidden
    expect(wrapper2.find('[data-cy="matrix-row-desk-1"]').exists()).toBe(false);
    // ig-2 should be expanded
    expect(wrapper2.find('[data-cy="matrix-row-desk-3"]').exists()).toBe(true);
  });

  it('renders free bookable cells', async () => {
    fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
    const wrapper = mountMatrix();
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-cell-free"]').exists()).toBe(true);
  });

  it('renders locked reserved free cells as non-clickable', async () => {
    const resp = makeMatrixResponse({
      groups: [{
        id: 'ig-1',
        name: 'Room',
        items: [{
          id: 'desk-1',
          name: 'VIP Desk',
          reserved: true,
          cells: [{ date: '2099-04-13', availability: 'free' }]
        }]
      }],
      days: [{ date: '2099-04-13', weekday: 'MO' }]
    });
    fetchMatrixMock.mockResolvedValue(resp);
    const wrapper = mountMatrix();
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-cell-locked"]').exists()).toBe(true);
  });

  it('renders occupied cells with initials and tooltip', async () => {
    const resp = makeMatrixResponse({
      groups: [{
        id: 'ig-1',
        name: 'Room',
        items: [{
          id: 'desk-1',
          name: 'Desk',
          cells: [{
            date: '2099-04-13',
            availability: 'occupied',
            booker_name: 'Ada Lovelace',
            booker_user_id: 'user-1',
            booked_by_me: false
          }]
        }]
      }],
      days: [{ date: '2099-04-13', weekday: 'MO' }]
    });
    fetchMatrixMock.mockResolvedValue(resp);
    const wrapper = mountMatrix();
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-cell-occupied"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="matrix-cell-initials"]').text()).toBe('A. Lovelace');
    expect(wrapper.find('[data-cy="matrix-cell-tooltip"]').text()).toBe('Ada Lovelace');
  });

  it('renders separate equipment icon with tooltip on desk label', async () => {
    fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
    const wrapper = mountMatrix();
    await flushPromises();

    // desk-1 has equipment
    expect(wrapper.find('[data-cy="matrix-equipment-icon"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="matrix-equipment-tooltip"]').text()).toBe('Dock, Monitor');
  });

  it('renders separate warning icon with tooltip showing warning text', async () => {
    fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
    const wrapper = mountMatrix();
    await flushPromises();

    // desk-1 has warning "Near window" — rendered via the shared ItemWarning component
    expect(wrapper.find('[data-cy="matrix-warning-icon"]').exists()).toBe(true);
    expect(wrapper.text()).toContain('Near window');
  });

  it('renders a heart for favorite matrix rows and removes the favorite on click', async () => {
    localStorage.setItem('sithub_favorite_items', JSON.stringify([{
      areaId: 'area-1',
      itemId: 'desk-1',
      itemName: 'Desk 1',
      itemGroupId: 'ig-1',
      itemGroupName: 'Room 101'
    }]));
    fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
    const wrapper = mountMatrix();
    await flushPromises();

    const heart = wrapper.get('[data-cy="matrix-favorite-heart-desk-1"]');
    expect(wrapper.find('[data-cy="matrix-favorite-heart-desk-2"]').exists()).toBe(false);

    await heart.trigger('click');
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-favorite-heart-desk-1"]').exists()).toBe(false);
    expect(localStorage.getItem('sithub_favorite_items')).toBe('[]');
  });

  it('does not render equipment or warning icons when absent', async () => {
    const resp = makeMatrixResponse({
      groups: [{
        id: 'ig-1',
        name: 'Room',
        items: [{
          id: 'desk-1',
          name: 'Desk',
          equipment: [],
          cells: [{ date: '2099-04-13', availability: 'free' }]
        }]
      }],
      days: [{ date: '2099-04-13', weekday: 'MO' }]
    });
    fetchMatrixMock.mockResolvedValue(resp);
    const wrapper = mountMatrix();
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-equipment-icon"]').exists()).toBe(false);
    expect(wrapper.find('[data-cy="matrix-warning-icon"]').exists()).toBe(false);
  });

  it('marks non-admin occupied cells as inert (no interaction)', async () => {
    fetchMatrixMock.mockResolvedValue(makeOccupiedResponse());
    const wrapper = mountMatrix();
    await flushPromises();

    const cell = wrapper.find('[data-cy="matrix-cell-occupied"]');
    expect(cell.classes()).toContain('cell-inert');
  });

  it('highlights booked-by-me cells', async () => {
    fetchMatrixMock.mockResolvedValue(
      makeOccupiedResponse({ booker_name: 'Me', booker_user_id: 'user-1', booked_by_me: true })
    );
    const { mount: mountAuth } = mountWithAuth('user-1');
    const wrapper = mountAuth();
    await flushPromises();

    expect(wrapper.find('.cell-booked-by-me').exists()).toBe(true);
  });

  it('marks own-booking occupied cells as interactive', async () => {
    fetchMatrixMock.mockResolvedValue(
      makeOccupiedResponse({ booker_name: 'Me', booker_user_id: 'user-1', booked_by_me: true, booking_id: 'b-123' })
    );
    const { mount: mountAuth } = mountWithAuth('user-1');
    const wrapper = mountAuth();
    await flushPromises();

    const cell = wrapper.find('[data-cy="matrix-cell-occupied"]');
    expect(cell.classes()).toContain('cell-interactive');
    expect(cell.classes()).not.toContain('cell-inert');
  });

  it('marks admin occupied cells as interactive for other users bookings', async () => {
    fetchMatrixMock.mockResolvedValue(
      makeOccupiedResponse({ booking_id: 'b-456' })
    );
    const { mount: mountAuth } = mountWithAuth('admin-1', { isAdmin: true });
    const wrapper = mountAuth();
    await flushPromises();

    const cell = wrapper.find('[data-cy="matrix-cell-occupied"]');
    expect(cell.classes()).toContain('cell-interactive');
  });

  it('reloads the matrix when booking conflict is reported', async () => {
    fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
    const requestAnimationFrameSpy = vi
      .spyOn(window, 'requestAnimationFrame')
      .mockImplementation((cb: FrameRequestCallback) => {
        cb(0);
        return 0;
      });

    const wrapper = mountMatrixWithStubs({
      MatrixBookingPopover: {
        template: '<button data-cy="matrix-booking-popover-stub" @click="$emit(\'bookingConflict\')" />',
        emits: ['bookingConflict']
      }
    });
    await flushPromises();

    await wrapper.find('[data-cy="matrix-cell-free"]').trigger('click');
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-popover-stub"]').trigger('click');
    await flushPromises();

    expect(fetchMatrixMock).toHaveBeenCalledTimes(2);
    requestAnimationFrameSpy.mockRestore();
  });

  it('refreshes silently after a booking so the table is not remounted', async () => {
    // A non-silent refresh swaps the table for the loading skeleton, recreating
    // the scroll container (scroll jumps to top) and every booker avatar
    // (images flicker). The post-booking refresh must keep the table mounted.
    let resolveRefresh: ((value: ReturnType<typeof makeMatrixResponse>) => void) | undefined;
    fetchMatrixMock
      .mockResolvedValueOnce(makeMatrixResponse())
      .mockImplementationOnce(() => new Promise(resolve => { resolveRefresh = resolve; }));
    const requestAnimationFrameSpy = vi
      .spyOn(window, 'requestAnimationFrame')
      .mockImplementation((cb: FrameRequestCallback) => {
        cb(0);
        return 0;
      });

    const wrapper = mountMatrixWithStubs({
      MatrixBookingPopover: {
        template: '<button data-cy="matrix-booking-popover-stub" @click="$emit(\'booked\')" />',
        emits: ['booked']
      }
    });
    await flushPromises();

    await wrapper.find('[data-cy="matrix-cell-free"]').trigger('click');
    await flushPromises();
    await wrapper.find('[data-cy="matrix-booking-popover-stub"]').trigger('click');
    await flushPromises();

    // Refresh is in-flight: the skeleton must NOT replace the mounted table.
    expect(wrapper.find('[data-cy="matrix-loading"]').exists()).toBe(false);
    expect(wrapper.find('[data-cy="matrix-container"]').exists()).toBe(true);

    resolveRefresh?.(makeMatrixResponse());
    await flushPromises();
    expect(wrapper.find('[data-cy="matrix-loading"]').exists()).toBe(false);
    expect(wrapper.find('[data-cy="matrix-container"]').exists()).toBe(true);
    requestAnimationFrameSpy.mockRestore();
  });

  describe('equipment filter', () => {
    it('does not mark any row as filtered when the parsed filter is empty', async () => {
      fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
      const wrapper = mountMatrix({ parsedEquipmentFilter: [] });
      await flushPromises();

      expect(wrapper.findAll('.matrix-row--filtered-out')).toHaveLength(0);
    });

    it('dims rows whose items do not match the parsed filter', async () => {
      fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
      const wrapper = mountMatrix({
        parsedEquipmentFilter: [{ exact: [], keywords: ['monitor'] }]
      });
      await flushPromises();

      // desk-1 has 'Dock', 'Monitor' → matches; desk-2 and desk-3 have [] → filtered out
      expect(wrapper.find('[data-cy="matrix-row-desk-1"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="matrix-row-desk-2"]').classes()).toContain('matrix-row--filtered-out');
      expect(wrapper.find('[data-cy="matrix-row-desk-3"]').classes()).toContain('matrix-row--filtered-out');
      expect(wrapper.findAll('[data-filtered-cy="matrix-row-filtered-out"]')).toHaveLength(2);
      expect(wrapper.findAll('.matrix-row--filtered-out')).toHaveLength(2);
    });

    it('marks all rows as filtered when no item matches', async () => {
      fetchMatrixMock.mockResolvedValue(makeMatrixResponse());
      const wrapper = mountMatrix({
        parsedEquipmentFilter: [{ exact: [], keywords: ['nonexistent'] }]
      });
      await flushPromises();

      expect(wrapper.findAll('.matrix-row--filtered-out')).toHaveLength(3);
    });
  });
});
