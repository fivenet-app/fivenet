import { useAuthStore } from '~/store/auth';
import slug from '~/utils/slugify';
import type { Perms } from '~~/gen/ts/perms';

export type canMode = 'oneof' | 'all';

/**
 *
 * @param perm one or more perms t ocheck
 * @param mode default 'oneof'
 * @returns boolean
 */
export const can = reactify((perm: Perms | Perms[], mode?: canMode): boolean => {
    const { permissions, isSuperuser } = storeToRefs(useAuthStore());

    if (isSuperuser.value === true) {
        return true;
    }

    return checkPermRef(permissions, perm, mode).value;
});

const checkPermRef = reactify(checkPerm);

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
    for (let index = 0; index < input.length; index++) {
        const val = slug(input[index] as string);
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

export const attr = reactify(attrRef);

function attrRef(perm: Perms, name: string, val: string): boolean {
    const { permissions } = storeToRefs(useAuthStore());

    return checkPermRef(permissions.value, perm + '.' + name + (val !== undefined ? '.' + val : '')).value;
}

export const attrList = reactify((perm: Perms, name: string): string[] => {
    const { permissions } = storeToRefs(useAuthStore());
    return attrListRef(permissions.value, perm, name).value;
});

const attrListRef = reactify(attrListFn);

function attrListFn(permissions: string[], perm: Perms, name: string): string[] {
    const key = slug(perm + '.' + name + '.');
    return permissions.filter((p) => p.startsWith(key)).map((p) => p.substring(key.length + 1));
}
