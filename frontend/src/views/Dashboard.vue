<template>
  <div class="dashboard">
    <h1 class="page-title">Dashboard</h1>

    <div v-if="loading" class="spinner"></div>

    <div v-else>
      <!-- Stats Cards -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon" style="background: #dbeafe">📦</div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total_products }}</div>
            <div class="stat-label">Total Products</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon" style="background: #dcfce7">📊</div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total_stock }}</div>
            <div class="stat-label">Total Stock</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon" style="background: #fef3c7">⚠️</div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.low_stock_count }}</div>
            <div class="stat-label">Low Stock Items</div>
          </div>
        </div>
      </div>

      <!-- Low Stock Products -->
      <div class="card mt-4">
        <h2 class="section-title">Low Stock Products</h2>

        <div
          v-if="stats.low_stock_products && stats.low_stock_products.length > 0"
        >
          <div class="table-container">
            <table>
              <thead>
                <tr>
                  <th>Product Name</th>
                  <th>SKU</th>
                  <th>Quantity</th>
                  <th>Location</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="product in stats.low_stock_products"
                  :key="product.id"
                >
                  <td>{{ product.name }}</td>
                  <td>
                    <code>{{ product.sku }}</code>
                  </td>
                  <td>{{ product.quantity }}</td>
                  <td>{{ product.location }}</td>
                  <td>
                    <span class="badge badge-warning">Low Stock</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div v-else class="empty-state">
          <p>✅ No low stock items. All products are well stocked!</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { productService } from "../services/productService";
import { useToast } from "../composables/useToast";

const toast = useToast();
const stats = ref({
  total_products: 0,
  total_stock: 0,
  low_stock_count: 0,
  low_stock_products: [],
});

const loading = ref(true);

const fetchStats = async () => {
  try {
    const response = await productService.getDashboardStats();
    stats.value = response.data;
  } catch (error) {
    console.error("Error fetching stats:", error);
    toast.error("Failed to load dashboard statistics");
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchStats();
});
</script>

<style scoped>
.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 2rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  border-radius: 0.5rem;
  padding: 1.5rem;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 1rem;
}

.stat-icon {
  width: 3rem;
  height: 3rem;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 1.875rem;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 1.5rem;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: var(--text-secondary);
}

code {
  background: var(--bg-gray);
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
  font-family: monospace;
  font-size: 0.875rem;
}

.mt-4 {
  margin-top: 2rem;
}
</style>
