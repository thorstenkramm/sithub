import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import App from './App.vue';
import router from './router';

describe('App', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('mounts with router and vuetify', async () => {
    const slotStub = { template: '<div><slot /></div>' };
    const wrapper = mount(App, {
      global: {
        plugins: [router, createPinia()],
        stubs: {
          'v-app': slotStub,
          'v-app-bar': slotStub,
          'v-main': slotStub,
          'v-navigation-drawer': slotStub,
          'v-menu': slotStub,
          'v-list': slotStub,
          'v-list-item': slotStub,
          'router-view': true,
          'router-link': true
        }
      }
    });

    await router.isReady();

    expect(wrapper.find('router-view-stub').exists()).toBe(true);
  });
});
