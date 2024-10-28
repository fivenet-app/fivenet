import type { AccessLevel, ThreadAccess } from '~~/gen/ts/resources/messenger/access';
import type { UserShort } from '~~/gen/ts/resources/users/users';

export function canAccessThread(access: ThreadAccess | undefined, creator: UserShort | undefined, level: AccessLevel): boolean {
    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) {
        return true;
    }

    if (activeChar.value === null) {
        return false;
    }

    if (creator !== undefined && activeChar.value.userId === creator.userId) {
        return true;
    }

    const ju = access?.users.find((ua) => ua.userId === activeChar.value?.userId && level <= ua.access);
    if (ju !== undefined) {
        return true;
    }

    return false;
}
