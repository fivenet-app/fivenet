import { ChevronDoubleDownIcon, ChevronDoubleUpIcon, ExitRunIcon, HelpIcon, IslandIcon, NewBoxIcon } from 'mdi-vue3';
import type { DefineComponent } from 'vue';
import { JobsUserActivityType } from '~~/gen/ts/resources/jobs/colleagues';

export function jobsUserActivityTypeBGColor(activityType: JobsUserActivityType | undefined): string {
    switch (activityType) {
        case JobsUserActivityType.ABSENCE_DATE:
            return 'fill-yellow-300';
        case JobsUserActivityType.HIRED:
            return 'fill-blue-500';
        case JobsUserActivityType.FIRED:
            return 'fill-red-500';
        case JobsUserActivityType.PROMOTED:
            return 'fill-green-500';
        case JobsUserActivityType.DEMOTED:
            return 'fill-orange-500';
        case JobsUserActivityType.UNSPECIFIED:
        default:
            return 'fill-info-500';
    }
}

export function jobsUserActivityTypeIcon(activityType: JobsUserActivityType | undefined): DefineComponent {
    switch (activityType) {
        case JobsUserActivityType.ABSENCE_DATE:
            return markRaw(IslandIcon);
        case JobsUserActivityType.HIRED:
            return markRaw(NewBoxIcon);
        case JobsUserActivityType.FIRED:
            return markRaw(ExitRunIcon);
        case JobsUserActivityType.PROMOTED:
            return markRaw(ChevronDoubleUpIcon);
        case JobsUserActivityType.DEMOTED:
            return markRaw(ChevronDoubleDownIcon);
        case JobsUserActivityType.UNSPECIFIED:
        default:
            return markRaw(HelpIcon);
    }
}
