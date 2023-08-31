import { DISPATCH_STATUS } from '~~/gen/ts/resources/dispatch/dispatches';
import { UNIT_STATUS } from '~~/gen/ts/resources/dispatch/units';

export function dispatchStatusToFillColor(status: DISPATCH_STATUS | undefined): string {
    switch (status) {
        case DISPATCH_STATUS.NEW:
            return 'fill-error-600';
        case DISPATCH_STATUS.UNASSIGNED:
            return 'fill-error-600';
        case DISPATCH_STATUS.EN_ROUTE:
            return 'fill-info-400';
        case DISPATCH_STATUS.ON_SCENE:
            return 'fill-primary-600';
        case DISPATCH_STATUS.NEED_ASSISTANCE:
            return 'fill-warn-600';
        case DISPATCH_STATUS.COMPLETED:
            return 'fill-success-600';
        case DISPATCH_STATUS.CANCELLED:
            return 'fill-success-600';
        case DISPATCH_STATUS.ARCHIVED:
            return 'fill-base-600';
        default:
            return 'fill-info-400';
    }
}

export function dispatchStatusToBGColor(status: DISPATCH_STATUS | undefined): string {
    switch (status) {
        case DISPATCH_STATUS.NEW:
            return 'bg-error-600';
        case DISPATCH_STATUS.UNASSIGNED:
            return 'bg-error-600';
        case DISPATCH_STATUS.EN_ROUTE:
            return 'bg-info-400';
        case DISPATCH_STATUS.ON_SCENE:
            return 'bg-primary-600';
        case DISPATCH_STATUS.NEED_ASSISTANCE:
            return 'bg-warn-600';
        case DISPATCH_STATUS.COMPLETED:
            return 'bg-success-600';
        case DISPATCH_STATUS.CANCELLED:
            return 'bg-success-600';
        case DISPATCH_STATUS.ARCHIVED:
            return 'bg-base-600';
        default:
            return 'bg-info-400';
    }
}

export function unitStatusToBGColor(status: UNIT_STATUS | undefined): string {
    switch (status) {
        case UNIT_STATUS.UNKNOWN:
            return 'bg-error-600';
        case UNIT_STATUS.UNAVAILABLE:
            return 'bg-error-600';
        case UNIT_STATUS.AVAILABLE:
            return 'bg-success-600';
        case UNIT_STATUS.ON_BREAK:
            return 'bg-warn-600';
        case UNIT_STATUS.BUSY:
            return 'bg-info-400';
        default:
            return 'bg-info-400';
    }
}

export const animateStates = [
    DISPATCH_STATUS.NEW.valueOf(),
    DISPATCH_STATUS.UNIT_UNASSIGNED.valueOf(),
    DISPATCH_STATUS.UNASSIGNED.valueOf(),
    DISPATCH_STATUS.NEED_ASSISTANCE.valueOf(),
];

export function dispatchStatusAnimate(status: DISPATCH_STATUS | undefined): boolean {
    return animateStates.includes((status ?? DISPATCH_STATUS.NEW).valueOf());
}
