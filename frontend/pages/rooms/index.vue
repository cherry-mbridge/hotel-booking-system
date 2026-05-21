<script setup>
const config = useRuntimeConfig();
const { data: rooms, pending, error } = await useFetch(`${config.public.apiBase}/rooms`);
</script>

<template>
  <div class="space-y-12">
    <header class="text-center sm:text-left">
      <h1 class="text-4xl font-extrabold tracking-tight text-slate-900">Explore Our Rooms</h1>
      <p class="mt-4 text-xl text-slate-500">Find the perfect sanctuary for your next stay.</p>
    </header>

    <div v-if="pending" class="flex items-center justify-center h-64">
      <div class="h-8 w-8 animate-spin rounded-full border-4 border-blue-600 border-t-transparent"></div>
    </div>

    <div v-else-if="error" class="p-8 text-center bg-red-50 text-red-600 rounded-3xl">
      Error loading rooms. Please check if the Go backend is running on 8080.
    </div>

    <div v-else class="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
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
  </div>
</template>
