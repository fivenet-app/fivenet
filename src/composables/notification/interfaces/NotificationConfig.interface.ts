import type { NotificationType } from './Notification.interface';

export interface NotificationConfig {
    title?: string;
    titleI18n?: boolean;
    content: string;
    contentI18n?: boolean;
    duration?: number;
    autoClose?: boolean;
    type?: NotificationType;
}
