import { StateTree, StoreDefinition, defineStore } from 'pinia';
import { v4 as uuidv4 } from 'uuid';
import { Notification } from '~/composables/notification/interfaces/Notification.interface';
import { NotificationConfig } from '~/composables/notification/interfaces/NotificationConfig.interface';
import { NOTIFICATION_CATEGORY } from '~~/gen/ts/resources/notifications/notifications';

class Serializer {
    serialize(value: StateTree): string {
        return JSON.stringify(value, (_, v) => (typeof v === 'bigint' ? v.toString() + 'n' : v));
    }

    deserialize(value: string): StateTree {
        return JSON.parse(value, (_, v) => {
            if (typeof v === 'string' && /^\d+n$/.test(v)) {
                return BigInt(v.substring(0, v.length - 1));
            }
            return value;
        }) as StateTree;
    }
}

export interface NotificationsState {
    lastId: bigint;
    notifications: Notification[];
}

export const useNotificationsStore = defineStore('notifications', {
    state: () =>
        ({
            lastId: 0n,
            notifications: [],
        } as NotificationsState),
    persist: {
        paths: ['lastId'],
        serializer: new Serializer(),
    },
    actions: {
        setLastId(lastId: bigint): void {
            this.lastId = lastId;
        },
        removeNotification(id: string): void {
            this.notifications = this.notifications.filter((notification) => notification.id !== id);
        },
        dispatchNotification({
            title,
            content,
            type,
            autoClose = true,
            duration = 6000,
            category = NOTIFICATION_CATEGORY.GENERAL,
            data = undefined,
        }: NotificationConfig) {
            const id = uuidv4();
            this.notifications.unshift({
                id,
                title,
                content,
                type,
                category,
                data,
            });

            if (autoClose) {
                setTimeout(() => {
                    this.removeNotification(id);
                }, duration);
            }
        },
    },
    getters: {
        getLastId(): bigint {
            return this.lastId;
        },
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
