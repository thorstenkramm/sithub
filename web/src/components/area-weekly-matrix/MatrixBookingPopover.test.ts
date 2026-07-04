import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import MatrixBookingPopover from './MatrixBookingPopover.vue';
import { createBooking } from '../../api/bookings';
import { fetchColleagues } from '../../api/users';
import { ApiError } from '../../api/client';
import { createTestI18n } from '../../__tests__/helpers/i18n';
import { popoverStubs } from './testHelpers';

vi.mock('../../api/bookings', () => ({
  createBooking: vi.fn()
}));

vi.mock('../../api/users', () => ({
  fetchColleagues: vi.fn()
}));

const createBookingMock = createBooking as unknown as ReturnType<typeof vi.fn>;
const fetchColleaguesMock = fetchColleagues as unknown as ReturnType<typeof vi.fn>;

const defaultItem = {
  item_id: 'desk-1',
  item_name: 'Desk 1',
  equipment: [],
  cells: [],
  warning: undefined,
  reserved: false
};

const defaultCell = {
  date: '2099-05-01',
  availability: 'free' as const,
  booked_by_me: false
};

const stubs = {
  ...popoverStubs,
  'v-radio-group': {
    template: '<div v-bind="$attrs"><slot /></div>',
    props: ['modelValue'],
    emits: ['update:modelValue']
  },
  'v-radio': {
    template: '<label v-bind="$attrs" @click="$emit(\'click\')"><input type="radio" /></label>',
    props: ['label', 'value'],
    emits: ['click']
  },
  'v-autocomplete': {
    template: '<select v-bind="$attrs" data-cy="matrix-colleague-select"></select>',
    props: ['modelValue', 'items', 'itemTitle', 'itemValue', 'label', 'density', 'loading', 'clearable', 'hideDetails']
  },
  'v-text-field': {
    template: '<input v-bind="$attrs" :value="modelValue" data-cy="matrix-booking-note" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'label', 'density', 'hideDetails'],
    emits: ['update:modelValue']
  },
  'v-expand-transition': { template: '<div><slot /></div>' },
  WarningConfirmDialog: {
    template: '<div v-if="modelValue" data-cy="warning-dialog">'
      + '<span data-cy="warning-message">{{ message }}</span>'
      + '<button data-cy="warning-confirm-btn" @click="$emit(\'confirm\')">confirm</button>'
      + '<button data-cy="warning-cancel-btn" @click="$emit(\'cancel\')">cancel</button></div>',
    props: ['modelValue', 'itemName', 'message', 'dontShowAgain']
  }
};

function mountPopover(overrides: {
  modelValue?: boolean;
  item?: typeof defaultItem;
  cell?: typeof defaultCell;
} = {}) {
  return mount(MatrixBookingPopover, {
    props: {
      modelValue: overrides.modelValue ?? true,
      activatorEl: document.createElement('td'),
      item: overrides.item ?? { ...defaultItem },
      cell: overrides.cell ?? { ...defaultCell }
    },
    global: {
      stubs,
      plugins: [createPinia(), createTestI18n()]
    }
  });
}

describe('MatrixBookingPopover', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    createBookingMock.mockReset();
    fetchColleaguesMock.mockReset();
    fetchColleaguesMock.mockResolvedValue({ data: [] });
    localStorage.clear();
  });

  it('renders booking card when open', async () => {
    const wrapper = mountPopover();
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-booking-card"]').exists()).toBe(true);
  });

  it('does not render when modelValue is false', async () => {
    const wrapper = mountPopover({ modelValue: false });
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-booking-card"]').exists()).toBe(false);
  });

  it('defaults to self-booking radio', async () => {
    const wrapper = mountPopover();
    await flushPromises();

    const selfRadio = wrapper.find('[data-cy="matrix-book-self-radio"]');
    expect(selfRadio.exists()).toBe(true);
  });

  it('shows colleague picker only when colleague radio is selected', async () => {
    const wrapper = mountPopover();
    await flushPromises();

    // Default: no colleague picker
    expect(wrapper.find('[data-cy="matrix-colleague-select"]').exists()).toBe(false);

    // Switch to colleague by setting internal ref directly
    const vm = wrapper.vm as unknown as { bookingType: string };
    vm.bookingType = 'colleague';
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-colleague-select"]').exists()).toBe(true);
  });

  it('shows note field', async () => {
    const wrapper = mountPopover();
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-booking-note"]').exists()).toBe(true);
  });

  it('shows the uniform warning confirmation on booking a warned item, before booking', async () => {
    localStorage.clear();
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    const itemWithWarning = { ...defaultItem, warning: 'Near window' };
    const wrapper = mountPopover({ item: itemWithWarning });
    await flushPromises();

    // No inline warning; clicking confirm opens the shared confirmation dialog.
    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();
    expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="warning-message"]').text()).toContain('Near window');
    expect(createBookingMock).not.toHaveBeenCalled();

    // Confirming the warning proceeds with the booking.
    await wrapper.find('[data-cy="warning-confirm-btn"]').trigger('click');
    await flushPromises();
    expect(createBookingMock).toHaveBeenCalled();
  });

  it('aborts the warning flow when the popover closes while the dialog is open', async () => {
    localStorage.clear();
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    const itemWithWarning = { ...defaultItem, warning: 'Near window' };
    const wrapper = mountPopover({ item: itemWithWarning });
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();
    expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);

    // The menu closes (outside click / Escape) while the confirmation is open.
    await wrapper.setProps({ modelValue: false });
    await flushPromises();

    // No orphaned dialog remains and no booking was made for the dismissed popover.
    expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
    expect(createBookingMock).not.toHaveBeenCalled();
  });

  it('books directly without a warning dialog when the item has no warning', async () => {
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    const wrapper = mountPopover();
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();

    expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
    expect(createBookingMock).toHaveBeenCalled();
  });

  it('calls createBooking on confirm with note', async () => {
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    const wrapper = mountPopover();
    await flushPromises();

    // Type a note
    const noteInput = wrapper.find('[data-cy="matrix-booking-note"]');
    await noteInput.setValue('Arriving at noon');
    await flushPromises();

    // Click confirm
    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();

    expect(createBookingMock).toHaveBeenCalledWith(
      'desk-1', '2099-05-01', undefined, undefined, 'Arriving at noon'
    );
  });

  it('emits booked and closes on success', async () => {
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    const wrapper = mountPopover();
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();

    expect(wrapper.emitted('booked')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([false]);
  });

  it('shows inline error on 409 conflict and stays open', async () => {
    createBookingMock.mockRejectedValue(new ApiError('Conflict', 409));
    const wrapper = mountPopover();
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-booking-error"]').exists()).toBe(true);
    // Should NOT have emitted close
    const closeEvents = wrapper.emitted('update:modelValue') ?? [];
    const lastClose = closeEvents[closeEvents.length - 1];
    expect(lastClose).not.toEqual([false]);
    expect(wrapper.emitted('bookingConflict')).toBeTruthy();
  });

  it('saves last colleague to localStorage on colleague booking', async () => {
    fetchColleaguesMock.mockResolvedValue({
      data: [{ id: 'c-1', type: 'users', attributes: { display_name: 'Alice' } }]
    });
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });

    const wrapper = mountPopover();
    await flushPromises();

    // Switch to colleague mode and select colleague via internal state
    const vm = wrapper.vm as unknown as { selectedColleagueId: string; bookingType: string };
    vm.bookingType = 'colleague';
    await flushPromises();
    vm.selectedColleagueId = 'c-1';
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();

    expect(localStorage.getItem('sithub_matrix_last_colleague')).toBe('c-1');
  });
});
