import type { ToastProps } from '@nuxt/ui';
import {
    NotificationCategory,
    NotificationKind,
    NotificationType,
    type Notification,
} from '~~/gen/ts/resources/notifications/notifications';

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

export function notificationTypeToColor(t?: NotificationType): ToastProps['color'] {
    switch (t) {
        case NotificationType.SUCCESS:
            return 'success';
        case NotificationType.WARNING:
            return 'amber';
        case NotificationType.ERROR:
            return 'error';
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
        case NotificationCategory.JOBS:
            return 'i-mdi-briefcase-outline';
        case NotificationCategory.QUALIFICATIONS:
            return 'i-mdi-certificate-outline';
        case NotificationCategory.MAILER:
            return 'i-mdi-email-outline';
        case NotificationCategory.SYSTEM:
            return 'i-mdi-cog-outline';
        default:
            return 'i-mdi-information-outline';
    }
}

/**
 * Kinds are more specific than categories. Unknown future kinds deliberately
 * fall back to their category so old clients still render a useful cue.
 */
export function notificationToIcon(notification: Pick<Notification, 'kind' | 'category'>): string {
    switch (notification.kind) {
        case NotificationKind.JOBS_GROUP_LEADERSHIP_ADDED:
            return 'i-mdi-account-star-outline';
        case NotificationKind.JOBS_GROUP_LEADERSHIP_REMOVED:
            return 'i-mdi-account-star-off-outline';
        default:
            return notificationCategoryToIcon(notification.category);
    }
}
