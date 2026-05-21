<script setup>
import { Plus, Pencil, Trash2, Bed } from 'lucide-vue-next';

definePageMeta({
  layout: 'admin',
  middleware: 'admin-auth'
});

const authStore = useAdminAuth();
const config = useRuntimeConfig();

const { data: rooms, refresh } = await useFetch(`${config.public.apiBase}/rooms`);

const deleteRoom = async (id) => {
  if (!confirm('Are you sure you want to delete this room?')) return;
  
  try {
    await $fetch(`${config.public.apiBase}/admin/rooms/${id}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${authStore.token}` }
    });
    refresh();
  } catch (err) {
    alert('Failed to delete room');
  }
};
</script>

<template>
  <div class="space-y-8">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-4xl font-black tracking-tight text-white">Rooms</h1>
        <p class="mt-2 text-slate-400">Manage your property's room inventory.</p>
      </div>
      <NuxtLink to="/admin/rooms/create" class="flex items-center gap-2 rounded-2xl bg-blue-600 hover:bg-blue-700 px-6 py-3 font-bold text-white shadow-lg shadow-blue-500/10 transition-all active:scale-95">
        <Plus class="h-5 w-5" />
        Add Room
      </NuxtLink>
    </div>

    <div class="overflow-hidden rounded-[2.5rem] bg-slate-950 border border-slate-800 shadow-md">
      <div class="overflow-x-auto">
        <table class="w-full text-left">
          <thead>
            <tr class="border-b border-slate-800 bg-slate-950/50 text-xs font-bold uppercase tracking-widest text-slate-500">
              <th class="px-8 py-5">Image</th>
              <th class="px-8 py-5">Name</th>
              <th class="px-8 py-5">Capacity</th>
              <th class="px-8 py-5">Price</th>
              <th class="px-8 py-5 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800">
            <tr v-for="room in rooms" :key="room.id" class="group transition-colors hover:bg-slate-900/40">
              <td class="px-8 py-4">
                <div class="h-16 w-16 overflow-hidden rounded-2xl bg-slate-900 border border-slate-800">
                  <img :src="room.image_url" class="h-full w-full object-cover opacity-85" referrerPolicy="no-referrer" />
                </div>
              </td>
              <td class="px-8 py-4 font-bold text-white">{{ room.name }}</td>
              <td class="px-8 py-4 text-slate-400">{{ room.capacity }} persons</td>
              <td class="px-8 py-4 font-mono font-bold text-blue-400">${{ room.price_per_night }}</td>
              <td class="px-8 py-4 text-right">
                <div class="flex justify-end gap-2">
                  <NuxtLink :to="`/admin/rooms/${room.id}`" class="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-900 border border-slate-800 text-slate-400 transition-all hover:bg-blue-500/10 hover:border-blue-500/30 hover:text-blue-400">
                    <Pencil class="h-4 w-4" />
                  </NuxtLink>
                  <button @click="deleteRoom(room.id)" class="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-900 border border-slate-800 text-slate-400 transition-all hover:bg-red-500/10 hover:border-red-500/30 hover:text-red-400">
                    <Trash2 class="h-4 w-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

