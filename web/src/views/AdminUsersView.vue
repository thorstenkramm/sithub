<template>
  <div class="page-container">
    <PageHeader
      :title="$t('adminUsers.title')"
      :subtitle="$t('adminUsers.subtitle')"
      :breadcrumbs="[{ text: $t('common.home'), to: '/' }, { text: $t('adminUsers.title') }]"
    />

    <!-- Load failure is a persistent page state (see .claude/rules/feedback.md) -->
    <v-alert v-if="loadError" type="error" class="mb-4 api-error-text" data-cy="users-load-error">
      {{ loadError }}
    </v-alert>

    <v-card>
      <!-- Order is fixed by the server (display name), so no column is sortable. -->
      <v-data-table-server
        :headers="headers"
        :items="rows"
        :items-length="total"
        :items-per-page="PAGE_SIZE"
        :page="page"
        :loading="loading"
        item-value="id"
        density="comfortable"
        hide-default-footer
        data-cy="users-table"
      >
        <template #[`item.email`]="{ item }">
          <span data-cy="user-email">{{ item.email }}</span>
        </template>
        <template #[`item.name`]="{ item }">
          <span data-cy="user-name">{{ item.name }}</span>
          <v-chip
            v-if="item.deleted"
            size="x-small"
            color="grey"
            class="ml-2"
            data-cy="user-deleted-chip"
          >
            {{ $t('adminUsers.deactivated') }}
          </v-chip>
        </template>
        <template #[`item.roles`]="{ item }">
          <div class="d-flex flex-wrap ga-1" data-cy="user-roles">
            <v-chip
              v-for="roleKey in item.roleKeys"
              :key="roleKey"
              size="small"
              :color="roleKey === 'adminUsers.roleAdmin' ? 'secondary' : undefined"
            >
              {{ $t(roleKey) }}
            </v-chip>
          </div>
        </template>
        <template #[`item.dataExporter`]="{ item }">
          <v-switch
            :model-value="item.isDataExporter"
            :disabled="togglingIds.has(item.id)"
            :aria-label="$t('adminUsers.colDataExporter')"
            density="compact"
            color="primary"
            hide-details
            data-cy="user-data-exporter-toggle"
            @update:model-value="toggleDataExporter(item, $event === true)"
          />
        </template>
        <template #[`item.lastLogin`]="{ item }">
          <span data-cy="user-last-login">{{ item.lastLogin ?? $t('adminUsers.never') }}</span>
        </template>
        <template #[`item.actions`]="{ item }">
          <v-tooltip location="top" :disabled="!isOwnRow(item)">
            <template #activator="{ props: tooltipProps }">
              <span v-bind="tooltipProps">
                <v-btn
                  icon="$delete"
                  variant="text"
                  size="small"
                  color="error"
                  :aria-label="$t('adminUsers.remove')"
                  :disabled="isOwnRow(item) || item.deleted"
                  data-cy="user-delete-btn"
                  @click="openDeleteDialog(item)"
                />
              </span>
            </template>
            <span>{{ $t('adminUsers.cannotRemoveSelf') }}</span>
          </v-tooltip>
        </template>
      </v-data-table-server>
    </v-card>

    <v-pagination
      v-model="page"
      :length="pageCount"
      class="mt-4"
      data-cy="users-pagination"
    />

    <ConfirmDialog
      v-model="deleteDialogOpen"
      :title="deleteDialogTitle"
      :message="deleteDialogMessage"
      :confirm-text="deleteDialogConfirmText"
      confirm-color="error"
      :loading="deleting"
      @confirm="confirmDelete"
    />

    <v-snackbar
      v-model="feedbackOpen"
      :timeout="feedbackTimeout"
      :color="feedbackColor"
      location="bottom"
      :closable="feedbackColor === 'error'"
      data-cy="users-feedback"
    >
      {{ feedbackMessage }}
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  deleteUser,
  fetchUserDeletionImpact,
  fetchUsers,
  updateUser,
  type UserDeletionImpactAttributes,
  type UserListAttributes
} from '../api/users';
import type { JsonApiResource } from '../api/types';
import { formatApiErrorMessage } from '../api/client';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { formatDateTime } from '../utils/date';
import { userRoleLabels } from '../utils/userRoles';
import { ConfirmDialog, PageHeader } from '../components';
import { useAuthStore } from '../stores/useAuthStore';

