<script setup>
import { Plus, Pencil, Trash2, Loader2, ChevronLeft, ChevronRight } from 'lucide-vue-next';

definePageMeta({
  layout: 'admin',
  middleware: 'admin-auth'
});

const authStore = useAdminAuth();
const config = useRuntimeConfig();
const toast = useToast();

const rooms = ref([]);
const loading = ref(false);

// Pagination state
const page = ref(1);
const perPage = ref(10);
const totalPages = ref(1);
const total = ref(0);

const fetchRooms = async () => {
  loading.value = true;
  try {
    const data = await $fetch(`${config.public.apiBase}/admin/rooms`, {
      query: { page: page.value, per_page: perPage.value },
      headers: { Authorization: `Bearer ${authStore.token}` }
    });
    rooms.value = data.data || [];
    total.value = data.total || 0;
    totalPages.value = data.total_pages || 1;
  } catch (err) {
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

// State
const deletingId = ref(null);
const showConfirm = ref(false);
const roomToDelete = ref(null);
const showErrorModal = ref(false);
const errorModalTitle = ref('Error');
const errorModalMessage = ref('Something went wrong.');

const confirmMessage = computed(() => {
  if (!roomToDelete.value) return 'Are you sure you want to delete this room? This action cannot be undone.';
  return `Are you sure you want to delete "${roomToDelete.value.name}"? This action cannot be undone.`;
});

const openConfirm = (room) => {
  roomToDelete.value = room;
  showConfirm.value = true;
};

const closeConfirm = () => {
  showConfirm.value = false;
  roomToDelete.value = null;
};

const closeErrorModal = () => {
  showErrorModal.value = false;
};

const handleDelete = async () => {
  if (!roomToDelete.value) return;

  const id = roomToDelete.value.id;
  deletingId.value = id;

  try {
    await $fetch(`${config.public.apiBase}/admin/rooms/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${authStore.token}` }
    });

    toast.success('Room Deleted', `"${roomToDelete.value.name}" has been deleted successfully.`);
    await fetchRooms();
    closeConfirm();
  } catch (err) {
    closeConfirm();
    const status = err?.statusCode || err?.response?.status;
    const data = err?.data || err?.response?._data || {};

    if (status === 409 && data.code === 'ROOM_HAS_ACTIVE_BOOKINGS') {
      errorModalTitle.value = 'Cannot Delete Room';
      errorModalMessage.value = 'This room still has active bookings.';
      showErrorModal.value = true;
    } else {
      errorModalTitle.value = 'Delete Failed';
      errorModalMessage.value = data.message || 'Unable to delete the room. Please try again later.';
      showErrorModal.value = true;
    }
  } finally {
    deletingId.value = null;
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

    <div v-if="loading" class="flex items-center justify-center h-48 bg-slate-950 rounded-[2rem] border border-slate-800/60">
      <div class="h-7 w-7 animate-spin rounded-full border-2 border-blue-500 border-t-transparent"></div>
    </div>

    <div v-else-if="rooms.length === 0" class="p-12 text-center bg-slate-950 rounded-[2rem] border border-slate-800/80 text-slate-500">
      No rooms found.
    </div>

    <div v-else class="overflow-hidden rounded-[2.5rem] bg-slate-950 border border-slate-800 shadow-md">
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
                  <button
                    @click="openConfirm(room)"
                    :disabled="deletingId === room.id"
                    class="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-900 border border-slate-800 text-slate-400 transition-all hover:bg-red-500/10 hover:border-red-500/30 hover:text-red-400 disabled:opacity-40 disabled:cursor-not-allowed"
                  >
                    <Loader2 v-if="deletingId === room.id" class="h-4 w-4 animate-spin" />
                    <Trash2 v-else class="h-4 w-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-slate-800 px-8 py-5">
        <div class="text-sm text-slate-400">
          Showing <span class="font-bold text-white">{{ rooms.length }}</span> of <span class="font-bold text-white">{{ total }}</span> rooms
        </div>
        <div class="flex items-center gap-2">
          <button
            @click="prevPage"
            :disabled="page === 1 || loading"
            class="flex items-center gap-1 px-4 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-300 hover:text-white hover:border-slate-700 disabled:opacity-40 disabled:cursor-not-allowed transition-all text-sm font-semibold"
          >
            <ChevronLeft class="h-4 w-4" />
            Prev
          </button>
          <span class="text-sm text-slate-400 font-medium px-2">Page {{ page }} of {{ totalPages }}</span>
          <button
            @click="nextPage"
            :disabled="page === totalPages || loading"
            class="flex items-center gap-1 px-4 py-2 rounded-xl bg-slate-900 border border-slate-800 text-slate-300 hover:text-white hover:border-slate-700 disabled:opacity-40 disabled:cursor-not-allowed transition-all text-sm font-semibold"
          >
            Next
            <ChevronRight class="h-4 w-4" />
          </button>
        </div>
      </div>
    </div>

    <ConfirmModal
      :is-open="showConfirm"
      title="Delete Room"
      :message="confirmMessage"
      confirm-text="Delete"
      :loading="deletingId !== null"
      @confirm="handleDelete"
      @cancel="closeConfirm"
    />

    <ErrorModal
      :is-open="showErrorModal"
      :title="errorModalTitle"
      :message="errorModalMessage"
      button-text="Got it"
      @close="closeErrorModal"
    />
  </div>
</template>
