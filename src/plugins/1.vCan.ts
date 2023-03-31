import { useAuthStore } from '../store/auth';
import slug from '../utils/slugify';

// Add `v-can` directive for easy client-side permission checking
export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.vueApp.directive('can', (el, binding, vnode) => {
        // Ignore empty v-can directives
        if (binding.value === "") {
            return;
        }

        const permissions = useAuthStore().$state.permissions;
        const val = slug(binding.value as string);
        if (permissions && (permissions.includes(val) || val === '')) {
            return (vnode.el.hidden = false);
        } else {
            return (vnode.el.hidden = true);
        }
    });
});
