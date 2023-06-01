import { TranslateItem } from '~~/gen/ts/resources/common/i18n';
import { Data, NOTIFICATION_CATEGORY } from '~~/gen/ts/resources/notifications/notifications';
import type { NotificationType } from './Notification.interface';

export interface NotificationConfig {
    title: TranslateItem;
    content: TranslateItem;
    duration?: number;
    autoClose?: boolean;
    type?: NotificationType;
    category?: NOTIFICATION_CATEGORY;
    data?: Data;
}
