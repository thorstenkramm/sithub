import { mount } from '@vue/test-utils';
import StatusChip from '../StatusChip.vue';

describe('StatusChip', () => {
  it('renders available status with success color', () => {
    const wrapper = mount(StatusChip, {
      props: { status: 'available' },
      global: {
        stubs: {
          'v-chip': {
            template: '<span :class="color"><slot /></span>',
            props: ['color', 'size', 'variant', 'prependIcon']
          }
        }
      }
    });
    expect(wrapper.text()).toContain('Available');
    expect(wrapper.find('span').classes()).toContain('success');
  });

  it('renders booked status with warning color', () => {
    const wrapper = mount(StatusChip, {
      props: { status: 'booked' },
      global: {
        stubs: {
          'v-chip': {
            template: '<span :class="color"><slot /></span>',
            props: ['color', 'size', 'variant', 'prependIcon']
          }
        }
      }
    });
    expect(wrapper.text()).toContain('Booked');
    expect(wrapper.find('span').classes()).toContain('warning');
  });

  it('renders mine status with primary color', () => {
    const wrapper = mount(StatusChip, {
      props: { status: 'mine' },
      global: {
        stubs: {
          'v-chip': {
            template: '<span :class="color"><slot /></span>',
            props: ['color', 'size', 'variant', 'prependIcon']
          }
        }
      }
    });
    expect(wrapper.text()).toContain('My Booking');
    expect(wrapper.find('span').classes()).toContain('primary');
  });

  it('renders custom label when provided', () => {
    const wrapper = mount(StatusChip, {
      props: { status: 'available', label: 'Free' },
      global: {
        stubs: {
          'v-chip': {
            template: '<span><slot /></span>',
            props: ['color', 'size', 'variant', 'prependIcon']
          }
        }
      }
    });
    expect(wrapper.text()).toContain('Free');
    expect(wrapper.text()).not.toContain('Available');
  });

  it('renders guest status with warning color', () => {
    const wrapper = mount(StatusChip, {
      props: { status: 'guest' },
      global: {
        stubs: {
          'v-chip': {
            template: '<span :class="color"><slot /></span>',
            props: ['color', 'size', 'variant', 'prependIcon']
          }
        }
      }
    });
    expect(wrapper.text()).toContain('Guest');
    expect(wrapper.find('span').classes()).toContain('warning');
  });
});