interface UserRow {
  id: string;
  email: string;
  name: string;
  roleKeys: string[];
  isDataExporter: boolean;
  deleted: boolean;
  authSource: string;
  /** Localized last login, null if the user never logged in. */
  lastLogin: string | null;
}

const PAGE_SIZE = 20;

const { t, locale } = useI18n();
const { loading, run } = useApi();
const { handleAuthError } = useAuthErrorHandler();
const authStore = useAuthStore();

const users = ref<JsonApiResource<UserListAttributes>[]>([]);
const total = ref(0);
const page = ref(1);
const loadError = ref<string | null>(null);
const deleteDialogOpen = ref(false);
const selectedUser = ref<UserRow | null>(null);
const selectedImpact = ref<UserDeletionImpactAttributes | null>(null);
const deleting = ref(false);
// Users whose data exporter toggle has a request in flight.
const togglingIds = ref(new Set<string>());
const feedbackOpen = ref(false);
const feedbackMessage = ref('');
const feedbackColor = ref<'success' | 'error'>('success');
const feedbackTimeout = ref(3000);

const rows = computed<UserRow[]>(() =>
  users.value.map((resource) => ({
    id: resource.id,
    email: resource.attributes.email,
    name: resource.attributes.display_name,
    roleKeys: userRoleLabels(resource.attributes),
    isDataExporter: resource.attributes.is_data_exporter,
    deleted: resource.attributes.deleted,
    authSource: resource.attributes.auth_source,
    lastLogin: formatDateTime(resource.attributes.last_login, locale.value || undefined)
  }))
);

const pageCount = computed(() => Math.max(1, Math.ceil(total.value / PAGE_SIZE)));
// Users with bookings are deactivated to keep the booking history, so the dialog never says "delete" for them.
const isDeactivation = computed(() => selectedImpact.value?.mode === 'deactivate');
const deleteDialogTitle = computed(() =>
  t(isDeactivation.value ? 'adminUsers.deactivateTitle' : 'adminUsers.deleteTitle')
);
const deleteDialogConfirmText = computed(() =>
  t(isDeactivation.value ? 'adminUsers.deactivateConfirm' : 'adminUsers.deleteConfirm')
);
const deleteDialogMessage = computed(() => {
  if (!selectedUser.value || !selectedImpact.value) return '';
  const name = selectedUser.value.name;
  if (!isDeactivation.value) {
    return t('adminUsers.deleteMessage', { name });
  }
  // Entra ID users are reactivated on their next Microsoft sign-in, so they are not locked out for good.
  const messageKey = selectedUser.value.authSource === 'entraid'
    ? 'adminUsers.deactivateMessageEntra'
    : 'adminUsers.deactivateMessage';
  const futureCount = selectedImpact.value.future_booking_count;
  const canceled = futureCount > 0 ? ` ${t('adminUsers.futureBookingsCanceled', { count: futureCount }, futureCount)}` : '';
  return t(messageKey, { name }) + canceled;
});

const headers = computed(() => [
  { title: t('adminUsers.colEmail'), key: 'email', sortable: false },
  { title: t('adminUsers.colName'), key: 'name', sortable: false },
  { title: t('adminUsers.colRoles'), key: 'roles', sortable: false },
  { title: t('adminUsers.colDataExporter'), key: 'dataExporter', sortable: false },
  { title: t('adminUsers.colLastLogin'), key: 'lastLogin', sortable: false },
  { title: '', key: 'actions', sortable: false, align: 'end' as const }
]);

