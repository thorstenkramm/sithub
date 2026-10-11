import { mount, flushPromises } from '@vue/test-utils';
import BookingCard from './BookingCard.vue';
import { updateBookingNote } from '../api/bookings';
import { ApiError } from '../api/client';
import { createTestI18n } from '../__tests__/helpers/i18n';
import type { MyBookingAttributes } from '../api/bookings';

vi.mock('../api/bookings', () => ({ updateBookingNote: vi.fn() }));

const baseAttributes = (): MyBookingAttributes => ({
  item_id: 'item-1',
  item_name: 'Corner Desk',
  item_group_id: 'ig-1',
  item_group_name: 'Room 101',
  area_id: 'area-1',
  area_name: 'Main Office',
  booking_date: '2026-01-20',
  created_at: '2026-01-19T10:00:00Z',
  booked_by_user_id: '',
  booked_by_user_name: '',
  booked_for_me: false,
  note: ''
});

const mountCard = (attrs: Partial<MyBookingAttributes>) =>
  mount(BookingCard, {
    props: {
      booking: { id: '1', type: 'bookings', attributes: { ...baseAttributes(), ...attrs } }
    },
    global: {
      plugins: [createTestI18n()],
      stubs: {
        'v-card': { template: '<div><slot /></div>' },
        'v-card-item': { template: '<div><slot name="prepend" /><slot /></div>' },
        'v-card-title': { template: '<div><slot /></div>' },
        'v-card-subtitle': { template: '<div><slot /></div>' },
        'v-card-text': { template: '<div><slot /></div>' },
        'v-card-actions': { template: '<div><slot /></div>' },
        'v-avatar': { template: '<div><slot /></div>' },
        'v-icon': { template: '<i><slot /></i>' },
        'v-btn': { template: '<button><slot /></button>' },
        'v-spacer': { template: '<div />' },
        'v-dialog': { template: '<div><slot /></div>' },
        'v-bottom-sheet': { template: '<div><slot /></div>' },
        'v-textarea': { template: '<textarea />' },
        'v-snackbar': { props: ['modelValue'], template: '<div v-if="modelValue" v-bind="$attrs"><slot /></div>' },
        StatusChip: { props: ['status'], template: '<span :data-cy-status="status">{{ status }}</span>' }
      }
    }
  });

describe('BookingCard on-behalf hint', () => {
  it('shows "On behalf of <name>" for an on-behalf booking made by me', () => {
    const wrapper = mountCard({
      booked_by_user_id: 'colleague-1',
      booked_for_me: false,
      for_user_name: 'John Smith'
    });

    const hint = wrapper.find('[data-cy="on-behalf-of"]');
    expect(hint.exists()).toBe(true);
    expect(hint.text()).toContain('On behalf of John Smith');
  });

  it('shows no on-behalf hint for a self-booking', () => {
    const wrapper = mountCard({});

    expect(wrapper.find('[data-cy="on-behalf-of"]').exists()).toBe(false);
    expect(wrapper.find('[data-cy="on-behalf-chip"]').exists()).toBe(false);
  });

  it('shows no on-behalf-of caption for a booking made FOR me', () => {
    const wrapper = mountCard({
      booked_by_user_id: 'colleague-1',
      booked_by_user_name: 'Jane Doe',
      booked_for_me: true
    });

    expect(wrapper.find('[data-cy="on-behalf-of"]').exists()).toBe(false);
    expect(wrapper.find('[data-cy="booked-by"]').text()).toContain('Booked by Jane Doe');
  });
});

describe('BookingCard note saving', () => {
  it('should show the backend detail when saving the note fails', async () => {
    vi.mocked(updateBookingNote).mockRejectedValue(new ApiError('Request failed: 400', 400, 'note is too long'));
    const wrapper = mountCard({});

    await wrapper.get('[data-cy="note-save-btn"]').trigger('click');
    await flushPromises();

    const error = wrapper.get('[data-cy="note-save-error"]');
    expect(error.text()).toBe('Unable to save note. API error:\nnote is too long');
    expect(wrapper.emitted('note-updated')).toBeFalsy();
  });
});
