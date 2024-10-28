import type { Perms } from '~~/gen/ts/perms';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { User } from '~~/gen/ts/resources/users/users';

export function checkIfCanAccessColleague(target: Colleague | User, perm: Perms): boolean {
    const { attrList, activeChar, isSuperuser } = useAuth();

    if (!activeChar.value) {
        return false;
    }

    if (isSuperuser.value) {
        return true;
    }

    const fields = attrList(perm, 'Access').value;
    if (fields.includes('any')) {
        return true;
    }
    if (fields.includes('lower_rank')) {
        if (target.jobGrade < activeChar.value.jobGrade) {
            return true;
        }
    }
    if (fields.includes('same_rank')) {
        if (target.jobGrade <= activeChar.value.jobGrade) {
            return true;
        }
    }
    if (fields.includes('own')) {
        if (target.userId === activeChar.value.userId) {
            return true;
        }
    }

    return target.userId === activeChar.value.userId;
}
