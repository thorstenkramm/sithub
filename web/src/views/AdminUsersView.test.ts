import { h } from 'vue';
import { mount, flushPromises } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import AdminUsersView from './AdminUsersView.vue';
import {
  deleteUser,
  fetchUserDeletionImpact,
  fetchUsers,
  updateUser,
  type UserListAttributes
} from '../api/users';
import { ApiError } from '../api/client';
import { buildViewStubs, createTestI18n, defineAuthRedirectTests } from './testHelpers';
import { useAuthStore } from '../stores/useAuthStore';

const pushMock = vi.fn();

vi.mock('../api/users', () => ({
  deleteUser: vi.fn(),
  fetchUserDeletionImpact: vi.fn(),
  fetchUsers: vi.fn(),
  updateUser: vi.fn()
}));
vi.mock('vue-router', () => ({ useRouter: () => ({ push: pushMock }) }));

type Row = { id: string };
type SlotFn = (args: { item: Row }) => unknown;

// Renders one <tr> per item with the named item.* slots, so cell content can be asserted.
const dataTableStub = {
  props: ['headers', 'items'],
  setup(
    props: { headers: Array<{ key: string; title: string }>; items: Row[] },
    { slots }: { slots: Record<string, SlotFn | undefined> }
  ) {
    return () =>
      h('table', { 'data-cy': 'users-table' }, [
        h('thead', {}, props.headers.map((header) => h('th', { key: header.key }, header.title))),
        h(
          'tbody',
          {},
          props.items.map((item) =>
            h(
              'tr',
              { key: item.id, 'data-cy': 'user-row' },
              props.headers.map((header) => h('td', {}, slots[`item.${header.key}`]?.({ item }) as never))
            )
          )
        )
      ]);
  }
};

const paginationStub = {
  props: ['modelValue', 'length'],
  emits: ['update:modelValue'],
  template: '<nav data-cy="users-pagination" :data-length="length"></nav>'
};

function user(id: string, overrides: Partial<UserListAttributes> = {}) {
  return {
    id,
    type: 'users',
    attributes: {
      email: `${id}@example.com`,
      display_name: `Name ${id}`,
      is_admin: false,
      is_data_exporter: false,
      auth_source: 'internal',
      role: 'user',
      last_login: '',
      deleted: false,
      created_at: '2026-01-01T00:00:00Z',
      updated_at: '2026-01-01T00:00:00Z',
      ...overrides
    }
  };
}

function page(data: ReturnType<typeof user>[], total: number, offset = 0) {
  return { data, meta: { total, limit: 20, offset }, links: {} };
}

