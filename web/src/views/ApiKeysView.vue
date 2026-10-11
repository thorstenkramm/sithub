<template>
  <div class="page-container">
    <PageHeader :title="$t('apiKeys.title')">
      <template #actions>
        <v-btn color="primary" prepend-icon="$plus" data-cy="api-keys-add" @click="createOpen = true">
          {{ $t('apiKeys.add') }}
        </v-btn>
      </template>
    </PageHeader>
    <v-alert v-if="loadError" type="error" class="mb-4 api-error-text" data-cy="api-keys-load-error">
      {{ loadError }}
    </v-alert>
    <EmptyState
      v-else-if="!loading && !rows.length" :title="$t('apiKeys.emptyTitle')"
      :message="$t('apiKeys.emptyMessage')" icon="$key" data-cy="api-keys-empty"
    />
    <v-card v-else>
      <v-data-table
        :headers="headers" :items="rows" :loading="loading" item-value="id"
        :row-props="{ 'data-cy': 'api-key-row' }" data-cy="api-keys-table"
      >
        <template #[`item.name`]="{ item }"><span data-cy="api-key-name">{{ item.name }}</span></template>
        <template #[`item.key_prefix`]="{ item }">
          <code data-cy="api-key-prefix">{{ item.key_prefix }}…</code>
        </template>
        <template #[`item.expires_at`]="{ item }">
          <span data-cy="api-key-expiry">
            {{ formatExpiry(item.expires_at) }}
            <v-chip v-if="isExpired(item.expires_at)" color="error" size="small" class="ml-2">
              {{ $t('apiKeys.expired') }}
            </v-chip>
          </span>
        </template>
        <template #[`item.scope`]><span data-cy="api-key-scope">{{ $t('apiKeys.scope') }}</span></template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            icon="$delete" color="error" variant="text" size="small" :aria-label="$t('apiKeys.revoke')"
            data-cy="api-key-revoke" @click="selected = item; revokeOpen = true"
          />
        </template>
      </v-data-table>
    </v-card>
    <ApiKeyCreateDialog v-if="createOpen" @close="closeCreate" @feedback="showFeedback" />
    <ConfirmDialog
      v-model="revokeOpen" :title="$t('apiKeys.revokeTitle')"
      :message="$t('apiKeys.revokeMessage', { name: selected?.name ?? '' })"
      :confirm-text="$t('apiKeys.revoke')" confirm-color="error" :loading="revoking" @confirm="revoke"
    />
    <v-snackbar
      v-model="feedbackOpen" :color="feedbackColor" :timeout="feedbackTimeout"
      :closable="feedbackColor === 'error'" location="bottom" data-cy="api-keys-snackbar"
    >
      <span class="api-error-text">{{ feedbackMessage }}</span>
      <template v-if="feedbackColor === 'error'" #actions>
        <v-btn icon="$close" :aria-label="$t('common.close')" @click="feedbackOpen = false" />
      </template>
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { deleteApiKey, fetchApiKeys, type ApiKeyAttributes } from '../api/apiKeys';
import { formatApiErrorMessage } from '../api/client';
import { ConfirmDialog, EmptyState, PageHeader } from '../components';
import ApiKeyCreateDialog from '../components/ApiKeyCreateDialog.vue';
import { useApi } from '../composables/useApi';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';
import { formatMediumDate } from '../utils/date';

type KeyRow = ApiKeyAttributes & { id: string };
const { t, locale } = useI18n();
const { loading, run } = useApi();
const { handleAuthError } = useAuthErrorHandler();
const rows = ref<KeyRow[]>([]);
const loadError = ref('');
const createOpen = ref(false);
const revokeOpen = ref(false);
const selected = ref<KeyRow | null>(null);
const revoking = ref(false);
const feedbackOpen = ref(false);
const feedbackMessage = ref('');
const feedbackColor = ref<'success' | 'warning' | 'error'>('success');
const feedbackTimeout = computed(() => ({ success: 3000, warning: 4000, error: 6000 })[feedbackColor.value]);

/**
 * Expiry is compared against the current time on every render; keeping `now` reactive
 * and updating it once per minute is enough for a date-granular expiry display.
 */
const now = ref(Date.now());
const timer = setInterval(() => { now.value = Date.now(); }, 60_000);
onBeforeUnmount(() => clearInterval(timer));

function isExpired(expiresAt: string): boolean {
  return Date.parse(expiresAt) <= now.value;
}

function formatExpiry(expiresAt: string): string {
  // expires_at is RFC 3339; show the calendar date without time of day.
  return formatMediumDate(expiresAt.slice(0, 10), locale.value || undefined);
}
const headers = computed(() => [
  { title: t('apiKeys.columns.name'), key: 'name' },
  { title: t('apiKeys.columns.prefix'), key: 'key_prefix' },
  { title: t('apiKeys.columns.expiry'), key: 'expires_at' },
  { title: t('apiKeys.columns.scope'), key: 'scope' },
  { title: '', key: 'actions', sortable: false }
]);

async function loadKeys() {
  loadError.value = '';
  try {
    const response = await run(fetchApiKeys);
    rows.value = response.data.map(resource => ({ id: resource.id, ...resource.attributes }));
  } catch (err) {
    if (await handleAuthError(err)) return;
    loadError.value = formatApiErrorMessage(t('apiKeys.loadError'), err);
  }
}

function showFeedback(message: string, color: 'success' | 'warning' | 'error') {
  feedbackMessage.value = message;
  feedbackColor.value = color;
  feedbackOpen.value = true;
}

/**
 * Always reload after the dialog closes. A create request may still commit server-side
 * when the admin closes the dialog or navigates mid-flight, so the list must reflect
 * any key whose plaintext was never displayed.
 */
async function closeCreate() {
  createOpen.value = false;
  await loadKeys();
}

async function revoke() {
  if (!selected.value || revoking.value) return;
  revoking.value = true;
  try {
    await deleteApiKey(selected.value.id);
    revokeOpen.value = false;
    selected.value = null;
    showFeedback(t('apiKeys.revoked'), 'success');
    await loadKeys();
  } catch (err) {
    if (await handleAuthError(err)) return;
    showFeedback(formatApiErrorMessage(t('apiKeys.revokeError'), err), 'error');
  } finally {
    revoking.value = false;
  }
}

onMounted(loadKeys);
</script>
