import type { NotificationColor } from '#ui/types';
import { NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';

export function notificationTypeToIcon(t?: NotificationType): string {
    switch (t) {
        case NotificationType.SUCCESS:
            return 'i-mdi-check-circle';
        case NotificationType.WARNING:
            return 'i-mdi-alert-circle';
        case NotificationType.ERROR:
            return 'i-mdi-close-circle';
        case NotificationType.INFO:
        default:
            return 'i-mdi-information-outline';
    }
}

export function notificationTypeToColor(t?: NotificationType): NotificationColor {
    switch (t) {
        case NotificationType.SUCCESS:
            return 'green';
        case NotificationType.WARNING:
            return 'amber';
        case NotificationType.ERROR:
            return 'red';
        case NotificationType.INFO:
        default:
            return 'blue';
    }
}

export function notificationCategoryToIcon(category: NotificationCategory): string {
    switch (category) {
        case NotificationCategory.DOCUMENT:
            return 'i-mdi-file-document-box-multiple-outline';
        case NotificationCategory.CALENDAR:
            return 'i-mdi-calendar-outline';
        default:
            return 'i-mdi-information-outline';
    }
}
