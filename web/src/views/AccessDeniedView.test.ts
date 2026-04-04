import { mount } from '@vue/test-utils';
import AccessDeniedView from './AccessDeniedView.vue';
import { createTestI18n } from '../__tests__/helpers/i18n';

describe('AccessDeniedView', () => {
  it('renders access denied copy', () => {
    const slotStub = {
      template: '<div><slot /></div>'
    };

    const wrapper = mount(AccessDeniedView, {
      global: {
        plugins: [createTestI18n()],
        stubs: {
          'v-container': slotStub,
          'v-row': slotStub,
          'v-col': slotStub,
          'v-card': slotStub,
          'v-card-title': slotStub,
          'v-card-text': slotStub,
          'v-card-actions': slotStub,
          'v-btn': slotStub
        }
      }
    });

    expect(wrapper.text()).toContain('Access denied');
  });
});
