import { TranslateItem } from '~/composables/i18n';
import { Data, NotificationCategory } from '~~/gen/ts/resources/notifications/notifications';

export type NotificationType = 'success' | 'info' | 'warning' | 'error';

export interface Notification {
    id: string;
    title: TranslateItem;
    content: TranslateItem;
    type?: NotificationType;
    category?: NotificationCategory;
    data?: Data;
}
