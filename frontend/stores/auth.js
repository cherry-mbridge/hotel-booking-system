import { defineStore } from 'pinia';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin',
  },
  actions: {
    async login(email, password) {
      const config = useRuntimeConfig();
      try {
        const data = await $fetch(`${config.public.apiBase}/login`, {
          method: 'POST',
          body: { email, password }
        });
        
        this.token = data.token;
        this.user = data.user;
        
        // Store token in cookie/localStorage is handled by pinia-plugin-persistedstate
      } catch (error) {
        console.error('Login error:', error);
        throw error;
      }
    },
    async register(name, email, password) {
      const config = useRuntimeConfig();
      try {
        await $fetch(`${config.public.apiBase}/register`, {
          method: 'POST',
          body: { name, email, password }
        });
        // After registration, usually we log them in
        await this.login(email, password);
      } catch (error) {
        console.error('Registration error:', error);
        throw error;
      }
    },
    logout() {
      this.user = null;
      this.token = null;
    }
  },
  persist: true
});
