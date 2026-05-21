export default defineNuxtRouteMiddleware((to, from) => {
  const adminAuth = useAdminAuth();
  if (!adminAuth.isLoggedIn) {
    return navigateTo('/admin/login');
  }
});
