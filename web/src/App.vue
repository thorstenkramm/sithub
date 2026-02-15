<template>
  <v-app>
    <template v-if="authStore.isAuthenticated">
      <v-app-bar color="primary" density="comfortable" elevation="0">
        <router-link to="/" class="d-flex align-center text-decoration-none ml-2">
          <img src="/logo.svg" alt="SitHub" height="28" class="logo-image" />
        </router-link>

        <v-spacer />

        <!-- Desktop Navigation -->
        <nav class="d-none d-md-flex align-center ga-1">
          <v-btn
            to="/"
            variant="text"
            :class="{ 'nav-active': isRouteActive('/') }"
            data-cy="nav-areas"
          >
            Areas
          </v-btn>
          <v-btn
            to="/my-bookings"
            variant="text"
            :class="{ 'nav-active': isRouteActive('/my-bookings') }"
            data-cy="nav-my-bookings"
          >
            My Bookings
          </v-btn>
          <v-btn
            to="/bookings/history"
            variant="text"
            :class="{ 'nav-active': isRouteActive('/bookings/history') }"
            data-cy="nav-booking-history"
          >
            History
          </v-btn>
        </nav>

        <!-- User Menu -->
        <v-menu location="bottom end" :offset="8">
          <template #activator="{ props }">
            <v-btn
              v-bind="props"
              variant="text"
              class="ml-2"
              data-cy="user-menu-trigger"
            >
              <v-avatar size="32" color="primary-lighten-1" class="mr-2">
                <span class="text-body-2 font-weight-medium">{{ userInitials }}</span>
              </v-avatar>
              <span class="d-none d-sm-inline">{{ authStore.userName }}</span>
            </v-btn>
          </template>
          <v-list density="compact" min-width="200">
            <v-list-item>
              <v-list-item-title class="font-weight-medium">{{ authStore.userName }}</v-list-item-title>
              <v-list-item-subtitle v-if="authStore.isAdmin">
                <v-chip size="x-small" color="secondary" class="mt-1">Admin</v-chip>
              </v-list-item-subtitle>
            </v-list-item>
            <v-divider class="my-1" />
            <v-list-item data-cy="theme-selector">
              <v-list-item-title class="text-caption text-medium-emphasis mb-1">Theme</v-list-item-title>
              <v-btn-toggle
                :model-value="themePreference"
                mandatory
                density="compact"
                color="primary"
                @update:model-value="setThemePreference($event)"
              >
                <v-btn
                  v-for="opt in themeOptions"
                  :key="opt.value"
                  :value="opt.value"
                  size="small"
                >
                  {{ opt.label }}
                </v-btn>
              </v-btn-toggle>
            </v-list-item>
            <v-list-item data-cy="show-weekends-toggle">
              <v-checkbox
                v-model="showWeekends"
                label="Show weekends"
                hide-details
                density="compact"
              />
            </v-list-item>
            <v-divider class="my-1" />
            <v-list-item
              v-if="authStore.authSource === 'internal'"
              data-cy="change-password-btn"
              @click="showPasswordDialog = true"
            >
              <template #prepend>
                <v-icon size="small" data-cy="change-password-icon">$lockReset</v-icon>
              </template>
              <v-list-item-title>Change Password</v-list-item-title>
            </v-list-item>
            <v-list-item data-cy="logout-btn" @click="handleLogout">
              <template #prepend>
                <v-icon size="small">$logout</v-icon>
              </template>
              <v-list-item-title>Sign out</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>

        <!-- Mobile Menu Button -->
        <v-app-bar-nav-icon
          class="d-md-none"
          @click="mobileDrawer = true"
          data-cy="mobile-menu-btn"
        />
      </v-app-bar>

      <!-- Mobile Navigation Drawer -->
      <v-navigation-drawer v-model="mobileDrawer" temporary location="right">
        <v-list nav>
          <v-list-item>
            <v-list-item-title class="font-weight-bold">{{ authStore.userName }}</v-list-item-title>
            <v-list-item-subtitle v-if="authStore.isAdmin">
              <v-chip size="x-small" color="secondary" class="mt-1">Admin</v-chip>
            </v-list-item-subtitle>
          </v-list-item>
          <v-divider class="my-2" />
          <v-list-item to="/" @click="mobileDrawer = false" data-cy="mobile-nav-areas">
            <template #prepend>
              <v-icon>$area</v-icon>
            </template>
            <v-list-item-title>Areas</v-list-item-title>
          </v-list-item>
          <v-list-item to="/my-bookings" @click="mobileDrawer = false" data-cy="mobile-nav-my-bookings">
            <template #prepend>
              <v-icon>$calendar</v-icon>
            </template>
            <v-list-item-title>My Bookings</v-list-item-title>
          </v-list-item>
          <v-list-item to="/bookings/history" @click="mobileDrawer = false" data-cy="mobile-nav-history">
            <template #prepend>
              <v-icon>$history</v-icon>
            </template>
            <v-list-item-title>History</v-list-item-title>
          </v-list-item>
          <v-divider class="my-2" />
          <v-list-item data-cy="mobile-theme-selector">
            <v-list-item-title class="text-caption text-medium-emphasis mb-1">Theme</v-list-item-title>
            <v-btn-toggle
              :model-value="themePreference"
              mandatory
              density="compact"
              color="primary"
              @update:model-value="setThemePreference($event)"
            >
              <v-btn
                v-for="opt in themeOptions"
                :key="opt.value"
                :value="opt.value"
                size="small"
              >
                {{ opt.label }}
              </v-btn>
            </v-btn-toggle>
          </v-list-item>
          <v-list-item data-cy="mobile-show-weekends-toggle">
            <v-checkbox
              v-model="showWeekends"
              label="Show weekends"
              hide-details
              density="compact"
            />
          </v-list-item>
          <v-divider class="my-2" />
          <v-list-item
            v-if="authStore.authSource === 'internal'"
            data-cy="mobile-change-password-btn"
            @click="showPasswordDialog = true; mobileDrawer = false"
          >
            <template #prepend>
              <v-icon data-cy="mobile-change-password-icon">$lockReset</v-icon>
            </template>
            <v-list-item-title>Change Password</v-list-item-title>
          </v-list-item>
          <v-list-item data-cy="mobile-logout-btn" @click="handleLogout">
            <template #prepend>
              <v-icon>$logout</v-icon>
            </template>
            <v-list-item-title>Sign out</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-navigation-drawer>

      <!-- Password Change Dialog -->
      <v-dialog v-model="showPasswordDialog" max-width="400" persistent>
        <v-card>
          <v-card-title>Change Password</v-card-title>
          <v-card-text>
            <v-text-field
              v-model="currentPassword"
              label="Current Password"
              type="password"
              autocomplete="current-password"
              data-cy="current-password"
              class="mb-2"
            />
            <v-text-field
              v-model="newPassword"
              label="New Password"
              type="password"
              autocomplete="new-password"
              hint="Minimum 14 characters"
              data-cy="new-password"
              class="mb-2"
            />
            <v-alert
              v-if="passwordError"
              type="error"
              variant="tonal"
              density="compact"
              class="mb-2"
              data-cy="password-error"
            >
              {{ passwordError }}
            </v-alert>
            <v-alert
              v-if="passwordSuccess"
              type="success"
              variant="tonal"
              density="compact"
              class="mb-2"
              data-cy="password-success"
            >
              Password changed successfully.
            </v-alert>
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn variant="text" @click="closePasswordDialog" data-cy="password-cancel">Cancel</v-btn>
            <v-btn
              color="primary"
              :loading="passwordLoading"
              :disabled="passwordSuccess"
              data-cy="password-submit"
              @click="handlePasswordChange"
            >
              Change
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </template>

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from './stores/useAuthStore';
import { logout } from './api/auth';
import { changePassword } from './api/me';
import { ApiError } from './api/client';
import { useThemePreference } from './composables/useThemePreference';
import { useWeekendPreference } from './composables/useWeekendPreference';

