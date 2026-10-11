import { h } from 'vue';
import { enableAutoUnmount, flushPromises, mount } from '@vue/test-utils';
import { createPinia } from 'pinia';
import ApiKeysView from './ApiKeysView.vue';
import { createApiKey, deleteApiKey, fetchApiKeys } from '../api/apiKeys';
import { ApiError } from '../api/client';
import { buildViewStubs, createTestI18n, defineAuthRedirectTests } from './testHelpers';

const pushMock = vi.fn();
enableAutoUnmount(afterEach);
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));
vi.mock('../api/apiKeys', () => ({ createApiKey: vi.fn(), deleteApiKey: vi.fn(), fetchApiKeys: vi.fn() }));

const key = '3f2b8c1e-9a4d-4c7e-b1f0-2d6e8a9c0b13';
const record = {
  id: 'record-id', type: 'api-keys', attributes: {
    name: 'Finance', key_prefix: key.slice(0, 10), scope: 'read-bookings' as const,
    expires_at: '2020-10-08T09:15:00Z', created_at: '2019-10-08T09:15:00Z', created_by: 'admin'
  }
};
const inputStub = {
  props: ['modelValue', 'readonly', 'disabled'], emits: ['update:modelValue'],
  template: '<div><input :value="modelValue" :readonly="readonly" :disabled="disabled" @input="$emit(\'update:modelValue\', $event.target.value)" /></div>'
};
const selectStub = {
  props: ['modelValue', 'items'], emits: ['update:modelValue'],
  template: '<select :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value)"><option v-for="item in items" :key="item.value" :value="item.value">{{ item.title }}</option></select>'
};
const stubs = {
  ...buildViewStubs(['v-card-actions', 'v-spacer', 'v-icon']),
  PageHeader: { props: ['title'], template: '<header>{{ title }}<slot name="actions" /></header>' },
  'v-btn': {
    props: ['disabled', 'loading'],
    template: '<button :disabled="disabled || loading"><slot /></button>'
  },
  'v-text-field': inputStub,
  'v-select': selectStub,
  'v-dialog': { props: ['modelValue'], template: '<div v-if="modelValue"><slot /></div>' },
  'v-snackbar': { props: ['modelValue'], template: '<div v-if="modelValue"><slot /><slot name="actions" /></div>' },
  'v-data-table': {
    props: ['items', 'headers'],
    setup(props: { items: { id: string }[]; headers: { key: string; title: string }[] },
      { slots }: { slots: Record<string, (args: unknown) => never> }) {
      return () => h('table', {}, [
        h('thead', {}, props.headers.map(header => h('th', {}, header.title))),
        h('tbody', {}, props.items.map(item => h('tr', { 'data-cy': 'api-key-row' },
          props.headers.map(header => h('td', {}, slots[`item.${header.key}`]?.({ item }))))))
      ]);
    }
  }
};

