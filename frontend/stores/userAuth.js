import { defineStore } from 'pinia';

export const useUserAuthStore = defineStore('userAuth', {
  state: () => ({
    user: null,
    token: null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
  },
  actions: {
    async login(email, password) {
      const config = useRuntimeConfig();
      try {
        const data = await $fetch(`${config.public.apiBase}/user/login`, {
          method: 'POST',
          body: { email, password }
        });
        
        this.token = data.token;
        this.user = data.user;
      } catch (error) {
        console.error('User login error:', error);
        throw error;
      }
    },
    async register(name, email, password) {
      const config = useRuntimeConfig();
      try {
        await $fetch(`${config.public.apiBase}/user/register`, {
          method: 'POST',
          body: { name, email, password }
        });
        await this.login(email, password);
      } catch (error) {
        console.error('User registration error:', error);
        throw error;
      }
    },
    async fetchProfile() {
      if (!this.token) return;
      const config = useRuntimeConfig();
      try {
        const data = await $fetch(`${config.public.apiBase}/user/profile`, {
          headers: {
            'Authorization': `Bearer ${this.token}`
          }
        });
        this.user = data;
      } catch (error) {
        console.error('Error fetching user profile:', error);
        this.logout();
      }
    },
    async logout() {
      const config = useRuntimeConfig();
      if (this.token) {
        try {
          await $fetch(`${config.public.apiBase}/user/logout`, {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${this.token}` }
          });
        } catch (e) {
          // ignore logout backend error
        }
      }
      this.user = null;
      this.token = null;
    }
  },
  persist: {
    key: 'user_token'
  }
});

export const useUserAuth = useUserAuthStore;
