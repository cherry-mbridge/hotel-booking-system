<script setup>
const authStore = useUserAuth();
const router = useRouter();

const name = ref('');
const email = ref('');
const password = ref('');
const loading = ref(false);
const error = ref('');

const handleRegister = async () => {
  loading.value = true;
  error.value = '';
  try {
    await authStore.register(name.value, email.value, password.value);
    router.push('/');
  } catch (e) {
    error.value = 'Registration failed. This email might already be taken.';
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="flex min-h-[70vh] items-center justify-center">
    <div class="w-full max-w-md space-y-8 rounded-[2.5rem] bg-white p-10 shadow-xl ring-1 ring-slate-100">
      <div class="text-center">
        <h1 class="text-3xl font-black tracking-tight text-slate-900">Create Account</h1>
        <p class="mt-2 text-slate-500">Join Lumina Hotel & Resorts</p>
      </div>

      <form @submit.prevent="handleRegister" class="space-y-5">
        <div class="space-y-1">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-400">Full Name</label>
          <input v-model="name" type="text" required class="w-full rounded-2xl border border-slate-200 bg-slate-50 p-4 outline-none focus:ring-2 focus:ring-blue-500/20 transition-all" placeholder="John Doe" />
        </div>
        
        <div class="space-y-1">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-400">Email Address</label>
          <input v-model="email" type="email" required class="w-full rounded-2xl border border-slate-200 bg-slate-50 p-4 outline-none focus:ring-2 focus:ring-blue-500/20 transition-all" placeholder="john@example.com" />
        </div>

        <div class="space-y-1">
          <label class="text-xs font-bold uppercase tracking-widest text-slate-400">Password</label>
          <input v-model="password" type="password" required class="w-full rounded-2xl border border-slate-200 bg-slate-50 p-4 outline-none focus:ring-2 focus:ring-blue-500/20 transition-all" placeholder="••••••••" />
        </div>

        <div v-if="error" class="text-sm font-medium text-red-500 text-center">{{ error }}</div>

        <button type="submit" :disabled="loading" class="w-full mt-4 rounded-2xl bg-blue-600 py-4 text-lg font-bold text-white shadow-lg shadow-blue-200 transition-all hover:bg-blue-700 active:scale-[0.98] disabled:opacity-50">
          {{ loading ? 'Creating account...' : 'Sign Up' }}
        </button>
      </form>
      
      <p class="text-center text-sm text-slate-400">
        Already have an account? <NuxtLink to="/login" class="font-bold text-blue-600">Login</NuxtLink>
      </p>
    </div>
  </div>
</template>
