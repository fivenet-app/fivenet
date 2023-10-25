import { attr } from '~/composables/can';
import { useAuthStore } from '~/store/auth';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { DocumentAccess } from '~~/gen/ts/resources/documents/documents';
import { UserShort } from '~~/gen/ts/resources/users/users';

export function checkDocAccess(
    docAccess: DocumentAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
    perm?: string,
): boolean {
    const authStore = useAuthStore();
    if (authStore.isSuperuser) {
        return true;
    }

    if (docAccess === undefined) {
        return false;
    }

    const activeChar = authStore.activeChar;
    if (activeChar === null) {
        return false;
    }
    if (activeChar.userId === creator?.userId) {
        return true;
    }

    const ju = docAccess.users.find((ua) => ua.userId === activeChar.userId);
    if (ju !== undefined && level <= ju.access) {
        return true;
    }

    const ja = docAccess.jobs.find((ja) => ja.job === activeChar.job && ja.minimumGrade <= activeChar.jobGrade);
    if (ja !== undefined && level <= ja.access) {
        if (creator?.job === activeChar.job) {
            if (perm === undefined) {
                return false;
            }

            if (activeChar.jobGrade === creator?.jobGrade) {
                return attr(perm, 'Access', 'Same_Rank');
            } else if (activeChar.jobGrade > creator?.jobGrade) {
                return attr(perm, 'Access', 'Lower_Rank');
            } else {
                return attr(perm, 'Access', 'Own');
            }
        }

        return true;
    }

    return false;
}
