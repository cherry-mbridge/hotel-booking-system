<script setup>
import { X, CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-vue-next';
import { useToast } from '~/composables/useToast';

const { toasts, removeToast } = useToast();

const typeConfig = {
  success: {
    icon: CheckCircle,
    bg: 'bg-emerald-500/10',
    border: 'border-emerald-500/20',
    text: 'text-emerald-400',
    iconColor: 'text-emerald-400',
  },
  error: {
    icon: AlertCircle,
    bg: 'bg-red-500/10',
    border: 'border-red-500/20',
    text: 'text-red-400',
    iconColor: 'text-red-400',
  },
  warning: {
    icon: AlertTriangle,
    bg: 'bg-amber-500/10',
    border: 'border-amber-500/20',
    text: 'text-amber-400',
    iconColor: 'text-amber-400',
  },
  info: {
    icon: Info,
    bg: 'bg-blue-500/10',
    border: 'border-blue-500/20',
    text: 'text-blue-400',
    iconColor: 'text-blue-400',
  },
};
</script>

<template>
  <Teleport to="body">
    <div class="fixed top-6 right-6 z-[100] flex flex-col gap-3 pointer-events-none">
      <TransitionGroup
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="translate-x-full opacity-0"
        enter-to-class="translate-x-0 opacity-100"
        leave-active-class="transition-all duration-200 ease-in"
        leave-from-class="translate-x-0 opacity-100"
        leave-to-class="translate-x-full opacity-0"
      >
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="pointer-events-auto w-80 sm:w-96 rounded-2xl border backdrop-blur-md shadow-xl p-4 flex items-start gap-3"
          :class="[typeConfig[toast.type].bg, typeConfig[toast.type].border]"
        >
          <component
            :is="typeConfig[toast.type].icon"
            class="h-5 w-5 shrink-0 mt-0.5"
            :class="typeConfig[toast.type].iconColor"
          />
          <div class="flex-1 min-w-0">
            <p v-if="toast.title" class="text-sm font-bold text-white">
              {{ toast.title }}
            </p>
            <p v-if="toast.message" class="text-sm text-slate-300 mt-0.5">
              {{ toast.message }}
            </p>
          </div>
          <button
            @click="removeToast(toast.id)"
            class="shrink-0 text-slate-500 hover:text-white transition-colors"
          >
            <X class="h-4 w-4" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>
