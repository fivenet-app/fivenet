import { UserActivityType, type UserActivity } from '~~/gen/ts/resources/users/activity/activity';

export function citizenUserActivityTypeBGColor(activityType: UserActivityType | undefined): string {
    switch (activityType) {
        case UserActivityType.NAME:
            return 'text-info-600!';
        case UserActivityType.DOCUMENT:
            return 'text-info-600!';
        case UserActivityType.WANTED:
            return 'text-error-400!';
        case UserActivityType.JOB:
            return 'text-gray-400!';
        case UserActivityType.TRAFFIC_INFRACTION_POINTS:
            return 'text-orange-400!';
        case UserActivityType.MUGSHOT:
            return 'text-amber-400!';
        case UserActivityType.LABELS:
            return 'text-amber-200!';
        case UserActivityType.LICENSES:
            return 'text-info-600!';
        case UserActivityType.JAIL:
            return 'text-base-600!';
        case UserActivityType.FINE:
            return 'text-red-400!';
        case UserActivityType.UNSPECIFIED:
        default:
            return '!text-info-500';
    }
}

export function citizenUserActivityTypeIcon(activityType: UserActivityType | undefined): string {
    switch (activityType) {
        case UserActivityType.NAME:
            return 'i-mdi-identification-card';
        case UserActivityType.DOCUMENT:
            return 'i-mdi-file-account';
        case UserActivityType.WANTED:
            return 'i-mdi-bell-alert';
        case UserActivityType.JOB:
            return 'i-mdi-briefcase';
        case UserActivityType.TRAFFIC_INFRACTION_POINTS:
            return 'i-mdi-traffic-cone';
        case UserActivityType.MUGSHOT:
            return 'i-mdi-camera-account';
        case UserActivityType.LABELS:
            return 'i-mdi-tag';
        case UserActivityType.LICENSES:
            return 'i-mdi-license';
        case UserActivityType.JAIL:
            return 'i-mdi-handcuffs';
        case UserActivityType.FINE:
            return 'i-mdi-receipt-text-remove';
        case UserActivityType.UNSPECIFIED:
        default:
            return 'i-mdi-help';
    }
}

export function citizenUserActivityIconColor(activity: UserActivity): string {
    switch (activity.type) {
        case UserActivityType.NAME:
            return 'text-info-600';
        case UserActivityType.DOCUMENT:
            if (activity.data?.data.oneofKind === 'documentRelation') {
                return activity.data.data.documentRelation.added ? 'text-info-600' : 'text-base-600';
            }
            break;
        case UserActivityType.WANTED:
            if (activity.data?.data.oneofKind === 'wantedChange') {
                return activity.data.data.wantedChange.wanted ? 'text-error-400' : 'text-success-400';
            }
            break;
        case UserActivityType.JOB:
            return 'text-gray-400';
        case UserActivityType.TRAFFIC_INFRACTION_POINTS:
            if (activity.data?.data.oneofKind === 'trafficInfractionPointsChange') {
                return activity.data.data.trafficInfractionPointsChange.old >
                    activity.data.data.trafficInfractionPointsChange.new
                    ? 'text-gray-400'
                    : 'text-orange-400';
            }
            break;
        case UserActivityType.MUGSHOT:
            if (activity.data?.data.oneofKind === 'mugshotChange') {
                return activity.data.data.mugshotChange.new ? 'text-gray-400' : 'text-amber-400';
            }
            break;
        case UserActivityType.LABELS:
            return 'text-amber-200';
        case UserActivityType.LICENSES:
            if (activity.data?.data.oneofKind === 'licensesChange') {
                return activity.data.data.licensesChange.added ? 'text-info-600' : 'text-amber-600';
            }
            break;
        case UserActivityType.FINE:
            if (activity.data?.data.oneofKind === 'fineChange') {
                if (activity.data.data.fineChange.removed) {
                    return 'text-red-400';
                }

                return activity.data.data.fineChange.amount < 0 ? 'text-success-400' : 'text-info-400';
            }
            break;
        case UserActivityType.JAIL:
        case UserActivityType.UNSPECIFIED:
        default:
            return '';
    }

    return '';
}
