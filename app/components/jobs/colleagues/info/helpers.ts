import { JobsUserActivityType } from '~~/gen/ts/resources/jobs/colleagues';

export function jobsUserActivityTypeBGColor(activityType: JobsUserActivityType | undefined): string {
    switch (activityType) {
        case JobsUserActivityType.ABSENCE_DATE:
            return '!text-yellow-300';
        case JobsUserActivityType.HIRED:
            return '!text-blue-500';
        case JobsUserActivityType.FIRED:
            return '!text-red-500';
        case JobsUserActivityType.PROMOTED:
            return '!text-green-500';
        case JobsUserActivityType.DEMOTED:
            return '!text-orange-500';
        case JobsUserActivityType.NOTE:
            return '!text-teal-500';
        case JobsUserActivityType.UNSPECIFIED:
        default:
            return '!text-info-500';
    }
}

export function jobsUserActivityTypeIcon(activityType: JobsUserActivityType | undefined): string {
    switch (activityType) {
        case JobsUserActivityType.ABSENCE_DATE:
            return 'i-mdi-island';
        case JobsUserActivityType.HIRED:
            return 'i-mdi-new-box';
        case JobsUserActivityType.FIRED:
            return 'i-mdi-exit-run';
        case JobsUserActivityType.PROMOTED:
            return 'i-mdi-chevron-double-up';
        case JobsUserActivityType.DEMOTED:
            return 'i-mdi-chevron-double-down';
        case JobsUserActivityType.NOTE:
            return 'i-mdi-note-edit';
        case JobsUserActivityType.LABELS:
            return 'i-mdi-tag-mutiple';
        case JobsUserActivityType.UNSPECIFIED:
        default:
            return 'i-mdi-help';
    }
}
