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

    <!-- Search and Filters -->
    <div class="card filters">
      <div class="filter-group">
        <label class="form-label">Search</label>
        <input
          v-model="searchQuery"
          type="text"
          class="input"
          placeholder="Search by name, SKU, or location..."
        />
      </div>
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
              <th class="sortable" @click="sortBy('index')">
                No
                <span class="sort-icon" v-if="sortColumn === 'index'">
                  {{ sortDirection === "asc" ? "↑" : "↓" }}
                </span>
              </th>
              <th class="sortable" @click="sortBy('name')">
                Product Name
                <span class="sort-icon" v-if="sortColumn === 'name'">
                  {{ sortDirection === "asc" ? "↑" : "↓" }}
                </span>
              </th>
              <th class="sortable" @click="sortBy('sku')">
                SKU
                <span class="sort-icon" v-if="sortColumn === 'sku'">
                  {{ sortDirection === "asc" ? "↑" : "↓" }}
                </span>
              </th>
              <th class="sortable" @click="sortBy('quantity')">
                Quantity
                <span class="sort-icon" v-if="sortColumn === 'quantity'">
                  {{ sortDirection === "asc" ? "↑" : "↓" }}
                </span>
              </th>
              <th class="sortable" @click="sortBy('location')">
                Location
                <span class="sort-icon" v-if="sortColumn === 'location'">
                  {{ sortDirection === "asc" ? "↑" : "↓" }}
                </span>
              </th>
              <th class="sortable" @click="sortBy('status')">
                Status
                <span class="sort-icon" v-if="sortColumn === 'status'">
                  {{ sortDirection === "asc" ? "↑" : "↓" }}
                </span>
              </th>
              <th class="sortable" @click="sortBy('created_at')">
                Created At
                <span class="sort-icon" v-if="sortColumn === 'created_at'">
                  {{ sortDirection === "asc" ? "↑" : "↓" }}
                </span>
              </th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(product, index) in paginatedProducts" :key="product.id">
              <td>{{ startIndex + index + 1 }}</td>
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
              <td>{{ formatDate(product.created_at) }}</td>
              <td>
                <div class="action-buttons">
                  <button
                    @click="viewBarcode(product.sku)"
                    class="icon-btn"
                    title="View Barcode"
                  >
                    🔲
                  </button>
                  <button
                    @click="openEditModal(product)"
                    class="icon-btn icon-btn-primary"
                    title="Edit Product"
                  >
                    ✏️
                  </button>
                  <button
                    @click="confirmDelete(product)"
                    class="icon-btn icon-btn-danger"
                    title="Delete Product"
                  >
                    🗑️
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="pagination-container">
        <div class="pagination-info">
          Showing {{ startIndex + 1 }} to {{ endIndex }} of
          {{ filteredProducts.length }} products
        </div>
        <div class="pagination">
          <button
            @click="currentPage = 1"
            :disabled="currentPage === 1"
            class="pagination-btn"
            title="First Page"
          >
            «
          </button>
          <button
            @click="currentPage--"
            :disabled="currentPage === 1"
            class="pagination-btn"
            title="Previous Page"
          >
            ‹
          </button>
          <button
            v-for="page in visiblePages"
            :key="page"
            @click="currentPage = page"
            :class="['pagination-btn', { active: currentPage === page }]"
          >
            {{ page }}
          </button>
          <button
            @click="currentPage++"
            :disabled="currentPage === totalPages"
            class="pagination-btn"
            title="Next Page"
          >
            ›
          </button>
          <button
            @click="currentPage = totalPages"
            :disabled="currentPage === totalPages"
            class="pagination-btn"
            title="Last Page"
          >
            »
          </button>
        </div>
        <div class="pagination-select">
          <label>Items per page:</label>
          <select v-model="itemsPerPage" class="select">
            <option :value="10">10</option>
            <option :value="25">25</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
          </select>
        </div>
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

            <div class="form-info">
              <p style="color: #6b7280; font-size: 0.875rem">
                ℹ️ Status will be automatically set based on quantity:
                <br />
                • 0 units = Out of Stock
                <br />
                • 1-5 units = Low Stock
                <br />
                • &gt;5 units = In Stock
              </p>
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

    <!-- Delete Confirmation Modal -->
    <div
      v-if="showDeleteModal"
      class="modal-overlay"
      @click.self="closeDeleteModal"
    >
      <div class="modal modal-small">
        <div class="modal-header">
          <h2 class="modal-title">Confirm Delete</h2>
          <button @click="closeDeleteModal" class="btn btn-sm btn-outline">
            ✕
          </button>
        </div>
        <div class="modal-body">
          <div class="delete-confirmation">
            <div class="delete-icon">🗑️</div>
            <p class="delete-message">
              Are you sure you want to delete
              <strong>{{ productToDelete?.name }}</strong
              >?
            </p>
            <p class="delete-warning">This action cannot be undone.</p>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeDeleteModal" class="btn btn-secondary">
            Cancel
          </button>
          <button
            @click="deleteProduct"
            class="btn btn-danger"
            :disabled="formLoading"
          >
            {{ formLoading ? "Deleting..." : "Delete" }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from "vue";
import { productService } from "../services/productService";
import { useToast } from "../composables/useToast";

const toast = useToast();
const products = ref([]);
const loading = ref(true);
const showModal = ref(false);
const showBarcodeModal = ref(false);
const showDeleteModal = ref(false);
const isEditing = ref(false);
const formLoading = ref(false);
const formError = ref("");
const currentBarcodeSKU = ref("");
const barcodeUrl = ref("");
const productToDelete = ref(null);

const filters = ref({
  status: "",
  lowStock: false,
});

// Search state
const searchQuery = ref("");

// Sorting state
const sortColumn = ref("created_at");
const sortDirection = ref("desc");

// Pagination state
const currentPage = ref(1);
const itemsPerPage = ref(10);

const formData = ref({
  id: null,
  name: "",
  sku: "",
  quantity: 0,
  location: "",
});

// Computed property for filtered products (search)
const filteredProducts = computed(() => {
  let filtered = [...products.value];

  // Apply search filter
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(
      (product) =>
        product.name.toLowerCase().includes(query) ||
        product.sku.toLowerCase().includes(query) ||
        product.location.toLowerCase().includes(query)
    );
  }

  return filtered;
});

