import { StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { StatusUnit } from '~~/gen/ts/resources/dispatch/units';

export function dispatchStatusToFillColor(status: StatusDispatch | undefined): string {
    switch (status) {
        case StatusDispatch.NEW:
            return 'fill-error-600';
        case StatusDispatch.UNASSIGNED:
            return 'fill-error-600';
        case StatusDispatch.EN_ROUTE:
            return 'fill-info-500';
        case StatusDispatch.ON_SCENE:
            return 'fill-primary-600';
        case StatusDispatch.NEED_ASSISTANCE:
            return 'fill-warn-600';
        case StatusDispatch.COMPLETED:
            return 'fill-success-600';
        case StatusDispatch.CANCELLED:
            return 'fill-success-600';
        case StatusDispatch.ARCHIVED:
            return 'fill-base-600';
        default:
            return 'fill-info-500';
    }
}

export function dispatchStatusToBGColor(status: StatusDispatch | undefined): string {
    switch (status) {
        case StatusDispatch.NEW:
            return 'bg-error-600';
        case StatusDispatch.UNASSIGNED:
            return 'bg-error-600';
        case StatusDispatch.EN_ROUTE:
            return 'bg-info-500';
        case StatusDispatch.ON_SCENE:
            return 'bg-primary-600';
        case StatusDispatch.NEED_ASSISTANCE:
            return 'bg-warn-600';
        case StatusDispatch.COMPLETED:
            return 'bg-success-600';
        case StatusDispatch.CANCELLED:
            return 'bg-success-600';
        case StatusDispatch.ARCHIVED:
            return 'bg-base-600';
        default:
            return 'bg-info-500';
    }
}

export function unitStatusToBGColor(status: StatusUnit | undefined): string {
    switch (status) {
        case StatusUnit.UNKNOWN:
            return 'bg-error-600';
        case StatusUnit.UNAVAILABLE:
            return 'bg-error-600';
        case StatusUnit.AVAILABLE:
            return 'bg-success-600';
        case StatusUnit.ON_BREAK:
            return 'bg-warn-600';
        case StatusUnit.BUSY:
            return 'bg-info-500';
        default:
            return 'bg-info-500';
    }
}

export const animateStates = [
    StatusDispatch.NEW.valueOf(),
    StatusDispatch.UNIT_UNASSIGNED.valueOf(),
    StatusDispatch.UNASSIGNED.valueOf(),
    StatusDispatch.NEED_ASSISTANCE.valueOf(),
];

export function dispatchStatusAnimate(status: StatusDispatch | undefined): boolean {
    return animateStates.includes((status ?? StatusDispatch.NEW).valueOf());
}
