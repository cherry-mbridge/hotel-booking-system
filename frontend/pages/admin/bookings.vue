<script setup>
import { Calendar, User, Tag, ArrowLeft, CheckCircle2, XCircle, Trash2 } from 'lucide-vue-next';

definePageMeta({
  layout: 'admin',
  middleware: 'admin-auth'
});

const authStore = useAdminAuth();
const config = useRuntimeConfig();

const { data: bookings, pending, error, refresh } = await useFetch(`${config.public.apiBase}/admin/bookings`, {
  headers: {
    'Authorization': `Bearer ${authStore.token}`
  }
});

const updateStatus = async (id, status) => {
  try {
    await $fetch(`${config.public.apiBase}/admin/bookings/${id}/status`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${authStore.token}` },
      body: { status }
    });
    refresh();
  } catch (err) {
    alert('Failed to update status');
  }
};

const getStatusColor = (status) => {
  switch (status?.toLowerCase()) {
    case 'confirmed': return 'bg-green-500/10 text-green-400 border border-green-500/20';
    case 'pending': return 'bg-amber-500/10 text-amber-400 border border-amber-500/20';
    case 'cancelled': case 'rejected': return 'bg-red-500/10 text-red-400 border border-red-500/20';
    default: return 'bg-slate-500/10 text-slate-400 border border-slate-500/20';
  }
};

const handleDelete = async (id) => {
  if (!confirm('Delete this booking?')) return;

  try {
    await $fetch(`${config.public.apiBase}/admin/bookings/${id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    });

    await refresh();
    setNotification('success', 'Booking deleted');

  } catch (err) {
    setNotification('error', err.data?.error || 'Delete failed');
  }
};

const isBookingFinished = (checkOutDate) => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const checkOut = new Date(checkOutDate)
  checkOut.setHours(0, 0, 0, 0)

  return checkOut < today
}
</script>

<template>
  <div class="space-y-8">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-4xl font-black tracking-tight text-white">All Bookings</h1>
        <p class="mt-2 text-slate-400">Master log of all reservations in the system.</p>
      </div>
      <NuxtLink to="/admin" class="flex items-center gap-2 rounded-2xl bg-slate-800 hover:bg-slate-700 px-6 py-3 font-bold text-slate-200 transition-all">
        <ArrowLeft class="h-5 w-5" />
        Overview
      </NuxtLink>
    </div>

    <div v-if="pending" class="flex h-64 items-center justify-center">
       <div class="h-10 w-10 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"></div>
    </div>

    <div v-else-if="error" class="rounded-3xl bg-red-500/10 border border-red-500/20 p-12 text-center text-red-400">
      Failed to load master booking list.
    </div>

    <div v-else class="overflow-hidden rounded-[2.5rem] bg-slate-950 border border-slate-800 shadow-md">
      <div class="overflow-x-auto">
        <table class="w-full text-left">
          <thead>
            <tr class="border-b border-slate-800 bg-slate-950/50 text-xs font-bold uppercase tracking-widest text-slate-500">
              <th class="px-8 py-5">Guest</th>
              <th class="px-8 py-5">Room</th>
              <th class="px-8 py-5">Dates</th>
              <th class="px-8 py-5">Total</th>
              <th class="px-8 py-5">Promo Code</th>
              <th class="px-8 py-5">Status</th>
              <th class="px-8 py-5 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800">
            <tr v-for="booking in bookings" :key="booking.id" class="transition-colors hover:bg-slate-900/40">
              <td class="px-8 py-4">
                <div class="flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-900 border border-slate-800 text-blue-400">
                    <User class="h-5 w-5" />
                  </div>
                  <div>
                     <p class="font-bold text-slate-200">{{ booking.user?.name || 'Guest' }}</p>
                     <p class="text-xs text-slate-500">{{ booking.user?.email }}</p>
                  </div>
                </div>
              </td>
              <td class="px-8 py-4 font-semibold text-slate-300">{{ booking.room?.name }}</td>
              <td class="px-8 py-4 text-sm text-slate-400 font-medium">
                 {{ new Date(booking.check_in).toLocaleDateString() }} - {{ new Date(booking.check_out).toLocaleDateString() }}
              </td>
              <td class="px-8 py-4 font-bold text-white font-mono">${{ booking.total_price }}</td>
              <td class="px-8 py-4 font-semibold text-slate-300">
                {{ booking.promo_code?.code || '-' }}
              </td>
              <td class="px-8 py-4">
                <span :class="['rounded-full px-3 py-1 text-xs font-bold uppercase tracking-widest', getStatusColor(booking.status)]">
                  {{ booking.status }}
                </span>
              </td>
              <td class="px-8 py-4 text-right">
                <div v-if="booking.status === 'pending' && !isBookingFinished(booking.check_out)" class="flex justify-end gap-2">
                  <button @click="updateStatus(booking.id, 'confirmed')" class="flex h-9 items-center gap-1 rounded-xl bg-green-500/10 border border-green-500/25 px-3 text-sm font-bold text-green-400 transition-all hover:bg-green-500/20">
                    <CheckCircle2 class="h-4 w-4" /> Approve
                  </button>
                  <button @click="updateStatus(booking.id, 'rejected')" class="flex h-9 items-center gap-1 rounded-xl bg-red-500/10 border border-red-500/25 px-3 text-sm font-bold text-red-400 transition-all hover:bg-red-500/20">
                    <XCircle class="h-4 w-4" /> Reject
                  </button>
                </div>
                <div v-else class="flex justify-end">
                  <button @click="handleDelete(booking.id)" class="flex h-9 items-center gap-1 rounded-xl bg-red-500/10 border border-red-500/25 px-3 text-sm font-bold text-red-400 transition-all hover:bg-red-500/20">
                    <Trash2 class="h-4 w-4" />Delete
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

