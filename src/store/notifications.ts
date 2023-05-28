import { StoreDefinition, defineStore } from 'pinia';
import { v4 as uuidv4 } from 'uuid';
import { Notification } from '~/composables/notification/interfaces/Notification.interface';
import { NotificationConfig } from '~/composables/notification/interfaces/NotificationConfig.interface';

export interface NotificationsState {
    notifications: Notification[];
}

export const useNotificationsStore = defineStore('notifications', {
    state: () =>
        ({
            notifications: [],
        } as NotificationsState),
    persist: false,
    actions: {
        removeNotification(id: string) {
            this.notifications = this.notifications.filter((notification) => notification.id !== id);
        },
        dispatchNotification({
            title,
            content,
            type,
            autoClose = true,
            duration = 6000,
            titleI18n = false,
            contentI18n = false,
        }: NotificationConfig) {
            const id = uuidv4();
            this.notifications.unshift({
                id,
                title,
                titleI18n,
                content,
                contentI18n,
                type,
            });

            if (autoClose) {
                setTimeout(() => {
                    this.removeNotification(id);
                }, duration);
            }
        },
    },
    getters: {
        getNotifications(): Notification[] {
            return [...this.notifications].slice(0, 4);
        },
        getNotificationsCount(): number {
            return this.notifications.length;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useNotificationsStore as unknown as StoreDefinition, import.meta.hot));
}