// Computed property for sorted products
const sortedProducts = computed(() => {
  const productsCopy = [...filteredProducts.value];

  productsCopy.sort((a, b) => {
    let aVal, bVal;

    if (sortColumn.value === "index") {
      // For index, just return as-is since we'll number them in the template
      return 0;
    } else if (sortColumn.value === "created_at") {
      aVal = new Date(a.created_at).getTime();
      bVal = new Date(b.created_at).getTime();
    } else if (sortColumn.value === "quantity") {
      aVal = a.quantity;
      bVal = b.quantity;
    } else {
      aVal = a[sortColumn.value]?.toString().toLowerCase() || "";
      bVal = b[sortColumn.value]?.toString().toLowerCase() || "";
    }

    if (sortDirection.value === "asc") {
      return aVal > bVal ? 1 : aVal < bVal ? -1 : 0;
    } else {
      return aVal < bVal ? 1 : aVal > bVal ? -1 : 0;
    }
  });

  return productsCopy;
});

// Computed property for total pages
const totalPages = computed(() => {
  return Math.ceil(sortedProducts.value.length / itemsPerPage.value) || 1;
});

// Computed property for paginated products
const paginatedProducts = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value;
  const end = start + itemsPerPage.value;
  return sortedProducts.value.slice(start, end);
});

// Computed property for visible page numbers
const visiblePages = computed(() => {
  const pages = [];
  const total = totalPages.value;
  const current = currentPage.value;
  const delta = 2; // Show 2 pages before and after current

  let start = Math.max(1, current - delta);
  let end = Math.min(total, current + delta);

  // Adjust if we're near the beginning or end
  if (current <= delta) {
    end = Math.min(total, delta * 2 + 1);
  }
  if (current >= total - delta) {
    start = Math.max(1, total - delta * 2);
  }

  for (let i = start; i <= end; i++) {
    pages.push(i);
  }

  return pages;
});

// Computed properties for pagination info
const startIndex = computed(() => {
  return (currentPage.value - 1) * itemsPerPage.value;
});

const endIndex = computed(() => {
  return Math.min(
    startIndex.value + itemsPerPage.value,
    filteredProducts.value.length
  );
});

// Watch for changes in itemsPerPage and reset to page 1
watch(itemsPerPage, () => {
  currentPage.value = 1;
});

// Watch for changes in search query and reset to page 1
watch(searchQuery, () => {
  currentPage.value = 1;
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
    toast.error("Failed to load products");
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
      toast.success("Product updated successfully!");
    } else {
      await productService.createProduct(formData.value);
      toast.success("Product created successfully!");
    }
    closeModal();
    fetchProducts();
  } catch (error) {
    const errorMsg = error.response?.data?.error || "Failed to save product";
    formError.value = errorMsg;
    toast.error(errorMsg);
  } finally {
    formLoading.value = false;
  }
};

const confirmDelete = (product) => {
  productToDelete.value = product;
  showDeleteModal.value = true;
};

const closeDeleteModal = () => {
  showDeleteModal.value = false;
  productToDelete.value = null;
};

const deleteProduct = async () => {
  if (!productToDelete.value) return;

  formLoading.value = true;

  try {
    await productService.deleteProduct(productToDelete.value.id);
    toast.success("Product deleted successfully!");
    closeDeleteModal();
    fetchProducts();
  } catch (error) {
    toast.error("Failed to delete product");
  } finally {
    formLoading.value = false;
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
    toast.success("Products exported to CSV successfully!");
  } catch (error) {
    toast.error("Failed to export CSV");
  }
};

