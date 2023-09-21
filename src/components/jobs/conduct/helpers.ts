import { ConductType } from '~~/gen/ts/resources/jobs/conduct';

export function conductTypesToBGColor(status: ConductType | undefined): string {
    switch (status) {
        case ConductType.NEUTRAL:
            return 'bg-base-600/10';
        case ConductType.POSITIVE:
            return 'bg-success-600/10';
        case ConductType.NEGATIVE:
            return 'bg-error-600/10';
        case ConductType.WARNING:
            return 'bg-warn-600/10';
        case ConductType.SUSPENSION:
            return 'bg-info-500/10';
        default:
            return 'bg-base-600/10';
    }
}

export function conductTypesToRingColor(status: ConductType | undefined): string {
    switch (status) {
        case ConductType.NEUTRAL:
            return 'ring-base-600/20';
        case ConductType.POSITIVE:
            return 'ring-success-600/20';
        case ConductType.NEGATIVE:
            return 'ring-error-600/20';
        case ConductType.WARNING:
            return 'ring-warn-600/20';
        case ConductType.SUSPENSION:
            return 'ring-info-500/20';
        default:
            return 'ring-base-500/20';
    }
}

export function conductTypesToTextColor(status: ConductType | undefined): string {
    switch (status) {
        case ConductType.NEUTRAL:
            return 'text-base-400';
        case ConductType.POSITIVE:
            return 'text-success-400';
        case ConductType.NEGATIVE:
            return 'text-error-400';
        case ConductType.WARNING:
            return 'text-warn-400';
        case ConductType.SUSPENSION:
            return 'text-info-400';
        default:
            return 'text-base-400';
    }
}
