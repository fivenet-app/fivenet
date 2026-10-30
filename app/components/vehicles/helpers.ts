import { VehicleActivityType, type VehicleActivity } from '~~/gen/ts/resources/vehicles/activity/activity';

export function vehicleActivityTypeColor(activityType: VehicleActivityType | undefined): string {
    switch (activityType) {
        case VehicleActivityType.WANTED:
            return 'text-error-400!';
        case VehicleActivityType.FILE:
            return 'text-primary-500!';
        case VehicleActivityType.PROPERTY:
            return 'text-secondary-500!';
        case VehicleActivityType.UNSPECIFIED:
        default:
            return '!text-info-500';
    }
}

export function vehicleActivityTypeIcon(activityType: VehicleActivityType | undefined): string {
    switch (activityType) {
        case VehicleActivityType.WANTED:
            return 'i-mdi-bell-alert';
        case VehicleActivityType.FILE:
            return 'i-mdi-file-image';
        case VehicleActivityType.PROPERTY:
            return 'i-mdi-form-textbox';
        case VehicleActivityType.UNSPECIFIED:
        default:
            return 'i-mdi-help';
    }
}

export function vehicleActivityIndicatorColor(activity: VehicleActivity): string {
    switch (activity.activityType) {
        case VehicleActivityType.WANTED:
            if (activity.data?.data.oneofKind === 'wantedChange') {
                return activity.data.data.wantedChange.wanted ? '!bg-error-400 text-white!' : '!bg-success-400 text-white!';
            }
            break;
        case VehicleActivityType.FILE:
            if (activity.data?.data.oneofKind === 'fileChange' && !activity.data.data.fileChange.added) {
                return '!bg-error-400 text-white!';
            }

            return '!bg-primary-500 text-white!';
        case VehicleActivityType.PROPERTY:
            return '!bg-secondary-500 text-white!';
        case VehicleActivityType.UNSPECIFIED:
        default:
            return '!bg-info-500 text-white!';
    }

    return '!bg-info-500 text-white!';
}
