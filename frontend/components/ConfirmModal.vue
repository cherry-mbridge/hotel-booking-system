<script setup>
import { AlertTriangle, X, Loader2 } from 'lucide-vue-next';

const props = defineProps({
  isOpen: { type: Boolean, required: true },
  title: { type: String, default: 'Confirm Deletion' },
  message: { type: String, default: 'Are you sure you want to delete this item? This action cannot be undone.' },
  confirmText: { type: String, default: 'Delete' },
  cancelText: { type: String, default: 'Cancel' },
  loading: { type: Boolean, default: false },
});

const emit = defineEmits(['confirm', 'cancel']);

const handleConfirm = () => {
  if (!props.loading) emit('confirm');
};

const handleCancel = () => {
  if (!props.loading) emit('cancel');
};
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-150"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
      >
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-black/60 backdrop-blur-sm"
          @click="handleCancel"
        />

        <!-- Modal -->
        <div class="relative w-full max-w-md rounded-3xl bg-slate-950 border border-slate-800 shadow-2xl p-6 sm:p-8">
          <div class="flex flex-col items-center text-center">
            <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-red-500/10 border border-red-500/20 mb-4">
              <AlertTriangle class="h-7 w-7 text-red-400" />
            </div>

            <h3 class="text-xl font-bold text-white">
              {{ title }}
            </h3>
            <p class="mt-2 text-sm text-slate-400 leading-relaxed">
              {{ message }}
            </p>

            <div class="mt-6 flex w-full gap-3">
              <button
                @click="handleCancel"
                :disabled="loading"
                class="flex-1 rounded-xl bg-slate-900 border border-slate-800 px-4 py-3 text-sm font-bold text-slate-300 hover:bg-slate-800 hover:text-white transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ cancelText }}
              </button>
              <button
                @click="handleConfirm"
                :disabled="loading"
                class="flex-1 rounded-xl bg-red-600 px-4 py-3 text-sm font-bold text-white hover:bg-red-700 shadow-lg shadow-red-500/10 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100 flex items-center justify-center gap-2"
              >
                <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
                {{ confirmText }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
