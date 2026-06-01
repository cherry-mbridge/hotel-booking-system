<script setup>
import { ArrowLeft, Save } from 'lucide-vue-next';

definePageMeta({
  layout: 'admin',
  middleware: 'admin-auth'
});

const route = useRoute();
const authStore = useAdminAuth();
const config = useRuntimeConfig();

const categories = ref([]);
const categoriesLoading = ref(true);

const form = reactive({
  name: '',
  description: '',
  price_per_night: 0,
  capacity: 1,
  image_url: '',
  category_id: ''
});

const loading = ref(false);
const fetching = ref(true);
const errors = ref({});

onMounted(async () => {
  try {
    const catsData = await $fetch(`${config.public.apiBase}/categories/all`);
    categories.value = catsData || [];

    const data = await $fetch(`${config.public.apiBase}/rooms/${route.params.id}`);
    form.name = data.name;
    form.description = data.description;
    form.price_per_night = data.price_per_night;
    form.capacity = data.capacity;
    form.image_url = data.image_url;
    form.category_id = data.category_id;
  } catch (err) {
    alert('Failed to fetch details');
  } finally {
    fetching.value = false;
    categoriesLoading.value = false;
  }
});

const handleSubmit = async () => {
  loading.value = true;
  try {
    await $fetch(`${config.public.apiBase}/admin/rooms/${route.params.id}`, {
      method: 'PUT',
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
      <h1 class="text-4xl font-black tracking-tight text-white">Edit Room</h1>
      <p class="mt-2 text-slate-400">Update the information for this property listing.</p>
    </header>

    <div v-if="fetching || categoriesLoading" class="flex h-64 items-center justify-center">
      <div class="h-10 w-10 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"></div>
    </div>

    <form v-else @submit.prevent="handleSubmit" class="space-y-6 rounded-[2.5rem] bg-slate-950 border border-slate-800 p-10 shadow-md">
      <div class="grid gap-6 sm:grid-cols-2">
        <div class="space-y-2 lg:col-span-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Room Name</label>
          <input v-model="form.name" required type="text" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 transition-all placeholder-slate-600" />
          <p v-if="errors.name" class="text-sm text-red-500 mt-1">{{ errors.name }}</p>
        </div>
        
        <div class="space-y-2 lg:col-span-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Image URL</label>
          <input v-model="form.image_url" required type="text" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 transition-all placeholder-slate-600" />
          <p v-if="errors.image_url" class="text-sm text-red-500 mt-1">{{ errors.image_url }}</p>
        </div>

        <div class="space-y-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Room Category</label>
          <select v-model="form.category_id" required class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 select-style outline-none text-white focus:ring-2 focus:ring-blue-500/20 transition-all">
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">
              {{ cat.name }}
            </option>
          </select>
          <p v-if="errors.category_id" class="text-sm text-red-500 mt-1">{{ errors.category_id }}</p>
        </div>

        <div class="space-y-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Price/Night ($)</label>
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
          <textarea v-model="form.description" required rows="4" class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 outline-none text-white focus:ring-2 focus:ring-blue-500/20 placeholder-slate-600 transition-all"></textarea>
          <p v-if="errors.description" class="text-sm text-red-500 mt-1">{{ errors.description }}</p>
        </div>
      </div>

      <button type="submit" :disabled="loading" class="flex w-full items-center justify-center gap-2 rounded-2xl bg-blue-600 py-4 text-xl font-bold text-white shadow-lg shadow-blue-500/10 transition-all hover:bg-blue-700 active:scale-95 disabled:opacity-50">
        <Save v-if="!loading" class="h-6 w-6" />
        {{ loading ? 'Saving...' : 'Update Room' }}
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
