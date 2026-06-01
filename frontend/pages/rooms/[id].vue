<script setup>
import { Users, Waves, Shield, Calendar, CreditCard, ChevronLeft, Percent } from 'lucide-vue-next';
const route = useRoute();
const config = useRuntimeConfig();
const authStore = useUserAuth();

const { data: room, pending, error } = await useFetch(`${config.public.apiBase}/rooms/${route.params.id}`);

const checkIn = ref('');
const checkOut = ref('');
const promoCode = ref('');
const priceBreakdown = ref(null);
const priceLoading = ref(false);
const bookingLoading = ref(false);
const bookingSuccess = ref(false);
const bookingError = ref('');

const calculateDynamicPrice = async () => {
  if (!checkIn.value || !checkOut.value) {
    priceBreakdown.value = null;
    return;
  }
  priceLoading.value = true;
  try {
    const data = await $fetch(`${config.public.apiBase}/rooms/${route.params.id}/price`, {
      params: {
        check_in: checkIn.value,
        check_out: checkOut.value,
        promo_code: promoCode.value
      }
    });
    priceBreakdown.value = data;
  } catch (err) {
    priceBreakdown.value = null;
  } finally {
    priceLoading.value = false;
  }
};

watch([checkIn, checkOut, promoCode], () => {
  calculateDynamicPrice();
});

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
        check_out: checkOut.value,
        promo_code: promoCode.value
      }
    });
    bookingSuccess.value = true;
  } catch (err) {
    bookingError.value = err.data?.error || 'Failed to book. Please ensure you are logged in and dates are valid.';
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
              <input v-model="checkIn" type="date" required class="w-full rounded-2xl border border-slate-200 bg-slate-50 p-4 outline-none focus:ring-2 focus:ring-blue-500/20 text-slate-800 font-semibold" />
            </div>
            <div class="space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-slate-400">Check Out</label>
              <input v-model="checkOut" type="date" required class="w-full rounded-2xl border border-slate-200 bg-slate-50 p-4 outline-none focus:ring-2 focus:ring-blue-500/20 text-slate-800 font-semibold" />
            </div>
          </div>

          <!-- Promo Code Input -->
          <div class="space-y-2">
            <label class="text-xs font-bold uppercase tracking-widest text-slate-400">Promo Code (Optional)</label>
            <input v-model="promoCode" type="text" placeholder="WELCOME10, SUMMER20, LUMINA30" class="w-full rounded-2xl border border-slate-200 bg-slate-50 p-4 outline-none focus:ring-2 focus:ring-blue-500/20 uppercase font-semibold placeholder-slate-300" />
          </div>

          <!-- Price Calculation Breakdown -->
          <div v-if="priceLoading" class="flex justify-center py-4">
            <div class="h-6 w-6 animate-spin rounded-full border-2 border-blue-600 border-t-transparent"></div>
          </div>

          <div v-else-if="priceBreakdown" class="space-y-3.5 rounded-2xl bg-slate-50 p-5 ring-1 ring-slate-100/80">
            <h4 class="text-xs font-bold uppercase tracking-widest text-slate-400 mb-2">Price Details</h4>
            
            <div class="flex justify-between text-sm text-slate-600">
              <span>Base Rate (Stay Total)</span>
              <span>${{ priceBreakdown.base_price.toFixed(2) }}</span>
            </div>

            <!-- Weekend Surcharge Label -->
            <div v-if="priceBreakdown.weekend_adjustment !== 0" class="flex justify-between items-center text-sm">
              <span class="inline-flex items-center gap-1.5 font-semibold" :class="priceBreakdown.weekend_adjustment > 0 ? 'text-red-600' : 'text-green-600'">
                <Percent class="h-3.5 w-3.5" />
                {{ priceBreakdown.weekend_adjustment > 0 ? 'Weekend Surcharge' : 'Weekend Discount' }}
              </span>
              <span class="font-bold shrink-0" :class="priceBreakdown.weekend_adjustment > 0 ? 'text-red-600' : 'text-green-600'">
                {{ priceBreakdown.weekend_adjustment > 0 ? '+' : '-' }}${{ Math.abs(priceBreakdown.weekend_adjustment).toFixed(2) }}
              </span>
            </div>

            <!-- Seasonal pricing markup -->
            <div v-if="priceBreakdown.seasonal_pricing !== 0" class="flex justify-between items-center text-sm text-amber-700">
              <span class="font-semibold">Seasonal Peak Adjustment (10%)</span>
              <span class="font-bold">+${{ priceBreakdown.seasonal_pricing.toFixed(2) }}</span>
            </div>

            <!-- Promotion discount -->
            <div v-if="priceBreakdown.promotion_discount !== 0" class="flex justify-between items-center text-sm text-emerald-600">
              <span class="font-semibold">Promo code discount</span>
              <span class="font-bold">-${{ priceBreakdown.promotion_discount.toFixed(2) }}</span>
            </div>

            <div class="border-t border-slate-200 pt-3 flex justify-between items-baseline gap-4">
              <div class="min-w-0">
                <p class="text-[10px] text-slate-400 font-bold uppercase tracking-wider">Total Stay Price</p>
                <p v-if="(priceBreakdown.weekend_adjustment + priceBreakdown.seasonal_pricing + priceBreakdown.promotion_discount) !== 0" class="text-xs text-slate-400 line-through">Original: ${{ (priceBreakdown.base_price + priceBreakdown.weekend_adjustment + priceBreakdown.seasonal_pricing).toFixed(2) }}</p>
              </div>
              <p class="text-2xl font-black text-slate-900 shrink-0">${{ priceBreakdown.final_price.toFixed(2) }}</p>
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
