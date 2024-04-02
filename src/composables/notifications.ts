import { type TranslateItem } from '~/composables/i18n';
import { Data, NotificationCategory } from '~~/gen/ts/resources/notifications/notifications';

export type NotificationType = 'success' | 'info' | 'warning' | 'error';

export interface Notification {
    id?: string;
    title: TranslateItem;
    description: TranslateItem;
    timeout?: number;
    type?: NotificationType;
    category?: NotificationCategory;
    data?: Data;
    onClick?: (data?: Data) => Promise<any>;
    onClickText?: TranslateItem;
    callback?: Function;
}
