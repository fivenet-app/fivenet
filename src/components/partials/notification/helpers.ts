import type { NotificationColor } from '#ui/types';
import type { NotificationType } from '~/composables/notifications';
import { NotificationCategory } from '~~/gen/ts/resources/notifications/notifications';

export function notificationTypeToIcon(t?: NotificationType): string {
    switch (t) {
        case 'success':
            return 'i-mdi-check-circle';
        case 'warning':
            return 'i-mdi-alert-circle';
        case 'error':
            return 'i-mdi-close-circle';
        case 'info':
        default:
            return 'i-mdi-information';
    }
}

export function notificationTypeToColor(t?: NotificationType): NotificationColor {
    switch (t) {
        case 'success':
            return 'green';
        case 'warning':
            return 'amber';
        case 'error':
            return 'red';
        case 'info':
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
            return '';
    }
}
