<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card elevation="2">
          <v-card-title class="text-center pt-6">
            <img src="/logo.svg" alt="SitHub" height="40" class="mb-2" />
            <div class="text-h6">Sign in to SitHub</div>
          </v-card-title>
          <v-card-text>
            <v-form @submit.prevent="handleLogin" data-cy="login-form">
              <v-text-field
                v-model="email"
                label="Email"
                type="email"
                autocomplete="email"
                :error-messages="errorMessage ? [] : []"
                data-cy="login-email"
                class="mb-2"
              />
              <v-text-field
                v-model="password"
                label="Password"
                type="password"
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
                Sign in
              </v-btn>
            </v-form>
            <v-divider class="my-4" />
            <v-btn
              href="/oauth/login"
              variant="outlined"
              block
              data-cy="login-entraid"
            >
              Sign in with Entra ID
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { loginLocal } from '../api/auth';
import { useAuthStore } from '../stores/useAuthStore';
import { ApiError } from '../api/client';

const router = useRouter();
const authStore = useAuthStore();

const email = ref('');
const password = ref('');
const loading = ref(false);
const errorMessage = ref('');

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
    if (err instanceof ApiError) {
      errorMessage.value = err.detail || 'Invalid email or password';
    } else {
      errorMessage.value = 'An error occurred. Please try again.';
    }
  } finally {
    loading.value = false;
  }
}
</script>
