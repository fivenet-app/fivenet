import { useAuthStore } from '~/store/auth';
import slug from '~/utils/slugify';

// Add `v-can` directive for easy client-side permission checking
export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.vueApp.directive('can', (el, binding, vnode) => {
        // Ignore undefined/ empty permissions
        if (!binding.value || binding.value === '') {
            vnode.el.hidden = false;
            return;
        }

        vnode.el.hidden = !can(binding.value);
        return;
    });
});

export function can(perm: string | string[]): boolean {
    const permissions = useAuthStore().permissions;
    if (permissions.includes('superuser')) {
        return true;
    } else {
        const input = new Array<String>();
        if (typeof perm === 'string') {
            input.push(perm);
        } else {
            const vals = perm as Array<String>;
            input.push(...vals);
        }

        // Iterate over permissions and check in "OR" condition manner
        for (let index = 0; index < input.length; index++) {
            const val = slug(input[index] as string);
            if (permissions && (permissions.includes(val) || val === '')) {
                return true;
            }
        }
    }

    return false;
}
