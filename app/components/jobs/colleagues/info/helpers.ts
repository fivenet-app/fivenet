import { ColleagueActivityType } from '~~/gen/ts/resources/jobs/activity';

export function jobsUserActivityTypeBGColor(activityType: ColleagueActivityType | undefined): string {
    switch (activityType) {
        case ColleagueActivityType.ABSENCE_DATE:
            return '!text-yellow-300';
        case ColleagueActivityType.HIRED:
            return '!text-blue-500';
        case ColleagueActivityType.FIRED:
            return '!text-red-500';
        case ColleagueActivityType.PROMOTED:
            return '!text-green-500';
        case ColleagueActivityType.DEMOTED:
            return '!text-orange-500';
        case ColleagueActivityType.NOTE:
            return '!text-teal-500';
        case ColleagueActivityType.NAME:
            return '!text-pink-500';
        case ColleagueActivityType.UNSPECIFIED:
        default:
            return '!text-info-500';
    }
}

export function jobsUserActivityTypeIcon(activityType: ColleagueActivityType | undefined): string {
    switch (activityType) {
        case ColleagueActivityType.ABSENCE_DATE:
            return 'i-mdi-island';
        case ColleagueActivityType.HIRED:
            return 'i-mdi-new-box';
        case ColleagueActivityType.FIRED:
            return 'i-mdi-exit-run';
        case ColleagueActivityType.PROMOTED:
            return 'i-mdi-chevron-double-up';
        case ColleagueActivityType.DEMOTED:
            return 'i-mdi-chevron-double-down';
        case ColleagueActivityType.NOTE:
            return 'i-mdi-note-edit';
        case ColleagueActivityType.LABELS:
            return 'i-mdi-tag-multiple';
        case ColleagueActivityType.NAME:
            return 'i-mdi-rename';
        case ColleagueActivityType.UNSPECIFIED:
        default:
            return 'i-mdi-help';
    }
}
