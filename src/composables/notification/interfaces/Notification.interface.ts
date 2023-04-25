export type NotificationType = 'success' | 'info' | 'warning' | 'error';

export interface Notification {
	id: string;
	title?: string;
    titleI18n?: boolean;
	content: string;
    contentI18n?: boolean;
	type?: NotificationType;
}
