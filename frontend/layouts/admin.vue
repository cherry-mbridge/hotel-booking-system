<script setup>
import { LayoutDashboard, Bed, CalendarCheck, LogOut, ShieldAlert, Hotel, Percent } from 'lucide-vue-next';

const adminAuth = useAdminAuth();
const router = useRouter();

const handleLogout = async () => {
  await adminAuth.logout();
  router.push('/admin/login');
};
</script>

<template>
  <div class="min-h-screen bg-slate-900 font-sans text-slate-100 flex flex-col md:flex-row">
    <!-- Sidebar Navigation -->
    <aside class="w-full md:w-64 bg-slate-950 border-b md:border-b-0 md:border-r border-slate-800 flex flex-col shrink-0">
      <!-- Title branding -->
      <div class="p-6 border-b border-slate-800 flex items-center justify-between">
        <NuxtLink to="/admin" class="flex items-center gap-2.5 text-xl font-black tracking-tight text-blue-400">
          <ShieldAlert class="h-6 w-6 text-blue-500" />
          <span>Lumina Admin</span>
        </NuxtLink>
      </div>

      <!-- Navigation Links -->
      <nav class="p-4 space-y-1.5 flex-1">
        <NuxtLink to="/admin" class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold transition-all hover:bg-slate-900 text-slate-300 hover:text-white" active-class="bg-blue-600/10 text-blue-400 border border-blue-500/20 font-bold">
          <LayoutDashboard class="h-5 w-5" />
          <span>Overview</span>
        </NuxtLink>

        <NuxtLink to="/admin/rooms" class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold transition-all hover:bg-slate-900 text-slate-300 hover:text-white" active-class="bg-blue-600/10 text-blue-400 border border-blue-500/20 font-bold">
          <Bed class="h-5 w-5" />
          <span>Room Inventory</span>
        </NuxtLink>

        <NuxtLink to="/admin/bookings" class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold transition-all hover:bg-slate-900 text-slate-300 hover:text-white" active-class="bg-blue-600/10 text-blue-400 border border-blue-500/20 font-bold">
          <CalendarCheck class="h-5 w-5" />
          <span>Guest Bookings</span>
        </NuxtLink>

        <NuxtLink to="/admin/weekend-pricing" class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold transition-all hover:bg-slate-900 text-slate-300 hover:text-white" active-class="bg-blue-600/10 text-blue-400 border border-blue-500/20 font-bold">
          <Percent class="h-5 w-5" />
          <span>Weekend Pricing</span>
        </NuxtLink>
      </nav>

      <!-- Admin profile summary & logout -->
      <div class="p-4 border-t border-slate-800 bg-slate-950/60 flex flex-col gap-3">
        <div class="flex items-center gap-3 px-2">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 text-sm font-extrabold text-white uppercase shadow-md">
            {{ adminAuth.admin?.name?.charAt(0) || 'A' }}
          </div>
          <div class="min-w-0 flex-1">
            <h4 class="text-sm font-bold truncate text-slate-200">{{ adminAuth.admin?.name || 'Administrator' }}</h4>
            <p class="text-xs text-slate-500 truncate">{{ adminAuth.admin?.email || 'admin@lumina.com' }}</p>
          </div>
        </div>

        <button @click="handleLogout" class="w-full flex items-center justify-center gap-2 px-4 py-3 rounded-xl text-sm font-bold text-red-400 border border-red-500/10 hover:bg-red-500/10 transition-colors">
          <LogOut class="h-4 w-4" />
          <span>Sign Out Admin</span>
        </button>
      </div>
    </aside>

    <!-- Main administrative content space -->
    <div class="flex-1 flex flex-col min-w-0">
      <header class="h-16 bg-slate-950 border-b border-slate-800 flex items-center justify-between px-6 sm:px-8 shrink-0">
        <div class="text-sm font-semibold text-slate-400 flex items-center gap-2">
          <span>Security Context:</span>
          <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-bold bg-blue-500/10 text-blue-400 border border-blue-500/20">ADMIN GUARD</span>
        </div>
        <div>
          <NuxtLink to="/" class="flex items-center gap-1.5 text-xs font-bold uppercase tracking-wider text-slate-400 hover:text-white transition-colors">
            <Hotel class="h-4 w-4" />
            <span>Go to Guest Site &rarr;</span>
          </NuxtLink>
        </div>
      </header>

      <main class="p-6 sm:p-8 flex-1 overflow-y-auto">
        <slot />
      </main>
    </div>
  </div>
</template>

<style scoped>
/* Scoped overrides to enforce high contrast dark colors for administrative space */
body {
  background-color: #0f172a;
}
</style>
