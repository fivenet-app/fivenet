import type { NotificationType } from './Notification.interface';

export interface NotificationConfig {
	title?: string;
	content: string;
	duration?: number;
	autoClose?: boolean;
	type?: NotificationType;
}
