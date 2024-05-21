import { useAuthStore } from '~/store/auth';
import type { AccessLevel, ThreadAccess } from '~~/gen/ts/resources/messenger/access';
import type { UserShort } from '~~/gen/ts/resources/users/users';

export function canAccessThread(access: ThreadAccess | undefined, creator: UserShort | undefined, level: AccessLevel): boolean {
    const authStore = useAuthStore();
    if (authStore.isSuperuser) {
        return true;
    }

    const activeChar = authStore.activeChar;
    if (activeChar === null) {
        return false;
    }

    if (creator !== undefined && activeChar.userId === creator.userId) {
        return true;
    }

    const ju = access?.users.find((ua) => ua.userId === activeChar.userId && level <= ua.access);
    if (ju !== undefined) {
        return true;
    }

    return false;
}
