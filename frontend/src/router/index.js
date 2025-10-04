import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "../stores/auth";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("../views/Login.vue"),
  },
  {
    path: "/",
    component: () => import("../views/Layout.vue"),
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        name: "Dashboard",
        component: () => import("../views/Dashboard.vue"),
      },
      {
        path: "/products",
        name: "Products",
        component: () => import("../views/Products.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next("/login");
  } else if (to.path === "/login" && authStore.isAuthenticated) {
    next("/");
  } else {
    next();
  }
});

export default router;
