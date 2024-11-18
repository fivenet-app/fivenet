import type { Access, AccessLevel } from '~~/gen/ts/resources/mailer/access';

export function canAccess(access: Access | undefined, creatorId: number | undefined, level: AccessLevel): boolean {
    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) {
        return true;
    }

    if (activeChar.value === null) {
        return false;
    }

    if (creatorId !== undefined && activeChar.value.userId === creatorId) {
        return true;
    }

    const ju = access?.users.find((ua) => ua.userId === activeChar.value?.userId && level <= ua.access);
    if (ju !== undefined) {
        return true;
    }

    return false;
}
