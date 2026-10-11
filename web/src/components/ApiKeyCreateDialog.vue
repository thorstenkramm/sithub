<template>
  <v-dialog :model-value="true" persistent max-width="560">
    <v-card data-cy="api-key-dialog">
      <v-card-title>{{ $t(createdKey ? 'apiKeys.createdTitle' : 'apiKeys.add') }}</v-card-title>
      <v-card-text>
        <template v-if="createdKey">
          <v-alert type="warning" class="mb-4" data-cy="api-key-once-warning">
            {{ $t('apiKeys.onceWarning') }}
          </v-alert>
          <v-text-field
            ref="keyField" :model-value="createdKey" readonly class="key-value"
            :label="$t('apiKeys.key')" data-cy="api-key-value"
          />
          <v-btn prepend-icon="$copy" data-cy="api-key-copy" @click="copyKey">
            {{ $t('apiKeys.copy') }}
          </v-btn>
        </template>
        <template v-else>
          <v-text-field
            v-model="name" :label="$t('apiKeys.columns.name')" :disabled="creating"
            maxlength="100" counter="100" required data-cy="api-key-name-input"
          />
          <v-select
            v-model="expiry" :items="expiryOptions" :label="$t('apiKeys.columns.expiry')"
            :disabled="creating" data-cy="api-key-expiry-select"
          />
          <v-select
            v-model="scope" :items="[{ title: $t('apiKeys.scope'), value: 'read-bookings' }]"
            :label="$t('apiKeys.columns.scope')" :disabled="creating" data-cy="api-key-scope-select"
          />
        </template>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn v-if="createdKey" data-cy="api-key-done" @click="close">{{ $t('apiKeys.done') }}</v-btn>
        <template v-else>
          <v-btn :disabled="creating" data-cy="api-key-cancel" @click="close">{{ $t('common.cancel') }}</v-btn>
          <v-btn
            :disabled="!validName || creating" :loading="creating" color="primary"
            data-cy="api-key-submit" @click="submit"
          >{{ $t('apiKeys.add') }}</v-btn>
        </template>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { createApiKey, type ApiKeyExpiryMonths, type ApiKeyScope } from '../api/apiKeys';
import { formatApiErrorMessage } from '../api/client';
import { useAuthErrorHandler } from '../composables/useAuthErrorHandler';

const emit = defineEmits<{
  close: [];
  feedback: [message: string, color: 'success' | 'warning' | 'error'];
}>();
const { t } = useI18n();
const { handleAuthError } = useAuthErrorHandler();
const name = ref('');
const expiry = ref<ApiKeyExpiryMonths>(12);
const scope = ref<ApiKeyScope>('read-bookings');
const creating = ref(false);
// Kept only in this dialog's memory, never in the list or a shared store.
const createdKey = ref('');
const keyField = ref<{ $el: HTMLElement } | null>(null);
const validName = computed(() => name.value.trim().length > 0 && [...name.value.trim()].length <= 100);
const expiryOptions = computed(() => [3, 12, 24].map(value => ({ title: t('apiKeys.months', { count: value }), value })));
let disposed = false;

async function submit() {
  if (!validName.value || creating.value) return;
  creating.value = true;
  try {
    const response = await createApiKey({ name: name.value.trim(), scope: scope.value, expires_in_months: expiry.value });
    if (!disposed) createdKey.value = response.data.attributes.key;
  } catch (err) {
    if (disposed || await handleAuthError(err)) return;
    emit('feedback', formatApiErrorMessage(t('apiKeys.createError'), err), 'error');
  } finally {
    creating.value = false;
  }
}

async function copyKey() {
  try {
    await navigator.clipboard.writeText(createdKey.value);
    if (!disposed) emit('feedback', t('apiKeys.copied'), 'success');
  } catch {
    if (disposed) return;
    const input = keyField.value?.$el.querySelector('input');
    input?.focus();
    input?.select();
    emit('feedback', t('apiKeys.copyManual'), 'warning');
  }
}

function close() {
  createdKey.value = '';
  emit('close');
}

onBeforeUnmount(() => {
  disposed = true;
  createdKey.value = '';
});
</script>

<style scoped>
.key-value :deep(input) { font-family: monospace; }
</style>
