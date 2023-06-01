import { TranslateItem } from '~~/gen/ts/resources/common/i18n';
import { Data, NOTIFICATION_CATEGORY } from '~~/gen/ts/resources/notifications/notifications';

export type NotificationType = 'success' | 'info' | 'warning' | 'error';

export interface Notification {
    id: string;
    title: TranslateItem;
    content: TranslateItem;
    type?: NotificationType;
    category?: NOTIFICATION_CATEGORY;
    data?: Data;
}
