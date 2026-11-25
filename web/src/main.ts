import { createApp } from 'vue';
import { createPinia } from 'pinia';
import 'vuetify/styles';
import '@fontsource/inter/400.css';
import '@fontsource/inter/500.css';
import '@fontsource/inter/600.css';
import '@fontsource/inter/700.css';
import '@mdi/font/css/materialdesignicons.css';

import App from './App.vue';
import router from './router';
import { i18n } from './plugins/i18n';
import { vuetify } from './plugins/vuetify';
import './styles/global.css';

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(i18n);
app.use(router);
app.use(vuetify);

app.mount('#app');
