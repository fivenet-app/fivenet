import { type TranslateItem } from '~/composables/i18n';
import { Data, NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';

export const NotificationTypes: NotificationType[] = [
    NotificationType.ERROR,
    NotificationType.WARNING,
    NotificationType.INFO,
    NotificationType.SUCCESS,
];

export interface Notification {
    id?: string;
    title: TranslateItem;
    description: TranslateItem;
    timeout?: number;
    type?: NotificationType;
    category?: NotificationCategory;
    data?: Data;
    onClick?: ((data?: Data) => any) | ((data?: Data) => Promise<any>);
    onClickText?: TranslateItem;
    callback?: Function;
}
