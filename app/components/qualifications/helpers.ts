import type { BadgeProps } from '@nuxt/ui';
import type { Perms } from '~~/gen/ts/perms';
import type { Access } from '~~/gen/ts/resources/access/access';
import type { AccessLevel } from '~~/gen/ts/resources/qualifications/access/access';
import { QualificationActivityType } from '~~/gen/ts/resources/qualifications/activity/activity';
import { type QualificationRequirement, RequestStatus, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import type { User } from '~~/gen/ts/resources/users/user';

export function qualificationActivityTypeColor(type: QualificationActivityType | undefined): string {
    switch (type) {
        case QualificationActivityType.CREATED:
        case QualificationActivityType.REQUEST_CREATED:
        case QualificationActivityType.RESULT_CREATED:
        case QualificationActivityType.RESULT_RESTORED:
            return 'text-success-500!';
        case QualificationActivityType.DELETED:
        case QualificationActivityType.REQUEST_DELETED:
        case QualificationActivityType.RESULT_DELETED:
        case QualificationActivityType.EXAM_CANCELLED:
        case QualificationActivityType.EXAM_EXPIRED:
            return 'text-error-500!';
        case QualificationActivityType.UPDATED:
        case QualificationActivityType.ACCESS_UPDATED:
        case QualificationActivityType.REQUEST_UPDATED:
        case QualificationActivityType.RESULT_UPDATED:
            return 'text-warning-500!';
        default:
            return 'text-primary-500!';
    }
}

export function qualificationActivityTypeIcon(type: QualificationActivityType | undefined): string {
    switch (type) {
        case QualificationActivityType.CREATED:
            return 'i-mdi-plus-circle-outline';
        case QualificationActivityType.DELETED:
            return 'i-mdi-delete-outline';
        case QualificationActivityType.RESTORED:
        case QualificationActivityType.RESULT_RESTORED:
            return 'i-mdi-restore';
        case QualificationActivityType.ACCESS_UPDATED:
            return 'i-mdi-lock-outline';
        case QualificationActivityType.REQUEST_CREATED:
        case QualificationActivityType.REQUEST_UPDATED:
        case QualificationActivityType.REQUEST_DELETED:
            return 'i-mdi-email-outline';
        case QualificationActivityType.EXAM_STARTED:
        case QualificationActivityType.EXAM_SUBMITTED:
        case QualificationActivityType.EXAM_CANCELLED:
        case QualificationActivityType.EXAM_EXPIRED:
            return 'i-mdi-test-tube';
        case QualificationActivityType.RESULT_CREATED:
        case QualificationActivityType.RESULT_UPDATED:
        case QualificationActivityType.RESULT_DELETED:
            return 'i-mdi-list-status';
        default:
            return 'i-mdi-pencil-outline';
    }
}

export function checkQualificationAccess(
    qualiAccess: Access | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
    perm?: Perms,
    creatorJob?: string,
): boolean {
    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) {
        return true;
    }

    if (activeChar.value === null) {
        return false;
    }

    if (!checkBaseQualificationAccess(activeChar.value, qualiAccess, creator, level, creatorJob)) {
        return false;
    }

    if (perm !== undefined && creator !== undefined && (creatorJob ?? creator?.job) === activeChar.value.job) {
        return checkIfCanAccessOwnJobQualification(activeChar.value, creator, perm);
    }

    return true;
}

function checkBaseQualificationAccess(
    activeChar: UserShort,
    access: Access | undefined,
    creator: UserShort | undefined,
    level: AccessLevel,
    creatorJob?: string,
): boolean {
    return checkAccess(activeChar, access, creator, level, creatorJob);
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
    if (fields.includes('Lower_Rank')) {
        if (creator?.jobGrade < activeChar.jobGrade) {
            return true;
        }
    }
    if (fields.includes('Same_Rank')) {
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

export function requestStatusToBadgeColor(status: RequestStatus): BadgeProps['color'] {
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

export function requestStatusToTextColor(status: RequestStatus | undefined): string {
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

export function resultStatusToBadgeColor(status: ResultStatus): BadgeProps['color'] {
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
