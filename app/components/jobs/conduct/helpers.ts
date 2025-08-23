import type { BadgeProps } from '@nuxt/ui';
import { ConductType } from '~~/gen/ts/resources/jobs/conduct';

export function conductTypesToBadgeColor(status: ConductType | undefined): BadgeProps['color'] {
    switch (status) {
        case ConductType.NEUTRAL:
            return 'neutral';
        case ConductType.POSITIVE:
            return 'success';
        case ConductType.NEGATIVE:
            return 'error';
        case ConductType.WARNING:
            return 'amber';
        case ConductType.SUSPENSION:
            return 'sky';
        default:
            return 'white';
    }
}

export function conductTypesToBGColor(status: ConductType | undefined): string {
    switch (status) {
        case ConductType.NEUTRAL:
            return 'bg-background/10';
        case ConductType.POSITIVE:
            return 'bg-success-600/10';
        case ConductType.NEGATIVE:
            return 'bg-error-600/10';
        case ConductType.WARNING:
            return 'bg-warn-600/10';
        case ConductType.SUSPENSION:
            return 'bg-info-600/10';
        default:
            return 'bg-background/10';
    }
}
