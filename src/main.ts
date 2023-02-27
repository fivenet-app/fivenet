import { createApp } from "vue";

import "./index.css";
import "@fontsource/inter";

import App from "./App.vue";
import { store } from "./store";

const app = createApp(App);
app.use(store);
app.mount("#app");
