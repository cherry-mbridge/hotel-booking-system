<script setup>
import { Calendar, Tag, Info, ChevronRight } from 'lucide-vue-next';

definePageMeta({
  middleware: 'user-auth'
});

const authStore = useUserAuth();
const config = useRuntimeConfig();

const { data: bookings, pending, error } = await useFetch(`${config.public.apiBase}/user/bookings`, {
  headers: {
    'Authorization': `Bearer ${authStore.token}`
  }
});

const getStatusColor = (status) => {
  switch (status.toLowerCase()) {
    case 'confirmed': return 'bg-green-100 text-green-700';
    case 'pending': return 'bg-yellow-100 text-yellow-700';
    case 'cancelled': return 'bg-red-100 text-red-700';
    default: return 'bg-slate-100 text-slate-700';
  }
};
</script>


<template>
  <div class="space-y-12">
    <header>
      <h1 class="text-4xl font-black tracking-tight text-slate-900">Your Bookings</h1>
      <p class="mt-2 text-slate-500">Manage and view your upcoming stays at Lumina.</p>
    </header>

    <div v-if="pending" class="flex h-64 items-center justify-center">
       <div class="h-10 w-10 animate-spin rounded-full border-4 border-blue-600 border-t-transparent"></div>
    </div>

    <div v-else-if="error" class="rounded-3xl bg-red-50 p-12 text-center text-red-600">
      Failed to load bookings. Please check your connection.
    </div>

    <div v-else-if="!bookings || bookings.length === 0" class="rounded-[2.5rem] bg-white p-20 text-center shadow-xl ring-1 ring-slate-100">
      <div class="mx-auto mb-6 flex h-20 w-20 items-center justify-center rounded-3xl bg-slate-100 text-slate-400">
        <Calendar class="h-10 w-10" />
      </div>
      <h3 class="text-2xl font-bold text-slate-900">No bookings yet</h3>
      <p class="mt-2 text-slate-500">Your sanctuary is just a few clicks away.</p>
      <NuxtLink to="/rooms" class="mt-8 inline-block rounded-2xl bg-blue-600 px-8 py-3 font-bold text-white transition-all hover:bg-blue-700">
        Explore Rooms
      </NuxtLink>
    </div>

    <div v-else class="space-y-6">
      <div v-for="booking in bookings" :key="booking.id" class="group flex flex-col gap-6 rounded-[2.5rem] bg-white p-8 shadow-sm ring-1 ring-slate-100 transition-all hover:shadow-xl sm:flex-row sm:items-center">
        <div class="h-32 w-32 shrink-0 overflow-hidden rounded-3xl bg-slate-200">
          <img :src="booking.room?.image_url || 'https://images.unsplash.com/photo-1590490360182-c33d57733427'" class="h-full w-full object-cover" />
        </div>
        
        <div class="flex-1 space-y-2">
          <div class="flex items-center justify-between">
            <h3 class="text-xl font-bold text-slate-900">{{ booking.room?.name }}</h3>
            <span :class="['rounded-full px-4 py-1 text-xs font-bold uppercase tracking-widest', getStatusColor(booking.status)]">
              {{ booking.status }}
            </span>
          </div>
          
          <div class="flex flex-wrap gap-4 text-sm text-slate-500">
            <div class="flex items-center gap-1">
              <Calendar class="h-4 w-4" />
              {{ new Date(booking.check_in).toLocaleDateString() }} - {{ new Date(booking.check_out).toLocaleDateString() }}
            </div>
            <div class="flex items-center gap-1 font-bold text-blue-600">
              <Tag class="h-4 w-4" />
              ${{ booking.total_price }} Total
            </div>
          </div>
        </div>

        <ChevronRight class="hidden h-6 w-6 text-slate-300 transition-transform group-hover:translate-x-2 sm:block" />
      </div>
    </div>
  </div>
</template>
