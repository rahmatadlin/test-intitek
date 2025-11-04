import axios from "axios";

// Detect if running in Tauri
const isTauri = typeof window !== "undefined" && "__TAURI_INTERNALS__" in window;

// API base URL - can be configured via environment variable
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api";

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor to add token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Debug logging for Tauri
    if (typeof window !== "undefined" && "__TAURI_INTERNALS__" in window) {
      console.error("API Error:", error);
      console.error("Error message:", error.message);
      console.error("Error response:", error.response);
      console.error("Request URL:", error.config?.url);
      console.error("Request baseURL:", error.config?.baseURL);
    }
    
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      // Navigate to login page
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export default api;
