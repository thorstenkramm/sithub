import { createApp } from 'vue';
import { createPinia } from 'pinia';
import 'vuetify/styles';
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
