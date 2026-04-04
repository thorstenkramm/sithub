<script setup lang="ts">
</script>

<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="8" lg="6">
        <v-card>
          <v-card-title data-cy="access-denied-title">{{ $t('accessDenied.title') }}</v-card-title>
          <v-card-text>
            {{ $t('accessDenied.message') }}
          </v-card-text>
          <v-card-actions>
            <v-btn data-cy="access-denied-login" href="/login" variant="tonal">
              {{ $t('accessDenied.backToLogin') }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
