import { createRouter, createWebHistory } from 'vue-router';

import AreasView from '../views/AreasView.vue';
import AccessDeniedView from '../views/AccessDeniedView.vue';
import AreaPresenceView from '../views/AreaPresenceView.vue';
import RoomsView from '../views/RoomsView.vue';
import DesksView from '../views/DesksView.vue';
import MyBookingsView from '../views/MyBookingsView.vue';
import RoomBookingsView from '../views/RoomBookingsView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'areas',
      component: AreasView
    },
    {
      path: '/areas/:areaId/rooms',
      name: 'rooms',
      component: RoomsView
    },
    {
      path: '/areas/:areaId/presence',
      name: 'area-presence',
      component: AreaPresenceView
    },
    {
      path: '/rooms/:roomId/desks',
      name: 'desks',
      component: DesksView
    },
    {
      path: '/rooms/:roomId/bookings',
      name: 'room-bookings',
      component: RoomBookingsView
    },
    {
      path: '/my-bookings',
      name: 'my-bookings',
      component: MyBookingsView
    },
    {
      path: '/access-denied',
      name: 'access-denied',
      component: AccessDeniedView
    }
  ]
});

export default router;
