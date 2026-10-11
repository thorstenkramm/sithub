import { h } from 'vue';
import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { createVuetify } from 'vuetify';
import { VSnackbar, VBtn } from 'vuetify/components';
import DataExportView from './DataExportView.vue';
import { downloadBookingsCsv, fetchBookingExport, type BookingExportAttributes } from '../api/bookings';
import { saveBlob } from '../utils/download';
import { fetchAreas } from '../api/areas';
import { ApiError } from '../api/client';
import { buildViewStubs, createTestI18n } from './testHelpers';

const pushMock = vi.fn();

vi.mock('../api/bookings', () => ({ fetchBookingExport: vi.fn(), downloadBookingsCsv: vi.fn() }));
vi.mock('../utils/download', () => ({ saveBlob: vi.fn() }));
vi.mock('../api/areas', () => ({ fetchAreas: vi.fn() }));
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));

type Row = { id: string };
type SlotFn = (args: { item: Row }) => unknown;

// Renders one <tr> per item with the named item.* slots, so cell content can be asserted.
const dataTableStub = {
  props: ['headers', 'items', 'noDataText'],
  setup(
    props: { headers: Array<{ key: string; title: string }>; items: Row[]; noDataText: string },
    { slots }: { slots: Record<string, SlotFn | undefined> }
  ) {
    return () =>
      h('table', { 'data-cy': 'export-table' }, [
        h('thead', {}, props.headers.map((header) => h('th', { key: header.key }, header.title))),
        h('tbody', {}, props.items.length === 0
          ? [h('tr', {}, [h('td', { 'data-cy': 'export-empty' }, props.noDataText)])]
          : props.items.map((item) =>
            h('tr', { key: item.id, 'data-cy': 'export-row' },
              props.headers.map((header) => h('td', {}, slots[`item.${header.key}`]?.({ item }) as never)))))
      ]);
  }
};

const paginationStub = {
  props: ['modelValue', 'length'],
  emits: ['update:modelValue'],
  template: '<nav data-cy="export-pagination" :data-length="length"></nav>'
};

const datePickerStub = {
  props: ['modelValue', 'label', 'min', 'max', 'dataCy'],
  emits: ['update:modelValue'],
  template: '<input :data-cy="dataCy" :value="modelValue" :min="min" :max="max"' +
    ' @change="$emit(\'update:modelValue\', $event.target.value)" />'
};

const selectStub = {
  props: ['modelValue', 'items'],
  emits: ['update:modelValue'],
  template: '<select data-cy="export-area" :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)">' +
    '<option v-for="item in items" :key="item.value" :value="item.value">{{ item.title }}</option></select>'
};

function booking(id: string, overrides: Partial<BookingExportAttributes> = {}) {
  return {
    id,
    type: 'bookings',
    attributes: {
      booking_date: '2026-10-08',
      user_id: `user-${id}`,
      user_name: `Name ${id}`,
      user_email: `${id}@example.com`,
      is_guest: false,
      item_id: `item-${id}`,
      item_name: 'Tisch 3',
      item_group_id: 'cube-1',
      item_group_name: 'Cube 1',
      area_id: 'office',
      area_name: 'Office',
      booked_item: 'Cube 1 / Tisch 3',
      ...overrides
    }
  };
}

function page(data: ReturnType<typeof booking>[], total: number, offset = 0) {
  return { data, meta: { total, limit: 20, offset }, links: {} };
}

