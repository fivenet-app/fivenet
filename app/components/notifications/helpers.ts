import type { ToastProps } from '@nuxt/ui';
import {
    NotificationCategory,
    NotificationType,
    type Notification,
} from '~~/gen/ts/resources/notifications/notifications';
import { notificationCategoryDefinition, notificationKindDefinition } from './definitions';

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
    return notificationCategoryDefinition(category)?.icon ?? 'i-mdi-information-outline';
}

/**
 * Kinds are more specific than categories. Unknown future kinds deliberately
 * fall back to their category so old clients still render a useful cue.
 */
export function notificationToIcon(notification: Pick<Notification, 'kind' | 'category'>): string {
    return notificationKindDefinition(notification.kind)?.icon ?? notificationCategoryToIcon(notification.category);
}
