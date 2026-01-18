import { createPinia, setActivePinia } from 'pinia';
import { useAuthStore } from './useAuthStore';

describe('useAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('defaults to empty userName', () => {
    const store = useAuthStore();
    expect(store.userName).toBe('');
  });

  it('allows updating userName', () => {
    const store = useAuthStore();
    store.userName = 'Ada Lovelace';
    expect(store.userName).toBe('Ada Lovelace');
  });
});
