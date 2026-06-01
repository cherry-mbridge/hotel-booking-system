<script setup>
import { FolderTree, Plus, Trash2, Edit, Save, X, RefreshCw, Loader2, ChevronLeft, ChevronRight } from 'lucide-vue-next';

definePageMeta({
  layout: 'admin',
  middleware: 'admin-auth'
});

const config = useRuntimeConfig();
const adminAuth = useAdminAuth();
const toast = useToast();

const categories = ref([]);
const loading = ref(false);
const formError = ref('');

// Pagination state
const page = ref(1);
const perPage = ref(10);
const totalPages = ref(1);
const total = ref(0);

// Form state
const formMode = ref('create'); // 'create' or 'edit'
const editId = ref(null);
const categoryForm = ref({
  name: '',
  description: ''
});

// Delete state
const deletingId = ref(null);
const showConfirm = ref(false);
const categoryToDelete = ref(null);
const showErrorModal = ref(false);
const errorModalTitle = ref('Error');
const errorModalMessage = ref('Something went wrong.');

const confirmMessage = computed(() => {
  if (!categoryToDelete.value) return 'Are you sure you want to delete this category? This action cannot be undone.';
  return `Are you sure you want to delete "${categoryToDelete.value.name}"? This action cannot be undone.`;
});

const fetchCategories = async () => {
  loading.value = true;
  formError.value = '';
  try {
    const data = await $fetch(`${config.public.apiBase}/admin/categories`, {
      query: { page: page.value, per_page: perPage.value },
      headers: {
        'Authorization': `Bearer ${adminAuth.token}`
      }
    });
    categories.value = data.data || [];
    total.value = data.total || 0;
    totalPages.value = data.total_pages || 1;
  } catch (err) {
    formError.value = 'Failed to load room categories';
    console.error(err);
  } finally {
    loading.value = false;
  }
};

const prevPage = () => {
  if (page.value > 1) {
    page.value--;
    fetchCategories();
  }
};

const nextPage = () => {
  if (page.value < totalPages.value) {
    page.value++;
    fetchCategories();
  }
};

