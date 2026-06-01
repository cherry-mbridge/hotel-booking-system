<script setup>
import { ChevronLeft, ChevronRight } from 'lucide-vue-next';

const config = useRuntimeConfig();

const rooms = ref([]);
const loading = ref(false);
const error = ref(null);

// Pagination state
const page = ref(1);
const perPage = ref(9);
const totalPages = ref(1);
const total = ref(0);

const fetchRooms = async () => {
  loading.value = true;
  error.value = null;
  try {
    const data = await $fetch(`${config.public.apiBase}/rooms`, {
      query: { page: page.value, per_page: perPage.value }
    });
    rooms.value = data.data || [];
    total.value = data.total || 0;
    totalPages.value = data.total_pages || 1;
  } catch (err) {
    error.value = err;
    console.error(err);
  } finally {
    loading.value = false;
  }
};

const prevPage = () => {
  if (page.value > 1) {
    page.value--;
    fetchRooms();
  }
};

const nextPage = () => {
  if (page.value < totalPages.value) {
    page.value++;
    fetchRooms();
  }
};

onMounted(() => {
  fetchRooms();
});
</script>

<template>
  <div class="space-y-12">
    <header class="text-center sm:text-left">
      <h1 class="text-4xl font-extrabold tracking-tight text-slate-900">Explore Our Rooms</h1>
      <p class="mt-4 text-xl text-slate-500">Find the perfect sanctuary for your next stay.</p>
    </header>

    <div v-if="loading" class="flex items-center justify-center h-64">
      <div class="h-8 w-8 animate-spin rounded-full border-4 border-blue-600 border-t-transparent"></div>
    </div>

    <div v-else-if="error" class="p-8 text-center bg-red-50 text-red-600 rounded-3xl">
      Error loading rooms. Please check if the Go backend is running on 8080.
    </div>

    <div v-else-if="rooms.length === 0" class="p-12 text-center bg-slate-50 rounded-3xl text-slate-500">
      No rooms available at the moment.
    </div>

    <div v-else class="space-y-10">
      <div class="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
        <div v-for="room in rooms" :key="room.id" 
             class="group overflow-hidden rounded-3xl bg-white shadow-sm ring-1 ring-slate-200 transition-all hover:shadow-xl">
          <div class="aspect-[4/3] overflow-hidden">
            <img :src="room.image_url || 'https://images.unsplash.com/photo-1566073771259-6a8506099945'" class="h-full w-full object-cover transition-transform group-hover:scale-105" />
          </div>
          <div class="p-6">
            <h3 class="text-xl font-bold">{{ room.name }}</h3>
            <p class="mt-2 text-slate-500 line-clamp-2">{{ room.description }}</p>
            <div class="mt-6 flex items-center justify-between">
              <span class="text-2xl font-black text-blue-600">${{ room.price_per_night }}<span class="text-xs text-slate-400 font-normal">/night</span></span>
              <NuxtLink :to="`/rooms/${room.id}`" class="rounded-xl bg-slate-900 px-4 py-2 text-sm font-bold text-white transition-colors hover:bg-blue-600">
                Details
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="text-sm text-slate-500">
          Showing <span class="font-bold text-slate-900">{{ rooms.length }}</span> of <span class="font-bold text-slate-900">{{ total }}</span> rooms
        </div>
        <div class="flex items-center gap-2">
          <button
            @click="prevPage"
            :disabled="page === 1 || loading"
            class="flex items-center gap-1 px-4 py-2 rounded-xl bg-white border border-slate-200 text-slate-700 hover:text-slate-900 hover:border-slate-300 disabled:opacity-40 disabled:cursor-not-allowed transition-all text-sm font-semibold shadow-sm"
          >
            <ChevronLeft class="h-4 w-4" />
            Prev
          </button>
          <span class="text-sm text-slate-500 font-medium px-2">Page {{ page }} of {{ totalPages }}</span>
          <button
            @click="nextPage"
            :disabled="page === totalPages || loading"
            class="flex items-center gap-1 px-4 py-2 rounded-xl bg-white border border-slate-200 text-slate-700 hover:text-slate-900 hover:border-slate-300 disabled:opacity-40 disabled:cursor-not-allowed transition-all text-sm font-semibold shadow-sm"
          >
            Next
            <ChevronRight class="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
