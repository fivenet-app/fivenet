import { useAuthStore } from '~/store/auth';
import slug from '~/utils/slugify';

// Add `v-can` directive for easy client-side permission checking
export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.vueApp.directive('can', (el, binding, vnode) => {
        // Ignore undefined/ empty permissions
        if (!binding.value || binding.value === '') {
            return;
        }

        const permissions = useAuthStore().getPermissions;
        if (permissions.includes('superuser')) {
            return (vnode.el.hidden = false);
        }

        const input = new Array<String>();
        if (typeof binding.value === 'string') {
            input.push(binding.value);
        } else {
            const vals = binding.value as Array<String>;
            input.push(...vals);
        }

        // Iterate over permissions and check in "OR" condition manner
        for (let index = 0; index < input.length; index++) {
            const val = slug(input[index] as string);
            if (permissions && (permissions.includes(val) || val === '')) {
                return (vnode.el.hidden = false);
            }
        }

        return (vnode.el.hidden = true);
    });
});
