import { reactive } from 'vue'

const toasts = reactive([])

let nextId = 1

export function useToast() {
  const addToast = (options) => {
    const id = nextId++
    const toast = {
      id,
      title: options.title || '',
      message: options.message || '',
      type: options.type || 'info',
      duration: options.duration ?? 4000,
    }
    toasts.push(toast)

    if (toast.duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, toast.duration)
    }

    return id
  }

  const removeToast = (id) => {
    const index = toasts.findIndex((t) => t.id === id)
    if (index > -1) {
      toasts.splice(index, 1)
    }
  }

  const success = (title, message, duration) =>
    addToast({ title, message, type: 'success', duration })

  const error = (title, message, duration) =>
    addToast({ title, message, type: 'error', duration })

  const warning = (title, message, duration) =>
    addToast({ title, message, type: 'warning', duration })

  const info = (title, message, duration) =>
    addToast({ title, message, type: 'info', duration })

  return {
    toasts,
    addToast,
    removeToast,
    success,
    error,
    warning,
    info,
  }
}
