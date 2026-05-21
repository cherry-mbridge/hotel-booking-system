import { defineStore } from 'pinia';

export const useAdminAuthStore = defineStore('adminAuth', {
  state: () => ({
    admin: null,
    token: null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
  },
  actions: {
    async login(email, password) {
      const config = useRuntimeConfig();
      try {
        const data = await $fetch(`${config.public.apiBase}/admin/login`, {
          method: 'POST',
          body: { email, password }
        });
        
        this.token = data.token;
        this.admin = data.user;
      } catch (error) {
        console.error('Admin login error:', error);
        throw error;
      }
    },
    async fetchProfile() {
      if (!this.token) return;
      const config = useRuntimeConfig();
      try {
        const data = await $fetch(`${config.public.apiBase}/admin/profile`, {
          headers: {
            'Authorization': `Bearer ${this.token}`
          }
        });
        this.admin = data;
      } catch (error) {
        console.error('Error fetching admin profile:', error);
        this.logout();
      }
    },
    async logout() {
      const config = useRuntimeConfig();
      if (this.token) {
        try {
          await $fetch(`${config.public.apiBase}/admin/logout`, {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${this.token}` }
          });
        } catch (e) {
          // ignore logout backend error
        }
      }
      this.admin = null;
      this.token = null;
    }
  },
  persist: {
    key: 'admin_token'
  }
});

export const useAdminAuth = useAdminAuthStore;
