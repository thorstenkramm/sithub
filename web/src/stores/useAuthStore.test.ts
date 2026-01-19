import { createPinia, setActivePinia } from 'pinia';
import { useAuthStore } from './useAuthStore';

describe('useAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('defaults to empty userName and isAdmin false', () => {
    const store = useAuthStore();
    expect(store.userName).toBe('');
    expect(store.isAdmin).toBe(false);
  });

  it('allows updating userName', () => {
    const store = useAuthStore();
    store.userName = 'Ada Lovelace';
    expect(store.userName).toBe('Ada Lovelace');
  });

  it('allows updating isAdmin', () => {
    const store = useAuthStore();
    store.isAdmin = true;
    expect(store.isAdmin).toBe(true);
  });
});
