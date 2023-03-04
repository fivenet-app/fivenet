export type NotificationType = 'success' | 'info' | 'warning' | 'error';

export interface Notification {
	id: string;
	title?: string;
	content: string;
	type?: NotificationType;
}
