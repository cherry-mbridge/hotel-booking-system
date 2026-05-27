<script setup>
import { Percent, Plus, Trash2, Edit2, Save, X, Calendar, AlertCircle, CheckCircle2 } from 'lucide-vue-next';

definePageMeta({
  layout: 'admin',
  middleware: 'admin-auth'
});

const config = useRuntimeConfig();
const authStore = useAdminAuth();

// Fetch rooms for the dropdown selector
const { data: rooms } = await useFetch(`${config.public.apiBase}/rooms/all`);

// Retrieve weekend pricing rules
const { data: rules, pending, error, refresh } = await useFetch(`${config.public.apiBase}/admin/weekend-pricing`, {
  headers: {
    'Authorization': `Bearer ${authStore.token}`
  }
});

// Toast / Notification states
const successMsg = ref('');
const errorMsg = ref('');
const setNotification = (type, msg) => {
  if (type === 'success') {
    successMsg.value = msg;
    errorMsg.value = '';
    setTimeout(() => { successMsg.value = ''; }, 4000);
  } else {
    errorMsg.value = msg;
    successMsg.value = '';
    setTimeout(() => { errorMsg.value = ''; }, 6000);
  }
};

// Form states
const isEditing = ref(false);
const editingRuleId = ref(null);
const loading = ref(false);

const form = reactive({
  room_type_id: '',
  adjustment_type: 'increase',
  value_type: 'percentage',
  adjustment_value: 0,
  days_of_week: [], // checkbox bound array: e.g. ["Friday", "Saturday", "Sunday"]
  start_date: '',
  end_date: '',
  is_active: true
});

const resetForm = () => {
  form.room_type_id = rooms.value?.length > 0 ? rooms.value[0].id : '';
  form.adjustment_type = 'increase';
  form.value_type = 'percentage';
  form.adjustment_value = 0;
  form.days_of_week = ['Friday', 'Saturday', 'Sunday'];
  form.start_date = '';
  form.end_date = '';
  form.is_active = true;
  isEditing.value = false;
  editingRuleId.value = null;
};

// Initialize form room selector if rooms are present
onMounted(() => {
  resetForm();
});

// Helper to map room_type_id to room name
const getRoomName = (id) => {
  const room = rooms.value?.find(r => r.id === Number(id));
  return room ? room.name : `Room ID: ${id}`;
};

// Save rule (Create or Update)
const handleSubmit = async () => {
  if (!form.room_type_id) {
    setNotification('error', 'Please select a room.');
    return;
  }
  if (form.days_of_week.length === 0) {
    setNotification('error', 'Select at least one day of week.');
    return;
  }
  if (form.adjustment_value <= 0) {
    setNotification('error', 'Adjustment value must be greater than 0.');
    return;
  }

  loading.value = true;
  try {
    const payload = {
      room_type_id: Number(form.room_type_id),
      adjustment_type: form.adjustment_type,
      value_type: form.value_type,
      adjustment_value: Number(form.adjustment_value),
      days_of_week: form.days_of_week.join(','),
      is_active: form.is_active,
      start_date: form.start_date ? new Date(form.start_date).toISOString() : null,
      end_date: form.end_date ? new Date(form.end_date).toISOString() : null
    };

    let url = `${config.public.apiBase}/admin/weekend-pricing`;
    let method = 'POST';

    if (isEditing.value && editingRuleId.value) {
      url = `${config.public.apiBase}/admin/weekend-pricing/${editingRuleId.value}`;
      method = 'PUT';
    }

    const response = await $fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: payload
    });

    setNotification('success', isEditing.value ? 'Pricing rule updated successfully' : 'Weekend pricing rule created successfully');
    resetForm();
    await refresh();
  } catch (err) {
    const backendErr = err.data?.error || err.message || 'Failure writing weekend pricing rule';
    setNotification('error', backendErr);
  } finally {
    loading.value = false;
  }
};

