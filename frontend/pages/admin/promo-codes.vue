<script setup>
import {
  Percent,
  Plus,
  Trash2,
  Edit2,
  Save,
  CheckCircle2,
  AlertCircle,
  Ticket
} from 'lucide-vue-next';

definePageMeta({
  layout: 'admin',
  middleware: 'admin-auth'
});

const config = useRuntimeConfig();
const authStore = useAdminAuth();

/*
|--------------------------------------------------------------------------
| Fetch Promotions
|--------------------------------------------------------------------------
*/
const {
  data: promotions,
  pending,
  refresh
} = await useFetch(`${config.public.apiBase}/admin/promotions`, {
  headers: {
    Authorization: `Bearer ${authStore.token}`
  }
});

/*
|--------------------------------------------------------------------------
| Notification States
|--------------------------------------------------------------------------
*/
const successMsg = ref('');
const errorMsg = ref('');

const setNotification = (type, msg) => {
  if (type === 'success') {
    successMsg.value = msg;
    errorMsg.value = '';

    setTimeout(() => {
      successMsg.value = '';
    }, 4000);
  } else {
    errorMsg.value = msg;
    successMsg.value = '';

    setTimeout(() => {
      errorMsg.value = '';
    }, 5000);
  }
};

/*
|--------------------------------------------------------------------------
| Form States
|--------------------------------------------------------------------------
*/
const loading = ref(false);
const isEditing = ref(false);
const editingId = ref(null);

const form = reactive({
  code: '',
  discount_type: 'percentage',
  value: 0,
  max_uses: 0,
  used_count: 0,
  start_date: '',
  end_date: '',
  is_active: true
});

/*
|--------------------------------------------------------------------------
| Reset Form
|--------------------------------------------------------------------------
*/
const resetForm = () => {
  form.code = '';
  form.discount_type = 'percentage';
  form.value = 0;
  form.max_uses = 0;
  form.used_count = 0;
  form.start_date = '';
  form.end_date = '';
  form.is_active = true;

  isEditing.value = false;
  editingId.value = null;
};

