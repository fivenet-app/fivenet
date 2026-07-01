import type { ButtonProps } from '@nuxt/ui';
import type { Data, NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';

export interface NotificationActionI18n {
    label: I18NItem;
    to?: string;
    external?: boolean;
    icon?: string;
    color?: ButtonProps['color'];
    onClick?: () => void | Promise<void>;
}

export interface Notification {
    id?: number;
    title: I18NItem | string;
    description: I18NItem | string;
    duration?: number;
    type?: NotificationType;
    category?: NotificationCategory;
    data?: Data;
    callback?: () => void;
    actions?: NotificationActionI18n[];
}
