<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1 class="login-title">📦 Warehouse Management</h1>
        <p class="login-subtitle">Sign in to manage your inventory</p>
      </div>

      <div v-if="error" class="alert alert-error">
        {{ error }}
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label class="form-label">Username</label>
          <input
            v-model="credentials.username"
            type="text"
            class="input"
            placeholder="Enter your username"
            required
          />
        </div>

        <div class="form-group">
          <label class="form-label">Password</label>
          <input
            v-model="credentials.password"
            type="password"
            class="input"
            placeholder="Enter your password"
            required
          />
        </div>

        <button
          type="submit"
          class="btn btn-primary btn-full"
          :disabled="loading"
        >
          {{ loading ? "Signing in..." : "Sign In" }}
        </button>
      </form>

      <div class="login-footer">
        <p class="demo-info">
          <strong>Demo Credentials:</strong><br />
          Username: <code>admin</code><br />
          Password: <code>admin123</code>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";
import { useToast } from "../composables/useToast";

const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

const credentials = ref({
  username: "",
  password: "",
});

const loading = ref(false);
const error = ref("");

const handleLogin = async () => {
  loading.value = true;
  error.value = "";

  try {
    await authStore.login(credentials.value);
    toast.success("Login successful! Welcome back.");
    router.push("/");
  } catch (err) {
    const errorMsg =
      err.response?.data?.error || "Failed to login. Please try again.";
    error.value = errorMsg;
    toast.error(errorMsg);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 1rem;
}

.login-card {
  background: white;
  border-radius: 1rem;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1),
    0 10px 10px -5px rgba(0, 0, 0, 0.04);
  width: 100%;
  max-width: 400px;
  padding: 2rem;
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.login-subtitle {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.login-form {
  margin-bottom: 1.5rem;
}

.btn-full {
  width: 100%;
  justify-content: center;
}

.login-footer {
  text-align: center;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
}

.demo-info {
  font-size: 0.875rem;
  color: var(--text-secondary);
  line-height: 1.6;
}

.demo-info code {
  background: var(--bg-gray);
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
  font-family: monospace;
  color: var(--primary-color);
}
</style>
