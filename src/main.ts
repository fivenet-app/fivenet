import { createApp } from "vue";
import App from "./App.vue";

import "./style.css";
import "@fontsource/inter";

import { store } from "./store";

import {
  createWebHistory,
  createRouter,
  setupDataFetchingGuard,
  RouteRecordRaw,
} from "vue-router/auto";

const router = createRouter({
  history: createWebHistory(),
  extendRoutes: (routes: RouteRecordRaw[]) => {
    return routes;
  },
});

setupDataFetchingGuard(router);

const app = createApp(App);
app.use(router);
app.use(store);
app.mount("#app");
