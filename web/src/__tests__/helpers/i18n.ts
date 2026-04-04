import { createI18n } from 'vue-i18n';
import en from '@/locales/en.json';

export function createTestI18n() {
  return createI18n({
    legacy: false,
    locale: 'en',
    messages: { en }
  });
}
