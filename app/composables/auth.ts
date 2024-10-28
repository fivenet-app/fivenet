import { useAuthStore } from '~/store/auth';
import slug from '~/utils/slugify';
import type { Perms } from '~~/gen/ts/perms';

// Wrapper around auth store to make live easier
const _useAuth = () => {
    const authStore = useAuthStore();
    const { activeChar, jobProps, username, isSuperuser, permissions } = storeToRefs(authStore);

    /**
     * @param perm one or more perms to check
     * @param mode default 'oneof'
     * @returns boolean
     */
    const can = (perm: Perms | Perms[], mode: canMode = 'oneof') => {
        return computed(() => {
            if (isSuperuser.value === true) {
                return true;
            }

            return checkPerm(permissions.value, perm, mode);
        });
    };

    const attr = (perm: Perms, name: string, val: string) =>
        computed(() => checkPerm(permissions.value, perm + '.' + name + (val !== undefined ? '.' + val : '')));

    const attrList = (perm: Perms, field: string) =>
        computed(() => {
            const key = slug(perm + '.' + field + '.');
            return permissions.value.filter((p) => p.startsWith(key)).map((p) => p.substring(key.length + 1));
        });

    return {
        // Getters
        activeChar,
        isSuperuser,
        jobProps,
        username,

        // Funcs
        can,
        attr,
        attrList,
    };
};

export const useAuth = createSharedComposable(_useAuth);

export type canMode = 'oneof' | 'all';

function checkPerm(permissions: string[], perm: string | string[], mode: canMode = 'oneof'): boolean {
    if (mode === undefined) {
        mode = 'oneof';
    }

    if (permissions.includes('superuser')) {
        return true;
    }

    const input: string[] = [];
    if (typeof perm === 'string') {
        input.push(perm);
    } else {
        const vals = perm as string[];
        input.push(...vals);
    }

    let ok = false;
    // Iterate over permissions and check in "OR" condition manner
    for (let idx = 0; idx < input.length; idx++) {
        const val = slug(input[idx] as string);
        if (permissions.includes(val) || val === '') {
            // Permission found
            if (mode === 'oneof') {
                return true;
            }

            ok = true;
        } else if (mode === 'all') {
            // Permission not found and mode requires all to be found
            return false;
        }
    }

    return ok;
}
