import type { NotificationAction } from '#ui/types';
import type { Data, NotificationCategory } from '~~/gen/ts/resources/notifications/notifications';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

export const NotificationTypes: NotificationType[] = [
    NotificationType.ERROR,
    NotificationType.WARNING,
    NotificationType.INFO,
    NotificationType.SUCCESS,
];

export interface NotificationActionI18n extends Omit<NotificationAction, 'label'> {
    label: TranslateItem;
}

export interface Notification {
    id?: number;
    title: TranslateItem;
    description: TranslateItem;
    timeout?: number;
    type?: NotificationType;
    category?: NotificationCategory;
    data?: Data;
    callback?: () => void;
    actions?: NotificationActionI18n[];
}
