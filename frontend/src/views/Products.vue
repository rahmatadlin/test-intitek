<template>
  <div class="products">
    <div class="page-header">
      <h1 class="page-title">Products</h1>
      <div class="page-actions">
        <button @click="exportCSV" class="btn btn-outline">
          📥 Export CSV
        </button>
        <button @click="openCreateModal" class="btn btn-primary">
          ➕ Add Product
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="card filters">
      <div class="filter-group">
        <label class="form-label">Filter by Status</label>
        <select v-model="filters.status" @change="fetchProducts" class="select">
          <option value="">All</option>
          <option value="in_stock">In Stock</option>
          <option value="low_stock">Low Stock</option>
          <option value="out_of_stock">Out of Stock</option>
        </select>
      </div>
      <div class="filter-group">
        <label class="form-label">Show Low Stock Only</label>
        <input
          type="checkbox"
          v-model="filters.lowStock"
          @change="fetchProducts"
          class="checkbox"
        />
      </div>
    </div>

    <div v-if="loading" class="spinner"></div>

    <!-- Products Table -->
    <div v-else class="card">
      <div v-if="products.length === 0" class="empty-state">
        <p>No products found. Add your first product to get started!</p>
      </div>

      <div v-else class="table-container">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Product Name</th>
              <th>SKU</th>
              <th>Quantity</th>
              <th>Location</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="product in products" :key="product.id">
              <td>{{ product.id }}</td>
              <td>{{ product.name }}</td>
              <td>
                <code>{{ product.sku }}</code>
              </td>
              <td>{{ product.quantity }}</td>
              <td>{{ product.location }}</td>
              <td>
                <span
                  class="badge"
                  :class="{
                    'badge-success': product.status === 'in_stock',
                    'badge-warning': product.status === 'low_stock',
                    'badge-danger': product.status === 'out_of_stock',
                  }"
                >
                  {{ formatStatus(product.status) }}
                </span>
              </td>
              <td>
                <div class="action-buttons">
                  <button
                    @click="viewBarcode(product.sku)"
                    class="btn btn-sm btn-outline"
                    title="View Barcode"
                  >
                    🔲
                  </button>
                  <button
                    @click="openEditModal(product)"
                    class="btn btn-sm btn-primary"
                  >
                    ✏️ Edit
                  </button>
                  <button
                    @click="deleteProduct(product.id)"
                    class="btn btn-sm btn-danger"
                  >
                    🗑️ Delete
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Product Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h2 class="modal-title">
            {{ isEditing ? "Edit Product" : "Add New Product" }}
          </h2>
          <button @click="closeModal" class="btn btn-sm btn-outline">✕</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveProduct">
            <div class="form-group">
              <label class="form-label">Product Name *</label>
              <input
                v-model="formData.name"
                type="text"
                class="input"
                required
              />
            </div>

            <div class="form-group">
              <label class="form-label">SKU *</label>
              <input
                v-model="formData.sku"
                type="text"
                class="input"
                required
              />
            </div>

            <div class="form-group">
              <label class="form-label">Quantity *</label>
              <input
                v-model.number="formData.quantity"
                type="number"
                class="input"
                min="0"
                required
              />
            </div>

            <div class="form-group">
              <label class="form-label">Location *</label>
              <input
                v-model="formData.location"
                type="text"
                class="input"
                placeholder="e.g., Warehouse A, Shelf 12"
                required
              />
            </div>

            <div class="form-group">
              <label class="form-label">Status *</label>
              <select v-model="formData.status" class="select" required>
                <option value="in_stock">In Stock</option>
                <option value="low_stock">Low Stock</option>
                <option value="out_of_stock">Out of Stock</option>
              </select>
            </div>

            <div v-if="formError" class="alert alert-error">
              {{ formError }}
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn btn-secondary">Cancel</button>
          <button
            @click="saveProduct"
            class="btn btn-primary"
            :disabled="formLoading"
          >
            {{ formLoading ? "Saving..." : "Save" }}
          </button>
        </div>
      </div>
    </div>

    <!-- Barcode Modal -->
    <div
      v-if="showBarcodeModal"
      class="modal-overlay"
      @click.self="closeBarcodeModal"
    >
      <div class="modal">
        <div class="modal-header">
          <h2 class="modal-title">Barcode - {{ currentBarcodeSKU }}</h2>
          <button @click="closeBarcodeModal" class="btn btn-sm btn-outline">
            ✕
          </button>
        </div>
        <div class="modal-body text-center">
          <img :src="barcodeUrl" alt="Barcode" style="max-width: 100%" />
        </div>
        <div class="modal-footer">
          <button @click="closeBarcodeModal" class="btn btn-secondary">
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { productService } from "../services/productService";

