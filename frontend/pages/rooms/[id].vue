<script setup>
import { Users, Waves, Shield, Calendar, CreditCard, ChevronLeft } from 'lucide-vue-next';
const route = useRoute();
const config = useRuntimeConfig();
const authStore = useUserAuth();

const { data: room, pending, error } = await useFetch(`${config.public.apiBase}/rooms/${route.params.id}`);

const checkIn = ref('');
const checkOut = ref('');
const bookingLoading = ref(false);
const bookingSuccess = ref(false);
const bookingError = ref('');

const handleBooking = async () => {
  if (!authStore.isLoggedIn) {
    navigateTo('/login');
    return;
  }

  bookingLoading.value = true;
  bookingError.value = '';
  
  try {
    const response = await $fetch(`${config.public.apiBase}/user/bookings`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      },
      body: {
        room_id: Number(route.params.id),
        check_in: checkIn.value,
        check_out: checkOut.value
      }
    });
    bookingSuccess.value = true;
  } catch (err) {
    bookingError.value = 'Failed to book. Please ensure you are logged in and dates are valid.';
  } finally {
    bookingLoading.value = false;
  }
};
</script>

<template>
  <div class="space-y-8">
    <NuxtLink to="/rooms" class="inline-flex items-center gap-2 text-sm font-bold text-slate-500 hover:text-blue-600 transition-colors">
      <ChevronLeft class="h-4 w-4" />
      Back to Rooms
    </NuxtLink>

    <div v-if="pending" class="flex h-96 items-center justify-center">
       <div class="h-10 w-10 animate-spin rounded-full border-4 border-blue-600 border-t-transparent"></div>
    </div>

    <div v-else-if="error" class="rounded-3xl bg-red-50 p-12 text-center text-red-600">
      Room not found.
    </div>

    <div v-else class="grid gap-12 lg:grid-cols-2">
      <!-- Image & Info -->
      <div class="space-y-8">
        <div class="aspect-[4/3] overflow-hidden rounded-[2.5rem] bg-slate-200 shadow-2xl">
          <img :src="room.image_url || 'https://images.unsplash.com/photo-1590490360182-c33d57733427'" class="h-full w-full object-cover" />
        </div>
        
        <div class="space-y-4">
          <h1 class="text-4xl font-black tracking-tight text-slate-900">{{ room.name }}</h1>
          <div class="flex items-center gap-4">
            <span class="inline-flex items-center gap-1 rounded-full bg-blue-50 px-3 py-1 text-sm font-bold text-blue-600">
              <Users class="h-4 w-4" />
              Up to {{ room.capacity }} Guests
            </span>
            <span class="inline-flex items-center gap-1 rounded-full bg-green-50 px-3 py-1 text-sm font-bold text-green-600">
              <Shield class="h-4 w-4" />
              Fully Flexible
            </span>
          </div>
          <p class="text-lg text-slate-500 leading-relaxed">{{ room.description }}</p>
        </div>
      </div>

      <!-- Booking Card -->
      <div class="sticky top-24 self-start space-y-6 rounded-[2.5rem] bg-white p-10 shadow-2xl ring-1 ring-slate-100">
        <div class="flex items-baseline justify-between">
          <span class="text-3xl font-black text-slate-900">${{ room.price_per_night }}</span>
          <span class="text-slate-400">per night</span>
        </div>

        <div v-if="bookingSuccess" class="rounded-2xl bg-green-50 p-6 text-center text-green-700">
          <p class="font-bold">Booking Request Sent!</p>
          <p class="text-sm mt-1">Check your bookings page for status updates.</p>
          <NuxtLink to="/bookings" class="mt-4 inline-block font-bold underline">Go to Bookings</NuxtLink>
        </div>

        <form v-else @submit.prevent="handleBooking" class="space-y-6">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-slate-400">Check In</label>
              <input v-model="checkIn" type="date" required class="w-full rounded-2xl border border-slate-200 bg-slate-50 p-4 outline-none focus:ring-2 focus:ring-blue-500/20" />
            </div>
            <div class="space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-slate-400">Check Out</label>
              <input v-model="checkOut" type="date" required class="w-full rounded-2xl border border-slate-200 bg-slate-50 p-4 outline-none focus:ring-2 focus:ring-blue-500/20" />
            </div>
          </div>

          <div v-if="bookingError" class="text-sm font-medium text-red-500">{{ bookingError }}</div>

          <button type="submit" :disabled="bookingLoading" class="flex w-full items-center justify-center gap-3 rounded-2xl bg-blue-600 py-5 text-xl font-bold text-white shadow-lg shadow-blue-500/20 transition-all hover:bg-blue-700 active:scale-[0.98] disabled:opacity-50">
            <CreditCard v-if="!bookingLoading" class="h-6 w-6" />
            {{ bookingLoading ? 'Processing...' : 'Reserve Now' }}
          </button>
          
          <p class="text-center text-sm text-slate-400">You won't be charged yet</p>
        </form>
      </div>
    </div>
  </div>
</template>
