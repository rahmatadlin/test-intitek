import { defineStore } from "pinia";
import api from "../services/api";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    user: null,
    token: localStorage.getItem("token") || null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
  },

  actions: {
    async login(credentials) {
      try {
        const response = await api.post("/auth/login", credentials);
        this.token = response.data.token;
        this.user = response.data.user;
        localStorage.setItem("token", this.token);
        return response.data;
      } catch (error) {
        throw error;
      }
    },

    async register(userData) {
      try {
        const response = await api.post("/auth/register", userData);
        return response.data;
      } catch (error) {
        throw error;
      }
    },

    logout() {
      this.user = null;
      this.token = null;
      localStorage.removeItem("token");
    },

    checkAuth() {
      const token = localStorage.getItem("token");
      if (token) {
        this.token = token;
      }
    },
  },
});
