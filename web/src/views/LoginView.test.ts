import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import LoginView from './LoginView.vue';
import { loginLocal } from '../api/auth';
import { ApiError, CONNECTION_LOST_MESSAGE } from '../api/client';

const pushMock = vi.fn();

vi.mock('../api/auth', () => ({ loginLocal: vi.fn() }));
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

  const mountView = () =>
    mount(LoginView, {
      global: {
        stubs,
        plugins: [createPinia()]
      }
    });

  beforeEach(() => {
    setActivePinia(createPinia());
    pushMock.mockReset();
    loginLocalMock.mockReset();
  });

  it('shows a connection lost error when login fails due to a network issue', async () => {
    loginLocalMock.mockRejectedValue(new ApiError(CONNECTION_LOST_MESSAGE, 0));
    const wrapper = mountView();

    await wrapper.get('[data-cy="login-email"]').setValue('ada@example.com');
    await wrapper.get('[data-cy="login-password"]').setValue('secret');
    await wrapper.get('[data-cy="login-form"]').trigger('submit');
    await flushPromises();

    expect(wrapper.text()).toContain(CONNECTION_LOST_MESSAGE);
  });
});