describe('DataExportView', () => {
  const fetchExportMock = fetchBookingExport as unknown as ReturnType<typeof vi.fn>;
  const fetchAreasMock = fetchAreas as unknown as ReturnType<typeof vi.fn>;
  const stubs = {
    ...buildViewStubs(['v-icon', 'v-breadcrumbs', 'v-breadcrumbs-item', 'v-breadcrumbs-divider']),
    'v-data-table-server': dataTableStub,
    'v-pagination': paginationStub,
    'v-select': selectStub,
    DatePickerField: datePickerStub,
    'v-snackbar': { props: ['modelValue'], template: '<div v-if="modelValue" v-bind="$attrs"><slot /></div>' }
  };

  const mountView = () =>
    mount(DataExportView, { global: { stubs, plugins: [createPinia(), createTestI18n()] } });

  beforeEach(() => {
    vi.useFakeTimers({ toFake: ['Date'] });
    vi.setSystemTime(new Date(2026, 9, 8, 12, 0, 0));
    setActivePinia(createPinia());
    pushMock.mockReset();
    fetchExportMock.mockReset();
    fetchAreasMock.mockReset();
    fetchExportMock.mockResolvedValue(page([booking('b1'), booking('b2')], 45));
    fetchAreasMock.mockResolvedValue({
      data: [
        { id: 'office', type: 'areas', attributes: { name: 'Office' } },
        { id: 'lab', type: 'areas', attributes: { name: 'Lab' } }
      ]
    });
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('should default to the last 5 days including today and load page 1', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(fetchExportMock).toHaveBeenCalledTimes(1);
    expect(fetchExportMock).toHaveBeenCalledWith({
      from: '2026-10-04', to: '2026-10-08', area: '', limit: 20, offset: 0
    });
    expect((wrapper.get('[data-cy="export-from"]').element as HTMLInputElement).value).toBe('2026-10-04');
    expect((wrapper.get('[data-cy="export-to"]').element as HTMLInputElement).value).toBe('2026-10-08');
  });

  it('should download applied filters, disable while pending and show success', async () => {
    let resolve!: (value: { blob: Blob; filename: string }) => void;
    vi.mocked(downloadBookingsCsv).mockReturnValue(new Promise((done) => { resolve = done; }));
    const wrapper = mountView();
    await flushPromises();
    await wrapper.get('[data-cy="export-area"]').setValue('office');
    await wrapper.get('[data-cy="data-export-csv-button"]').trigger('click');
    expect(downloadBookingsCsv).toHaveBeenCalledWith({ from: '2026-10-04', to: '2026-10-08', area: 'office' });
    expect(wrapper.get('[data-cy="data-export-csv-button"]').attributes()).toHaveProperty('disabled');
    const blob = new Blob(['csv']);
    resolve({ blob, filename: 'export.csv' });
    await flushPromises();
    expect(saveBlob).toHaveBeenCalledWith(blob, 'export.csv');
    expect(wrapper.get('[data-cy="data-export-csv-success"]').text()).toContain('CSV file downloaded');
  });

  it('should show the backend detail when CSV fails and disable invalid ranges', async () => {
    vi.mocked(downloadBookingsCsv).mockRejectedValue(new ApiError('failed', 500, 'database unavailable'));
    const wrapper = mountView();
    await flushPromises();
    await wrapper.get('[data-cy="data-export-csv-button"]').trigger('click');
    await flushPromises();
    expect(wrapper.get('[data-cy="data-export-csv-error"]').text()).toContain('database unavailable');
    await wrapper.get('[data-cy="export-from"]').setValue('2026-10-09');
    expect(wrapper.get('[data-cy="data-export-csv-button"]').attributes()).toHaveProperty('disabled');
  });

  it('should dismiss CSV errors with the real Vuetify snackbar close action', async () => {
    vi.stubGlobal('visualViewport', undefined);
    vi.mocked(downloadBookingsCsv).mockRejectedValue(new ApiError('failed', 500, 'database unavailable'));
    const wrapper = mount(DataExportView, {
      attachTo: document.body,
      global: {
        stubs: { ...stubs, 'v-snackbar': VSnackbar, 'v-btn': VBtn, 'router-link': true },
        plugins: [createPinia(), createTestI18n(), createVuetify({ components: { VSnackbar, VBtn } })]
      }
    });
    try {
      await flushPromises();
      await wrapper.get('[data-cy="data-export-csv-button"]').trigger('click');
      await flushPromises();
      const snackbar = wrapper.findAllComponents(VSnackbar).find((entry) =>
        entry.vm.$attrs['data-cy'] === 'data-export-csv-error')!;
      expect(snackbar.props('modelValue')).toBe(true);
      const close = document.querySelector<HTMLButtonElement>('[data-cy="data-export-csv-error-close"]')!;
      expect(close.textContent).toContain('Close');
      close.click();
      await flushPromises();
      expect(snackbar.props('modelValue')).toBe(false);
    } finally {
      wrapper.unmount();
      vi.unstubAllGlobals();
    }
  });

  it('should use local time for the default range just after midnight', async () => {
    vi.setSystemTime(new Date(2026, 2, 1, 0, 5, 0));
    mountView();
    await flushPromises();

    expect(fetchExportMock).toHaveBeenCalledWith(
      expect.objectContaining({ from: '2026-02-25', to: '2026-03-01' })
    );
  });

  it('should offer "All areas" followed by every area', async () => {
    const wrapper = mountView();
    await flushPromises();

    const options = wrapper.findAll('[data-cy="export-area"] option').map((o) => [o.text(), o.attributes('value')]);
    expect(options).toEqual([['All areas', ''], ['Office', 'office'], ['Lab', 'lab']]);
  });

  it('should render name, email, booked item and date per row', async () => {
    fetchExportMock.mockResolvedValue(page([
      booking('b1', { is_guest: true, user_name: 'Guest Person', user_email: '' }),
      booking('b2', { user_name: '', user_email: '' })
    ], 2));
    const wrapper = mountView();
    await flushPromises();

    const rows = wrapper.findAll('[data-cy="export-row"]');
    expect(rows).toHaveLength(2);
    expect(rows[0].get('[data-cy="export-name"]').text()).toContain('Guest Person');
    expect(rows[0].get('[data-cy="export-item"]').text()).toBe('Cube 1 / Tisch 3');
    expect(rows[0].get('[data-cy="export-date"]').text()).toContain('2026');
    expect(rows[1].get('[data-cy="export-name"]').text()).toBe('Unknown user');
  });

  it('should show email next to the name', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.get('[data-cy="export-email"]').text()).toBe('b1@example.com');
  });

  it('should render the four column headers', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.findAll('th').map((th) => th.text())).toEqual(['Name', 'Email', 'Booked item', 'Booking date']);
  });

  it('should compute the page count from the total', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.get('[data-cy="export-pagination"]').attributes('data-length')).toBe('3');
  });

  it('should fetch offset 40 when selecting page 3', async () => {
    const wrapper = mountView();
    await flushPromises();

    await wrapper.getComponent(paginationStub).vm.$emit('update:modelValue', 3);
    await flushPromises();

    expect(fetchExportMock).toHaveBeenLastCalledWith(expect.objectContaining({ limit: 20, offset: 40 }));
  });

  it.each([
    { name: 'the start date', cy: 'export-from', value: '2026-10-01', expected: { from: '2026-10-01' } },
    { name: 'the end date', cy: 'export-to', value: '2026-10-20', expected: { to: '2026-10-20' } },
    { name: 'the area', cy: 'export-area', value: 'lab', expected: { area: 'lab' } }
  ])('should reload from page 1 when changing $name', async ({ cy, value, expected }) => {
    const wrapper = mountView();
    await flushPromises();
    await wrapper.getComponent(paginationStub).vm.$emit('update:modelValue', 3);
    await flushPromises();
    fetchExportMock.mockClear();

    await wrapper.get(`[data-cy="${cy}"]`).setValue(value);
    await flushPromises();

    expect(fetchExportMock).toHaveBeenCalledTimes(1);
    expect(fetchExportMock).toHaveBeenCalledWith(expect.objectContaining({ ...expected, offset: 0 }));
  });

  it('should limit the date pickers to a consistent range', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.get('[data-cy="export-from"]').attributes('max')).toBe('2026-10-08');
    expect(wrapper.get('[data-cy="export-to"]').attributes('min')).toBe('2026-10-04');
  });

  it('should show the no-data text for an empty result', async () => {
    fetchExportMock.mockResolvedValue(page([], 0));
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.get('[data-cy="export-empty"]').text()).toBe('No bookings found for this range and area.');
  });

  it('should ignore a stale response that resolves after a newer one', async () => {
    let resolveFirst: (value: unknown) => void = () => {};
    fetchExportMock.mockReturnValueOnce(new Promise((resolve) => { resolveFirst = resolve; }));
    const wrapper = mountView();
    await flushPromises(); // area options are loaded, the first export request is still pending

    fetchExportMock.mockResolvedValueOnce(page([booking('fresh')], 1));
    await wrapper.get('[data-cy="export-area"]').setValue('lab');
    await flushPromises();
    resolveFirst(page([booking('stale')], 1));
    await flushPromises();

    expect(wrapper.get('[data-cy="export-email"]').text()).toBe('fresh@example.com');
  });

  it('should show the backend detail in an error snackbar', async () => {
    fetchExportMock.mockRejectedValue(new ApiError('Request failed: 400', 400, '"filter[from]" must not be after "filter[to]"'));
    const wrapper = mountView();
    await flushPromises();

    const snackbar = wrapper.get('[data-cy="export-error-snackbar"]');
    expect(snackbar.text()).toContain('Unable to load the bookings. API error:');
    expect(snackbar.text()).toContain('"filter[from]" must not be after "filter[to]"');
  });

  it('should redirect to access denied on 403', async () => {
    fetchExportMock.mockRejectedValue(new ApiError('Forbidden', 403));
    mountView();
    await flushPromises();

    expect(pushMock).toHaveBeenCalledWith('/access-denied');
  });

  it('should still render when the areas cannot be loaded', async () => {
    fetchAreasMock.mockRejectedValue(new ApiError('Request failed: 500', 500, 'boom'));
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.findAll('[data-cy="export-area"] option')).toHaveLength(1);
    expect(wrapper.get('[data-cy="export-error-snackbar"]').text()).toContain('boom');
  });
});
