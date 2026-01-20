import { createApp } from 'vue';
import { createPinia } from 'pinia';

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

// Mock vuetify styles import
vi.mock('vuetify/styles', () => ({}));

// Mock the vuetify plugin module
vi.mock('./plugins/vuetify', () => ({
  vuetify: { name: 'vuetify' }
}));

// Mock global CSS
vi.mock('./styles/global.css', () => ({}));

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

    expect(useMock).toHaveBeenCalledWith({ name: 'pinia' });
    expect(useMock).toHaveBeenCalledWith({ name: 'router' });
    expect(useMock).toHaveBeenCalledWith({ name: 'vuetify' });
    expect(mountMock).toHaveBeenCalledWith('#app');
  });
});
