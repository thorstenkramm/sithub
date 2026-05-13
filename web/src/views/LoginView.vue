<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card elevation="2">
          <v-card-title class="text-center pt-6 d-flex flex-column align-center">
            <img src="/sithub_logo.svg" alt="SitHub" class="login-logo mb-2" />
            <div class="text-h6">{{ $t('auth.signInTitle') }}</div>
          </v-card-title>
          <v-card-text>
            <!-- Entra ID primary action (rendered only when configured on the server) -->
            <v-btn
              v-if="entraIdAvailable"
              block
              size="large"
              variant="outlined"
              :loading="entraIdLoading"
              :disabled="entraIdLoading"
              data-cy="login-entraid"
              class="login-entraid-btn"
              @click="handleEntraIdLogin"
            >
              <template #prepend>
                <img src="/entra-id-icon.svg" alt="Entra ID" class="entra-id-icon" />
              </template>
              {{ $t('auth.signInWithEntraId') }}
            </v-btn>

            <div v-if="entraIdAvailable" class="text-center mt-3">
              <a
                href="#"
                class="text-caption text-medium-emphasis login-more-options"
                data-cy="login-toggle-local"
                @click.prevent="showLocalForm = !showLocalForm"
              >
                {{ showLocalForm ? $t('auth.lessLoginOptions') : $t('auth.moreLoginOptions') }}
              </a>
            </div>

            <v-expand-transition>
              <div v-if="!entraIdAvailable || showLocalForm">
                <v-divider v-if="entraIdAvailable" class="my-4" />
                <v-form action="/api/v1/auth/login" method="post" data-cy="login-form" @submit.prevent="handleLogin">
                  <v-text-field
                    v-model="email"
                    :label="$t('auth.email')"
                    type="email"
                    name="email"
                    autocomplete="username"
                    :error-messages="errorMessage ? [] : []"
                    data-cy="login-email"
                    class="mb-2"
                  />
                  <v-text-field
                    v-model="password"
                    :label="$t('auth.password')"
                    type="password"
                    name="password"
                    autocomplete="current-password"
                    data-cy="login-password"
                    class="mb-2"
                  />
                  <v-alert
                    v-if="errorMessage"
                    type="error"
                    variant="tonal"
                    density="compact"
                    class="mb-4"
                    data-cy="login-error"
                  >
                    {{ errorMessage }}
                  </v-alert>
                  <v-btn
                    type="submit"
                    color="primary"
                    block
                    :loading="loading"
                    data-cy="login-submit"
                  >
                    {{ $t('auth.signIn') }}
                  </v-btn>
                </v-form>
              </div>
            </v-expand-transition>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { nextTick, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { fetchAuthProviders, loginLocal } from '../api/auth';
import { useAuthStore } from '../stores/useAuthStore';
import { ApiError, isConnectionError, CONNECTION_LOST_MESSAGE } from '../api/client';

const router = useRouter();
const authStore = useAuthStore();

const email = ref('');
const password = ref('');
const loading = ref(false);
const entraIdLoading = ref(false);
const errorMessage = ref('');
const entraIdAvailable = ref(false);
const showLocalForm = ref(false);
const { t } = useI18n();

onMounted(async () => {
  try {
    const resp = await fetchAuthProviders();
    entraIdAvailable.value = resp.data.attributes.entraid;
    // When Entra ID is unavailable, show the local form by default so users
    // are not locked out. When it is available, keep the local form collapsed
    // behind the "more login options" link.
    showLocalForm.value = !entraIdAvailable.value;
  } catch {
    // If the providers endpoint fails (older server, network error, etc.)
    // fall back to showing both options so users can still authenticate.
    entraIdAvailable.value = true;
    showLocalForm.value = true;
  }
});

function getLoginErrorMessage(err: ApiError): string {
  if (err.status === 401) {
    return t('auth.invalidCredentials');
  }
  if (err.status === 400) {
    return t('auth.requiredFields');
  }
  return t('auth.genericError');
}

async function handleLogin() {
  errorMessage.value = '';
  loading.value = true;
  try {
    const response = await loginLocal(email.value, password.value);
    authStore.setUser({
      id: response.data.id,
      display_name: response.data.attributes.display_name,
      email: response.data.attributes.email,
      is_admin: response.data.attributes.is_admin,
      auth_source: response.data.attributes.auth_source
    });
    router.push('/');
  } catch (err) {
    if (isConnectionError(err)) {
      errorMessage.value = CONNECTION_LOST_MESSAGE;
    } else if (err instanceof ApiError) {
      errorMessage.value = getLoginErrorMessage(err);
    } else {
      errorMessage.value = t('auth.genericError');
    }
  } finally {
    loading.value = false;
  }
}

async function handleEntraIdLogin() {
  entraIdLoading.value = true;
  await nextTick();
  await new Promise<void>((resolve) => {
    if (typeof window.requestAnimationFrame === 'function') {
      window.requestAnimationFrame(() => resolve());
      return;
    }
    window.setTimeout(resolve, 0);
  });
  window.location.assign('/oauth/login');
}
</script>

<style scoped>
.login-logo {
  max-width: 220px;
  height: auto;
}

.login-entraid-btn {
  text-transform: none;
  font-weight: 500;
}

.entra-id-icon {
  width: 20px;
  height: 20px;
  display: inline-block;
}

.login-more-options {
  text-decoration: none;
  cursor: pointer;
}

.login-more-options:hover {
  text-decoration: underline;
}
</style>
