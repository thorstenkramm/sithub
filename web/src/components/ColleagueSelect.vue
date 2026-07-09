<script setup lang="ts">
/**
 * ColleagueSelect is the shared "book for myself / for a colleague" fragment
 * used by every booking dialog (tiles, floor plan, weekly table pattern). It
 * exposes a v-model carrying the selected colleague id (null = for myself) and
 * emits nothing else. Colleague loading and name resolution live in
 * useColleagues so the pattern stays DRY across surfaces.
 */
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useColleagues } from '../composables/useColleagues';

const modelValue = defineModel<string | null>({ default: null });

const props = withDefaults(
  defineProps<{
    /** Prefix for the two radio data-cy hooks, e.g. "tile" -> tile-book-self-radio. */
    dataCyPrefix?: string;
  }>(),
  { dataCyPrefix: 'colleague' }
);

useI18n(); // expose $t in template

const { colleagueList, colleaguesLoading, loadColleagues } = useColleagues();

const bookingType = ref<'self' | 'colleague'>(modelValue.value ? 'colleague' : 'self');

// Switching to colleague mode lazily loads the list; switching back to self
// clears any selection so the booking defaults to the current user.
watch(bookingType, (type) => {
  if (type === 'colleague') {
    loadColleagues();
  } else {
    modelValue.value = null;
  }
});

// Load eagerly so the autocomplete is ready when the dialog is already in
// colleague mode (e.g. reopened with a preselection).
if (bookingType.value === 'colleague') {
  loadColleagues();
}
</script>

<template>
  <div :data-cy="`${props.dataCyPrefix}-colleague-select`">
    <v-radio-group
      v-model="bookingType"
      inline
      density="compact"
      hide-details
      class="mb-1"
    >
      <v-radio
        :label="$t('items.bookForMyself')"
        value="self"
        :data-cy="`${props.dataCyPrefix}-book-self-radio`"
      />
      <v-radio
        :label="$t('items.bookForColleague')"
        value="colleague"
        :data-cy="`${props.dataCyPrefix}-book-colleague-radio`"
      />
    </v-radio-group>

    <v-expand-transition>
      <v-autocomplete
        v-if="bookingType === 'colleague'"
        v-model="modelValue"
        :items="colleagueList"
        item-title="displayName"
        item-value="id"
        :label="$t('items.selectColleague')"
        density="compact"
        :loading="colleaguesLoading"
        clearable
        hide-details
        :data-cy="`${props.dataCyPrefix}-colleague-autocomplete`"
      />
    </v-expand-transition>
  </div>
</template>
