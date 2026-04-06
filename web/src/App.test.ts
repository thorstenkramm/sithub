import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { nextTick } from 'vue';
import { createVuetify } from 'vuetify';
import App from './App.vue';
import router from './router';
import { createTestI18n } from './__tests__/helpers/i18n';
import { useAuthStore } from './stores/useAuthStore';

describe('App', () => {
  const slotStub = { template: '<div><slot /></div>' };
  const vImgStub = {
    emits: ['error'],
    template: '<img data-cy="stub-v-img" v-bind="$attrs" @error="$emit(\'error\')" />',
  };
  const dialogStub = {
    props: ['modelValue'],
    template: '<div v-if="modelValue"><slot /></div>',
  };
  const menuStub = {
    template: '<div><slot name="activator" :props="{}" /><slot /></div>',
  };
  let pinia: ReturnType<typeof createPinia>;

  function mountApp() {
    return mount(App, {
      global: {
        plugins: [router, pinia, createVuetify(), createTestI18n()],
        stubs: {
          'v-app': slotStub,
          'v-app-bar': slotStub,
          'v-app-bar-nav-icon': slotStub,
          'v-main': slotStub,
          'v-navigation-drawer': slotStub,
          'v-menu': menuStub,
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
          'v-snackbar': slotStub,
          'v-dialog': dialogStub,
          'v-img': vImgStub,
          'router-view': true,
          'router-link': true,
        },
      },
    });
  }

  beforeEach(async () => {
    pinia = createPinia();
    setActivePinia(pinia);
    await router.push('/');
    await router.isReady();
  });

  it('mounts with router and vuetify', () => {
    const wrapper = mountApp();
    expect(wrapper.find('router-view-stub').exists()).toBe(true);
  });

  it('renders the current user avatar and retries after a user switch', async () => {
    const authStore = useAuthStore(pinia);
    authStore.setUser({
      id: 'user-1',
      display_name: 'Alice Admin',
      email: 'alice@example.com',
      is_admin: true,
      auth_source: 'internal',
    });

    const wrapper = mountApp();
    await nextTick();

    const firstAvatar = wrapper.get('[data-cy="stub-v-img"]');
    expect(firstAvatar.attributes('src')).toMatch(/^\/api\/v1\/avatars\/user-1\?t=\d+$/);

    await firstAvatar.trigger('error');
    await nextTick();

    expect(wrapper.find('[data-cy="stub-v-img"]').exists()).toBe(false);
    expect(wrapper.text()).toContain('AA');

    authStore.setUser({
      id: 'user-2',
      display_name: 'Bob Builder',
      email: 'bob@example.com',
      is_admin: false,
      auth_source: 'entra',
    });
    await nextTick();

    const secondAvatar = wrapper.get('[data-cy="stub-v-img"]');
    expect(secondAvatar.attributes('src')).toMatch(/^\/api\/v1\/avatars\/user-2\?t=\d+$/);
  });
});
