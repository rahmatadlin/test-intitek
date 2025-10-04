<template>
  <div class="layout">
    <nav class="navbar">
      <div class="navbar-container">
        <div class="navbar-brand">
          <h1 class="navbar-title">📦 Warehouse Management</h1>
        </div>
        <div class="navbar-menu">
          <router-link to="/" class="nav-link">Dashboard</router-link>
          <router-link to="/products" class="nav-link">Products</router-link>
          <button @click="handleLogout" class="btn btn-secondary btn-sm">
            Logout
          </button>
        </div>
      </div>
    </nav>

    <main class="main-content">
      <div class="container">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup>
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";

const router = useRouter();
const authStore = useAuthStore();

const handleLogout = () => {
  authStore.logout();
  router.push("/login");
};
</script>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar {
  background: white;
  border-bottom: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
}

.navbar-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 1rem 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.navbar-brand {
  display: flex;
  align-items: center;
}

.navbar-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--text-primary);
}

.navbar-menu {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.nav-link {
  color: var(--text-secondary);
  text-decoration: none;
  font-weight: 500;
  font-size: 0.875rem;
  transition: color 0.2s;
}

.nav-link:hover {
  color: var(--primary-color);
}

.nav-link.router-link-active {
  color: var(--primary-color);
}

.main-content {
  flex: 1;
  padding: 2rem 0;
}

.container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 1.5rem;
}

@media (max-width: 768px) {
  .navbar-container {
    flex-direction: column;
    gap: 1rem;
  }

  .navbar-menu {
    width: 100%;
    justify-content: center;
  }
}
</style>