const themeOptions = [
  { label: 'Auto', value: 'auto' as const },
  { label: 'Light', value: 'light' as const },
  { label: 'Dark', value: 'dark' as const }
];

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const mobileDrawer = ref(false);
const { preference: themePreference, setPreference: setThemePreference } = useThemePreference();
const { showWeekends } = useWeekendPreference();

const showPasswordDialog = ref(false);
const currentPassword = ref('');
const newPassword = ref('');
const passwordError = ref('');
const passwordSuccess = ref(false);
const passwordLoading = ref(false);

const userInitials = computed(() => {
  const name = authStore.userName || 'U';
  const parts = name.split(' ');
  const first = parts[0];
  const last = parts[parts.length - 1];
  if (parts.length >= 2 && first && last) {
    return (first.charAt(0) + last.charAt(0)).toUpperCase();
  }
  return name.substring(0, 2).toUpperCase();
});

function isRouteActive(path: string): boolean {
  if (path === '/') {
    return route.path === '/' || route.path.startsWith('/areas') || route.path.startsWith('/item-groups');
  }
  return route.path.startsWith(path);
}

async function handleLogout() {
  await logout();
  authStore.clearUser();
  router.push('/login');
}

async function handlePasswordChange() {
  passwordError.value = '';
  passwordSuccess.value = false;
  passwordLoading.value = true;
  try {
    await changePassword(currentPassword.value, newPassword.value);
    passwordSuccess.value = true;
  } catch (err) {
    if (err instanceof ApiError) {
      passwordError.value = err.detail || 'Failed to change password';
    } else {
      passwordError.value = 'An error occurred. Please try again.';
    }
  } finally {
    passwordLoading.value = false;
  }
}

function closePasswordDialog() {
  showPasswordDialog.value = false;
  currentPassword.value = '';
  newPassword.value = '';
  passwordError.value = '';
  passwordSuccess.value = false;
}
</script>

<style scoped>
.logo-image {
  filter: brightness(0) invert(1);
}

.nav-active {
  background-color: rgba(255, 255, 255, 0.15) !important;
}

.nav-active::before {
  opacity: 0 !important;
}
</style>
