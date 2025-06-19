import type { BadgeColor } from '#ui/types';
import type { Perms } from '~~/gen/ts/perms';
import type { AccessLevel, QualificationAccess } from '~~/gen/ts/resources/qualifications/access';
import type { QualificationRequirement } from '~~/gen/ts/resources/qualifications/qualifications';
import { RequestStatus, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { User, UserShort } from '~~/gen/ts/resources/users/users';

export function checkQualificationAccess(
    qualiAccess: QualificationAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
    perm?: Perms,
): boolean {
    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) {
        return true;
    }

    if (activeChar.value === null) {
        return false;
    }

    if (!checkBaseQualificationAccess(activeChar.value, qualiAccess, creator, level)) {
        return false;
    }

    if (perm !== undefined && creator !== undefined && creator?.job === activeChar.value.job) {
        return checkIfCanAccessOwnJobQualification(activeChar.value, creator, perm);
    }

    return true;
}

function checkBaseQualificationAccess(
    activeChar: UserShort,
    access: QualificationAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
): boolean {
    return checkAccess(activeChar, access, creator, level);
}

function checkIfCanAccessOwnJobQualification(activeChar: User, creator: UserShort, perm: Perms): boolean {
    const { attrStringList } = useAuth();

    const fields = attrStringList(perm, 'Access').value;
    if (fields.length === 0) {
        return creator?.userId === activeChar.userId;
    }

    if (fields.includes('Any')) {
        return true;
    }
    if (fields.includes('LowerRank')) {
        if (creator?.jobGrade < activeChar.jobGrade) {
            return true;
        }
    }
    if (fields.includes('SameRank')) {
        if (creator?.jobGrade <= activeChar.jobGrade) {
            return true;
        }
    }
    if (fields.includes('Own')) {
        if (creator?.userId === activeChar.userId) {
            return true;
        }
    }

    return false;
}

export function requestStatusToBadgeColor(status: RequestStatus): BadgeColor {
    switch (status) {
        case RequestStatus.ACCEPTED:
        case RequestStatus.COMPLETED:
            return 'success';
        case RequestStatus.DENIED:
            return 'error';
        default:
            return 'primary';
    }
}

export function requestStatusToTextColor(status: RequestStatus): string {
    switch (status) {
        case RequestStatus.ACCEPTED:
        case RequestStatus.COMPLETED:
            return 'text-success-400';
        case RequestStatus.DENIED:
            return 'text-error-400';
        default:
            return 'text-info-400';
    }
}

export function requestStatusToBgColor(status: RequestStatus): string {
    switch (status) {
        case RequestStatus.ACCEPTED:
        case RequestStatus.COMPLETED:
            return 'bg-success-400';
        case RequestStatus.DENIED:
            return 'bg-error-400';
        default:
            return 'bg-info-400';
    }
}

export function resultStatusToBadgeColor(status: ResultStatus): BadgeColor {
    switch (status) {
        case ResultStatus.FAILED:
            return 'error';
        case ResultStatus.SUCCESSFUL:
            return 'success';
        default:
            return 'primary';
    }
}

export function resultStatusToTextColor(status: ResultStatus): string {
    switch (status) {
        case ResultStatus.FAILED:
            return 'text-error-400';
        case ResultStatus.SUCCESSFUL:
            return 'text-success-400';
        default:
            return 'text-info-400';
    }
}

export function resultStatusToBgColor(status: ResultStatus): string {
    switch (status) {
        case ResultStatus.FAILED:
            return 'bg-error-400';
        case ResultStatus.SUCCESSFUL:
            return 'bg-success-400';
        default:
            return 'bg-info-400';
    }
}

export function requirementsFullfilled(reqs: QualificationRequirement[]): boolean {
    for (let i = 0; i < reqs.length; i++) {
        const req = reqs[i];
        if (req?.targetQualification?.result?.status !== ResultStatus.SUCCESSFUL) {
            return false;
        }
    }

    return true;
}
