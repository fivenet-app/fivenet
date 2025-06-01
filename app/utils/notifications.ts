import type { NotificationAction } from '#ui/types';
import type { TranslateItem } from '~/types/i18n';
import type { Data, NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';

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
