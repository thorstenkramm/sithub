<template>
  <div class="page-header mb-6">
    <!-- Breadcrumbs -->
    <nav v-if="breadcrumbs?.length" class="breadcrumbs mb-2" aria-label="Breadcrumb" data-cy="breadcrumbs">
      <ol class="d-flex align-center ga-1 text-body-2">
        <li
          v-for="(crumb, index) in breadcrumbs"
          :key="index"
          class="d-flex align-center"
          :data-cy="`breadcrumb-item-${index}`"
        >
          <router-link
            v-if="crumb.to && index < breadcrumbs.length - 1"
            :to="crumb.to"
            class="breadcrumb-link text-primary"
            :data-cy="`breadcrumb-link-${index}`"
          >
            {{ crumb.text }}
          </router-link>
          <span v-else class="text-medium-emphasis" :data-cy="`breadcrumb-text-${index}`">
            {{ crumb.text }}
          </span>
          <v-icon
            v-if="index < breadcrumbs.length - 1"
            size="small"
            class="mx-1 text-disabled"
          >
            $chevronRight
          </v-icon>
        </li>
      </ol>
    </nav>

    <!-- Title row -->
    <div class="d-flex align-center justify-space-between flex-wrap ga-4">
      <div>
        <h1 v-if="title" class="text-h4 font-weight-bold mb-1">{{ title }}</h1>
        <p v-if="subtitle" class="text-body-1 text-medium-emphasis ma-0">{{ subtitle }}</p>
      </div>
      <div v-if="$slots.actions" class="d-flex align-center ga-2">
        <slot name="actions" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { RouteLocationRaw } from 'vue-router';

export interface BreadcrumbItem {
  text: string;
  to?: RouteLocationRaw;
}

defineProps<{
  title: string;
  subtitle?: string;
  breadcrumbs?: BreadcrumbItem[];
}>();
</script>

<style scoped>
.breadcrumbs ol {
  list-style: none;
  padding: 0;
  margin: 0;
}

.breadcrumb-link {
  text-decoration: none;
}

.breadcrumb-link:hover {
  text-decoration: underline;
}
</style>
