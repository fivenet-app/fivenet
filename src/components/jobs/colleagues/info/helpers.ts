import { HelpIcon, IslandIcon } from 'mdi-vue3';
import type { DefineComponent } from 'vue';
import { JobsUserActivityType } from '~~/gen/ts/resources/jobs/colleagues';

export function jobsUserActivityTypeBGColor(activityType: JobsUserActivityType | undefined): string {
    switch (activityType) {
        case JobsUserActivityType.ABSENCE_DATE:
            return 'fill-green-600';
        case JobsUserActivityType.UNSPECIFIED:
        default:
            return 'fill-info-500';
    }
}

export function jobsUserActivityTypeIcon(activityType: JobsUserActivityType | undefined): DefineComponent {
    switch (activityType) {
        case JobsUserActivityType.ABSENCE_DATE:
            return markRaw(IslandIcon);
        case JobsUserActivityType.UNSPECIFIED:
        default:
            return markRaw(HelpIcon);
    }
}
