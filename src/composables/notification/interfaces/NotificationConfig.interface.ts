import { TranslateItem } from '~/composables/i18n';
import { Data, NotificationCategory } from '~~/gen/ts/resources/notifications/notifications';
import type { NotificationPosition, NotificationType } from './Notification.interface';

export interface NotificationConfig {
    title: TranslateItem;
    content: TranslateItem;
    duration?: number;
    autoClose?: boolean;
    type?: NotificationType;
    category?: NotificationCategory;
    data?: Data;
    position?: NotificationPosition;
}
