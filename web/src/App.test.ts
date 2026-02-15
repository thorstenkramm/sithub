import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { createVuetify } from 'vuetify';
import App from './App.vue';
import router from './router';

describe('App', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('mounts with router and vuetify', async () => {
    const vuetify = createVuetify();
    const slotStub = { template: '<div><slot /></div>' };
    const wrapper = mount(App, {
      global: {
        plugins: [router, createPinia(), vuetify],
        stubs: {
          'v-app': slotStub,
          'v-app-bar': slotStub,
          'v-app-bar-nav-icon': slotStub,
          'v-main': slotStub,
          'v-navigation-drawer': slotStub,
          'v-menu': slotStub,
          'v-list': slotStub,
          'v-list-item': slotStub,
          'v-list-item-title': slotStub,
          'v-list-item-subtitle': slotStub,
          'v-divider': slotStub,
          'v-spacer': slotStub,
          'v-btn': slotStub,
          'v-btn-toggle': slotStub,
          'v-checkbox': slotStub,
          'v-icon': slotStub,
          'v-avatar': slotStub,
          'v-chip': slotStub,
          'v-card': slotStub,
          'v-card-title': slotStub,
          'v-card-text': slotStub,
          'v-card-actions': slotStub,
          'v-text-field': slotStub,
          'v-alert': slotStub,
          'v-dialog': slotStub,
          'router-view': true,
          'router-link': true
        }
      }
    });

    await router.isReady();

    expect(wrapper.find('router-view-stub').exists()).toBe(true);
  });
});
