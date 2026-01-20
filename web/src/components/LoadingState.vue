<template>
  <div class="loading-state">
    <!-- Card grid skeleton -->
    <template v-if="type === 'cards'">
      <div class="card-grid">
        <v-skeleton-loader
          v-for="i in count"
          :key="i"
          type="card"
          class="rounded-lg"
        />
      </div>
    </template>

    <!-- List skeleton -->
    <template v-else-if="type === 'list'">
      <v-skeleton-loader
        v-for="i in count"
        :key="i"
        type="list-item-two-line"
        class="mb-2"
      />
    </template>

    <!-- Table skeleton -->
    <template v-else-if="type === 'table'">
      <v-skeleton-loader type="table-heading" class="mb-2" />
      <v-skeleton-loader
        v-for="i in count"
        :key="i"
        type="table-row"
      />
    </template>

    <!-- Detail skeleton -->
    <template v-else-if="type === 'detail'">
      <v-skeleton-loader type="heading" class="mb-4" />
      <v-skeleton-loader type="paragraph" class="mb-4" />
      <v-skeleton-loader type="paragraph" />
    </template>

    <!-- Default: simple list -->
    <template v-else>
      <v-skeleton-loader
        v-for="i in count"
        :key="i"
        type="list-item"
        class="mb-2"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  type?: 'list' | 'cards' | 'table' | 'detail';
  count?: number;
}>(), {
  type: 'list',
  count: 3
});
</script>

<style scoped>
.card-grid {
  display: grid;
  gap: var(--space-4, 16px);
  grid-template-columns: 1fr;
}

@media (min-width: 600px) {
  .card-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 960px) {
  .card-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>
