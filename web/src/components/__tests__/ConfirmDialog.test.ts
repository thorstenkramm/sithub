import { mount } from '@vue/test-utils';
import ConfirmDialog from '../ConfirmDialog.vue';

/* jscpd:ignore-start */
describe('ConfirmDialog', () => {
  const defaultProps = {
    modelValue: true,
    title: 'Confirm Action',
    message: 'Are you sure?'
  };

  it('renders title and message', () => {
    const wrapper = mount(ConfirmDialog, {
      props: defaultProps,
      global: {
        stubs: {
          'v-dialog': {
            template: '<div v-if="modelValue"><slot /></div>',
            props: ['modelValue', 'maxWidth']
          },
          'v-card': { template: '<div><slot /></div>' },
          'v-card-title': { template: '<div class="title"><slot /></div>' },
          'v-card-text': { template: '<div class="text"><slot /></div>' },
          'v-card-actions': { template: '<div class="actions"><slot /></div>' },
          'v-spacer': { template: '<span></span>' },
          'v-btn': {
            template: '<button :data-cy="$attrs[\'data-cy\']" @click="$emit(\'click\')"><slot /></button>',
            props: ['variant', 'color', 'loading', 'disabled']
          }
        }
      }
    });
    expect(wrapper.find('.title').text()).toContain('Confirm Action');
    expect(wrapper.find('.text').text()).toContain('Are you sure?');
  });

  it('emits confirm when confirm button clicked', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: defaultProps,
      global: {
        stubs: {
          'v-dialog': {
            template: '<div v-if="modelValue"><slot /></div>',
            props: ['modelValue', 'maxWidth']
          },
          'v-card': { template: '<div><slot /></div>' },
          'v-card-title': { template: '<div><slot /></div>' },
          'v-card-text': { template: '<div><slot /></div>' },
          'v-card-actions': { template: '<div><slot /></div>' },
          'v-spacer': { template: '<span></span>' },
          'v-btn': {
            template: '<button :data-cy="$attrs[\'data-cy\']" @click="$emit(\'click\')"><slot /></button>',
            props: ['variant', 'color', 'loading', 'disabled']
          }
        }
      }
    });
    await wrapper.find('[data-cy="confirm-dialog-confirm"]').trigger('click');
    expect(wrapper.emitted('confirm')).toBeTruthy();
  });

  it('emits cancel and closes when cancel button clicked', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: defaultProps,
      global: {
        stubs: {
          'v-dialog': {
            template: '<div v-if="modelValue"><slot /></div>',
            props: ['modelValue', 'maxWidth']
          },
          'v-card': { template: '<div><slot /></div>' },
          'v-card-title': { template: '<div><slot /></div>' },
          'v-card-text': { template: '<div><slot /></div>' },
          'v-card-actions': { template: '<div><slot /></div>' },
          'v-spacer': { template: '<span></span>' },
          'v-btn': {
            template: '<button :data-cy="$attrs[\'data-cy\']" @click="$emit(\'click\')"><slot /></button>',
            props: ['variant', 'color', 'loading', 'disabled']
          }
        }
      }
    });
    await wrapper.find('[data-cy="confirm-dialog-cancel"]').trigger('click');
    expect(wrapper.emitted('cancel')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('uses custom button text', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        ...defaultProps,
        confirmText: 'Delete',
        cancelText: 'Keep'
      },
      global: {
        stubs: {
          'v-dialog': {
            template: '<div v-if="modelValue"><slot /></div>',
            props: ['modelValue', 'maxWidth']
          },
          'v-card': { template: '<div><slot /></div>' },
          'v-card-title': { template: '<div><slot /></div>' },
          'v-card-text': { template: '<div><slot /></div>' },
          'v-card-actions': { template: '<div><slot /></div>' },
          'v-spacer': { template: '<span></span>' },
          'v-btn': {
            template: '<button><slot /></button>',
            props: ['variant', 'color', 'loading', 'disabled']
          }
        }
      }
    });
    expect(wrapper.text()).toContain('Delete');
    expect(wrapper.text()).toContain('Keep');
  });
});
/* jscpd:ignore-end */