describe('ApiKeysView', () => {
  const mountView = () => mount(ApiKeysView, {
    attachTo: document.body,
    global: { stubs, plugins: [createPinia(), createTestI18n()] }
  });
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(fetchApiKeys).mockResolvedValue({ data: [record] });
    vi.mocked(createApiKey).mockResolvedValue({ data: { ...record, attributes: { ...record.attributes, key } } });
    vi.mocked(deleteApiKey).mockResolvedValue(undefined);
  });
  afterEach(() => {
    document.body.innerHTML = '';
    vi.restoreAllMocks();
  });

  async function openAndCreate(wrapper: ReturnType<typeof mountView>) {
    await flushPromises();
    await wrapper.get('[data-cy="api-keys-add"]').trigger('click');
    await wrapper.get('[data-cy="api-key-name-input"] input').setValue(' Finance ');
    await wrapper.get('[data-cy="api-key-submit"]').trigger('click');
    await flushPromises();
  }

  it('renders public metadata, localized scope and the expired chip', async () => {
    const wrapper = mountView();
    await flushPromises();
    expect(wrapper.get('[data-cy="api-key-name"]').text()).toBe('Finance');
    expect(wrapper.get('[data-cy="api-key-prefix"]').text()).toBe(`${key.slice(0, 10)}…`);
    expect(wrapper.get('[data-cy="api-key-expiry"]').text()).toContain('Expired');
    expect(wrapper.get('[data-cy="api-key-scope"]').text()).toBe('Read bookings');
    expect(wrapper.text()).not.toContain(key);
  });

  it('validates the name and defaults, shows the key once, then clears on Done', async () => {
    const wrapper = mountView();
    await flushPromises();
    await wrapper.get('[data-cy="api-keys-add"]').trigger('click');
    const submit = wrapper.get('[data-cy="api-key-submit"]');
    expect(submit.attributes('disabled')).toBeDefined();
    await wrapper.get('[data-cy="api-key-name-input"] input').setValue('  ');
    expect(submit.attributes('disabled')).toBeDefined();
    await wrapper.get('[data-cy="api-key-name-input"] input').setValue('x'.repeat(101));
    expect(submit.attributes('disabled')).toBeDefined();
    await wrapper.get('[data-cy="api-key-name-input"] input').setValue(' Finance ');
    await submit.trigger('click');
    await flushPromises();
    expect(createApiKey).toHaveBeenCalledWith({ name: 'Finance', scope: 'read-bookings', expires_in_months: 12 });
    expect((wrapper.get('[data-cy="api-key-value"] input').element as HTMLInputElement).value).toBe(key);
    expect(wrapper.find('[data-cy="api-key-once-warning"]').exists()).toBe(true);
    expect(wrapper.find('[data-cy="api-key-cancel"]').exists()).toBe(false);
    await wrapper.get('[data-cy="api-key-done"]').trigger('click');
    await flushPromises();
    expect(fetchApiKeys).toHaveBeenCalledTimes(2);
    expect(wrapper.find('[data-cy="api-key-value"]').exists()).toBe(false);
    await wrapper.get('[data-cy="api-keys-add"]').trigger('click');
    expect(wrapper.find('[data-cy="api-key-value"]').exists()).toBe(false);
    expect((wrapper.get('[data-cy="api-key-name-input"] input').element as HTMLInputElement).value).toBe('');
  });

  it.each(['success', 'missing', 'rejected'])('handles clipboard %s', async (mode) => {
    const writeText = vi.fn().mockImplementation(() => mode === 'rejected'
      ? Promise.reject(new Error('Denied')) : Promise.resolve());
    Object.defineProperty(navigator, 'clipboard', { configurable: true, value: mode === 'missing' ? undefined : { writeText } });
    const wrapper = mountView();
    await openAndCreate(wrapper);
    await wrapper.get('[data-cy="api-key-copy"]').trigger('click');
    await flushPromises();
    const snackbar = wrapper.get('[data-cy="api-keys-snackbar"]');
    expect(snackbar.attributes('color')).toBe(mode === 'success' ? 'success' : 'warning');
    expect(snackbar.attributes('timeout')).toBe(mode === 'success' ? '3000' : '4000');
    if (mode === 'success') expect(writeText).toHaveBeenCalledWith(key);
    else {
      const input = wrapper.get('[data-cy="api-key-value"] input').element as HTMLInputElement;
      expect(input.selectionEnd! - input.selectionStart!).toBe(key.length);
    }
  });

  it('requires confirmation to revoke, then reloads the list', async () => {
    const wrapper = mountView();
    await flushPromises();
    await wrapper.get('[data-cy="api-key-revoke"]').trigger('click');
    expect(deleteApiKey).not.toHaveBeenCalled();
    expect(wrapper.get('[data-cy="confirm-dialog"]').text()).toContain('Finance');
    await wrapper.get('[data-cy="confirm-dialog-confirm"]').trigger('click');
    await flushPromises();
    expect(deleteApiKey).toHaveBeenCalledWith('record-id');
    expect(fetchApiKeys).toHaveBeenCalledTimes(2);
    expect(wrapper.get('[data-cy="api-keys-snackbar"]').attributes('timeout')).toBe('3000');
  });

  it.each(['create', 'revoke'])('shows backend detail on %s failure', async (operation) => {
    const err = new ApiError('Failed', 400, 'Precise backend detail');
    const wrapper = mountView();
    if (operation === 'create') {
      vi.mocked(createApiKey).mockRejectedValueOnce(err);
      await openAndCreate(wrapper);
    } else {
      vi.mocked(deleteApiKey).mockRejectedValueOnce(err);
      await flushPromises();
      await wrapper.get('[data-cy="api-key-revoke"]').trigger('click');
      await wrapper.get('[data-cy="confirm-dialog-confirm"]').trigger('click');
      await flushPromises();
    }
    const snackbar = wrapper.get('[data-cy="api-keys-snackbar"]');
    expect(snackbar.text()).toContain('Precise backend detail');
    expect(snackbar.attributes('timeout')).toBe('6000');
    expect(snackbar.attributes('closable')).toBe('true');
  });

  it('renders an empty state and persistent load failures', async () => {
    vi.mocked(fetchApiKeys).mockResolvedValueOnce({ data: [] });
    const empty = mountView();
    await flushPromises();
    expect(empty.find('[data-cy="api-keys-empty"]').exists()).toBe(true);
    empty.unmount();
    vi.mocked(fetchApiKeys).mockRejectedValueOnce(new ApiError('Fail', 500, 'Database unavailable'));
    const failed = mountView();
    await flushPromises();
    expect(failed.get('[data-cy="api-keys-load-error"]').text()).toContain('Database unavailable');
  });

  it('cancels creation without a request and still refreshes the list', async () => {
    const wrapper = mountView();
    await flushPromises();
    await wrapper.get('[data-cy="api-keys-add"]').trigger('click');
    await wrapper.get('[data-cy="api-key-cancel"]').trigger('click');
    expect(wrapper.find('[data-cy="api-key-dialog"]').exists()).toBe(false);
    expect(createApiKey).not.toHaveBeenCalled();
    // Always reloaded: a closed dialog cannot rule out a committed mid-flight create.
    expect(fetchApiKeys).toHaveBeenCalledTimes(2);
  });

  it('disables the form while creation is in flight', async () => {
    let resolve!: (response: Awaited<ReturnType<typeof createApiKey>>) => void;
    vi.mocked(createApiKey).mockReturnValueOnce(new Promise(done => { resolve = done; }));
    const wrapper = mountView();
    await openAndCreate(wrapper);
    expect(wrapper.get('[data-cy="api-key-submit"]').attributes('disabled')).toBeDefined();
    expect(wrapper.get('[data-cy="api-key-cancel"]').attributes('disabled')).toBeDefined();
    await wrapper.get('[data-cy="api-key-submit"]').trigger('click');
    expect(createApiKey).toHaveBeenCalledTimes(1);
    resolve({ data: { ...record, attributes: { ...record.attributes, key } } });
    await flushPromises();
    expect(wrapper.find('[data-cy="api-key-value"]').exists()).toBe(true);
  });

  it('does not retain a key from a request that finishes after leaving the page', async () => {
    let resolve!: (response: Awaited<ReturnType<typeof createApiKey>>) => void;
    vi.mocked(createApiKey).mockReturnValueOnce(new Promise(done => { resolve = done; }));
    const wrapper = mountView();
    await openAndCreate(wrapper);
    wrapper.unmount();
    resolve({ data: { ...record, attributes: { ...record.attributes, key } } });
    await flushPromises();
    expect(document.body.innerHTML).not.toContain(key);
  });

  it.each(['create', 'revoke'])('redirects on authorization failure during %s', async (operation) => {
    const wrapper = mountView();
    const err = new ApiError('Forbidden', 403);
    if (operation === 'create') {
      vi.mocked(createApiKey).mockRejectedValueOnce(err);
      await openAndCreate(wrapper);
    } else {
      await flushPromises();
      vi.mocked(deleteApiKey).mockRejectedValueOnce(err);
      await wrapper.get('[data-cy="api-key-revoke"]').trigger('click');
      await wrapper.get('[data-cy="confirm-dialog-confirm"]').trigger('click');
      await flushPromises();
    }
    expect(pushMock).toHaveBeenCalledWith('/access-denied');
    expect(wrapper.find('[data-cy="api-keys-snackbar"]').exists()).toBe(false);
  });

  it('does not mark an unexpired key as expired', async () => {
    vi.mocked(fetchApiKeys).mockResolvedValueOnce({ data: [{
      ...record, attributes: { ...record.attributes, expires_at: '2099-01-01T00:00:00Z' }
    }] });
    const wrapper = mountView();
    await flushPromises();
    expect(wrapper.get('[data-cy="api-key-expiry"]').text()).not.toContain('Expired');
  });

  defineAuthRedirectTests(vi.mocked(fetchApiKeys), mountView, pushMock);
});