// Select a rule to edit
const startEdit = (rule) => {
  isEditing.value = true;
  editingRuleId.value = rule.id;
  form.room_type_id = rule.room_type_id;
  form.adjustment_type = rule.adjustment_type;
  form.value_type = rule.value_type;
  form.adjustment_value = rule.adjustment_value;
  form.days_of_week = rule.days_of_week ? rule.days_of_week.split(',') : [];
  form.is_active = rule.is_active;
  form.start_date = rule.start_date ? rule.start_date.substring(0, 10) : '';
  form.end_date = rule.end_date ? rule.end_date.substring(0, 10) : '';
  
  // Smooth scroll up to form space
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

// Delete pricing rule
const handleDelete = async (id) => {
  if (!confirm('Are you sure you want to delete this weekend pricing rule?')) return;
  
  try {
    await $fetch(`${config.public.apiBase}/admin/weekend-pricing/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    });
    setNotification('success', 'Weekend pricing rule deleted successfully');
    await refresh();
    if (editingRuleId.value === id) {
      resetForm();
    }
  } catch (err) {
    setNotification('error', err.data?.error || 'Error deleting rule');
  }
};

// Toggle active status directly
const toggleActive = async (rule) => {
  try {
    const payload = {
      ...rule,
      is_active: !rule.is_active
    };
    await $fetch(`${config.public.apiBase}/admin/weekend-pricing/${rule.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      },
      body: payload
    });
    setNotification('success', `Rule toggled ${!rule.is_active ? 'Active' : 'Inactive'}`);
    await refresh();
  } catch (err) {
    setNotification('error', err.data?.error || 'Error toggling rule status');
  }
};
</script>

