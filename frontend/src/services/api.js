import axios from "axios";

// Detect if running in Wails (desktop app) or browser
// In Wails, API server runs on localhost:8080 in background
// Check if we're in Wails by checking for wails runtime
const isWails = typeof window !== "undefined" && window.runtime;

// Always use localhost:8080 for API calls
// In Wails desktop app, the Gin server runs on localhost:8080
// In browser development, it also runs on localhost:8080
const baseURL = "http://localhost:8080/api";

const api = axios.create({
  baseURL: baseURL,
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
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      // Use router for navigation instead of window.location
      // Import router dynamically to avoid circular dependency
      import("../router").then(({ default: router }) => {
        router.push("/login");
      }).catch(() => {
        // Fallback to window.location if router import fails
        window.location.hash = "#/login";
      });
    }
    return Promise.reject(error);
  }
);

export default api;
