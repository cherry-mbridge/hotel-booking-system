<script setup>
import { AlertCircle, X } from 'lucide-vue-next';

const props = defineProps({
  isOpen: { type: Boolean, required: true },
  title: { type: String, default: 'Error' },
  message: { type: String, default: 'Something went wrong.' },
  buttonText: { type: String, default: 'Got it' },
});

const emit = defineEmits(['close']);
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
          @click="emit('close')"
        />

        <!-- Modal -->
        <div class="relative w-full max-w-md rounded-3xl bg-slate-950 border border-slate-800 shadow-2xl p-6 sm:p-8">
          <div class="flex flex-col items-center text-center">
            <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-amber-500/10 border border-amber-500/20 mb-4">
              <AlertCircle class="h-7 w-7 text-amber-400" />
            </div>

            <h3 class="text-xl font-bold text-white">
              {{ title }}
            </h3>
            <p class="mt-2 text-sm text-slate-400 leading-relaxed">
              {{ message }}
            </p>

            <div class="mt-6 flex w-full">
              <button
                @click="emit('close')"
                class="w-full rounded-xl bg-amber-600 px-4 py-3 text-sm font-bold text-white hover:bg-amber-700 shadow-lg shadow-amber-500/10 transition-all active:scale-95"
              >
                {{ buttonText }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
