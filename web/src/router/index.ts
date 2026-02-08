import { createRouter, createWebHistory } from 'vue-router';

import AreasView from '../views/AreasView.vue';
import AccessDeniedView from '../views/AccessDeniedView.vue';
import AreaPresenceView from '../views/AreaPresenceView.vue';
import RoomsView from '../views/RoomsView.vue';
import DesksView from '../views/DesksView.vue';
import MyBookingsView from '../views/MyBookingsView.vue';
import BookingHistoryView from '../views/BookingHistoryView.vue';
import RoomBookingsView from '../views/RoomBookingsView.vue';
import LoginView from '../views/LoginView.vue';

import { useAuthStore } from '../stores/useAuthStore';
import { fetchMe } from '../api/me';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { public: true }
    },
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
      path: '/bookings/history',
      name: 'booking-history',
      component: BookingHistoryView
    },
    {
      path: '/access-denied',
      name: 'access-denied',
      component: AccessDeniedView,
      meta: { public: true }
    }
  ]
});

router.beforeEach(async (to) => {
  if (to.meta.public) return true;

  const authStore = useAuthStore();

  if (!authStore.isAuthenticated) {
    try {
      const response = await fetchMe();
      authStore.setUser({
        id: response.data.id,
        display_name: response.data.attributes.display_name,
        email: response.data.attributes.email,
        is_admin: response.data.attributes.is_admin,
        auth_source: response.data.attributes.auth_source
      });
    } catch {
      return { name: 'login' };
    }
  }

  return true;
});

export default router;
