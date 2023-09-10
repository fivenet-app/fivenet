import { CONDUCT_TYPE } from '~~/gen/ts/resources/jobs/conduct';

export function conductTypesToBGColor(status: CONDUCT_TYPE | undefined): string {
    switch (status) {
        case CONDUCT_TYPE.NEUTRAL:
            return 'bg-base-600/10';
        case CONDUCT_TYPE.POSITIVE:
            return 'bg-success-600/10';
        case CONDUCT_TYPE.NEGATIVE:
            return 'bg-error-600/10';
        case CONDUCT_TYPE.WARN:
            return 'bg-warn-600/10';
        case CONDUCT_TYPE.SUSPENSION:
            return 'bg-info-500/10';
        default:
            return 'bg-base-600/10';
    }
}

export function conductTypesToRingColor(status: CONDUCT_TYPE | undefined): string {
    switch (status) {
        case CONDUCT_TYPE.NEUTRAL:
            return 'ring-base-600/20';
        case CONDUCT_TYPE.POSITIVE:
            return 'ring-success-600/20';
        case CONDUCT_TYPE.NEGATIVE:
            return 'ring-error-600/20';
        case CONDUCT_TYPE.WARN:
            return 'ring-warn-600/20';
        case CONDUCT_TYPE.SUSPENSION:
            return 'ring-info-500/20';
        default:
            return 'ring-base-500/20';
    }
}

export function conductTypesToTextColor(status: CONDUCT_TYPE | undefined): string {
    switch (status) {
        case CONDUCT_TYPE.NEUTRAL:
            return 'text-base-600';
        case CONDUCT_TYPE.POSITIVE:
            return 'text-success-600';
        case CONDUCT_TYPE.NEGATIVE:
            return 'text-error-600';
        case CONDUCT_TYPE.WARN:
            return 'text-warn-600';
        case CONDUCT_TYPE.SUSPENSION:
            return 'text-info-500';
        default:
            return 'text-base-500';
    }
}
