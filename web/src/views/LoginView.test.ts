import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import LoginView from './LoginView.vue';
import { loginLocal, fetchAuthProviders } from '../api/auth';
import { ApiError, CONNECTION_LOST_MESSAGE } from '../api/client';
import { createTestI18n } from '../__tests__/helpers/i18n';

const pushMock = vi.fn();

vi.mock('../api/auth', () => ({
  loginLocal: vi.fn(),
  fetchAuthProviders: vi.fn()
}));
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));

describe('LoginView', () => {
  const stubs = {
    'v-container': { template: '<div><slot /></div>' },
    'v-row': { template: '<div><slot /></div>' },
    'v-col': { template: '<div><slot /></div>' },
    'v-card': { template: '<div><slot /></div>' },
    'v-card-title': { template: '<div><slot /></div>' },
    'v-card-text': { template: '<div><slot /></div>' },
    'v-divider': { template: '<div />' },
    'v-alert': { template: '<div><slot /></div>' },
    'v-expand-transition': { template: '<div><slot /></div>' },
    'v-form': {
      template: '<form v-bind="$attrs" @submit.prevent="$emit(\'submit\', $event)"><slot /></form>'
    },
    'v-text-field': {
      props: ['modelValue'],
      template: '<input v-bind="$attrs" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />'
    },
    'v-btn': {
      template: '<button v-bind="$attrs" @click="$emit(\'click\', $event)"><slot /></button>'
    }
  };

  const loginLocalMock = loginLocal as unknown as ReturnType<typeof vi.fn>;
  const fetchAuthProvidersMock = fetchAuthProviders as unknown as ReturnType<typeof vi.fn>;

  const providersResponse = (entraid: boolean) => ({
    data: {
      type: 'auth-providers',
      id: 'current',
      attributes: { entraid, local: true }
    }
  });

  const decodeSvgAsset = (src: string) => {
    if (src.startsWith('data:image/svg+xml;base64,')) {
      return Buffer.from(src.split(',')[1] ?? '', 'base64').toString('utf8');
    }
    return src;
  };

  const mountView = () =>
    mount(LoginView, {
      global: {
        stubs,
        plugins: [createPinia(), createTestI18n()]
      }
    });

  beforeEach(() => {
    setActivePinia(createPinia());
    pushMock.mockReset();
    loginLocalMock.mockReset();
    fetchAuthProvidersMock.mockReset();
    // Default: Entra ID is configured; the local form is hidden behind the toggle.
    fetchAuthProvidersMock.mockResolvedValue(providersResponse(true));
  });

  it('shows a connection lost error when login fails due to a network issue', async () => {
    loginLocalMock.mockRejectedValue(new ApiError(CONNECTION_LOST_MESSAGE, 0));
    // Allow the local form to render (force-expand it via the toggle).
    const wrapper = mountView();
    await flushPromises();
    await wrapper.get('[data-cy="login-toggle-local"]').trigger('click');
    await flushPromises();

    await wrapper.get('[data-cy="login-email"]').setValue('ada@example.com');
    await wrapper.get('[data-cy="login-password"]').setValue('secret');
    await wrapper.get('[data-cy="login-form"]').trigger('submit');
    await flushPromises();

    expect(wrapper.text()).toContain(CONNECTION_LOST_MESSAGE);
  });

  it('disables the Entra ID button immediately after click', async () => {
    const requestAnimationFrameMock = vi
      .spyOn(window, 'requestAnimationFrame')
      .mockImplementation(() => 0);
    const wrapper = mountView();
    await flushPromises();

    await wrapper.get('[data-cy="login-entraid"]').trigger('click');

    expect(wrapper.get('[data-cy="login-entraid"]').attributes('disabled')).toBeDefined();

    requestAnimationFrameMock.mockRestore();
  });

  describe('auth providers gating', () => {
    it('renders only the Entra ID button and the toggle link when Entra ID is configured', async () => {
      fetchAuthProvidersMock.mockResolvedValue(providersResponse(true));
      const wrapper = mountView();
      await flushPromises();

      const loginLogo = wrapper.get('img.login-logo');
      expect(decodeSvgAsset(loginLogo.attributes('src') ?? '')).toContain('width="1200" height="700"');
      expect(wrapper.find('[data-cy="login-entraid"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="login-toggle-local"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="login-form"]').exists()).toBe(false);
    });

    it('toggles the local form open and closed via the more/less link', async () => {
      const wrapper = mountView();
      await flushPromises();

      const toggle = wrapper.get('[data-cy="login-toggle-local"]');
      await toggle.trigger('click');
      await flushPromises();
      expect(wrapper.find('[data-cy="login-form"]').exists()).toBe(true);

      await toggle.trigger('click');
      await flushPromises();
      expect(wrapper.find('[data-cy="login-form"]').exists()).toBe(false);
    });

    it('shows the local form by default when Entra ID is NOT configured and hides the Entra ID button', async () => {
      fetchAuthProvidersMock.mockResolvedValue(providersResponse(false));
      const wrapper = mountView();
      await flushPromises();

      expect(wrapper.find('[data-cy="login-entraid"]').exists()).toBe(false);
      expect(wrapper.find('[data-cy="login-toggle-local"]').exists()).toBe(false);
      expect(wrapper.find('[data-cy="login-form"]').exists()).toBe(true);
    });

    it('falls back to showing both options when the providers endpoint errors', async () => {
      fetchAuthProvidersMock.mockRejectedValue(new Error('network'));
      const wrapper = mountView();
      await flushPromises();

      expect(wrapper.find('[data-cy="login-entraid"]').exists()).toBe(true);
      expect(wrapper.find('[data-cy="login-form"]').exists()).toBe(true);
    });
  });
});
