import { createI18n } from 'vue-i18n';
import en from '@/locales/en.json';

export function createTestI18n(options: { locale?: string; messages?: Record<string, unknown> } = {}) {
  const { locale = 'en', messages = { en } } = options;
  return createI18n({
    legacy: false,
    locale,
    messages
  });
}
