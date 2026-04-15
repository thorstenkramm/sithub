import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import MatrixCancelPopover from './MatrixCancelPopover.vue';
import { cancelBooking } from '../../api/bookings';
import { createTestI18n } from '../../__tests__/helpers/i18n';
import { popoverStubs } from './testHelpers';

vi.mock('../../api/bookings', () => ({
  cancelBooking: vi.fn()
}));

const cancelBookingMock = cancelBooking as unknown as ReturnType<typeof vi.fn>;

const defaultItem = {
  item_id: 'desk-1',
  item_name: 'Desk 1',
  equipment: [],
  cells: [],
  reserved: false
};

const defaultCell = {
  date: '2099-05-01',
  availability: 'occupied' as const,
  booker_name: 'Ada Lovelace',
  booker_user_id: 'user-1',
  booked_by_me: true,
  booking_id: 'b-123'
};

const stubs = popoverStubs;

function mountPopover(overrides: {
  modelValue?: boolean;
  item?: typeof defaultItem;
  cell?: typeof defaultCell;
} = {}) {
  return mount(MatrixCancelPopover, {
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

describe('MatrixCancelPopover', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    cancelBookingMock.mockReset();
  });

  it('renders cancel card with person, desk, and date', async () => {
    const wrapper = mountPopover();
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-cancel-card"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="matrix-cancel-person"]').text()).toContain('Ada Lovelace');
    expect(wrapper.find('[data-cy="matrix-cancel-desk"]').text()).toContain('Desk 1');
    expect(wrapper.find('[data-cy="matrix-cancel-date"]').text()).toContain('2099-05-01');
  });

  it('does not render when modelValue is false', async () => {
    const wrapper = mountPopover({ modelValue: false });
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-cancel-card"]').exists()).toBe(false);
  });

  it('calls cancelBooking and emits cancelled on success', async () => {
    cancelBookingMock.mockResolvedValue(undefined);
    const wrapper = mountPopover();
    await flushPromises();

    await wrapper.find('[data-cy="matrix-cancel-confirm"]').trigger('click');
    await flushPromises();

    expect(cancelBookingMock).toHaveBeenCalledWith('b-123');
    expect(wrapper.emitted('cancelled')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([false]);
  });

  it('shows error on cancel failure', async () => {
    cancelBookingMock.mockRejectedValue(new Error('fail'));
    const wrapper = mountPopover();
    await flushPromises();

    await wrapper.find('[data-cy="matrix-cancel-confirm"]').trigger('click');
    await flushPromises();

    expect(wrapper.find('[data-cy="matrix-cancel-error"]').exists()).toBe(true);
  });

  it('closes without cancelling when close is clicked', async () => {
    const wrapper = mountPopover();
    await flushPromises();

    await wrapper.find('[data-cy="matrix-cancel-close"]').trigger('click');
    await flushPromises();

    expect(cancelBookingMock).not.toHaveBeenCalled();
  });
});
