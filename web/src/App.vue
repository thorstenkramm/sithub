<template>
  <v-app>
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
          <v-list-item href="/oauth/logout" data-cy="logout-btn">
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
        <v-list-item href="/oauth/logout" data-cy="mobile-logout-btn">
          <template #prepend>
            <v-icon>$logout</v-icon>
          </template>
          <v-list-item-title>Sign out</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from './stores/useAuthStore';

const route = useRoute();
const authStore = useAuthStore();
const mobileDrawer = ref(false);

const userInitials = computed(() => {
  const name = authStore.userName || 'U';
  const parts = name.split(' ');
  if (parts.length >= 2) {
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
  }
  return name.substring(0, 2).toUpperCase();
});

function isRouteActive(path: string): boolean {
  if (path === '/') {
    return route.path === '/' || route.path.startsWith('/areas') || route.path.startsWith('/rooms');
  }
  return route.path.startsWith(path);
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