<template>
  <div class="space-y-10">
    <header class="flex items-center justify-between gap-4">
      <div>
        <h1 class="text-4xl font-black tracking-tight text-white flex items-center gap-3">
          <Percent class="h-9 w-9 text-blue-500" />
          Weekend Pricing Rules
        </h1>
        <p class="mt-2 text-slate-400">Configure customized surcharge or discounts for rooms on weekends (Friday, Saturday, and Sunday).</p>
      </div>
    </header>

    <!-- Success & Error Banners -->
    <div v-if="successMsg" class="flex p-4 rounded-2xl bg-green-500/10 border border-green-500/20 text-green-400 gap-3 items-center">
      <CheckCircle2 class="h-5 w-5 shrink-0" />
      <span class="text-sm font-bold">{{ successMsg }}</span>
    </div>

    <div v-if="errorMsg" class="flex p-4 rounded-2xl bg-red-500/10 border border-red-500/20 text-red-400 gap-3 items-center">
      <AlertCircle class="h-5 w-5 shrink-0" />
      <span class="text-sm font-bold">{{ errorMsg }}</span>
    </div>

    <div class="grid gap-10 lg:grid-cols-3">
      <!-- Form Panel -->
      <div class="lg:col-span-1 space-y-6">
        <div class="rounded-[2.5rem] bg-slate-950 border border-slate-800 p-8 shadow-xl">
          <h2 class="text-xl font-bold text-white mb-6 flex items-center gap-2">
            <Plus v-if="!isEditing" class="h-5 w-5 text-blue-500" />
            <Edit2 v-else class="h-5 w-5 text-blue-500" />
            {{ isEditing ? 'Edit Pricing Rule' : 'Create Pricing Rule' }}
          </h2>

          <form @submit.prevent="handleSubmit" class="space-y-5">
            <!-- Room Selector -->
            <div class="space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Select Room / Category</label>
              <select v-model="form.room_type_id" required class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 text-white outline-none focus:ring-2 focus:ring-blue-500/20 transition-all">
                <option value="" disabled>Choose a listing...</option>
                <option v-for="room in rooms" :key="room.id" :value="room.id">
                  {{ room.name }} (${{ room.price_per_night }}/night)
                </option>
              </select>
            </div>

            <!-- Days of Week check -->
            <div class="space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Applicable Days</label>
              <div class="grid grid-cols-3 gap-2 pt-1">
                <label v-for="day in ['Friday', 'Saturday', 'Sunday']" :key="day" 
                       class="flex flex-col items-center justify-center p-3 rounded-xl border border-slate-800 cursor-pointer transition-all hover:bg-slate-900"
                       :class="form.days_of_week.includes(day) ? 'bg-blue-600/10 border-blue-500/40 text-blue-400' : 'bg-slate-900 text-slate-400'">
                  <input type="checkbox" :value="day" v-model="form.days_of_week" class="sr-only" />
                  <span class="text-xs font-bold select-none">{{ day }}</span>
                </label>
              </div>
            </div>

            <!-- Adjustment Type -->
            <div class="space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Adjustment Type</label>
              <select v-model="form.adjustment_type" required class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 text-white outline-none focus:ring-2 focus:ring-blue-500/20 transition-all">
                <option value="increase">Surcharge (Increase)</option>
                <option value="decrease">Discount (Decrease)</option>
              </select>
            </div>

            <!-- Value Type & Adjustment Value -->
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Value Type</label>
                <select v-model="form.value_type" required class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 text-white outline-none focus:ring-2 focus:ring-blue-500/20 transition-all">
                  <option value="percentage">Percentage (%)</option>
                  <option value="fixed">Fixed ($)</option>
                </select>
              </div>

              <div class="space-y-2">
                <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Value</label>
                <input v-model="form.adjustment_value" type="number" step="any" min="0.01" required 
                       class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 text-white outline-none focus:ring-2 focus:ring-blue-500/20 transition-all text-center font-bold" />
              </div>
            </div>

            <!-- Date range (optional) -->
            <div class="space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-slate-500">Applicable Date Range (Optional)</label>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="text-[10px] uppercase font-bold text-slate-600 block mb-1">Start Date</label>
                  <input v-model="form.start_date" type="date" class="w-full rounded-xl border border-slate-800 bg-slate-900 p-3 text-sm text-slate-300 outline-none focus:ring-2 focus:ring-blue-500/25" />
                </div>
                <div>
                  <label class="text-[10px] uppercase font-bold text-slate-600 block mb-1">End Date</label>
                  <input v-model="form.end_date" type="date" class="w-full rounded-xl border border-slate-800 bg-slate-900 p-3 text-sm text-slate-300 outline-none focus:ring-2 focus:ring-blue-500/25" />
                </div>
              </div>
            </div>

            <!-- Is Active -->
            <div class="flex items-center justify-between py-2">
              <span class="text-xs font-bold uppercase tracking-widest text-slate-500">Enable Right Away</span>
              <button type="button" @click="form.is_active = !form.is_active" 
                      class="relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out"
                      :class="form.is_active ? 'bg-blue-600' : 'bg-slate-800'">
                <span class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"
                      :class="form.is_active ? 'translate-x-5' : 'translate-x-0'"></span>
              </button>
            </div>

            <!-- Submit & Reset buttons -->
            <div class="flex gap-3 pt-4">
              <button v-if="isEditing" type="button" @click="resetForm" :disabled="loading" 
                      class="flex-1 py-4 text-sm font-bold text-slate-400 bg-slate-900 border border-slate-800 rounded-xl transition-all hover:bg-slate-850 active:scale-95 disabled:opacity-50">
                Cancel
              </button>
              <button type="submit" :disabled="loading" 
                      class="flex-[2] flex items-center justify-center gap-2 rounded-xl bg-blue-600 py-4 text-sm font-bold text-white transition-all hover:bg-blue-700 active:scale-95 disabled:opacity-50 shadow-md">
                <Save class="h-4 w-4" />
                {{ loading ? 'Saving...' : isEditing ? 'Update Rule' : 'Create Rule' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Rules List Panel -->
      <div class="lg:col-span-2 space-y-6">
        <div class="rounded-[2.5rem] bg-slate-950 border border-slate-800 p-8 shadow-xl">
          <h2 class="text-xl font-bold text-white mb-6">Configured Pricing Adjustments</h2>

          <div v-if="pending" class="flex h-64 items-center justify-center">
            <div class="h-10 w-10 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"></div>
          </div>

          <div v-else-if="rules?.length === 0" class="flex flex-col items-center justify-center h-64 border-2 border-dashed border-slate-800 rounded-3xl p-6 text-center">
            <Percent class="h-12 w-12 text-slate-600 mb-3" />
            <h3 class="text-lg font-bold text-slate-300">No general rules set yet</h3>
            <p class="mt-1 text-sm text-slate-500 max-w-sm">Create high-season markups or custom discounts using the creator form on the left.</p>
          </div>

          <div v-else class="space-y-4">
            <div v-for="rule in rules" :key="rule.id" 
                 class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-5 rounded-3xl border border-slate-850 bg-slate-900/40 hover:border-slate-800 transition-all"
                 :class="!rule.is_active ? 'opacity-60' : ''">
              <div class="space-y-1.5 flex-1 min-w-0">
                <div class="flex items-center gap-2.5 flex-wrap">
                  <h4 class="font-bold text-md text-white truncate">{{ getRoomName(rule.room_type_id) }}</h4>
                  <span class="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-black uppercase tracking-wider bg-slate-800"
                        :class="rule.adjustment_type === 'increase' ? 'text-rose-400 bg-rose-500/10' : 'text-emerald-400 bg-emerald-500/10'">
                    {{ rule.adjustment_type === 'increase' ? 'Surcharge' : 'Discount' }}
                  </span>
                </div>

                <div class="flex items-center gap-1.5 text-xs text-slate-400">
                  <span class="font-bold text-slate-300">Days:</span>
                  <span class="inline-flex gap-1">
                    <span v-for="d in rule.days_of_week.split(',')" :key="d" class="px-1.5 py-0.5 rounded bg-slate-800 text-[10px] text-slate-300">{{ d }}</span>
                  </span>
                </div>

                <div class="flex items-center gap-1.5 text-xs text-slate-400">
                  <Calendar class="h-3.5 w-3.5 text-slate-500" />
                  <span v-if="rule.start_date || rule.end_date" class="font-mono text-[11px]">
                    {{ rule.start_date ? rule.start_date.substring(0, 10) : 'Any Start' }} &rarr; {{ rule.end_date ? rule.end_date.substring(0, 10) : 'Any End' }}
                  </span>
                  <span v-else class="text-slate-500 italic">Always Active Date Range</span>
                </div>
              </div>

              <!-- Value & Actions -->
              <div class="flex items-center justify-between md:justify-end gap-6 border-t md:border-t-0 pt-3 md:pt-0 border-slate-800/40">
                <div class="text-right">
                  <p class="text-xs font-bold uppercase text-slate-500 tracking-wider">Adjustment</p>
                  <p class="text-xl font-black" :class="rule.adjustment_type === 'increase' ? 'text-rose-400' : 'text-emerald-400'">
                    {{ rule.adjustment_type === 'increase' ? '+' : '-' }}
                    {{ rule.value_type === 'percentage' ? rule.adjustment_value + '%' : '$' + rule.adjustment_value }}
                  </p>
                </div>

                <!-- Control Buttons -->
                <div class="flex items-center gap-1.5 shrink-0">
                  <button @click="toggleActive(rule)" 
                          class="p-2.5 rounded-xl border border-slate-800 bg-slate-900 text-xs font-bold hover:bg-slate-850 text-slate-400 hover:text-white transition-colors"
                          :title="rule.is_active ? 'Deactivate' : 'Activate'">
                    {{ rule.is_active ? 'Active' : 'Disabled' }}
                  </button>
                  <button @click="startEdit(rule)" 
                          class="p-2.5 rounded-xl border border-slate-800 bg-slate-900 text-slate-400 hover:text-white transition-colors"
                          title="Edit">
                    <Edit2 class="h-4 w-4" />
                  </button>
                  <button @click="handleDelete(rule.id)" 
                          class="p-2.5 rounded-xl border border-slate-800 bg-slate-900 text-red-400 hover:bg-red-500/10 hover:border-red-500/20 transition-all"
                          title="Delete">
                    <Trash2 class="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
