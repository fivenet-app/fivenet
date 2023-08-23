import { useAuthStore } from '~/store/auth';
import slug from '~/utils/slugify';

export function can(perm: string): boolean {
    const permissions = useAuthStore().permissions;
    if (permissions.includes('superuser')) {
        return true;
    } else {
        const val = slug(perm as string);
        if (permissions && (permissions.includes(val) || val === '')) {
            return true;
        }
    }

    return false;
}
