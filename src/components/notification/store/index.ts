import type { NotificationsState, NotificationConfig } from '../interfaces';
import { reactive, readonly } from 'vue';
import { v4 as uuidv4 } from 'uuid';

const state = reactive<NotificationsState>({
	notifications: [],
});

const actions = {
	removeNotification(id: string) {
		state.notifications = state.notifications.filter((notification) => notification.id !== id);
	},
	dispatchNotification({ title, content, type, autoClose = true, duration = 6000 }: NotificationConfig) {
		const id = uuidv4();
		const notifications = [
			{
				id,
				title,
				content,
				type,
			},
			...state.notifications,
		];
		state.notifications = notifications;

		if (autoClose) {
			setTimeout(() => {
				actions.removeNotification(id);
			}, duration);
		}
	},
};

const getters = {
	getNotifications() {
		return [...state.notifications].slice(0, 4);
	},
	getNotificationsCount() {
		return state.notifications.length;
	},
};

export default {
	state: readonly(state),
	actions,
	getters,
};
