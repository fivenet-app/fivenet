import { type BadgeColor } from '#ui/types';
import { useAuthStore } from '~/store/auth';
import type { Perms } from '~~/gen/ts/perms';
import { AccessLevel, QualificationAccess } from '~~/gen/ts/resources/qualifications/access';
import { QualificationRequirement, RequestStatus, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { User, UserShort } from '~~/gen/ts/resources/users/users';

export function checkQualificationAccess(
    qualiAccess: QualificationAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
    perm?: Perms,
): boolean {
    const authStore = useAuthStore();
    if (authStore.isSuperuser) {
        return true;
    }

    const activeChar = authStore.activeChar;
    if (activeChar === null) {
        return false;
    }

    if (!checkBaseQualificationAccess(activeChar, qualiAccess, creator, level)) {
        return false;
    }

    if (perm !== undefined && creator !== undefined && creator?.job === activeChar.job) {
        return checkIfCanAccessOwnJobQualification(activeChar, creator, perm);
    }

    return true;
}

function checkBaseQualificationAccess(
    activeChar: User,
    docAccess: QualificationAccess | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
): boolean {
    if (docAccess === undefined) {
        return false;
    }

    if (creator !== undefined && activeChar.userId === creator.userId) {
        return true;
    }

    const ja = docAccess.jobs.find(
        (ja) => ja.job === activeChar.job && ja.minimumGrade <= activeChar.jobGrade && level <= ja.access,
    );
    if (ja !== undefined) {
        return true;
    }

    return false;
}

function checkIfCanAccessOwnJobQualification(activeChar: User, creator: UserShort, perm: Perms): boolean {
    const authStore = useAuthStore();
    if (authStore.isSuperuser) {
        return true;
    }

    const fields = attrList(perm, 'Access').value;
    if (fields.length === 0) {
        return creator?.userId === activeChar.userId;
    }

    if (fields.includes('any')) {
        return true;
    }
    if (fields.includes('lower_rank')) {
        if (creator?.jobGrade < activeChar.jobGrade) {
            return true;
        }
    }
    if (fields.includes('same_rank')) {
        if (creator?.jobGrade <= activeChar.jobGrade) {
            return true;
        }
    }
    if (fields.includes('own')) {
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
            return 'green';
        case RequestStatus.DENIED:
            return 'red';
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

export function resultStatusToBadgeColor(status: ResultStatus): BadgeColor {
    switch (status) {
        case ResultStatus.FAILED:
            return 'red';
        case ResultStatus.SUCCESSFUL:
            return 'green';
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

export function requirementsFullfilled(reqs: QualificationRequirement[]): boolean {
    for (let i = 0; i < reqs.length; i++) {
        const req = reqs[i];
        if (req?.targetQualification?.result?.status !== ResultStatus.SUCCESSFUL) {
            return false;
        }
    }

    return true;
}
