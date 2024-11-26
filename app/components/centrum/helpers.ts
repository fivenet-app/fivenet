import type { UnitAccess, UnitAccessLevel } from '~~/gen/ts/resources/centrum/access';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import type { Unit } from '~~/gen/ts/resources/centrum/units';
import { StatusUnit } from '~~/gen/ts/resources/centrum/units';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

export type GroupedUnits = { status: StatusUnit; key: string; units: Unit[] }[];

export function dispatchStatusToFillColor(status: StatusDispatch | undefined): string {
    switch (status) {
        case StatusDispatch.UNSPECIFIED:
        case StatusDispatch.NEW:
        case StatusDispatch.UNASSIGNED:
        case StatusDispatch.UNIT_DECLINED:
            return '!text-error-600';
        case StatusDispatch.EN_ROUTE:
            return '!text-info-500';
        case StatusDispatch.ON_SCENE:
            return '!text-info-700';
        case StatusDispatch.NEED_ASSISTANCE:
            return '!text-warn-600';
        case StatusDispatch.COMPLETED:
            return '!text-success-600';
        case StatusDispatch.CANCELLED:
            return '!text-success-800';
        case StatusDispatch.ARCHIVED:
            return '!text-base-600';
        case StatusDispatch.UNIT_ACCEPTED:
            return '!text-info-600';
        default:
            return '!text-info-500';
    }
}

export function dispatchStatusToBGColor(status: StatusDispatch | undefined): string {
    switch (status) {
        case StatusDispatch.UNSPECIFIED:
        case StatusDispatch.NEW:
        case StatusDispatch.UNASSIGNED:
        case StatusDispatch.UNIT_DECLINED:
            return '!bg-error-600';
        case StatusDispatch.EN_ROUTE:
            return '!bg-info-500';
        case StatusDispatch.ON_SCENE:
            return '!bg-info-700';
        case StatusDispatch.NEED_ASSISTANCE:
            return '!bg-warn-600';
        case StatusDispatch.COMPLETED:
            return '!bg-success-600';
        case StatusDispatch.CANCELLED:
            return '!bg-success-800';
        case StatusDispatch.ARCHIVED:
            return '!bg-background';
        case StatusDispatch.UNIT_ACCEPTED:
            return '!bg-info-600';
        default:
            return '!bg-info-500';
    }
}

export const animateStates = [
    StatusDispatch.NEW.valueOf(),
    StatusDispatch.UNASSIGNED.valueOf(),
    StatusDispatch.NEED_ASSISTANCE.valueOf(),
];

export function dispatchStatusAnimate(status: StatusDispatch | undefined): boolean {
    return animateStates.includes((status ?? StatusDispatch.NEW).valueOf());
}

export function unitStatusToBGColor(status: StatusUnit | undefined): string {
    switch (status) {
        case StatusUnit.ON_BREAK:
        case StatusUnit.USER_ADDED:
        case StatusUnit.USER_REMOVED:
            return '!bg-info-500';
        case StatusUnit.AVAILABLE:
            return '!bg-success-600';
        case StatusUnit.BUSY:
            return '!bg-warn-600';
        case StatusUnit.UNSPECIFIED:
        case StatusUnit.UNKNOWN:
        case StatusUnit.UNAVAILABLE:
        default:
            return '!bg-error-600';
    }
}

export const statusOrder = [
    StatusUnit.AVAILABLE,
    StatusUnit.ON_BREAK,
    StatusUnit.BUSY,
    StatusUnit.USER_ADDED,
    StatusUnit.USER_REMOVED,
    StatusUnit.UNAVAILABLE,
    StatusUnit.UNKNOWN,
    StatusUnit.UNSPECIFIED,
];

export const unitStatuses: {
    icon: string;
    name: string;
    status?: StatusUnit;
}[] = [
    { icon: 'i-mdi-cancel', name: 'Unavailable', status: StatusUnit.UNAVAILABLE },
    { icon: 'i-mdi-calendar-check', name: 'Available', status: StatusUnit.AVAILABLE },
    { icon: 'i-mdi-coffee', name: 'On Break', status: StatusUnit.ON_BREAK },
    { icon: 'i-mdi-calendar-remove', name: 'Busy', status: StatusUnit.BUSY },
];

export const dispatchStatuses: {
    icon: string;
    name: string;
    status?: StatusDispatch;
}[] = [
    { icon: 'i-mdi-car-back', name: 'En Route', status: StatusDispatch.EN_ROUTE },
    { icon: 'i-mdi-marker-check', name: 'On Scene', status: StatusDispatch.ON_SCENE },
    { icon: 'i-mdi-help-circle', name: 'Need Assistance', status: StatusDispatch.NEED_ASSISTANCE },
    { icon: 'i-mdi-check-bold', name: 'Completed', status: StatusDispatch.COMPLETED },
    { icon: 'i-mdi-cancel', name: 'Cancelled', status: StatusDispatch.CANCELLED },
];

export function isStatusDispatchCompleted(status: StatusDispatch): boolean {
    return status === StatusDispatch.ARCHIVED || status === StatusDispatch.CANCELLED || status === StatusDispatch.COMPLETED;
}

export function dispatchTimeToTextColor(
    date: Timestamp | undefined,
    status: StatusDispatch = StatusDispatch.UNSPECIFIED,
    maxTime: number = 600,
): string {
    if (isStatusDispatchCompleted(status)) {
        return 'text-success-300';
    }

    // Get passed time in minutes
    const time = (Date.now() - toDate(date).getTime()) / 1000;

    const over = time / maxTime;
    if (over >= 0.85) {
        return 'text-red-700 animate-bounce';
    } else if (over >= 0.7) {
        return 'text-red-400';
    } else if (over >= 0.55) {
        return 'text-orange-400';
    } else if (over >= 0.35) {
        return 'text-orange-300';
    } else if (over >= 0.2) {
        return 'text-yellow-300';
    } else if (over >= 0.1) {
        return 'text-yellow-100';
    }

    return '';
}

export function dispatchTimeToTextColorSidebar(
    date: Timestamp | undefined,
    status: StatusDispatch = StatusDispatch.UNSPECIFIED,
    maxTime: number = 900,
): { ping: boolean; class: string } {
    const time = (Date.now() - toDate(date).getTime()) / 1000;

    if (isStatusDispatchCompleted(status)) {
        return { ping: false, class: '' };
    }

    const over = time / maxTime;
    if (over <= 0.15) {
        return { ping: false, class: '' };
    } else if (over <= 0.2) {
        return { ping: false, class: '!bg-orange-300' };
    } else if (over <= 0.3) {
        return { ping: false, class: '!bg-yellow-300' };
    } else if (over <= 0.5) {
        return { ping: false, class: '!bg-orange-500' };
    } else if (over <= 0.8) {
        return { ping: true, class: '!bg-red-400' };
    }

    return { ping: true, class: '!bg-red-700' };
}

export function checkUnitAccess(unitAccess: UnitAccess | undefined, level: UnitAccessLevel): boolean {
    if (unitAccess === undefined || (unitAccess.jobs.length === 0 && unitAccess.qualifications.length === 0)) {
        return true;
    }

    const { activeChar, isSuperuser } = useAuth();
    if (isSuperuser.value) {
        return true;
    }

    if (activeChar.value === null) {
        return false;
    }

    if (!checkAccess(activeChar.value, unitAccess, undefined, level)) {
        return false;
    }

    return true;
}
