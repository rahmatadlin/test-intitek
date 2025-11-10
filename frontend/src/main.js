import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import "./assets/styles.css";
import logger from "./utils/logger";

// Initialize logger untuk Tauri
if (typeof window !== 'undefined' && '__TAURI_INTERNALS__' in window) {
  logger.init().catch(console.error);
}

const app = createApp(App);

app.use(createPinia());
app.use(router);

app.mount("#app");
