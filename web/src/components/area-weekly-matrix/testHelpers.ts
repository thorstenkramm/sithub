/** Shared Vuetify stubs for matrix popover tests. */
export const popoverStubs = {
  'v-menu': {
    template: '<div v-if="modelValue" v-bind="$attrs"><slot /></div>',
    props: ['modelValue', 'activator', 'closeOnContentClick', 'location', 'maxWidth']
  },
  'v-card': { template: '<div v-bind="$attrs"><slot /></div>' },
  'v-card-title': { template: '<div v-bind="$attrs"><slot /></div>' },
  'v-card-actions': { template: '<div v-bind="$attrs"><slot /></div>' },
  'v-alert': { template: '<div v-bind="$attrs"><slot /><slot name="append" /></div>' },
  'v-btn': {
    template: '<button v-bind="$attrs" @click="$emit(\'click\', $event)"><slot /></button>',
    emits: ['click'],
    props: ['variant', 'color', 'size', 'loading']
  },
  'v-spacer': { template: '<div />' }
};
