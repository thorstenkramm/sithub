import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import MatrixBookingPopover from './MatrixBookingPopover.vue';
import { createBooking, cancelBooking, fetchMyBookings } from '../../api/bookings';
import { fetchColleagues } from '../../api/users';
import { ApiError } from '../../api/client';
import { createTestI18n } from '../../__tests__/helpers/i18n';
import { popoverStubs } from './testHelpers';

vi.mock('../../api/bookings', () => ({
  createBooking: vi.fn(),
  cancelBooking: vi.fn(),
  fetchMyBookings: vi.fn()
}));

vi.mock('../../api/users', () => ({
  fetchColleagues: vi.fn()
}));

const createBookingMock = createBooking as unknown as ReturnType<typeof vi.fn>;
const cancelBookingMock = cancelBooking as unknown as ReturnType<typeof vi.fn>;
const fetchMyBookingsMock = fetchMyBookings as unknown as ReturnType<typeof vi.fn>;
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
  },
  ConfirmDialog: {
    template: '<div v-if="modelValue" data-cy="confirm-dialog">'
      + '<span data-cy="confirm-dialog-message">{{ message }}</span>'
      + '<button data-cy="confirm-dialog-confirm" @click="$emit(\'confirm\')">confirm</button>'
      + '<button data-cy="confirm-dialog-cancel" @click="$emit(\'update:modelValue\', false); $emit(\'cancel\')">cancel</button></div>',
    props: ['modelValue', 'title', 'message', 'confirmText', 'confirmColor', 'loading']
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
      cell: overrides.cell ?? { ...defaultCell },
      areaId: 'area-1'
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
    cancelBookingMock.mockReset();
    cancelBookingMock.mockResolvedValue(undefined);
    fetchMyBookingsMock.mockReset();
    fetchMyBookingsMock.mockResolvedValue({ data: [] });
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

  it('prompts to swap an existing same area/day booking, then books on confirm (story 36.9)', async () => {
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    fetchMyBookingsMock.mockResolvedValue({
      data: [{
        id: 'existing-1',
        type: 'bookings',
        attributes: {
          item_id: 'other-desk',
          item_name: 'Desk Z',
          item_group_id: 'ig',
          item_group_name: 'IG',
          area_id: 'area-1',
          area_name: 'Area',
          booking_date: defaultCell.date,
          created_at: '',
          booked_by_user_id: 'me',
          booked_by_user_name: 'Me',
          booked_for_me: true,
          note: ''
        }
      }]
    });
    const wrapper = mountPopover();
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();

    expect(wrapper.find('[data-cy="confirm-dialog"]').exists()).toBe(true);
    expect(createBookingMock).not.toHaveBeenCalled();

    await wrapper.find('[data-cy="confirm-dialog-confirm"]').trigger('click');
    await flushPromises();

    // Create-then-cancel: the new booking is created BEFORE the old one is
    // cancelled (story 36.9 D2).
    expect(createBookingMock).toHaveBeenCalled();
    expect(cancelBookingMock).toHaveBeenCalledWith('existing-1');
    const createOrder = createBookingMock.mock.invocationCallOrder[0]!;
    const cancelOrder = cancelBookingMock.mock.invocationCallOrder[0]!;
    expect(createOrder).toBeLessThan(cancelOrder);
  });

  it('keeps the new booking and warns when the post-create cancel fails (story 36.9)', async () => {
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    cancelBookingMock.mockRejectedValueOnce(new Error('cancel failed'));
    fetchMyBookingsMock.mockResolvedValue({
      data: [{
        id: 'existing-1',
        type: 'bookings',
        attributes: {
          item_id: 'other-desk',
          item_name: 'Desk Z',
          item_group_id: 'ig',
          item_group_name: 'IG',
          area_id: 'area-1',
          area_name: 'Area',
          booking_date: defaultCell.date,
          created_at: '',
          booked_by_user_id: 'me',
          booked_by_user_name: 'Me',
          booked_for_me: true,
          note: ''
        }
      }]
    });
    const wrapper = mountPopover();
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();
    await wrapper.find('[data-cy="confirm-dialog-confirm"]').trigger('click');
    await flushPromises();

    expect(createBookingMock).toHaveBeenCalled();
    expect(cancelBookingMock).toHaveBeenCalledWith('existing-1');
    // 'booked' still emitted since the new booking succeeded.
    expect(wrapper.emitted('booked')).toBeTruthy();
  });

  it('ignores an on-behalf booking in the same area/day (self-scoped guard, story 36.9)', async () => {
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    // for_user_name set: made for a colleague; must not be swapped.
    fetchMyBookingsMock.mockResolvedValue({
      data: [{
        id: 'existing-1',
        type: 'bookings',
        attributes: {
          item_id: 'other-desk',
          item_name: 'Desk Z',
          item_group_id: 'ig',
          item_group_name: 'IG',
          area_id: 'area-1',
          area_name: 'Area',
          booking_date: defaultCell.date,
          created_at: '',
          booked_by_user_id: 'me',
          booked_by_user_name: 'Me',
          booked_for_me: false,
          for_user_name: 'Colleague',
          note: ''
        }
      }]
    });
    const wrapper = mountPopover();
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();

    // No swap prompt; books directly without cancelling the colleague's seat.
    expect(wrapper.find('[data-cy="confirm-dialog"]').exists()).toBe(false);
    expect(cancelBookingMock).not.toHaveBeenCalled();
    expect(createBookingMock).toHaveBeenCalled();
  });

  it('does not prompt for a colleague booking when only the user has a conflict (story 36.9)', async () => {
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    // The user has an own conflicting booking, but the NEW booking is for a
    // colleague — it never occupies the user's own seat, so no guard prompt.
    fetchMyBookingsMock.mockResolvedValue({
      data: [{
        id: 'existing-1',
        type: 'bookings',
        attributes: {
          item_id: 'other-desk',
          item_name: 'Desk Z',
          item_group_id: 'ig',
          item_group_name: 'IG',
          area_id: 'area-1',
          area_name: 'Area',
          booking_date: defaultCell.date,
          created_at: '',
          note: ''
        }
      }]
    });
    const wrapper = mountPopover();
    await flushPromises();
    const vm = wrapper.vm as unknown as { bookingType: string; selectedColleagueId: string | null };
    vm.bookingType = 'colleague';
    vm.selectedColleagueId = 'u-1';
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();

    expect(wrapper.find('[data-cy="confirm-dialog-confirm"]').exists()).toBe(false);
    expect(cancelBookingMock).not.toHaveBeenCalled();
    expect(createBookingMock).toHaveBeenCalled();
  });

  it('prompts when the colleague already has a booking in the same area/day (story 36.9)', async () => {
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    cancelBookingMock.mockResolvedValue(undefined);
    // The user already booked Desk Z for colleague u-1 on the same area/day.
    fetchMyBookingsMock.mockResolvedValue({
      data: [{
        id: 'existing-1',
        type: 'bookings',
        attributes: {
          item_id: 'other-desk',
          item_name: 'Desk Z',
          item_group_id: 'ig',
          item_group_name: 'IG',
          area_id: 'area-1',
          area_name: 'Area',
          booking_date: defaultCell.date,
          created_at: '',
          booked_by_user_id: 'me',
          for_user_id: 'u-1',
          for_user_name: 'Jane Doe',
          note: ''
        }
      }]
    });
    const wrapper = mountPopover();
    await flushPromises();
    const vm = wrapper.vm as unknown as { bookingType: string; selectedColleagueId: string | null };
    vm.bookingType = 'colleague';
    vm.selectedColleagueId = 'u-1';
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();

    // The colleague-variant swap prompt names the colleague; nothing created yet.
    const dialog = wrapper.find('[data-cy="confirm-dialog-confirm"]');
    expect(dialog.exists()).toBe(true);
    expect(wrapper.text()).toContain('Jane Doe');
    expect(createBookingMock).not.toHaveBeenCalled();

    // Confirming swaps: create the new booking first, then cancel the old one.
    await dialog.trigger('click');
    await flushPromises();
    expect(createBookingMock).toHaveBeenCalled();
    expect(cancelBookingMock).toHaveBeenCalledWith('existing-1');
  });

  it('does not book when the swap prompt is cancelled (story 36.9)', async () => {
    fetchMyBookingsMock.mockResolvedValue({
      data: [{
        id: 'existing-1',
        type: 'bookings',
        attributes: {
          item_id: 'other-desk',
          item_name: 'Desk Z',
          item_group_id: 'ig',
          item_group_name: 'IG',
          area_id: 'area-1',
          area_name: 'Area',
          booking_date: defaultCell.date,
          created_at: '',
          booked_by_user_id: 'me',
          booked_by_user_name: 'Me',
          booked_for_me: true,
          note: ''
        }
      }]
    });
    const wrapper = mountPopover();
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();
    await wrapper.find('[data-cy="confirm-dialog-cancel"]').trigger('click');
    await flushPromises();

    expect(cancelBookingMock).not.toHaveBeenCalled();
    expect(createBookingMock).not.toHaveBeenCalled();
  });

  it('keeps the warning dialog open when the popover menu closes (controlled by CANCEL/CONFIRM only)', async () => {
    // The v-menu closes on any interaction outside its content — including
    // ticking "Don't show again" inside the confirmation dialog. Closing the
    // menu must NOT abort the flow; the dialog stays until CANCEL/CONFIRM.
    localStorage.clear();
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    const itemWithWarning = { ...defaultItem, warning: 'Near window' };
    const wrapper = mountPopover({ item: itemWithWarning });
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();
    expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);

    // The menu closes (e.g. the checkbox click propagates as an outside click).
    await wrapper.setProps({ modelValue: false });
    await flushPromises();

    // The confirmation dialog persists; nothing was booked or aborted.
    expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);
    expect(createBookingMock).not.toHaveBeenCalled();

    // Confirming still completes the booking.
    await wrapper.find('[data-cy="warning-confirm-btn"]').trigger('click');
    await flushPromises();
    expect(createBookingMock).toHaveBeenCalled();
  });

  it('aborts the booking when the warning confirmation is cancelled', async () => {
    createBookingMock.mockResolvedValue({ data: { id: 'b-1', type: 'bookings', attributes: {} } });
    const itemWithWarning = { ...defaultItem, warning: 'Near window' };
    const wrapper = mountPopover({ item: itemWithWarning });
    await flushPromises();

    await wrapper.find('[data-cy="matrix-booking-confirm"]').trigger('click');
    await flushPromises();
    expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(true);

    // Cancelling the warning must abort: no booking is created and the dialog closes.
    await wrapper.find('[data-cy="warning-cancel-btn"]').trigger('click');
    await flushPromises();
    expect(createBookingMock).not.toHaveBeenCalled();
    expect(wrapper.find('[data-cy="warning-dialog"]').exists()).toBe(false);
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
