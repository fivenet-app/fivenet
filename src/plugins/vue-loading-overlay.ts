import { LoadingPlugin } from 'vue-loading-overlay';

export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.vueApp.use(LoadingPlugin);
});
