<script setup>
import { Hotel, User, LogOut, History } from 'lucide-vue-next';
const userAuth = useUserAuth();
const router = useRouter();

const handleUserLogout = () => {
  userAuth.logout();
  router.push('/');
};
</script>

<template>
  <div class="min-h-screen bg-slate-50 font-sans text-slate-900">
    <nav class="sticky top-0 z-50 border-b border-slate-200 bg-white/80 backdrop-blur-md">
      <div class="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <NuxtLink to="/" class="flex items-center gap-2 text-xl font-bold tracking-tight text-blue-600">
          <Hotel class="h-6 w-6" />
          <span>Lumina</span>
        </NuxtLink>

        <div class="flex items-center gap-6">
          <NuxtLink to="/rooms" class="text-sm font-medium text-slate-600 transition-colors hover:text-blue-600">Rooms</NuxtLink>
          
          <!-- User/Guest Session Controls (Simultaneous Guard) -->
          <template v-if="userAuth.isLoggedIn">
            <NuxtLink to="/bookings" class="flex items-center gap-1 text-sm font-medium text-slate-600 hover:text-blue-600">
              <History class="h-4 w-4" />
              <span class="hidden sm:inline">My Bookings</span>
            </NuxtLink>
            <button @click="handleUserLogout" class="flex items-center gap-1 text-sm font-medium text-slate-650 hover:text-red-650" title="Sign out guest session">
              <LogOut class="h-4 w-4" />
              <span class="hidden sm:inline">Logout</span>
            </button>
          </template>
          
          <template v-else>
            <NuxtLink to="/login" class="rounded-full bg-blue-600 px-5 py-2 text-sm font-semibold text-white shadow-lg transition-all hover:bg-blue-700">
              Sign In
            </NuxtLink>
          </template>
        </div>
      </div>
    </nav>

    <main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <slot />
    </main>

    <footer class="border-t border-slate-200 bg-white py-12">
      <div class="mx-auto max-w-7xl px-4 text-center sm:flex sm:items-center sm:justify-between">
        <p class="text-sm text-slate-500">© 2026 Lumina Hotel & Resorts. Guest Portal.</p>
        <NuxtLink to="/admin/login" class="mt-4 sm:mt-0 inline-flex items-center text-xs font-bold uppercase tracking-widest text-slate-400 hover:text-blue-600 transition-colors">
          Admin Gateway
        </NuxtLink>
      </div>
    </footer>
  </div>
</template>

