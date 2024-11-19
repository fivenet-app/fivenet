import type { Access, AccessLevel } from '~~/gen/ts/resources/mailer/access';

export const defaultEmailDomain = '@fivenet.app';

export const defaultEmptyContent = '<p><br></p><p><br></p>';

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

    return checkAccess(activeChar.value, access, undefined, level);
}
