<script setup>
import { ArrowLeft, Save } from 'lucide-vue-next';

definePageMeta({
  layout: 'admin',
  middleware: 'admin-auth'
});

const authStore = useAdminAuth();
const config = useRuntimeConfig();

const form = reactive({
  name: '',
  description: '',
  price_per_night: 0,
  capacity: 1,
  image_url: 'https://images.unsplash.com/photo-1590490360182-c33d57733427',
  category_id: 1
});

const loading = ref(false);

const handleSubmit = async () => {
  loading.value = true;
  try {
    await $fetch(`${config.public.apiBase}/admin/rooms`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${authStore.token}` },
      body: {
        ...form,
        price_per_night: Number(form.price_per_night),
        capacity: Number(form.capacity),
        category_id: Number(form.category_id)
      }
    });
    navigateTo('/admin/rooms');
  } catch (err) {
    alert('Failed to create room');
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="max-w-3xl space-y-8">
    <NuxtLink to="/admin/rooms" class="inline-flex items-center gap-2 text-sm font-bold text-slate-400 hover:text-white transition-colors">
      <ArrowLeft class="h-4 w-4" />
      Back to Inventory
    </NuxtLink>

    <header>
      <h1 class="text-4xl font-black tracking-tight text-white">Create New Room</h1>
      <p class="mt-2 text-slate-400">Fill in the details to add a new room to your property.</p>
    </header>

    <form @submit.prevent="handleSubmit" class="space-y-6 rounded-[2.5rem] bg-slate-950 border border-slate-800 p-10 shadow-md">
      <div class="grid gap-6 sm:grid-cols-2">
        <div class="space-y-2 lg:col-span-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Room Name</label>
          <input v-model="form.name" required type="text" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 placeholder-slate-600 transition-all" placeholder="Deluxe Ocean View" />
        </div>
        
        <div class="space-y-2 lg:col-span-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Image URL</label>
          <input v-model="form.image_url" required type="text" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 placeholder-slate-600 transition-all" placeholder="https://..." />
        </div>

        <div class="space-y-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Price/Night ($)</label>
          <input v-model="form.price_per_night" required type="number" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 transition-all" />
        </div>

        <div class="space-y-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Max Capacity</label>
          <input v-model="form.capacity" required type="number" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 transition-all" />
        </div>

        <div class="space-y-2 lg:col-span-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Description</label>
          <textarea v-model="form.description" required rows="4" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 placeholder-slate-600 transition-all" placeholder="Describe the luxury details..."></textarea>
        </div>
      </div>

      <button type="submit" :disabled="loading" class="flex w-full items-center justify-center gap-2 rounded-2xl bg-blue-600 py-4 text-xl font-bold text-white shadow-lg shadow-blue-500/10 transition-all hover:bg-blue-700 active:scale-95 disabled:opacity-50">
        <Save v-if="!loading" class="h-6 w-6" />
        {{ loading ? 'Creating...' : 'Save Room' }}
      </button>
    </form>
  </div>
</template>

