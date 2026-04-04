import { createI18n } from 'vue-i18n';
import { en as vuetifyEn, de as vuetifyDe, es as vuetifyEs, fr as vuetifyFr, uk as vuetifyUk } from 'vuetify/locale';
import en from '@/locales/en.json';
import de from '@/locales/de.json';
import es from '@/locales/es.json';
import fr from '@/locales/fr.json';
import uk from '@/locales/uk.json';

export const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: { ...en, $vuetify: vuetifyEn },
    de: { ...de, $vuetify: vuetifyDe },
    es: { ...es, $vuetify: vuetifyEs },
    fr: { ...fr, $vuetify: vuetifyFr },
    uk: { ...uk, $vuetify: vuetifyUk }
  }
});
