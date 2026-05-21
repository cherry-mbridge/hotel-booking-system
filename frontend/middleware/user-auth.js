export default defineNuxtRouteMiddleware((to, from) => {
  const userAuth = useUserAuth();
  if (!userAuth.isLoggedIn) {
     return navigateTo('/login');
  }
});
