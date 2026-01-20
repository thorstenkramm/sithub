import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';
import { aliases, mdi } from 'vuetify/iconsets/mdi-svg';
import {
  mdiOfficeBuilding,
  mdiDoorOpen,
  mdiDesk,
  mdiCalendar,
  mdiCalendarClock,
  mdiAccount,
  mdiAccountGroup,
  mdiAccountPlus,
  mdiCheck,
  mdiClose,
  mdiAlert,
  mdiInformation,
  mdiChevronRight,
  mdiHome,
  mdiMenu,
  mdiLogout,
  mdiMagnify,
  mdiPlus,
  mdiPencil,
  mdiDelete,
  mdiMonitor,
  mdiKeyboard,
  mdiMouse,
  mdiPhone,
  mdiHeadphones,
  mdiSeatOutline,
  mdiAlertCircle,
  mdiCheckCircle,
  mdiClockOutline,
  mdiHistory,
  mdiFilterVariant,
  mdiRefresh,
  mdiMapMarker
} from '@mdi/js';

// Custom icon aliases for the app
const customAliases = {
  ...aliases,
  area: mdiOfficeBuilding,
  room: mdiDoorOpen,
  desk: mdiDesk,
  calendar: mdiCalendar,
  calendarClock: mdiCalendarClock,
  user: mdiAccount,
  users: mdiAccountGroup,
  userPlus: mdiAccountPlus,
  check: mdiCheck,
  close: mdiClose,
  alert: mdiAlert,
  info: mdiInformation,
  chevronRight: mdiChevronRight,
  home: mdiHome,
  menu: mdiMenu,
  logout: mdiLogout,
  search: mdiMagnify,
  plus: mdiPlus,
  edit: mdiPencil,
  delete: mdiDelete,
  monitor: mdiMonitor,
  keyboard: mdiKeyboard,
  mouse: mdiMouse,
  phone: mdiPhone,
  headphones: mdiHeadphones,
  chair: mdiSeatOutline,
  warning: mdiAlertCircle,
  success: mdiCheckCircle,
  clock: mdiClockOutline,
  history: mdiHistory,
  filter: mdiFilterVariant,
  refresh: mdiRefresh,
  location: mdiMapMarker
};

// Light theme colors
const lightTheme = {
  dark: false,
  colors: {
    // Primary palette - Professional blue
    primary: '#2563EB',
    'primary-darken-1': '#1D4ED8',
    'primary-lighten-1': '#3B82F6',

    // Secondary palette - Accent violet
    secondary: '#7C3AED',
    'secondary-darken-1': '#6D28D9',
    'secondary-lighten-1': '#8B5CF6',

    // Semantic colors
    success: '#059669',
    'success-darken-1': '#047857',
    'success-lighten-1': '#10B981',

    warning: '#D97706',
    'warning-darken-1': '#B45309',
    'warning-lighten-1': '#F59E0B',

    error: '#DC2626',
    'error-darken-1': '#B91C1C',
    'error-lighten-1': '#EF4444',

    info: '#0891B2',
    'info-darken-1': '#0E7490',
    'info-lighten-1': '#06B6D4',

    // Surface colors
    background: '#F8FAFC',
    surface: '#FFFFFF',
    'surface-variant': '#F1F5F9',
    'surface-bright': '#FFFFFF',

    // Text colors
    'on-background': '#1E293B',
    'on-surface': '#1E293B',
    'on-surface-variant': '#64748B',
    'on-primary': '#FFFFFF',
    'on-secondary': '#FFFFFF',
    'on-success': '#FFFFFF',
    'on-warning': '#FFFFFF',
    'on-error': '#FFFFFF',
    'on-info': '#FFFFFF',

    // Border color
    'outline': '#E2E8F0',
    'outline-variant': '#CBD5E1'
  }
};

// Dark theme colors
const darkTheme = {
  dark: true,
  colors: {
    // Primary palette
    primary: '#3B82F6',
    'primary-darken-1': '#2563EB',
    'primary-lighten-1': '#60A5FA',

    // Secondary palette
    secondary: '#8B5CF6',
    'secondary-darken-1': '#7C3AED',
    'secondary-lighten-1': '#A78BFA',

    // Semantic colors
    success: '#10B981',
    'success-darken-1': '#059669',
    'success-lighten-1': '#34D399',

    warning: '#F59E0B',
    'warning-darken-1': '#D97706',
    'warning-lighten-1': '#FBBF24',

    error: '#EF4444',
    'error-darken-1': '#DC2626',
    'error-lighten-1': '#F87171',

    info: '#06B6D4',
    'info-darken-1': '#0891B2',
    'info-lighten-1': '#22D3EE',

    // Surface colors
    background: '#0F172A',
    surface: '#1E293B',
    'surface-variant': '#334155',
    'surface-bright': '#475569',

    // Text colors
    'on-background': '#F1F5F9',
    'on-surface': '#F1F5F9',
    'on-surface-variant': '#94A3B8',
    'on-primary': '#FFFFFF',
    'on-secondary': '#FFFFFF',
    'on-success': '#FFFFFF',
    'on-warning': '#000000',
    'on-error': '#FFFFFF',
    'on-info': '#000000',

    // Border color
    'outline': '#475569',
    'outline-variant': '#334155'
  }
};

export const vuetify = createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'mdi',
    aliases: customAliases,
    sets: {
      mdi
    }
  },
  theme: {
    defaultTheme: 'light',
    themes: {
      light: lightTheme,
      dark: darkTheme
    }
  },
  defaults: {
    // Global component defaults
    VBtn: {
      variant: 'flat',
      rounded: 'lg'
    },
    VCard: {
      rounded: 'lg',
      elevation: 1
    },
    VTextField: {
      variant: 'outlined',
      density: 'comfortable',
      rounded: 'lg'
    },
    VSelect: {
      variant: 'outlined',
      density: 'comfortable',
      rounded: 'lg'
    },
    VChip: {
      rounded: 'lg'
    },
    VAlert: {
      variant: 'tonal',
      rounded: 'lg'
    },
    VDialog: {
      rounded: 'lg'
    },
    VMenu: {
      rounded: 'lg'
    },
    VList: {
      rounded: 'lg'
    },
    VListItem: {
      rounded: 'lg'
    }
  }
});

export default vuetify;
