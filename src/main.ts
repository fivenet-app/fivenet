import { createApp } from 'vue';
import App from './App.vue';
import store from './store';
import { createWebHistory, createRouter, setupDataFetchingGuard, RouteRecordRaw } from 'vue-router/auto';
import jwt_decode from 'jwt-decode';

import './style.css';
import '@fontsource/inter';

const router = createRouter({
	history: createWebHistory(),
	extendRoutes: (routes: RouteRecordRaw[]) => {
		return routes;
	},
});
setupDataFetchingGuard(router);

router.beforeEach((to, from) => {
	if (to.meta.requiresAuth && !store.state.accessToken) {
		var decoded = jwt_decode(store.state.accessToken as string);
		console.log(JSON.stringify(decoded));

		return {
			path: '/login',
			// save the location we were at to come back later
			query: { redirect: to.fullPath },
		};
	}
});
export default router;

const app = createApp(App);
app.use(router);
app.use(store);
app.mount('#app');
