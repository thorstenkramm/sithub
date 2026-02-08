import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useAuthStore = defineStore('auth', () => {
  const userName = ref('');
  const email = ref('');
  const userId = ref('');
  const isAdmin = ref(false);
  const authSource = ref('');
  const isAuthenticated = ref(false);

  function setUser(data: {
    id: string;
    display_name: string;
    email: string;
    is_admin: boolean;
    auth_source: string;
  }) {
    userId.value = data.id;
    userName.value = data.display_name;
    email.value = data.email;
    isAdmin.value = data.is_admin;
    authSource.value = data.auth_source;
    isAuthenticated.value = true;
  }

  function clearUser() {
    userId.value = '';
    userName.value = '';
    email.value = '';
    isAdmin.value = false;
    authSource.value = '';
    isAuthenticated.value = false;
  }

  return { userName, email, userId, isAdmin, authSource, isAuthenticated, setUser, clearUser };
});
