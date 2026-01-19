<template>
  <v-app>
    <v-app-bar color="primary" density="compact">
      <v-app-bar-title>
        <router-link to="/" class="text-white text-decoration-none">SitHub</router-link>
      </v-app-bar-title>
      <v-btn to="/my-bookings" variant="text" data-cy="nav-my-bookings">My Bookings</v-btn>
    </v-app-bar>
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>
