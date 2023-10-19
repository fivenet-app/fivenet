import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { StoreDefinition, defineStore } from 'pinia';
import { v4 as uuidv4 } from 'uuid';
import { Notification, NotificationType } from '~/composables/notification/interfaces/Notification.interface';
import { NotificationConfig } from '~/composables/notification/interfaces/NotificationConfig.interface';
import { NotificationCategory } from '~~/gen/ts/resources/notifications/notifications';
import { useAuthStore } from './auth';

// In seconds
const initialBackoffTime = 2;

export interface NotificationsState {
    lastId: string;
    notifications: Notification[];
    abort: AbortController | undefined;
    restarting: boolean;
    restartBackoffTime: number;
}

export const useNotificatorStore = defineStore('notifications', {
    state: () =>
        ({
            lastId: '0',
            notifications: [],
            abort: undefined,
            restarting: false,
            restartBackoffTime: 0,
        }) as NotificationsState,
    persist: {
        paths: ['lastId'],
    },
    actions: {
        setLastId(lastId: bigint): void {
            this.lastId = lastId.toString();
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
            category = NotificationCategory.GENERAL,
            data = undefined,
            position = 'top-right',
        }: NotificationConfig) {
            const id = uuidv4();
            this.notifications.unshift({
                id,
                title,
                content,
                type,
                category,
                data,
                position,
            });

            if (autoClose) {
                setTimeout(() => {
                    this.removeNotification(id);
                }, duration);
            }
        },

        async startStream(): Promise<void> {
            if (this.abort !== undefined) return;

            console.debug('Notificator: Stream starting, starting at ID:', this.getLastId);

            this.abort = new AbortController();
            this.restarting = false;
            const { $grpc } = useNuxtApp();

            try {
                this.abort = new AbortController();

                const call = $grpc.getNotificatorClient().stream(
                    {
                        lastId: this.getLastId,
                    },
                    {
                        abort: this.abort.signal,
                    },
                );

                for await (let resp of call.responses) {
                    if (resp.lastId > this.getLastId) this.setLastId(resp.lastId);

                    if (resp.data.oneofKind !== undefined) {
                        if (resp.data.oneofKind === 'ping') {
                        } else if (resp.data.oneofKind === 'token') {
                            const tokenUpdate = resp.data.token;

                            const authStore = useAuthStore();

                            // Update active char when updated user info is received
                            if (tokenUpdate.userInfo) {
                                console.debug('Notificator: Updated UserInfo received');

                                authStore.setActiveChar(tokenUpdate.userInfo);
                                authStore.setPermissions(tokenUpdate.permissions);
                                if (tokenUpdate.jobProps) {
                                    authStore.setJobProps(tokenUpdate.jobProps!);
                                } else {
                                    authStore.setJobProps(null);
                                }
                            }

                            if (tokenUpdate.newToken && tokenUpdate.expires) {
                                console.debug('Notificator: New Token received');

                                authStore.setAccessToken(tokenUpdate.newToken, toDate(tokenUpdate.expires) as null | Date);

                                this.dispatchNotification({
                                    title: { key: 'notifications.renewed_token.title', parameters: {} },
                                    content: { key: 'notifications.renewed_token.content', parameters: {} },
                                    type: 'info',
                                });
                            }
                        } else if (resp.data.oneofKind === 'notifications') {
                            resp.data.notifications.notifications.forEach((n) => {
                                let nType: NotificationType = (n.type as NotificationType) ?? 'info';

                                switch (n.category) {
                                    case NotificationCategory.GENERAL:
                                        this.dispatchNotification({
                                            title: { key: n.title?.key!, parameters: {} },
                                            content: { key: n.content?.key!, parameters: {} },
                                            type: nType,
                                            category: n.category,
                                            data: n.data,
                                        });
                                        break;

                                    default:
                                        this.dispatchNotification({
                                            title: { key: n.title?.key!, parameters: {} },
                                            content: { key: n.content?.key!, parameters: {} },
                                            type: nType,
                                            category: n.category,
                                            data: n.data,
                                        });
                                        break;
                                }
                            });
                        } else {
                            // @ts-ignore this is a catch all "unknown", so okay if it is technically "never" reached till it is
                            console.warn('Notificator: Unknown data received - Kind: ', resp.data.oneofKind, resp.data);
                        }
                    }

                    if (resp.restart) {
                        console.debug('Notificator: Server requested stream to be restarted');
                        this.restartBackoffTime = 0;
                        this.restartStream();
                        break;
                    }
                }
            } catch (e) {
                const error = e as RpcError;
                if (error.code !== 'CANCELLED' && error.code !== 'ABORTED') {
                    console.debug('Notificator: Stream failed', error.code, error.message, error.cause);
                    this.restartStream();
                }
            }

            console.debug('Notificator: Stream ended');
        },

        async stopStream(): Promise<void> {
            if (this.abort !== undefined) this.abort?.abort();
            this.abort = undefined;

            console.debug('Notificator: Stopping Data Stream');
        },

        async restartStream(): Promise<void> {
            // Reset back off time when over 3 minutes
            if (this.restartBackoffTime > 180) {
                this.restartBackoffTime = initialBackoffTime;
            } else {
                this.restartBackoffTime += initialBackoffTime;
            }

            console.debug('Notificator: Restart back off time in', this.restartBackoffTime, 'seconds');
            await this.stopStream();

            setTimeout(async () => {
                if (this.restarting) this.startStream();
            }, this.restartBackoffTime * 1000);
        },
    },
    getters: {
        getLastId(): bigint {
            return BigInt(this.lastId);
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
    import.meta.hot.accept(acceptHMRUpdate(useNotificatorStore as unknown as StoreDefinition, import.meta.hot));
}
