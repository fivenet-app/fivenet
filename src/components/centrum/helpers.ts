import { StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { StatusUnit } from '~~/gen/ts/resources/dispatch/units';

export function dispatchStatusToFillColor(status: StatusDispatch | undefined): string {
    switch (status) {
        case StatusDispatch.UNSPECIFIED:
        case StatusDispatch.NEW:
        case StatusDispatch.UNASSIGNED:
            return 'fill-error-600';
        case StatusDispatch.ON_SCENE:
            return 'fill-primary-600';
        case StatusDispatch.NEED_ASSISTANCE:
            return 'fill-warn-600';
        case StatusDispatch.COMPLETED:
        case StatusDispatch.CANCELLED:
            return 'fill-success-600';
        case StatusDispatch.ARCHIVED:
            return 'fill-base-600';
        case StatusDispatch.EN_ROUTE:
        default:
            return 'fill-info-500';
    }
}

export function dispatchStatusToBGColor(status: StatusDispatch | undefined): string {
    switch (status) {
        case StatusDispatch.UNSPECIFIED:
        case StatusDispatch.NEW:
        case StatusDispatch.UNASSIGNED:
            return 'bg-error-600';
        case StatusDispatch.ON_SCENE:
            return 'bg-primary-600';
        case StatusDispatch.NEED_ASSISTANCE:
            return 'bg-warn-600';
        case StatusDispatch.COMPLETED:
        case StatusDispatch.CANCELLED:
            return 'bg-success-600';
        case StatusDispatch.ARCHIVED:
            return 'bg-base-600';
        case StatusDispatch.EN_ROUTE:
        default:
            return 'bg-info-500';
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
            return 'bg-info-500';
        case StatusUnit.AVAILABLE:
            return 'bg-success-600';
        case StatusUnit.BUSY:
            return 'bg-warn-600';
        case StatusUnit.UNSPECIFIED:
        case StatusUnit.UNKNOWN:
        case StatusUnit.UNAVAILABLE:
        default:
            return 'bg-error-600';
    }
}

export const statusOrder = [
    StatusUnit.AVAILABLE,
    StatusUnit.USER_ADDED,
    StatusUnit.USER_REMOVED,
    StatusUnit.ON_BREAK,
    StatusUnit.BUSY,
    StatusUnit.UNAVAILABLE,
    StatusUnit.UNKNOWN,
    StatusUnit.UNSPECIFIED,
];
