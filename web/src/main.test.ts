import { createApp } from 'vue';
import { createPinia } from 'pinia';
import { createVuetify } from 'vuetify';

const useMock = vi.fn().mockReturnThis();
const mountMock = vi.fn();

vi.mock('vue', () => ({
  createApp: vi.fn(() => ({
    use: useMock,
    mount: mountMock
  }))
}));

vi.mock('pinia', () => ({
  createPinia: vi.fn(() => ({ name: 'pinia' }))
}));

vi.mock('vuetify', () => ({
  createVuetify: vi.fn(() => ({ name: 'vuetify' }))
}));

vi.mock('./router', () => ({
  default: { name: 'router' }
}));

vi.mock('./App.vue', () => ({
  default: { name: 'App' }
}));

describe('main', () => {
  it('creates app, installs plugins, and mounts', async () => {
    await import('./main');

    expect(createApp).toHaveBeenCalled();
    expect(createPinia).toHaveBeenCalled();
    expect(createVuetify).toHaveBeenCalled();

    expect(useMock).toHaveBeenCalledWith({ name: 'pinia' });
    expect(useMock).toHaveBeenCalledWith({ name: 'router' });
    expect(useMock).toHaveBeenCalledWith({ name: 'vuetify' });
    expect(mountMock).toHaveBeenCalledWith('#app');
  });
});
