import type { BadgeProps } from '@nuxt/ui';
import { EventAction, EventResult } from '~~/gen/ts/resources/audit/audit';

export function eventActionToBadgeColor(et: EventAction): BadgeProps['color'] {
    switch (et) {
        case EventAction.VIEWED:
            return 'info';
        case EventAction.CREATED:
            return 'success';
        case EventAction.UPDATED:
            return 'warning';
        default:
            return 'error';
    }
}

export function eventResultToBadgeColor(er: EventResult): BadgeProps['color'] {
    switch (er) {
        case EventResult.SUCCEEDED:
            return 'success';
        case EventResult.FAILED:
            return 'warning';
        case EventResult.ERRORED:
            return 'error';
        default:
            return 'info';
    }
}
