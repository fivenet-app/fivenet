import { useAuthStore } from '~/store/auth';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { DocumentAccess } from '~~/gen/ts/resources/documents/documents';

export function checkDocAccess(docAccess: DocumentAccess | undefined, level: AccessLevel): boolean {
    if (docAccess === undefined) {
        return false;
    }

    const activeChar = useAuthStore().activeChar;
    if (activeChar === null) {
        return false;
    }

    const ja = docAccess.jobs.find((ja) => ja.job === activeChar.job && ja.minimumGrade <= activeChar.jobGrade);
    if (ja !== undefined && level <= ja.access) {
        return true;
    }

    const ju = docAccess.users.find((ua) => ua.userId === activeChar.userId);
    if (ju !== undefined && level <= ju.access) {
        return true;
    }

    return false;
}
