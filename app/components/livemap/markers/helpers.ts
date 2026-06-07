import type { UserShort } from '~~/gen/ts/resources/users/short/user';

export function checkIfCanEditMarker(creator: UserShort | null | undefined): boolean {
    const { activeChar, attrStringList, isSuperuser } = useAuth();
    if (isSuperuser.value) return true;

    if (!activeChar.value || !creator) return false;

    const fields = attrStringList('livemap.LivemapService/CreateOrUpdateMarker', 'Access').value;
    if (fields.length === 0) {
        return creator.userId === activeChar.value.userId;
    }

    if (fields.includes('Any')) return true;

    if (fields.includes('Lower_Rank')) {
        if (creator.jobGrade < activeChar.value.jobGrade) return true;
    }
    if (fields.includes('Same_Rank')) {
        if (creator.jobGrade <= activeChar.value.jobGrade) return true;
    }
    if (fields.includes('Own')) {
        if (creator.userId === activeChar.value.userId) return true;
    }

    return false;
}
