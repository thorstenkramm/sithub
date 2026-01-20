import { mount } from '@vue/test-utils';
import EmptyState from '../EmptyState.vue';

describe('EmptyState', () => {
  it('renders title', () => {
    const wrapper = mount(EmptyState, {
      props: { title: 'No items found' },
      global: {
        stubs: {
          'v-icon': true,
          'v-btn': true
        }
      }
    });
    expect(wrapper.text()).toContain('No items found');
  });

  it('renders message when provided', () => {
    const wrapper = mount(EmptyState, {
      props: { title: 'Empty', message: 'Try adding some items' },
      global: {
        stubs: {
          'v-icon': true,
          'v-btn': true
        }
      }
    });
    expect(wrapper.text()).toContain('Try adding some items');
  });

  it('renders action button when actionText provided', () => {
    const wrapper = mount(EmptyState, {
      props: { title: 'Empty', actionText: 'Add Item' },
      global: {
        stubs: {
          'v-icon': true,
          'v-btn': {
            template: '<button data-cy="empty-state-action"><slot /></button>',
            props: ['to', 'color', 'variant']
          }
        }
      }
    });
    expect(wrapper.find('[data-cy="empty-state-action"]').exists()).toBe(true);
    expect(wrapper.text()).toContain('Add Item');
  });

  it('emits action event when button clicked without actionTo', async () => {
    const wrapper = mount(EmptyState, {
      props: { title: 'Empty', actionText: 'Add Item' },
      global: {
        stubs: {
          'v-icon': true,
          'v-btn': {
            template: '<button @click="$emit(\'click\')"><slot /></button>',
            props: ['to', 'color', 'variant']
          }
        }
      }
    });
    await wrapper.find('button').trigger('click');
    expect(wrapper.emitted('action')).toBeTruthy();
  });
});
