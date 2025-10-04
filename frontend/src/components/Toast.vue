<template>
  <div class="toast-container">
    <transition-group name="toast">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        :class="['toast-item', `toast-${toast.type}`]"
        @click="removeToast(toast.id)"
      >
        <div class="toast-icon">
          <span v-if="toast.type === 'success'">✓</span>
          <span v-else-if="toast.type === 'error'">✕</span>
          <span v-else-if="toast.type === 'warning'">⚠</span>
          <span v-else>ℹ</span>
        </div>
        <div class="toast-message">{{ toast.message }}</div>
        <button class="toast-close" @click.stop="removeToast(toast.id)">
          ✕
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script setup>
import { useToast } from "../composables/useToast";

const { toasts, removeToast } = useToast();
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-width: 400px;
}

.toast-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 300px;
  cursor: pointer;
  transition: all 0.3s ease;
  animation: slideIn 0.3s ease;
}

.toast-item:hover {
  transform: translateX(-5px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.toast-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 14px;
  flex-shrink: 0;
}

.toast-message {
  flex: 1;
  font-size: 14px;
  line-height: 1.4;
}

.toast-close {
  width: 20px;
  height: 20px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 16px;
  opacity: 0.6;
  transition: opacity 0.2s;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.toast-close:hover {
  opacity: 1;
}

/* Success Toast */
.toast-success {
  background: #10b981;
  color: white;
}

.toast-success .toast-icon {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

/* Error Toast */
.toast-error {
  background: #ef4444;
  color: white;
}

.toast-error .toast-icon {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

/* Warning Toast */
.toast-warning {
  background: #f59e0b;
  color: white;
}

.toast-warning .toast-icon {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

/* Info Toast */
.toast-info {
  background: #3b82f6;
  color: white;
}

.toast-info .toast-icon {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

/* Animations */
@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.toast-enter-active {
  animation: slideIn 0.3s ease;
}

.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

@media (max-width: 768px) {
  .toast-container {
    left: 10px;
    right: 10px;
    max-width: none;
  }

  .toast-item {
    min-width: auto;
  }
}
</style>
