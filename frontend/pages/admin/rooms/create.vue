<script setup>
import { ArrowLeft, Save } from 'lucide-vue-next';

definePageMeta({
  layout: 'admin',
  middleware: 'admin-auth'
});

const authStore = useAdminAuth();
const config = useRuntimeConfig();

const categories = ref([]);
const categoriesLoading = ref(true);

const form = reactive({
  name: '',
  description: '',
  price_per_night: 150,
  capacity: 2,
  image_url: 'https://images.unsplash.com/photo-1590490360182-c33d57733427',
  category_id: ''
});

onMounted(async () => {
  try {
    const data = await $fetch(`${config.public.apiBase}/categories/all`);
    categories.value = data || [];
    if (categories.value.length > 0) {
      form.category_id = categories.value[0].id; // Default to first category
    }
  } catch (err) {
    console.error('Failed to load categories', err);
  } finally {
    categoriesLoading.value = false;
  }
});

const loading = ref(false);
const errors = ref({});

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
    if (err.data?.errors) {
      errors.value = err.data.errors;
    }
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

    <div v-if="categoriesLoading" class="flex h-48 items-center justify-center bg-slate-950 rounded-[2.5rem] border border-slate-800">
      <div class="h-8 w-8 animate-spin rounded-full border-2 border-blue-500 border-t-transparent"></div>
    </div>

    <form v-else @submit.prevent="handleSubmit" class="space-y-6 rounded-[2.5rem] bg-slate-950 border border-slate-800 p-10 shadow-md">
      <div class="grid gap-6 sm:grid-cols-2">
        <div class="space-y-2 lg:col-span-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Room Name</label>
          <input v-model="form.name" required type="text" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 placeholder-slate-600 transition-all" placeholder="Deluxe Ocean View" />
          <p v-if="errors.name" class="text-sm text-red-500 mt-1">{{ errors.name }}</p>
        </div>
        
        <div class="space-y-2 lg:col-span-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Image URL</label>
          <input v-model="form.image_url" required type="text" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 placeholder-slate-600 transition-all" placeholder="https://..." />
          <p v-if="errors.image_url" class="text-sm text-red-500 mt-1">{{ errors.image_url }}</p>
        </div>

        <div class="space-y-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500 font-sans">Room Category</label>
          <select v-model="form.category_id" required class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 select-style transition-all">
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">
              {{ cat.name }}
            </option>
          </select>
          <p v-if="errors.category_id" class="text-sm text-red-500 mt-1">{{ errors.category_id }}</p>
        </div>

        <div class="space-y-2 mb-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Price Per Night ($/night)</label>
          <input v-model="form.price_per_night" required type="number" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 transition-all" />
          <p v-if="errors.price_per_night" class="text-sm text-red-500 mt-1">{{ errors.price_per_night }}</p>
        </div>

        <div class="space-y-2 lg:col-span-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Max Capacity</label>
          <input v-model="form.capacity" required type="number" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 transition-all" />
          <p v-if="errors.capacity" class="text-sm text-red-500 mt-1">{{ errors.capacity }}</p>
        </div>

        <div class="space-y-2 lg:col-span-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Description</label>
          <textarea v-model="form.description" required rows="4" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 placeholder-slate-600 transition-all" placeholder="Describe the luxury details..."></textarea>
          <p v-if="errors.description" class="text-sm text-red-500 mt-1">{{ errors.description }}</p>
        </div>
      </div>

      <button type="submit" :disabled="loading" class="flex w-full items-center justify-center gap-2 rounded-2xl bg-blue-600 py-4 text-xl font-bold text-white shadow-lg shadow-blue-500/10 transition-all hover:bg-blue-700 active:scale-95 disabled:opacity-50">
        <Save v-if="!loading" class="h-6 w-6" />
        {{ loading ? 'Creating...' : 'Save Room' }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.select-style {
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='white' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 1rem center;
  background-size: 1em;
}
</style>
