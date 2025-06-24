import type { Perms } from '~~/gen/ts/perms';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import type { User } from '~~/gen/ts/resources/users/users';

export function checkIfCanAccessColleague(target: Colleague | User, perm: Perms): boolean {
    const { attrStringList, activeChar, isSuperuser } = useAuth();

    if (!activeChar.value) {
        return false;
    }

    if (isSuperuser.value) {
        return true;
    }
    console.log('checkIfCanAccessColleague 3');
    const fields = attrStringList(perm, 'Access').value;
    console.log('checkIfCanAccessColleague 4', fields);
    if (fields.includes('Any')) {
        return true;
    }
    if (fields.includes('Lower_Rank')) {
        if (target.jobGrade < activeChar.value.jobGrade) {
            return true;
        }
    }
    if (fields.includes('Same_Rank')) {
        if (target.jobGrade <= activeChar.value.jobGrade) {
            return true;
        }
    }
    if (fields.includes('Own')) {
        if (target.userId === activeChar.value.userId) {
            return true;
        }
    }

    return target.userId === activeChar.value.userId;
}
