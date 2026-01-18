import { mount } from '@vue/test-utils';
import App from './App.vue';
import router from './router';

describe('App', () => {
  it('mounts with router and vuetify', async () => {
    const slotStub = { template: '<div><slot /></div>' };
    const wrapper = mount(App, {
      global: {
        plugins: [router],
        stubs: {
          'v-app': slotStub,
          'v-main': slotStub,
          'router-view': true
        }
      }
    });

    await router.isReady();

    expect(wrapper.find('router-view-stub').exists()).toBe(true);
  });
});
