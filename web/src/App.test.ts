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
  const routerLinkStub = {
    props: ['to'],
    template: '<a><slot /></a>',
  };
  const decodeSvgAsset = (src: string) => {
    if (src.startsWith('data:image/svg+xml;base64,')) {
      return Buffer.from(src.split(',')[1] ?? '', 'base64').toString('utf8');
    }
    return src;
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
          'router-link': routerLinkStub,
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

  it('renders the compact SitHub logo in the authenticated header', () => {
    const authStore = useAuthStore(pinia);
    authStore.setUser({
      id: 'user-1',
      display_name: 'Test User',
      email: 'test@example.com',
      is_admin: false,
      auth_source: 'internal',
    });

    const wrapper = mountApp();

    const logo = wrapper.get('img.logo-image');
    expect(decodeSvgAsset(logo.attributes('src') ?? '')).toContain('viewBox="0 0 320 80"');
    expect(logo.attributes('alt')).toBe('SitHub');
  });

  describe('consolidated profile menu', () => {
    beforeEach(() => {
      const authStore = useAuthStore(pinia);
      authStore.setUser({
        id: 'user-1',
        display_name: 'Test User',
        email: 'test@example.com',
        is_admin: true,
        auth_source: 'internal',
      });
    });

    it('desktop menu contains all settings and profile items', async () => {
      const wrapper = mountApp();
      await nextTick();

      // Profile actions present
      expect(wrapper.find('[data-cy="floor-plan-editor-btn"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="avatar-btn"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="change-password-btn"]').exists()).toBe(true);

      // Settings controls present
      expect(wrapper.find('[data-cy="theme-selector"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="language-selector"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="show-weekends-toggle"]').exists()).toBe(true);

      // Logout present
      expect(wrapper.find('[data-cy="logout-btn"]').exists()).toBe(true);
    });

    it('mobile menu contains all settings and profile items', async () => {
      const wrapper = mountApp();
      await nextTick();

      // Profile actions present
      expect(wrapper.find('[data-cy="mobile-floor-plan-editor-btn"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="mobile-avatar-btn"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="mobile-change-password-btn"]').exists()).toBe(true);

      // Settings controls present
      expect(wrapper.find('[data-cy="mobile-theme-selector"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="mobile-language-selector"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="mobile-show-weekends-toggle"]').exists()).toBe(true);

      // Logout present
      expect(wrapper.find('[data-cy="mobile-logout-btn"]').exists()).toBe(true);
    });
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
