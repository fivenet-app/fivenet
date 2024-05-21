import { defineStore, type StoreDefinition } from 'pinia';
import { v4 as uuidv4 } from 'uuid';
import { type Notification } from '~/composables/notifications';
import { useAuthStore } from '~/store/auth';
import { NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { MarkNotificationsRequest } from '~~/gen/ts/services/notificator/notificator';
import { useCalendarStore } from './calendar';
import { messengerStore } from './messenger';

// In seconds
const initialReconnectBackoffTime = 2;
const maxBackoffTime = 40;

export interface NotificationsState {
    notifications: Notification[];
    notificationsCount: number;

    abort: AbortController | undefined;
    reconnecting: boolean;
    reconnectBackoffTime: number;
}

export const useNotificatorStore = defineStore('notifications', {
    state: () =>
        ({
            notifications: [],
            notificationsCount: 0,

            abort: undefined,
            reconnecting: false,
            reconnectBackoffTime: 0,
        }) as NotificationsState,
    persist: false,
    actions: {
        remove(id: string): void {
            this.notifications = this.notifications.filter((notification) => notification.id !== id);
        },
        add(notification: Notification): void {
            notification.id = uuidv4();

            if (notification.timeout === undefined) {
                notification.timeout = 3500;
            }

            this.notifications.push(notification);
        },

        async startStream(): Promise<void> {
            if (this.abort !== undefined) {
                return;
            }

            console.debug('Notificator: Starting Data Stream');

            this.abort = new AbortController();
            this.reconnecting = false;
            const authStore = useAuthStore();

            try {
                this.abort = new AbortController();

                const call = getGRPCNotificatorClient().stream(
                    {},
                    {
                        abort: this.abort.signal,
                    },
                );

                for await (const resp of call.responses) {
                    this.notificationsCount = resp.notificationCount;

                    if (resp.data.oneofKind !== undefined) {
                        if (resp.data.oneofKind === 'userEvent') {
                            if (resp.data.userEvent.data.oneofKind === 'refreshToken') {
                                console.info('Notificator: Refreshing token...');
                                await authStore.chooseCharacter(undefined);
                                continue;
                            } else if (resp.data.userEvent.data.oneofKind === 'notification') {
                                const n = resp.data.userEvent.data.notification;
                                const nType: NotificationType =
                                    n.type !== NotificationType.UNSPECIFIED ? n.type : NotificationType.INFO;

                                if (n.title === undefined || n.content === undefined) {
                                    continue;
                                }

                                const not: Notification = {
                                    title: { key: n.title.key, parameters: n.title.parameters },
                                    description: {
                                        key: n.content.key,
                                        parameters: n.content.parameters,
                                    },
                                    type: nType,
                                    category: n.category,
                                    data: n.data,
                                };

                                if (n.data?.link !== undefined) {
                                    if (n.data.link.external === true) {
                                        not.onClick = async () => navigateTo(n.data!.link!.to, { external: true });
                                    } else {
                                        // @ts-ignore route from a notification is a string
                                        const route = useRouter().resolve(n.data!.link!.to);

                                        not.onClick = async () => {
                                            navigateTo(route);
                                        };
                                    }
                                }

                                if (n.category === NotificationCategory.CALENDAR) {
                                    useSound().play({ name: 'notification' });

                                    if (n.data?.calendar !== undefined) {
                                        const calendarStore = useCalendarStore();

                                        try {
                                            if (n.data?.calendar.calendarId !== undefined) {
                                                calendarStore.getCalendar({
                                                    calendarId: n.data?.calendar.calendarId,
                                                });
                                            } else if (n.data?.calendar.calendarEntryId !== undefined) {
                                                calendarStore.getCalendarEntry({
                                                    entryId: n.data?.calendar.calendarEntryId,
                                                });
                                            }
                                        } catch (e) {}
                                    }
                                }

                                this.add(not);
                                continue;
                            } else if (resp.data.userEvent.data.oneofKind === 'messenger') {
                                if (can('MessengerService.ListThreads')) {
                                    messengerStore.handleEvent(resp.data.userEvent.data.messenger);
                                }
                            }
                        } else if (resp.data.oneofKind === 'jobEvent') {
                            if (resp.data.jobEvent.data.oneofKind === 'jobProps') {
                                authStore.setJobProps(resp.data.jobEvent.data.jobProps);
                            } else {
                                console.warn(
                                    'Notificator: Unknown job event data received - Kind: ',
                                    resp.data.oneofKind,
                                    resp.data,
                                );
                            }
                            continue;
                        } else if (resp.data.oneofKind === 'systemEvent') {
                        } else {
                            // @ts-ignore this is a catch all "unknown", so okay if it is technically "never" reached till it is..
                            console.warn('Notificator: Unknown data received - Kind: ', resp.data.oneofKind, resp.data);
                        }
                    }

                    if (resp.restart) {
                        console.debug('Notificator: Server requested stream to be restarted');
                        this.reconnectBackoffTime = 0;
                        this.restartStream();
                        break;
                    }
                }
            } catch (e) {
                const error = e as RpcError;
                if (error.code !== 'CANCELLED' && error.code !== 'ABORTED') {
                    console.debug('Notificator: Stream failed', error.code, error.message, error.cause);

                    if (error.message.includes('ErrCharLock')) {
                        handleGRPCError(error);
                    } else {
                        this.restartStream();
                    }
                }
            }

            console.debug('Notificator: Stream ended');
        },

        async stopStream(): Promise<void> {
            if (this.abort === undefined) {
                return;
            }

            this.abort?.abort();
            this.abort = undefined;
            console.debug('Notificator: Stopping Data Stream');
        },

        async restartStream(): Promise<void> {
            this.reconnecting = true;

            // Reset back off time when over the max back off time
            if (this.reconnectBackoffTime > maxBackoffTime) {
                this.reconnectBackoffTime = initialReconnectBackoffTime;
            } else {
                this.reconnectBackoffTime += initialReconnectBackoffTime;
            }

            console.debug('Notificator: Restart back off time in', this.reconnectBackoffTime, 'seconds');
            await this.stopStream();

            setTimeout(async () => {
                if (this.reconnecting) {
                    this.startStream();
                }
            }, this.reconnectBackoffTime * 1000);
        },

        // Notification Actions
        async markNotifications(req: MarkNotificationsRequest): Promise<void> {
            try {
                await getGRPCNotificatorClient().markNotifications(req);
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }

            if (req.all === true || req.ids.length >= this.notificationsCount) {
                this.notificationsCount = 0;
            } else {
                this.notificationsCount = this.notificationsCount - req.ids.length;
            }
        },
    },
    getters: {
        getNotificationsCount(): number {
            return this.notificationsCount;
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useNotificatorStore as unknown as StoreDefinition, import.meta.hot));
}
