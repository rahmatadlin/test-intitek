import { ref } from "vue";

const toasts = ref([]);
let toastId = 0;

export function useToast() {
  const addToast = (message, type = "success", duration = 3000) => {
    const id = toastId++;
    const toast = {
      id,
      message,
      type, // 'success', 'error', 'warning', 'info'
      visible: true,
    };

    toasts.value.push(toast);

    // Auto remove after duration
    setTimeout(() => {
      removeToast(id);
    }, duration);
  };

  const removeToast = (id) => {
    const index = toasts.value.findIndex((t) => t.id === id);
    if (index > -1) {
      toasts.value.splice(index, 1);
    }
  };

  const success = (message, duration) => {
    addToast(message, "success", duration);
  };

  const error = (message, duration) => {
    addToast(message, "error", duration);
  };

  const warning = (message, duration) => {
    addToast(message, "warning", duration);
  };

  const info = (message, duration) => {
    addToast(message, "info", duration);
  };

  return {
    toasts,
    addToast,
    removeToast,
    success,
    error,
    warning,
    info,
  };
}