const products = ref([]);
const loading = ref(true);
const showModal = ref(false);
const showBarcodeModal = ref(false);
const isEditing = ref(false);
const formLoading = ref(false);
const formError = ref("");
const currentBarcodeSKU = ref("");
const barcodeUrl = ref("");

const filters = ref({
  status: "",
  lowStock: false,
});

const formData = ref({
  id: null,
  name: "",
  sku: "",
  quantity: 0,
  location: "",
  status: "in_stock",
});

const fetchProducts = async () => {
  loading.value = true;
  try {
    const params = {};
    if (filters.value.status) params.status = filters.value.status;
    if (filters.value.lowStock) params.low_stock = "true";

    const response = await productService.getProducts(params);
    products.value = response.data.data || [];
  } catch (error) {
    console.error("Error fetching products:", error);
  } finally {
    loading.value = false;
  }
};

const openCreateModal = () => {
  isEditing.value = false;
  formData.value = {
    id: null,
    name: "",
    sku: "",
    quantity: 0,
    location: "",
    status: "in_stock",
  };
  formError.value = "";
  showModal.value = true;
};

const openEditModal = (product) => {
  isEditing.value = true;
  formData.value = { ...product };
  formError.value = "";
  showModal.value = true;
};

const closeModal = () => {
  showModal.value = false;
};

const saveProduct = async () => {
  formLoading.value = true;
  formError.value = "";

  try {
    if (isEditing.value) {
      await productService.updateProduct(formData.value.id, formData.value);
    } else {
      await productService.createProduct(formData.value);
    }
    closeModal();
    fetchProducts();
  } catch (error) {
    formError.value = error.response?.data?.error || "Failed to save product";
  } finally {
    formLoading.value = false;
  }
};

const deleteProduct = async (id) => {
  if (!confirm("Are you sure you want to delete this product?")) return;

  try {
    await productService.deleteProduct(id);
    fetchProducts();
  } catch (error) {
    alert("Failed to delete product");
  }
};

const exportCSV = async () => {
  try {
    const response = await productService.exportCSV();
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", "products.csv");
    document.body.appendChild(link);
    link.click();
    link.remove();
  } catch (error) {
    alert("Failed to export CSV");
  }
};

const viewBarcode = async (sku) => {
  try {
    const response = await productService.getBarcode(sku);
    barcodeUrl.value = window.URL.createObjectURL(new Blob([response.data]));
    currentBarcodeSKU.value = sku;
    showBarcodeModal.value = true;
  } catch (error) {
    alert("Failed to generate barcode");
  }
};

const closeBarcodeModal = () => {
  showBarcodeModal.value = false;
  barcodeUrl.value = "";
  currentBarcodeSKU.value = "";
};

const formatStatus = (status) => {
  return status.replace("_", " ").replace(/\b\w/g, (l) => l.toUpperCase());
};

onMounted(() => {
  fetchProducts();
});
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
}

.page-actions {
  display: flex;
  gap: 1rem;
}

.filters {
  display: flex;
  gap: 2rem;
  align-items: flex-end;
  margin-bottom: 1.5rem;
}

.filter-group {
  flex: 1;
  max-width: 250px;
}

.checkbox {
  width: 1.25rem;
  height: 1.25rem;
  cursor: pointer;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
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

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .page-actions {
    width: 100%;
  }

  .filters {
    flex-direction: column;
  }

  .filter-group {
    max-width: 100%;
  }

  .action-buttons {
    flex-direction: column;
  }
}
</style>
