import { Data, NotificationCategory } from '~~/gen/ts/resources/notifications/notifications';
import type { NotificationType } from './Notification.interface';

export interface NotificationConfig {
    title: TranslateItem;
    content: TranslateItem;
    duration?: number;
    autoClose?: boolean;
    type?: NotificationType;
    category?: NotificationCategory;
    data?: Data;
}
