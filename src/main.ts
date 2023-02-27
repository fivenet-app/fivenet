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

router.beforeEach((to, from ) => {
  store.state
  if (to.meta.requiresAuth && !store.state.accessToken) {
    return {
      path: '/login',
      // save the location we were at to come back later
      query: { redirect: to.fullPath },
    }
  }
});

const app = createApp(App);
app.use(router);
app.use(store);
app.mount("#app");
