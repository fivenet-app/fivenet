import type { BadgeProps } from '@nuxt/ui';
import { EventType } from '~~/gen/ts/resources/audit/audit';

export function eventTypeToBadgeColor(et: EventType): BadgeProps['color'] {
    switch (et) {
        case EventType.ERRORED:
            return 'error';
        case EventType.VIEWED:
            return 'info';
        case EventType.CREATED:
            return 'success';
        case EventType.UPDATED:
            return 'warning';
        default:
            return 'error';
    }
}
