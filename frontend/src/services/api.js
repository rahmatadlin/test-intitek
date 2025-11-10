import axios from "axios";
import logger from "../utils/logger";

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
    
    // Log API request
    logger.apiRequest(config.method.toUpperCase(), config.url, config.data);
    
    return config;
  },
  (error) => {
    logger.apiError('REQUEST', error.config?.url || 'unknown', error);
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    // Log successful API response
    logger.apiResponse(
      response.config.method.toUpperCase(),
      response.config.url,
      response.status,
      response.data
    );
    return response;
  },
  (error) => {
    // Log API error
    logger.apiError(
      error.config?.method?.toUpperCase() || 'UNKNOWN',
      error.config?.url || 'unknown',
      error
    );
    
    if (error.response?.status === 401) {
      logger.warn("Unauthorized - redirecting to login");
      localStorage.removeItem("token");
      // Navigate to login page
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export default api;