describe('AdminUsersView', () => {
  const fetchUsersMock = fetchUsers as unknown as ReturnType<typeof vi.fn>;
  const fetchUserDeletionImpactMock = fetchUserDeletionImpact as unknown as ReturnType<typeof vi.fn>;
  const deleteUserMock = deleteUser as unknown as ReturnType<typeof vi.fn>;
  const updateUserMock = updateUser as unknown as ReturnType<typeof vi.fn>;
  const stubs = {
    ...buildViewStubs(['v-icon', 'v-breadcrumbs', 'v-breadcrumbs-item', 'v-breadcrumbs-divider']),
    'v-data-table-server': dataTableStub,
    'v-pagination': paginationStub,
    'v-chip': { template: '<span v-bind="$attrs"><slot /></span>' },
    'v-btn': {
      props: ['disabled'],
      emits: ['click'],
      template: '<button v-bind="$attrs" :disabled="disabled" @click="$emit(\'click\', $event)"><slot /></button>'
    },
    // Minimal switch: a checkbox that emits update:modelValue like v-switch.
    'v-switch': {
      props: ['modelValue', 'disabled'],
      emits: ['update:modelValue'],
      methods: {
        // Like a controlled v-switch: the DOM follows modelValue, not the click.
        onChange(event: Event) {
          const input = event.target as HTMLInputElement;
          this.$emit('update:modelValue', input.checked);
          void this.$nextTick(() => { input.checked = this.modelValue; });
        }
      },
      template: '<input type="checkbox" v-bind="$attrs" :checked="modelValue" :disabled="disabled" @change="onChange" />'
    },
    'v-tooltip': { template: '<div><slot name="activator" :props="{}" /><slot /></div>' },
    ConfirmDialog: {
      props: ['modelValue', 'title', 'message', 'confirmText', 'loading'],
      emits: ['confirm'],
      template: '<div v-if="modelValue" data-cy="confirm-dialog"><h2 data-cy="confirm-dialog-title">{{ title }}</h2><p data-cy="confirm-dialog-message">{{ message }}</p><button data-cy="confirm-dialog-confirm" @click="$emit(\'confirm\')">{{ confirmText }}</button></div>'
    },
    'v-snackbar': { props: ['modelValue'], template: '<div v-if="modelValue" v-bind="$attrs"><slot /></div>' }
  };

  const mountView = () =>
    mount(AdminUsersView, {
      global: { stubs, plugins: [createPinia(), createTestI18n()] }
    });

  beforeEach(() => {
    setActivePinia(createPinia());
    pushMock.mockReset();
    fetchUsersMock.mockReset();
    fetchUserDeletionImpactMock.mockReset();
    deleteUserMock.mockReset();
    updateUserMock.mockReset();
    fetchUsersMock.mockResolvedValue(page([user('u1'), user('u2')], 45));
    fetchUserDeletionImpactMock.mockResolvedValue({
      data: { id: 'u1', type: 'user-deletion-impacts', attributes: { has_bookings: false, future_booking_count: 0, mode: 'delete' } }
    });
    deleteUserMock.mockResolvedValue(null);
  });

  it('should load the first page with limit 20 and offset 0', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(fetchUsersMock).toHaveBeenCalledTimes(1);
    expect(fetchUsersMock).toHaveBeenCalledWith({ limit: 20, offset: 0 });
    expect(wrapper.findAll('[data-cy="user-row"]')).toHaveLength(2);
    expect(wrapper.get('[data-cy="users-pagination"]').attributes('data-length')).toBe('3');
  });

  it('should render the six column headers', async () => {
    const wrapper = mountView();
    await flushPromises();

    const headers = wrapper.findAll('th').map((th) => th.text());
    expect(headers).toEqual(['Email', 'Name', 'Role(s)', 'Data exporter', 'Last login', '']);
  });

  it('should fetch offset 40 when selecting page 3', async () => {
    const wrapper = mountView();
    await flushPromises();

    fetchUsersMock.mockResolvedValue(page([user('u41')], 45, 40));
    await wrapper.getComponent(paginationStub).vm.$emit('update:modelValue', 3);
    await flushPromises();

    expect(fetchUsersMock).toHaveBeenLastCalledWith({ limit: 20, offset: 40 });
    expect(wrapper.get('[data-cy="user-email"]').text()).toBe('u41@example.com');
  });

  it('should ignore a stale response that resolves after a newer page', async () => {
    let resolveFirst: (value: unknown) => void = () => {};
    fetchUsersMock.mockReturnValueOnce(new Promise((resolve) => { resolveFirst = resolve; }));
    const wrapper = mountView();

    fetchUsersMock.mockResolvedValueOnce(page([user('u21')], 45, 20));
    await wrapper.getComponent(paginationStub).vm.$emit('update:modelValue', 2);
    await flushPromises();
    resolveFirst(page([user('u1')], 45));
    await flushPromises();

    expect(wrapper.get('[data-cy="user-email"]').text()).toBe('u21@example.com');
  });

  it('should render email, name and "never" for a user without login', async () => {
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.get('[data-cy="user-email"]').text()).toBe('u1@example.com');
    expect(wrapper.get('[data-cy="user-name"]').text()).toBe('Name u1');
    expect(wrapper.get('[data-cy="user-last-login"]').text()).toBe('never');
  });

  it('should render a localized date-time for a set last login', async () => {
    const lastLogin = '2026-03-05T14:30:00Z';
    fetchUsersMock.mockResolvedValue(page([user('u1', { last_login: lastLogin })], 1));
    const wrapper = mountView();
    await flushPromises();

    const expected = new Intl.DateTimeFormat('en', { dateStyle: 'medium', timeStyle: 'short' })
      .format(new Date(lastLogin));
    expect(wrapper.get('[data-cy="user-last-login"]').text()).toBe(expected);
  });

  it.each([
    { isAdmin: true, isExporter: false, label: 'Admin' },
    { isAdmin: false, isExporter: false, label: 'User' },
    { isAdmin: false, isExporter: true, label: 'Data exporter' },
    { isAdmin: true, isExporter: true, label: 'AdminData exporter' }
  ])('should show the "$label" role chips', async ({ isAdmin, isExporter, label }) => {
    fetchUsersMock.mockResolvedValue(
      page([user('u1', { is_admin: isAdmin, is_data_exporter: isExporter })], 1)
    );
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.get('[data-cy="user-roles"]').text()).toBe(label);
  });

  it('should grant the data exporter role with the toggle and update the row', async () => {
    updateUserMock.mockResolvedValue({ data: user('u1', { is_data_exporter: true }) });
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.findAll('[data-cy="user-data-exporter-toggle"]')[0].element).toHaveProperty('checked', false);
    await wrapper.findAll('[data-cy="user-data-exporter-toggle"]')[0].setValue(true);
    await flushPromises();

    expect(updateUserMock).toHaveBeenCalledWith('u1', { is_data_exporter: true });
    expect(wrapper.findAll('[data-cy="user-roles"]')[0].text()).toBe('Data exporter');
    expect(wrapper.get('[data-cy="users-feedback"]').text()).toContain('Data exporter role granted');
  });

  it('should revoke the role and report it', async () => {
    fetchUsersMock.mockResolvedValue(page([user('u1', { is_data_exporter: true })], 1));
    updateUserMock.mockResolvedValue({ data: user('u1', { is_data_exporter: false }) });
    const wrapper = mountView();
    await flushPromises();

    await wrapper.get('[data-cy="user-data-exporter-toggle"]').setValue(false);
    await flushPromises();

    expect(updateUserMock).toHaveBeenCalledWith('u1', { is_data_exporter: false });
    expect(wrapper.get('[data-cy="user-roles"]').text()).toBe('User');
    expect(wrapper.get('[data-cy="users-feedback"]').text()).toContain('revoked');
  });

  it('should revert the toggle and show the backend detail when the update fails', async () => {
    updateUserMock.mockRejectedValue(new ApiError('Request failed: 500', 500, 'database is locked'));
    const wrapper = mountView();
    await flushPromises();

    await wrapper.findAll('[data-cy="user-data-exporter-toggle"]')[0].setValue(true);
    await flushPromises();

    expect(wrapper.findAll('[data-cy="user-data-exporter-toggle"]')[0].element).toHaveProperty('checked', false);
    expect(wrapper.findAll('[data-cy="user-roles"]')[0].text()).toBe('User');
    expect(wrapper.get('[data-cy="users-feedback"]').text()).toContain('database is locked');
  });

  it('should keep the signed-in admin menu state in sync when toggling their own row', async () => {
    updateUserMock.mockResolvedValue({ data: user('u1', { is_data_exporter: true }) });
    const pinia = createPinia();
    setActivePinia(pinia);
    const wrapper = mount(AdminUsersView, { global: { stubs, plugins: [pinia, createTestI18n()] } });
    const authStore = useAuthStore();
    authStore.userId = 'u1';
    await flushPromises();

    await wrapper.findAll('[data-cy="user-data-exporter-toggle"]')[0].setValue(true);
    await flushPromises();

    expect(authStore.isDataExporter).toBe(true);
  });

  it('should render a deleted chip and disable the delete button for deleted rows', async () => {
    fetchUsersMock.mockResolvedValue(page([user('u1', { deleted: true })], 1));
    const wrapper = mountView();
    await flushPromises();

    expect(wrapper.get('[data-cy="user-deleted-chip"]').text()).toBe('deactivated');
    expect(wrapper.get('[data-cy="user-delete-btn"]').attributes()).toHaveProperty('disabled');
  });

  it('should disable the delete button for the current user', async () => {
    const pinia = createPinia();
    setActivePinia(pinia);
    const wrapper = mount(AdminUsersView, {
      global: { stubs, plugins: [pinia, createTestI18n()] }
    });
    useAuthStore().userId = 'u1';
    await flushPromises();

    expect(wrapper.findAll('[data-cy="user-delete-btn"]')[0].attributes()).toHaveProperty('disabled');
  });

  function impact(mode: 'delete' | 'deactivate', futureBookingCount = 0) {
    return {
      data: {
        id: 'u1',
        type: 'user-deletion-impacts',
        attributes: { has_bookings: mode === 'deactivate', future_booking_count: futureBookingCount, mode }
      }
    };
  }

  async function openDialogFor(rows: ReturnType<typeof user>[]) {
    fetchUsersMock.mockResolvedValue(page(rows, rows.length));
    const wrapper = mountView();
    await flushPromises();
    await wrapper.findAll('[data-cy="user-delete-btn"]')[0].trigger('click');
    await flushPromises();
    return wrapper;
  }

  const deactivationText = 'To keep the booking history, Name u1 will be deactivated instead of deleted.';

  it.each([
    {
      name: 'a user without any bookings is deleted',
      impact: impact('delete'),
      authSource: 'internal',
      title: 'Delete user',
      confirm: 'Delete',
      message: 'Name u1 has no bookings and will be permanently deleted. This cannot be undone.',
      success: 'Name u1 was deleted.'
    },
    {
      name: 'a local user with past bookings only is deactivated without mentioning cancellations',
      impact: impact('deactivate'),
      authSource: 'internal',
      title: 'Deactivate user',
      confirm: 'Deactivate',
      message: `${deactivationText} They will no longer be able to sign in.`,
      success: 'Name u1 was deactivated.'
    },
    {
      name: 'a local user with one future booking is deactivated and the cancellation is announced',
      impact: impact('deactivate', 1),
      authSource: 'internal',
      title: 'Deactivate user',
      confirm: 'Deactivate',
      message: `${deactivationText} They will no longer be able to sign in. One upcoming booking will be canceled.`,
      success: 'Name u1 was deactivated and one booking was canceled.'
    },
    {
      name: 'an Entra ID user is told that a Microsoft sign-in reactivates the account',
      impact: impact('deactivate', 3),
      authSource: 'entraid',
      title: 'Deactivate user',
      confirm: 'Deactivate',
      message: `${deactivationText} If Name u1 signs in with Microsoft again, the account is reactivated automatically. `
        + 'To block access permanently, remove them from the Entra ID group. 3 upcoming bookings will be canceled.',
      success: 'Name u1 was deactivated and 3 bookings were canceled.'
    }
  ])('should word the dialog consistently when $name', async (tc) => {
    fetchUserDeletionImpactMock.mockResolvedValue(tc.impact);
    deleteUserMock.mockResolvedValue(tc.impact.data.attributes.mode === 'delete' ? null : { data: user('u1') });
    const wrapper = await openDialogFor([user('u1', { auth_source: tc.authSource }), user('u2')]);

    expect(fetchUserDeletionImpactMock).toHaveBeenCalledWith('u1');
    expect(wrapper.get('[data-cy="confirm-dialog-title"]').text()).toBe(tc.title);
    expect(wrapper.get('[data-cy="confirm-dialog-message"]').text()).toBe(tc.message);
    expect(wrapper.get('[data-cy="confirm-dialog-confirm"]').text()).toBe(tc.confirm);

    await wrapper.get('[data-cy="confirm-dialog-confirm"]').trigger('click');
    await flushPromises();

    expect(deleteUserMock).toHaveBeenCalledWith('u1');
    expect(wrapper.get('[data-cy="users-feedback"]').text()).toBe(tc.success);
  });

  it.each([
    { mode: 'delete' as const, error: 'User could not be deleted.' },
    { mode: 'deactivate' as const, error: 'User could not be deactivated.' }
  ])('should name the failed action when $mode fails', async ({ mode, error }) => {
    fetchUserDeletionImpactMock.mockResolvedValue(impact(mode));
    deleteUserMock.mockRejectedValue(new ApiError('Request failed: 500', 500, 'database is locked'));
    const wrapper = await openDialogFor([user('u1')]);

    await wrapper.get('[data-cy="confirm-dialog-confirm"]').trigger('click');
    await flushPromises();

    const feedback = wrapper.get('[data-cy="users-feedback"]').text();
    expect(feedback).toContain(error);
    expect(feedback).toContain('database is locked');
  });

  it('should show backend detail when delete preview fails', async () => {
    fetchUserDeletionImpactMock.mockRejectedValue(new ApiError('Request failed: 500', 500, 'database is locked'));
    const wrapper = mountView();
    await flushPromises();

    await wrapper.findAll('[data-cy="user-delete-btn"]')[0].trigger('click');
    await flushPromises();

    const feedback = wrapper.get('[data-cy="users-feedback"]').text();
    expect(feedback).toContain('User could not be removed.');
    expect(feedback).toContain('database is locked');
  });

  it('should show the backend detail when loading fails', async () => {
    fetchUsersMock.mockRejectedValue(new ApiError('Request failed: 500', 500, 'database is locked'));
    const wrapper = mountView();
    await flushPromises();

    const alert = wrapper.get('[data-cy="users-load-error"]');
    expect(alert.text()).toContain('Users could not be loaded.');
    expect(alert.text()).toContain('database is locked');
  });

  defineAuthRedirectTests(fetchUsersMock, () => mountView(), pushMock);
});
