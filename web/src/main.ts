import { createApp } from 'vue';
import { createPinia } from 'pinia';
import { createVuetify } from 'vuetify';
import 'vuetify/styles';

import App from './App.vue';
import router from './router';

const app = createApp(App);
const pinia = createPinia();
const vuetify = createVuetify();

app.use(pinia);
app.use(router);
app.use(vuetify);

app.mount('#app');
