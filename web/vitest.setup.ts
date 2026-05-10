import { beforeEach } from "vitest";
import { createPinia, setActivePinia } from "pinia";

class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}

globalThis.ResizeObserver = ResizeObserver;

// Ensure every test starts with a fresh Pinia instance. Components in the
// app use Pinia stores (auth, live feed) and would otherwise crash on first
// store access in tests that don't manually wire it up.
beforeEach(() => {
  setActivePinia(createPinia());
});

function installLocalStorageShim() {
  try {
    const existing = globalThis.localStorage;
    if (
      existing &&
      typeof existing.getItem === "function" &&
      typeof existing.setItem === "function" &&
      typeof existing.removeItem === "function" &&
      typeof existing.clear === "function"
    ) {
      return;
    }
  } catch {
    // Fall through to install the shim.
  }

  const store = new Map<string, string>();
  const localStorageShim = {
    get length() {
      return store.size;
    },
    clear() {
      store.clear();
    },
    getItem(key: string) {
      return store.has(key) ? store.get(key)! : null;
    },
    key(index: number) {
      return Array.from(store.keys())[index] ?? null;
    },
    removeItem(key: string) {
      store.delete(key);
    },
    setItem(key: string, value: string) {
      store.set(key, String(value));
    },
  } satisfies Storage;

  Object.defineProperty(globalThis, "localStorage", {
    configurable: true,
    value: localStorageShim,
  });

  if (typeof window !== "undefined") {
    Object.defineProperty(window, "localStorage", {
      configurable: true,
      value: localStorageShim,
    });
  }
}

installLocalStorageShim();
