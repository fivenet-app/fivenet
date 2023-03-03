import { createApp } from 'vue';
import App from './App.vue';
import store from './store';
import { createWebHistory, createRouter, setupDataFetchingGuard, RouteRecordRaw } from 'vue-router/auto';
import './style.css';
import '@fontsource/inter';
import { dispatchNotification } from './components/Notification';

const router = createRouter({
	history: createWebHistory(),
	extendRoutes: (routes: RouteRecordRaw[]) => {
		return routes;
	},
});
setupDataFetchingGuard(router);
router.beforeEach((to, from) => {
	// Default is that a page requires authentication
	if (!to.meta.hasOwnProperty('requiresAuth') || to.meta.requiresAuth) {
		// Check if user has access token
		if (store.state.accessToken && store.state.activeChar) {
			// Route has permission attached to it, check if user has required permission
			if (to.meta.permission) {
				if (store.state.permissions.includes(to.meta.permission)) {
					// User has permission
					return;
				} else {
					dispatchNotification({
						title: "You don't have permission!",
						content: "You don't have permission to go to " + to.path + '.',
						type: 'warning',
					});

                    return {
                        path: '/overview'
                    };
				}
			} else {
                // No route permission check, so we can go ahead and return
				return;
			}
		}

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

// Add `v-can` directive for easy permission checking
app.directive('can', (el, binding, vnode) => {
	var permissions = store.state.permissions;
	if (permissions.includes(binding.value)) {
		return (vnode.el.hidden = false);
	} else {
		return (vnode.el.hidden = true);
	}
});

app.mount('#app');
