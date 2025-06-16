import type { NotificationAction } from '#ui/types';
import type { I18NItem } from '~/types/i18n';
import type { Data, NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';

export interface NotificationActionI18n extends Omit<NotificationAction, 'label'> {
    label: I18NItem;
}

export interface Notification {
    id?: number;
    title: I18NItem;
    description: I18NItem;
    timeout?: number;
    type?: NotificationType;
    category?: NotificationCategory;
    data?: Data;
    callback?: () => void;
    actions?: NotificationActionI18n[];
}
