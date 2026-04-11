<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card elevation="2">
          <v-card-title class="text-center pt-6">
            <img src="/logo.svg" alt="SitHub" height="40" class="mb-2" />
            <div class="text-h6">{{ $t('auth.signInTitle') }}</div>
          </v-card-title>
          <v-card-text>
            <v-form action="/api/v1/auth/login" method="post" @submit.prevent="handleLogin" data-cy="login-form">
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
            <v-divider class="my-4" />
            <v-btn
              variant="outlined"
              block
              :loading="entraIdLoading"
              :disabled="entraIdLoading"
              data-cy="login-entraid"
              @click="handleEntraIdLogin"
            >
              {{ $t('auth.signInWithEntraId') }}
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { loginLocal } from '../api/auth';
import { useAuthStore } from '../stores/useAuthStore';
import { ApiError, isConnectionError, CONNECTION_LOST_MESSAGE } from '../api/client';

const router = useRouter();
const authStore = useAuthStore();

const email = ref('');
const password = ref('');
const loading = ref(false);
const entraIdLoading = ref(false);
const errorMessage = ref('');
const { t } = useI18n();

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

function handleEntraIdLogin() {
  entraIdLoading.value = true;
  window.location.href = '/oauth/login';
}
</script>
