import { useAuthStore } from '~/store/auth';
import type { Perms } from '~~/gen/ts/perms';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { User, UserShort } from '~~/gen/ts/resources/users/users';

export function checkIfCanAccessColleague(activeChar: UserShort | User, target: Colleague | User, perm: Perms): boolean {
    const authStore = useAuthStore();
    if (authStore.isSuperuser) {
        return true;
    }

    const fields = attrList(perm, 'Access').value;
    if (fields.includes('any')) {
        return true;
    }
    if (fields.includes('lower_rank')) {
        if (target.jobGrade < activeChar.jobGrade) {
            return true;
        }
    }
    if (fields.includes('same_rank')) {
        if (target.jobGrade <= activeChar.jobGrade) {
            return true;
        }
    }
    if (fields.includes('own')) {
        if (target.userId === activeChar.userId) {
            return true;
        }
    }

    return target.userId === activeChar.userId;
}
