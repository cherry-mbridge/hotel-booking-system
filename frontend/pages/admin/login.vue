<script setup>
import { ShieldCheck } from 'lucide-vue-next';

definePageMeta({
  layout: false
});

// Use Admin Auth guard store
const adminAuthStore = useAdminAuth();
const router = useRouter();

const email = ref('admin@lumina.com');
const password = ref('admin123');
const loading = ref(false);
const error = ref('');

const handleLogin = async () => {
  loading.value = true;
  error.value = '';
  try {
    await adminAuthStore.login(email.value, password.value);
    router.push('/admin');
  } catch (e) {
    error.value = 'Invalid admin credentials or unauthorized user.';
  } finally {
    loading.value = false;
  }
};
</script>


<template>
  <div class="flex min-h-[60vh] items-center justify-center">
    <div class="w-full max-w-md space-y-8 rounded-[2.5rem] bg-slate-930 text-white p-10 shadow-2xl ring-1 ring-slate-800">
      <div class="text-center">
        <div class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-blue-500/10 text-blue-400">
          <ShieldCheck class="h-8 w-8" />
        </div>
        <h1 class="text-3xl font-black tracking-tight">Admin Gateway</h1>
        <p class="mt-2 text-slate-400">Identify yourself to access the administration suite</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-6">
        <div class="space-y-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-400">Control Panel Email</label>
          <input v-model="email" type="email" required class="w-full rounded-2xl border border-slate-850 bg-slate-900 p-4 text-white outline-none focus:ring-2 focus:ring-blue-500/20 transition-all placeholder-slate-600" placeholder="admin@lumina.com" />
        </div>
        <div class="space-y-2">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-400">Security Phrase</label>
          <input v-model="password" type="password" required class="w-full rounded-2xl border border-slate-850 bg-slate-900 p-4 text-white outline-none focus:ring-2 focus:ring-blue-500/20 transition-all placeholder-slate-600" placeholder="••••••••" />
        </div>

        <div v-if="error" class="text-sm font-medium text-red-400 text-center bg-red-500/10 py-3 rounded-xl border border-red-500/20">{{ error }}</div>

        <button type="submit" :disabled="loading" class="w-full rounded-2xl bg-blue-600 py-4 text-lg font-bold text-white shadow-lg shadow-blue-500/20 transition-all hover:bg-blue-700 active:scale-[0.98] disabled:opacity-50">
          {{ loading ? 'Authenticating...' : 'Enter Dashboard' }}
        </button>
      </form>
      
      <p class="text-center text-sm text-slate-500">
        Standard customer? <NuxtLink to="/login" class="font-bold text-blue-400 underline">Customer Login</NuxtLink>
      </p>
    </div>
  </div>
</template>

<style scoped>
.bg-slate-930 {
  background-color: #0f172a;
}
.border-slate-850 {
  border-color: #1e293b;
}
</style>
