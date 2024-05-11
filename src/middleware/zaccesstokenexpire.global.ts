import { type RouteLocationNormalized } from 'vue-router';
import { useAuthStore } from '~/store/auth';

export default defineNuxtRouteMiddleware(async (to: RouteLocationNormalized, _: RouteLocationNormalized) => {
    const authStore = useAuthStore();

    const expiration = authStore.getAccessTokenExpiration;

    // Check if we have an expiration time, make sure the token isn't expired (yet)
    if (expiration !== null) {
        const now = new Date();
        // Token expired, redirect to login
        if (expiration <= now) {
            console.info('Auth: Token expired, redirecting to login.');
            authStore.setAccessTokenExpiration(null);

            // Only update the redirect query param if it isn't set already
            const redirect = to.query.redirect ?? to.fullPath;
            return navigateTo({
                name: 'auth-login',
                // save the location we were at to come back later
                query: { redirect },
            });
        }
    }
});