// Guards against out-of-order responses when the page changes quickly.
let latestRequest = 0;

async function loadUsers() {
  const requestId = ++latestRequest;
  loadError.value = null;
  try {
    const resp = await run(() => fetchUsers({ limit: PAGE_SIZE, offset: (page.value - 1) * PAGE_SIZE }));
    if (requestId !== latestRequest) return;
    total.value = resp.meta.total;
    users.value = resp.data;
  } catch (err) {
    if (requestId !== latestRequest) return;
    if (await handleAuthError(err)) return;
    loadError.value = formatApiErrorMessage(t('adminUsers.loadError'), err);
  }
}

function isOwnRow(item: UserRow) {
  return item.id === authStore.userId;
}

function showFeedback(message: string, color: 'success' | 'error') {
  feedbackMessage.value = message;
  feedbackColor.value = color;
  feedbackTimeout.value = color === 'success' ? 3000 : 6000;
  feedbackOpen.value = true;
}

/**
 * Grants or revokes the data exporter role. The switch is bound to the server state, so a failed
 * request leaves it in its previous position without extra revert logic.
 */
async function toggleDataExporter(item: UserRow, granted: boolean) {
  togglingIds.value = new Set(togglingIds.value).add(item.id);
  try {
    const resp = await updateUser(item.id, { is_data_exporter: granted });
    users.value = users.value.map((resource) =>
      resource.id === item.id ? { ...resource, attributes: resp.data.attributes } : resource
    );
    if (isOwnRow(item)) {
      authStore.isDataExporter = granted;
    }
    showFeedback(t(granted ? 'adminUsers.dataExporterGranted' : 'adminUsers.dataExporterRevoked'), 'success');
  } catch (err) {
    if (await handleAuthError(err)) return;
    showFeedback(formatApiErrorMessage(t('adminUsers.dataExporterError'), err), 'error');
  } finally {
    const remaining = new Set(togglingIds.value);
    remaining.delete(item.id);
    togglingIds.value = remaining;
  }
}

async function openDeleteDialog(item: UserRow) {
  if (isOwnRow(item) || item.deleted) return;
  selectedUser.value = item;
  selectedImpact.value = null;
  try {
    const resp = await fetchUserDeletionImpact(item.id);
    selectedImpact.value = resp.data.attributes;
    deleteDialogOpen.value = true;
  } catch (err) {
    if (await handleAuthError(err)) return;
    showFeedback(formatApiErrorMessage(t('adminUsers.removeError'), err), 'error');
  }
}

function removalSuccessMessage(name: string, deleted: boolean, canceledCount: number) {
  if (deleted) return t('adminUsers.deleteSuccess', { name });
  if (canceledCount === 0) return t('adminUsers.deactivateSuccess', { name });
  return t('adminUsers.deactivateSuccessCanceled', { name, count: canceledCount }, canceledCount);
}

async function confirmDelete() {
  if (!selectedUser.value || !selectedImpact.value) return;
  deleting.value = true;
  const { name } = selectedUser.value;
  const deactivation = isDeactivation.value;
  const canceledCount = selectedImpact.value.future_booking_count;
  try {
    const result = await deleteUser(selectedUser.value.id);
    deleteDialogOpen.value = false;
    // A 204 means the server hard-deleted the user, even if the preview predicted a deactivation.
    showFeedback(removalSuccessMessage(name, result === null || !deactivation, canceledCount), 'success');
    if (result === null && rows.value.length === 1 && page.value > 1) {
      page.value -= 1;
    } else {
      await loadUsers();
    }
  } catch (err) {
    if (await handleAuthError(err)) return;
    const errorKey = deactivation ? 'adminUsers.deactivateError' : 'adminUsers.deleteError';
    showFeedback(formatApiErrorMessage(t(errorKey), err), 'error');
  } finally {
    deleting.value = false;
  }
}

watch(page, loadUsers);
onMounted(loadUsers);
</script>
