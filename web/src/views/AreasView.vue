<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title data-cy="areas-title">
            Areas
            <span v-if="userName" class="text-caption ml-2">(Signed in as {{ userName }})</span>
          </v-card-title>
        <v-card-text>
          Area list will render here.
            <div v-if="isAdmin" class="mt-2">
              <div class="text-caption">Admin-only cancellation controls</div>
              <v-btn data-cy="admin-cancel" size="small" variant="tonal">Cancel booking (admin)</v-btn>
            </div>
        </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ApiError } from '../api/client';
import { fetchMe } from '../api/me';

const userName = ref('');
const isAdmin = ref(false);

onMounted(async () => {
  try {
    const resp = await fetchMe();
    userName.value = resp.data.attributes.display_name;
    isAdmin.value = resp.data.attributes.is_admin;
  } catch (err) {
    if (err instanceof ApiError && err.status === 401) {
      window.location.href = '/oauth/login';
      return;
    }
    throw err;
  }
});
</script>
