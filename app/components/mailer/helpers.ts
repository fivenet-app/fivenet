import type { Access, AccessLevel } from '~~/gen/ts/resources/mailer/access';
import type { Thread } from '~~/gen/ts/resources/mailer/thread';

export const defaultEmailDomain = 'fivenet.ls';

export const threadResponseTitlePrefix = 'RE: ';
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

export function generateResponseTitle(thread?: Thread): string {
    if (!thread) {
        return '';
    }

    if (thread.title.startsWith(threadResponseTitlePrefix)) {
        return thread.title;
    }

    return threadResponseTitlePrefix + thread.title;
}
