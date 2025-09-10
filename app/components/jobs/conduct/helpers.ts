import type { BadgeProps } from '@nuxt/ui';
import { ConductType } from '~~/gen/ts/resources/jobs/conduct';

export function conductTypesToBadgeColor(status: ConductType | undefined): BadgeProps['color'] {
    switch (status) {
        case ConductType.NOTE:
            return 'gray';

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
