import { createRouter, createWebHistory } from 'vue-router';

import { useAuthStore } from '../stores/useAuthStore';
import { fetchMe } from '../api/me';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { public: true }
    },
    {
      path: '/',
      name: 'areas',
      component: () => import('../views/AreasView.vue')
    },
    {
      path: '/areas/:areaId/item-groups',
      name: 'item-groups',
      component: () => import('../views/ItemGroupsView.vue')
    },
    {
      // Favorites is a virtual room: it reuses ItemsView so users get the
      // exact same day/week toggle, date picker, equipment filter, and
      // booking flow as a real room. ItemsView branches on
      // `route.meta.favoritesMode` to skip the per-item-group fetch and
      // aggregate items across all favorited (areaId, itemGroupId) pairs.
      path: '/favorites',
      name: 'favorites',
      component: () => import('../views/ItemsView.vue'),
      meta: { favoritesMode: true }
    },
    {
      path: '/areas/:areaId/presence',
      name: 'area-presence',
      component: () => import('../views/AreaPresenceView.vue')
    },
    {
      path: '/item-groups/:itemGroupId/items',
      name: 'items',
      component: () => import('../views/ItemsView.vue')
    },
    {
      path: '/item-groups/:itemGroupId/bookings',
      name: 'item-group-bookings',
      component: () => import('../views/ItemGroupBookingsView.vue')
    },
    {
      path: '/my-bookings',
      name: 'my-bookings',
      component: () => import('../views/MyBookingsView.vue')
    },
    {
      path: '/bookings/history',
      name: 'booking-history',
      component: () => import('../views/BookingHistoryView.vue')
    },
    {
      path: '/admin/floor-plan-editor',
      name: 'floor-plan-editor',
      component: () => import('../views/FloorPlanEditorView.vue'),
      meta: { requiresAdmin: true }
    },
    {
      path: '/access-denied',
      name: 'access-denied',
      component: () => import('../views/AccessDeniedView.vue'),
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

  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    return { name: 'access-denied' };
  }

  return true;
});

export default router;