/*
|--------------------------------------------------------------------------
| Submit Promotion
|--------------------------------------------------------------------------
*/
const handleSubmit = async () => {
  if (!form.code) {
    setNotification('error', 'Promo code is required');
    return;
  }

  if (form.value <= 0) {
    setNotification('error', 'Discount value must be greater than 0');
    return;
  }

  loading.value = true;

  try {
    const payload = {
      code: form.code.toUpperCase(),
      discount_type: form.discount_type,
      value: Number(form.value),
      max_uses: Number(form.max_uses),
      used_count: Number(form.used_count),
      is_active: form.is_active,
      start_date: form.start_date
        ? new Date(form.start_date).toISOString()
        : null,
      end_date: form.end_date
        ? new Date(form.end_date).toISOString()
        : null
    };

    let url = `${config.public.apiBase}/admin/promotions`;
    let method = 'POST';

    if (isEditing.value && editingId.value) {
      url = `${config.public.apiBase}/admin/promotions/${editingId.value}`;
      method = 'PUT';
    }

    await $fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${authStore.token}`
      },
      body: payload
    });

    setNotification(
      'success',
      isEditing.value
        ? 'Promotion updated successfully'
        : 'Promotion created successfully'
    );

    resetForm();
    await refresh();
  } catch (err) {
    setNotification(
      'error',
      err.data?.error || err.message || 'Something went wrong'
    );
  } finally {
    loading.value = false;
  }
};

/*
|--------------------------------------------------------------------------
| Edit Promotion
|--------------------------------------------------------------------------
*/
const startEdit = (promo) => {
  isEditing.value = true;
  editingId.value = promo.id;

  form.code = promo.code;
  form.discount_type = promo.discount_type;
  form.value = promo.value;
  form.max_uses = promo.max_uses;
  form.used_count = promo.used_count;
  form.is_active = promo.is_active;

  form.start_date = promo.start_date
    ? promo.start_date.substring(0, 10)
    : '';

  form.end_date = promo.end_date
    ? promo.end_date.substring(0, 10)
    : '';

  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  });
};

/*
|--------------------------------------------------------------------------
| Delete Promotion
|--------------------------------------------------------------------------
*/
const handleDelete = async (promo) => {
  if (!confirm('Delete this promotion?')) {
    return;
  }

  const isBlocked =
  Number(promo.pending_count) > 0 ||
  (promo.is_active === true && Number(promo.used_count) > 0);

  if (isBlocked) {
    setNotification(
      'error',
      'This promo code cannot be deleted because it is currently in use.'
    );
    return;
  }

  try {
    await $fetch(`${config.public.apiBase}/admin/promotions/${promo.id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${authStore.token}`
      }
    });

    setNotification('success', 'Promotion deleted');

    await refresh();

    if (editingId.value === promo.id) {
      resetForm();
    }
  } catch (err) {
    setNotification(
      'error',
      err.data?.error || 'Delete failed'
    );
  }
};
/*
|--------------------------------------------------------------------------
| Expired edite
|--------------------------------------------------------------------------
*/
const isExpired = (promo) => {
  if (!promo.end_date) return false

  const endDate = new Date(promo.end_date)
  endDate.setHours(23, 59, 59, 999)

  return endDate < new Date()
}

const isUsedUp = (promo) => {
  return promo.max_uses > 0 && promo.used_count >= promo.max_uses
}

const isDisabled = (promo) => {
  return !promo.is_active || isExpired(promo) || isUsedUp(promo)
}

const canEditPromo = (promo) => {
  return !isDisabled(promo)
}
/*
|--------------------------------------------------------------------------
| Toggle Active
|--------------------------------------------------------------------------
*/
const toggleActive = async (promo) => {
  try {
    await $fetch(`${config.public.apiBase}/admin/promotions/${promo.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${authStore.token}`
      },
      body: {
        ...promo,
        is_active: !promo.is_active
      }
    });

    setNotification(
      'success',
      `Promotion ${!promo.is_active ? 'Activated' : 'Disabled'}`
    );

    await refresh();
  } catch (err) {
    setNotification(
      'error',
      err.data?.error || 'Status update failed'
    );
  }
};
</script>

<template>
  <div class="space-y-10">
    <!-- Header -->
    <header class="flex items-center justify-between">
      <div>
        <h1 class="flex items-center gap-3 text-4xl font-black text-white">
          <Ticket class="h-9 w-9 text-blue-500" />
          Promotions
        </h1>

        <p class="mt-2 text-slate-400">
          Create and manage hotel promo discount codes.
        </p>
      </div>
    </header>

    <!-- Success -->
    <div
      v-if="successMsg"
      class="flex items-center gap-3 rounded-2xl border border-green-500/20 bg-green-500/10 p-4 text-green-400"
    >
      <CheckCircle2 class="h-5 w-5" />
      <span class="text-sm font-bold">{{ successMsg }}</span>
    </div>

    <!-- Error -->
    <div
      v-if="errorMsg"
      class="flex items-center gap-3 rounded-2xl border border-red-500/20 bg-red-500/10 p-4 text-red-400"
    >
      <AlertCircle class="h-5 w-5" />
      <span class="text-sm font-bold">{{ errorMsg }}</span>
    </div>

    <div class="grid gap-10 lg:grid-cols-3">
      <!-- Form -->
      <div class="lg:col-span-1">
        <div class="rounded-[2.5rem] border border-slate-800 bg-slate-950 p-8">
          <h2 class="mb-6 flex items-center gap-2 text-xl font-bold text-white">
            <Plus v-if="!isEditing" class="h-5 w-5 text-blue-500" />
            <Edit2 v-else class="h-5 w-5 text-blue-500" />

            {{ isEditing ? 'Edit Promotion' : 'Create Promotion' }}
          </h2>

          <form
            class="space-y-5"
            @submit.prevent="handleSubmit"
          >
            <!-- Promo Code -->
            <div class="space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-slate-500">
                Promo Code
              </label>

              <input
                v-model="form.code"
                type="text"
                placeholder="WELCOME10"
                class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 text-white outline-none"
              />
            </div>

            <!-- Discount Type -->
            <div class="space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-slate-500">
                Discount Type
              </label>

              <select
                v-model="form.discount_type"
                class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 text-white"
              >
                <option value="percentage">
                  Percentage (%)
                </option>

                <option value="fixed">
                  Fixed Amount ($)
                </option>
              </select>
            </div>

            <!-- Value -->
            <div class="space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-slate-500">
                Discount Value
              </label>

              <input
                v-model="form.value"
                type="number"
                min="0"
                step="0.01"
                class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 text-white"
              />
            </div>

            <!-- Max Uses -->
            <div class="space-y-2">
              <label class="text-xs font-bold uppercase tracking-widest text-slate-500">
                Max Uses
              </label>

              <input
                v-model="form.max_uses"
                type="number"
                min="0"
                class="w-full rounded-2xl border border-slate-800 bg-slate-900 p-4 text-white"
              />
            </div>

            <!-- Dates -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="mb-1 block text-[10px] font-bold uppercase text-slate-500">
                  Start Date
                </label>

                <input
                  v-model="form.start_date"
                  type="date"
                  class="w-full rounded-xl border border-slate-800 bg-slate-900 p-3 text-white"
                />
              </div>

              <div>
                <label class="mb-1 block text-[10px] font-bold uppercase text-slate-500">
                  End Date
                </label>

                <input
                  v-model="form.end_date"
                  type="date"
                  class="w-full rounded-xl border border-slate-800 bg-slate-900 p-3 text-white"
                />
              </div>
            </div>

            <!-- Active -->
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold uppercase tracking-widest text-slate-500">
                Active
              </span>

              <button
                type="button"
                @click="form.is_active = !form.is_active"
                class="relative inline-flex h-6 w-11 rounded-full"
                :class="form.is_active ? 'bg-blue-600' : 'bg-slate-700'"
              >
                <span
                  class="inline-block h-5 w-5 transform rounded-full bg-white transition"
                  :class="form.is_active ? 'translate-x-5' : 'translate-x-0'"
                />
              </button>
            </div>

            <!-- Buttons -->
            <div class="flex gap-3 pt-4">
              <button
                v-if="isEditing"
                type="button"
                @click="resetForm"
                class="flex-1 rounded-xl border border-slate-800 bg-slate-900 py-4 text-sm font-bold text-slate-300"
              >
                Cancel
              </button>

              <button
                type="submit"
                :disabled="loading"
                class="flex flex-[2] items-center justify-center gap-2 rounded-xl bg-blue-600 py-4 text-sm font-bold text-white"
              >
                <Save class="h-4 w-4" />

                {{
                  loading
                    ? 'Saving...'
                    : isEditing
                    ? 'Update Promotion'
                    : 'Create Promotion'
                }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- List -->
      <div class="lg:col-span-2">
        <div class="rounded-[2.5rem] border border-slate-800 bg-slate-950 p-8">
          <h2 class="mb-6 text-xl font-bold text-white">
            Promotions List
          </h2>

          <div
            v-if="pending"
            class="flex h-40 items-center justify-center"
          >
            <div class="h-10 w-10 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"></div>
          </div>

          <div
            v-else-if="promotions?.length === 0"
            class="text-center text-slate-500"
          >
            No promotions created yet.
          </div>

          <div
            v-else
            class="space-y-4"
          >
            <div
              v-for="promo in promotions"
              :key="promo.id"
              class="flex flex-col justify-between gap-4 rounded-3xl border border-slate-800 bg-slate-900/40 p-5 md:flex-row md:items-center"
            >
              <div class="space-y-2">
                <div class="flex items-center gap-3">
                  <h3 class="text-lg font-black text-white">
                    {{ promo.code }}
                  </h3>

                  <span
                    class="rounded-md px-2 py-1 text-[10px] font-black uppercase"
                    :class="
                      isDisabled(promo)
                        ? 'bg-red-500/10 text-red-400'
                        : 'bg-green-500/10 text-green-400'
                    "
                    >
                    {{ isDisabled(promo) ? 'Disabled' : 'Active' }}
                  </span>
                </div>

                <p class="text-sm text-slate-400">
                  {{
                    promo.discount_type === 'percentage'
                      ? promo.value + '% OFF'
                      : '$' + promo.value + ' OFF'
                  }}
                </p>

                <p class="text-xs text-slate-500">
                  Max Uses: {{ promo.max_uses }}
                </p>

                <p class="text-xs text-slate-500">
                  Start Date:
                  {{ promo.start_date ? new Date(promo.start_date).toLocaleDateString() : 'N/A' }}
                </p>

                <p class="text-xs text-slate-500">
                  End Date:
                  {{ promo.end_date ? new Date(promo.end_date).toLocaleDateString() : 'N/A' }}
                </p>
              </div>

              <div class="flex items-center gap-2">
                <button
                  v-if="promo.is_active && !isExpired(promo) && !isUsedUp(promo)"
                  @click="toggleActive(promo)"
                  class="rounded-xl border border-slate-800 bg-slate-900 px-4 py-2 text-xs font-bold text-slate-300"
                >
                  Disable
                </button>

                <button
                  v-if="canEditPromo(promo)"
                  @click="startEdit(promo)"
                  class="rounded-xl border border-slate-800 bg-slate-900 p-2.5 text-slate-300"
                >
                  <Edit2 class="h-4 w-4" />
                </button>

                <button @click="handleDelete(promo)" class="rounded-xl border border-red-500/20 bg-red-500/10 p-2.5 text-red-400">
                  <Trash2 class="h-4 w-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>