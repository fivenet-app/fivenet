import { type TranslateItem } from '~/composables/i18n';
import type { NotificationPosition, NotificationType } from '~/composables/notification/interfaces/Notification.interface';
import { Data, NotificationCategory } from '~~/gen/ts/resources/notifications/notifications';

export interface NotificationConfig {
    title: TranslateItem;
    content: TranslateItem;
    duration?: number;
    autoClose?: boolean;
    type?: NotificationType;
    category?: NotificationCategory;
    data?: Data;
    position?: NotificationPosition;
    callback?: () => Promise<any>;
}
