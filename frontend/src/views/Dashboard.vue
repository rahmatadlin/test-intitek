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

      <!-- Inventory Chart -->
      <div class="card mt-4">
        <h2 class="section-title">Top 10 Products by Quantity</h2>
        <p class="section-subtitle">
          Products are color-coded by status:
          <span class="legend-item">
            <span
              class="legend-dot"
              style="background: rgb(16, 185, 129)"
            ></span>
            In Stock
          </span>
          <span class="legend-item">
            <span
              class="legend-dot"
              style="background: rgb(245, 158, 11)"
            ></span>
            Low Stock
          </span>
          <span class="legend-item">
            <span
              class="legend-dot"
              style="background: rgb(239, 68, 68)"
            ></span>
            Out of Stock
          </span>
        </p>
        <div class="chart-wrapper">
          <BarChart :data="chartData" :horizontal="true" />
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
import { ref, computed, onMounted } from "vue";
import { productService } from "../services/productService";
import { useToast } from "../composables/useToast";
import BarChart from "../components/BarChart.vue";

const toast = useToast();
const stats = ref({
  total_products: 0,
  total_stock: 0,
  low_stock_count: 0,
  low_stock_products: [],
});

const allProducts = ref([]);
const loading = ref(true);

// Computed property for chart data
const chartData = computed(() => {
  // Sort products by quantity (descending) and take top 10
  const sortedProducts = [...allProducts.value]
    .sort((a, b) => b.quantity - a.quantity)
    .slice(0, 10);

  // Extract product names and quantities
  const productNames = sortedProducts.map((p) => p.name);
  const quantities = sortedProducts.map((p) => p.quantity);

  // Color based on status
  const backgroundColors = sortedProducts.map((p) => {
    if (p.status === "in_stock") return "rgba(16, 185, 129, 0.8)"; // Green
    if (p.status === "low_stock") return "rgba(245, 158, 11, 0.8)"; // Orange
    return "rgba(239, 68, 68, 0.8)"; // Red for out_of_stock
  });

  const borderColors = sortedProducts.map((p) => {
    if (p.status === "in_stock") return "rgb(16, 185, 129)";
    if (p.status === "low_stock") return "rgb(245, 158, 11)";
    return "rgb(239, 68, 68)";
  });

  return {
    labels: productNames,
    datasets: [
      {
        label: "Quantity",
        data: quantities,
        backgroundColor: backgroundColors,
        borderColor: borderColors,
        borderWidth: 2,
        borderRadius: 6,
      },
    ],
  };
});

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

const fetchAllProducts = async () => {
  try {
    const response = await productService.getProducts();
    allProducts.value = response.data.data || [];
  } catch (error) {
    console.error("Error fetching products:", error);
  }
};

onMounted(() => {
  fetchStats();
  fetchAllProducts();
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
  margin-bottom: 0.5rem;
}

.section-subtitle {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-bottom: 1.5rem;
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  align-items: center;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
}

.legend-dot {
  width: 12px;
  height: 12px;
  border-radius: 3px;
  display: inline-block;
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

.chart-wrapper {
  height: 500px;
  padding: 1rem;
}

.mt-4 {
  margin-top: 2rem;
}

@media (max-width: 768px) {
  .chart-wrapper {
    height: 450px;
  }

  .section-subtitle {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
}
</style>
