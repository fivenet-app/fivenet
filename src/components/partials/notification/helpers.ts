import type { NotificationType } from '~/composables/notifications/notifications';

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
