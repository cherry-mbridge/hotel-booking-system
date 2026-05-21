<script setup>
import { LayoutDashboard, Bed, CalendarCheck, Users, ArrowUpRight } from 'lucide-vue-next';

definePageMeta({
  layout: 'admin',
  middleware: 'admin-auth'
});

const authStore = useAdminAuth();
const config = useRuntimeConfig();

const { data: rooms } = await useFetch(`${config.public.apiBase}/rooms`);
const { data: bookings } = await useFetch(`${config.public.apiBase}/admin/bookings`, {
  headers: { 'Authorization': `Bearer ${authStore.token}` }
});

const stats = computed(() => [
  { label: 'Total Rooms', value: rooms.value?.length || 0, icon: Bed, color: 'text-blue-400', bg: 'bg-blue-500/10' },
  { label: 'Total Bookings', value: bookings.value?.length || 0, icon: CalendarCheck, color: 'text-purple-400', bg: 'bg-purple-500/10' },
  { label: 'Revenue (Est)', value: `$${bookings.value?.reduce((acc, b) => acc + b.total_price, 0).toFixed(2) || '0.00'}`, icon: ArrowUpRight, color: 'text-green-400', bg: 'bg-green-500/10' },
]);
</script>

<template>
  <div class="space-y-10">
    <header>
      <h1 class="text-4xl font-black tracking-tight text-white">Admin Dashboard</h1>
      <p class="mt-2 text-slate-400">System overview and key metrics.</p>
    </header>

    <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="stat in stats" :key="stat.label" class="rounded-[2.5rem] bg-slate-950 p-8 shadow-md ring-1 ring-slate-800 transition-all hover:ring-slate-700">
        <div class="flex items-center justify-between">
          <div :class="['flex h-14 w-14 items-center justify-center rounded-2xl', stat.bg, stat.color]">
            <component :is="stat.icon" class="h-7 w-7" />
          </div>
        </div>
        <div class="mt-6">
          <p class="text-sm font-bold uppercase tracking-widest text-slate-500">{{ stat.label }}</p>
          <p class="mt-1 text-3xl font-black text-white">{{ stat.value }}</p>
        </div>
      </div>
    </div>

    <div class="grid gap-10 lg:grid-cols-2">
      <NuxtLink to="/admin/rooms" class="group relative overflow-hidden rounded-[2.5rem] bg-slate-950 hover:bg-slate-900 border border-slate-800 p-10 text-white shadow-xl transition-all hover:scale-[1.01]">
        <div class="relative z-10">
          <h3 class="text-2xl font-bold">Room Management</h3>
          <p class="mt-2 text-slate-400">Add, edit or remove rooms from the inventory.</p>
          <div class="mt-8 inline-flex items-center gap-2 font-bold text-blue-400">
            Manage Rooms <ArrowUpRight class="h-4 w-4" />
          </div>
        </div>
        <Bed class="absolute -bottom-10 -right-10 h-64 w-64 opacity-10 transition-transform group-hover:scale-110" />
      </NuxtLink>

      <NuxtLink to="/admin/bookings" class="group relative overflow-hidden rounded-[2.5rem] bg-slate-950 hover:bg-slate-900 border border-slate-800 p-10 text-white shadow-xl transition-all hover:scale-[1.01]">
        <div class="relative z-10">
          <h3 class="text-2xl font-bold">Booking Logs</h3>
          <p class="mt-2 text-slate-400">Monitor and track all guest reservations.</p>
          <div class="mt-8 inline-flex items-center gap-2 font-bold text-blue-400">
            View Bookings <ArrowUpRight class="h-4 w-4" />
          </div>
        </div>
        <CalendarCheck class="absolute -bottom-10 -right-10 h-64 w-64 opacity-10 transition-transform group-hover:scale-110" />
      </NuxtLink>
    </div>
  </div>
</template>

