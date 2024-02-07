import { useAuthStore } from '~/store/auth';
import slug from '~/utils/slugify';
import type { Perms } from '~~/gen/ts/perms';

export function can(perm: Perms | Perms[], mode?: 'oneof' | 'all'): boolean {
    return checkPerm(perm, mode);
}

function checkPerm(perm: string | string[], mode?: 'oneof' | 'all'): boolean {
    if (mode === undefined) {
        mode = 'oneof';
    }

    const permissions = useAuthStore().permissions;
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

    let can = false;
    // Iterate over permissions and check in "OR" condition manner
    for (let index = 0; index < input.length; index++) {
        const val = slug(input[index] as string);
        if (permissions && (permissions.includes(val) || val === '')) {
            // Permission found
            if (mode === 'oneof') {
                return true;
            }

            can = true;
        } else if (mode === 'all') {
            // Permission not found and mode requires all to be found
            return false;
        }
    }

    return can;
}

export function attr(perm: Perms, name: string, val: string): boolean {
    return checkPerm(perm + '.' + name + (val !== undefined ? '.' + val : ''));
}

export function attrList(perm: Perms, name: string): string[] {
    const key = slug(perm + '.' + name + '.');
    const permissions = useAuthStore().permissions;
    return permissions.filter((p) => p.startsWith(key)).map((p) => p.substring(key.length + 1));
}
