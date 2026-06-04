import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import type { User } from '~~/gen/ts/resources/users/user';

export function checkIfCanEditMarker(activeChar: User | null, creator: UserShort | null | undefined): boolean {
    if (!activeChar || !creator) return false;

    const { attrStringList } = useAuth();

    const fields = attrStringList('livemap.LivemapService/CreateOrUpdateMarker', 'Access').value;
    if (fields.length === 0) {
        return creator.userId === activeChar.userId;
    }

    if (fields.includes('Any')) return true;

    if (fields.includes('Lower_Rank')) {
        if (creator.jobGrade < activeChar.jobGrade) return true;
    }
    if (fields.includes('Same_Rank')) {
        if (creator.jobGrade <= activeChar.jobGrade) return true;
    }
    if (fields.includes('Own')) {
        if (creator.userId === activeChar.userId) return true;
    }

    return false;
}
