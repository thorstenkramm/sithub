import { createRouter, createWebHistory } from 'vue-router';

import AreasView from '../views/AreasView.vue';
import AccessDeniedView from '../views/AccessDeniedView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'areas',
      component: AreasView
    },
    {
      path: '/access-denied',
      name: 'access-denied',
      component: AccessDeniedView
    }
  ]
});

export default router;