const viewBarcode = async (sku) => {
  try {
    const response = await productService.getBarcode(sku);
    barcodeUrl.value = window.URL.createObjectURL(new Blob([response.data]));
    currentBarcodeSKU.value = sku;
    showBarcodeModal.value = true;
  } catch (error) {
    toast.error("Failed to generate barcode");
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

const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const sortBy = (column) => {
  if (sortColumn.value === column) {
    // Toggle direction if clicking same column
    sortDirection.value = sortDirection.value === "asc" ? "desc" : "asc";
  } else {
    // Set new column and default to ascending
    sortColumn.value = column;
    sortDirection.value = "asc";
  }
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

/* Table improvements */
.table-container {
  overflow-x: auto;
  border-radius: 0.5rem;
  border: 1px solid var(--border-color);
}

table {
  width: 100%;
  min-width: 1200px; /* Ensure table is wide enough */
  border-collapse: collapse;
  table-layout: auto; /* Allow columns to adjust based on content */
}

thead {
  background: var(--bg-gray);
  position: sticky;
  top: 0;
  z-index: 10;
}

th {
  padding: 0.875rem 1rem;
  text-align: left;
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-color);
  white-space: nowrap; /* Prevent header text from wrapping */
}

/* Specific column widths */
th:nth-child(1) {
  width: 60px; /* No */
}

th:nth-child(2) {
  width: 250px; /* Product Name */
  min-width: 200px;
}

th:nth-child(3) {
  width: 140px; /* SKU */
}

th:nth-child(4) {
  width: 100px; /* Quantity */
}

th:nth-child(5) {
  width: 200px; /* Location */
  min-width: 180px;
}

th:nth-child(6) {
  width: 120px; /* Status */
}

th:nth-child(7) {
  width: 180px; /* Created At */
}

th:nth-child(8) {
  width: 180px; /* Actions */
  min-width: 180px;
}

td {
  padding: 0.875rem 1rem;
  border-bottom: 1px solid var(--border-color);
  font-size: 0.875rem;
  vertical-align: middle;
}

/* Prevent text wrapping in specific columns */
td:nth-child(1),
td:nth-child(3),
td:nth-child(4),
td:nth-child(6),
td:nth-child(7) {
  white-space: nowrap;
}

/* Allow wrapping for longer text columns with ellipsis */
td:nth-child(2),
td:nth-child(5) {
  max-width: 250px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

tbody tr {
  transition: background-color 0.2s;
}

tbody tr:hover {
  background: var(--bg-gray);
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
  flex-wrap: nowrap;
  justify-content: center;
}

/* Icon-only buttons */
.icon-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.25rem;
  padding: 0.25rem;
  transition: all 0.2s;
  border-radius: 0.25rem;
  opacity: 0.7;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
}

.icon-btn:hover {
  opacity: 1;
  transform: scale(1.1);
}

.icon-btn-primary:hover {
  background: rgba(79, 70, 229, 0.1);
}

.icon-btn-danger:hover {
  background: rgba(239, 68, 68, 0.1);
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
  white-space: nowrap;
}

/* Sortable table headers */
.sortable {
  cursor: pointer;
  user-select: none;
  position: relative;
  transition: background-color 0.2s;
}

.sortable:hover {
  background-color: #e5e7eb;
}

.sort-icon {
  margin-left: 0.5rem;
  font-size: 0.875rem;
  color: var(--primary-color);
  font-weight: bold;
}

/* Pagination styles */
.pagination-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-top: 1px solid var(--border-color);
  flex-wrap: wrap;
  gap: 1rem;
}

.pagination-info {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.pagination {
  display: flex;
  gap: 0.25rem;
}

.pagination-btn {
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--border-color);
  background: white;
  color: var(--text-primary);
  cursor: pointer;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  transition: all 0.2s;
  min-width: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pagination-btn:hover:not(:disabled) {
  background: var(--bg-gray);
  border-color: var(--primary-color);
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-btn.active {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.pagination-select {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
}

.pagination-select label {
  color: var(--text-secondary);
  white-space: nowrap;
}

.pagination-select .select {
  width: auto;
  padding: 0.375rem 0.5rem;
}

/* Delete Confirmation Modal */
.modal-small {
  max-width: 400px;
}

.delete-confirmation {
  text-align: center;
  padding: 1rem 0;
}

.delete-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.delete-message {
  font-size: 1rem;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
  line-height: 1.5;
}

.delete-message strong {
  color: var(--danger-color);
}

.delete-warning {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-style: italic;
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
    flex-direction: row;
    gap: 0.25rem;
  }

  .pagination-container {
    flex-direction: column;
    gap: 1rem;
  }

  .pagination {
    justify-content: center;
  }

  .pagination-info,
  .pagination-select {
    width: 100%;
    justify-content: center;
  }
}
</style>
