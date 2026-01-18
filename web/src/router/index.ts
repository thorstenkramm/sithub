import { createRouter, createWebHistory } from 'vue-router';

import AreasView from '../views/AreasView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'areas',
      component: AreasView
    }
  ]
});

export default router;
