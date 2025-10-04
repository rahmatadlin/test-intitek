import api from "./api";

export const productService = {
  // Get all products with optional filters
  getProducts(params = {}) {
    return api.get("/products", { params });
  },

  // Get single product
  getProduct(id) {
    return api.get(`/products/${id}`);
  },

  // Create new product
  createProduct(productData) {
    return api.post("/products", productData);
  },

  // Update product
  updateProduct(id, productData) {
    return api.put(`/products/${id}`, productData);
  },

  // Delete product
  deleteProduct(id) {
    return api.delete(`/products/${id}`);
  },

  // Get dashboard statistics
  getDashboardStats() {
    return api.get("/dashboard/stats");
  },

  // Export products as CSV
  exportCSV() {
    return api.get("/export/csv", {
      responseType: "blob",
    });
  },

  // Generate barcode
  getBarcode(sku) {
    return api.get(`/barcode/${sku}`, {
      responseType: "blob",
    });
  },
};
