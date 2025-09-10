import type { ButtonProps } from '@nuxt/ui';
import type { Data, NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';

export interface NotificationActionI18n extends Omit<ButtonProps, 'label'> {
    label: I18NItem;
}

export interface Notification {
    id?: number;
    title: I18NItem;
    description: I18NItem;
    duration?: number;
    type?: NotificationType;
    category?: NotificationCategory;
    data?: Data;
    callback?: () => void;
    actions?: NotificationActionI18n[];
}
