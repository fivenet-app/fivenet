import { VehicleActivityType, type VehicleActivity } from '~~/gen/ts/resources/vehicles/activity/activity';

export function vehicleActivityTypeBGColor(activityType: VehicleActivityType | undefined): string {
    switch (activityType) {
        case VehicleActivityType.WANTED:
            return 'text-error-400!';
        case VehicleActivityType.UNSPECIFIED:
        default:
            return '!text-info-500';
    }
}

export function vehicleActivityTypeIcon(activityType: VehicleActivityType | undefined): string {
    switch (activityType) {
        case VehicleActivityType.WANTED:
            return 'i-mdi-bell-alert';
        case VehicleActivityType.UNSPECIFIED:
        default:
            return 'i-mdi-help';
    }
}

export function vehicleActivityIconColor(activity: VehicleActivity): string {
    switch (activity.activityType) {
        case VehicleActivityType.WANTED:
            if (activity.data?.data.oneofKind === 'wantedChange') {
                return activity.data.data.wantedChange.wanted ? 'text-error-400' : 'text-success-400';
            }
            break;
        case VehicleActivityType.UNSPECIFIED:
        default:
            return '';
    }

    return '';
}