const handleCreateOrUpdate = async () => {
  formError.value = '';

  if (!categoryForm.value.name) {
    formError.value = 'Category Name is required';
    return;
  }

  try {
    if (formMode.value === 'create') {
      await $fetch(`${config.public.apiBase}/admin/categories`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${adminAuth.token}`
        },
        body: categoryForm.value
      });
      toast.success('Category Created', `"${categoryForm.value.name}" has been created successfully.`);
    } else {
      await $fetch(`${config.public.apiBase}/admin/categories/${editId.value}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${adminAuth.token}`
        },
        body: categoryForm.value
      });
      toast.success('Category Updated', `"${categoryForm.value.name}" has been updated successfully.`);
    }

    resetForm();
    await fetchCategories();
  } catch (err) {
    formError.value = err.data?.message || err.data?.error || 'Failed to sync category';
    console.error(err);
  }
};

const handleEdit = (cat) => {
  formMode.value = 'edit';
  editId.value = cat.id;
  categoryForm.value = {
    name: cat.name,
    description: cat.description || ''
  };
};

const openConfirm = (cat) => {
  categoryToDelete.value = cat;
  showConfirm.value = true;
};

const closeConfirm = () => {
  showConfirm.value = false;
  categoryToDelete.value = null;
};

const closeErrorModal = () => {
  showErrorModal.value = false;
};

const handleDelete = async () => {
  if (!categoryToDelete.value) return;

  const id = categoryToDelete.value.id;
  deletingId.value = id;

  try {
    await $fetch(`${config.public.apiBase}/admin/categories/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${adminAuth.token}`
      }
    });

    toast.success('Category Deleted', `"${categoryToDelete.value.name}" has been deleted successfully.`);
    await fetchCategories();
    closeConfirm();
  } catch (err) {
    closeConfirm();
    const status = err?.statusCode || err?.response?.status;
    const data = err?.data || err?.response?._data || {};

    if (status === 409 && data.code === 'CATEGORY_HAS_ROOMS') {
      errorModalTitle.value = 'Cannot Delete Category';
      errorModalMessage.value = 'This category still has rooms assigned to it.';
      showErrorModal.value = true;
    } else {
      errorModalTitle.value = 'Delete Failed';
      errorModalMessage.value = data.message || 'Unable to delete the category. Please try again later.';
      showErrorModal.value = true;
    }
  } finally {
    deletingId.value = null;
  }
};

const resetForm = () => {
  formMode.value = 'create';
  editId.value = null;
  categoryForm.value = {
    name: '',
    description: ''
  };
};

onMounted(() => {
  fetchCategories();
});
</script>

<template>
  <div class="space-y-8">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-black text-white flex items-center gap-3">
          <FolderTree class="text-blue-500 h-8 w-8 animate-pulse" />
          <span>Room Category Management</span>
        </h1>
        <p class="text-slate-400 mt-1">Configure room categories for your property</p>
      </div>
      <button @click="fetchCategories" class="p-3 text-slate-400 hover:text-white bg-slate-800 rounded-xl transition-all">
        <RefreshCw class="h-5 w-5" :class="{ 'animate-spin': loading }" />
      </button>
    </div>

    <!-- Layout Grid -->
    <div class="grid gap-8 lg:grid-cols-3">
      <!-- Form Panel -->
      <div class="bg-slate-950 p-6 rounded-[2rem] border border-slate-800/80 h-fit space-y-6">
        <h2 class="text-lg font-bold text-white flex items-center gap-2">
          <span>{{ formMode === 'create' ? 'Define New Category' : 'Edit Category' }}</span>
        </h2>

        <form @submit.prevent="handleCreateOrUpdate" class="space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-400 uppercase tracking-wider">Category Name</label>
            <input v-model="categoryForm.name" type="text" placeholder="e.g. Luxury Beachfront" required
                   class="w-full bg-slate-900 border border-slate-800 rounded-xl px-4 py-3 outline-none focus:border-blue-500 text-white placeholder-slate-650" />
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-slate-400 uppercase tracking-wider">Description</label>
            <textarea v-model="categoryForm.description" rows="3" placeholder="Describe inclusions or qualities of this room style..."
                      class="w-full bg-slate-900 border border-slate-800 rounded-xl px-4 py-3 outline-none focus:border-blue-500 text-white text-sm placeholder-slate-650"></textarea>
          </div>

          <div v-if="formError" class="text-xs font-medium text-red-400 bg-red-500/10 border border-red-500/20 p-3 rounded-lg">
            {{ formError }}
          </div>

          <div class="flex items-center gap-3 pt-2">
            <button type="submit" class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 rounded-xl transition-all flex items-center justify-center gap-2">
              <Save class="h-4 w-4" />
              <span>Save Category</span>
            </button>
            <button v-if="formMode === 'edit'" type="button" @click="resetForm" class="bg-slate-800 hover:bg-slate-700 text-slate-350 px-4 py-3 rounded-xl transition-all">
              <X class="h-4 w-4" />
            </button>
          </div>
        </form>
      </div>

      <!-- Categories Table List -->
      <div class="lg:col-span-2 space-y-4">
        <div v-if="loading" class="flex items-center justify-center h-48 bg-slate-950 rounded-[2rem] border border-slate-800/60">
          <div class="h-7 w-7 animate-spin rounded-full border-2 border-blue-500 border-t-transparent"></div>
        </div>

        <div v-else-if="categories.length === 0" class="p-12 text-center bg-slate-950 rounded-[2rem] border border-slate-800/80 text-slate-500">
          No categories configured yet. Add your first room category above!
        </div>

        <div v-else class="grid gap-4 sm:grid-cols-2">
          <div v-for="cat in categories" :key="cat.id"
               class="bg-slate-950 p-6 rounded-[2rem] border border-slate-800 flex flex-col justify-between hover:border-slate-750 transition-all group">
            <div class="space-y-2">
              <div class="flex items-start justify-between">
                <h3 class="text-lg font-bold text-white group-hover:text-blue-400 transition-colors">{{ cat.name }}</h3>
              </div>
              <p class="text-sm text-slate-400 line-clamp-3 leading-relaxed">{{ cat.description || 'No description provided.' }}</p>

              <div class="mt-4 text-[10px] font-bold uppercase tracking-wider text-slate-600">
                Rooms Assigned: {{ cat.rooms?.length || 0 }}
              </div>
            </div>

            <div class="flex items-center gap-2 border-t border-slate-800/60 pt-4 mt-6">
              <button @click="handleEdit(cat)" class="flex-1 bg-slate-900 hover:bg-blue-600/10 text-slate-300 hover:text-blue-400 font-semibold py-2 rounded-xl transition-all text-sm flex items-center justify-center gap-1.5 border border-slate-800 hover:border-blue-500/20">
                <Edit class="h-3.5 w-3.5" />
                <span>Modify</span>
              </button>
              <button
                @click="openConfirm(cat)"
                :disabled="deletingId === cat.id"
                class="bg-slate-900 hover:bg-red-600/15 text-slate-400 hover:text-red-400 p-2 rounded-xl transition-all border border-slate-800 hover:border-red-500/25 disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center h-10 w-10"
              >
                <Loader2 v-if="deletingId === cat.id" class="h-4 w-4 animate-spin" />
                <Trash2 v-else class="h-3.5 w-3.5" />
              </button>
            </div>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex items-center justify-between bg-slate-950 rounded-[2rem] border border-slate-800 px-6 py-4">
          <div class="text-sm text-slate-400">
            Showing <span class="font-bold text-white">{{ categories.length }}</span> of <span class="font-bold text-white">{{ total }}</span> categories
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
    </div>

    <ConfirmModal
      :is-open="showConfirm"
      title="Delete Category"
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
