import type { BadgeColor } from '#ui/types';
import { EventType } from '~~/gen/ts/resources/audit/audit';

export function eventTypeToBadgeColor(et: EventType): BadgeColor {
    switch (et) {
        case EventType.ERRORED:
            return 'orange';
        case EventType.VIEWED:
            return 'blue';
        case EventType.CREATED:
            return 'success';
        case EventType.UPDATED:
            return 'amber';
        default:
            return 'error';
    }
}
