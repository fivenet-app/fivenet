import { defineStore } from 'pinia';
import { v4 as uuidv4 } from 'uuid';
import { useGRPCWebsocketTransport } from '~/composables/grpcws';
import type { Notification } from '~/composables/notifications';
import { useAuthStore } from '~/store/auth';
import { NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { MarkNotificationsRequest } from '~~/gen/ts/services/notificator/notificator';
import { useCalendarStore } from './calendar';
import { useMailerStore } from './mailer';

const logger = useLogger('📣 Notificator');

// In seconds
const maxBackoffTime = 30;
const initialReconnectBackoffTime = 2;

export interface NotificationsState {
    doNotDisturb: boolean;

    notifications: Notification[];
    notificationsCount: number;

    abort: AbortController | undefined;
    reconnecting: boolean;
    reconnectBackoffTime: number;
}

export const useNotificatorStore = defineStore('notifications', {
    state: () =>
        ({
            doNotDisturb: false,

            notifications: [],
            notificationsCount: 0,

            abort: undefined,
            reconnecting: false,
            reconnectBackoffTime: 0,
        }) as NotificationsState,
    persist: {
        pick: ['doNotDisturb'],
    },
    actions: {
        remove(id: string): void {
            this.notifications = this.notifications.filter((notification) => notification.id !== id);
        },
        add(notification: Notification): void {
            if (notification.id === undefined) {
                notification.id = uuidv4();
            }

            if (notification.timeout === undefined) {
                notification.timeout = useAppConfig().timeouts.notification;
            }

            this.notifications.push(notification);
        },

        async startStream(): Promise<void> {
            if (this.abort !== undefined) {
                return;
            }

            logger.debug('Starting Stream');

            this.abort = new AbortController();
            this.reconnecting = false;

            const authStore = useAuthStore();
            const { can } = useAuth();

            try {
                const call = getGRPCNotificatorClient().stream(
                    {},
                    {
                        abort: this.abort.signal,
                    },
                );

                for await (const resp of call.responses) {
                    this.notificationsCount = resp.notificationCount;

                    if (resp === undefined || !resp.data || resp.data.oneofKind === undefined) {
                        continue;
                    }

                    if (resp.data.oneofKind === 'userEvent') {
                        if (resp.data.userEvent.data.oneofKind === 'refreshToken') {
                            logger.info('Refreshing token...');
                            await authStore.chooseCharacter(undefined);
                        } else if (resp.data.userEvent.data.oneofKind === 'notification') {
                            // Don't display server notifications when do not disturb is on
                            if (this.doNotDisturb) {
                                continue;
                            }

                            const n = resp.data.userEvent.data.notification;
                            const nType: NotificationType =
                                n.type !== NotificationType.UNSPECIFIED ? n.type : NotificationType.INFO;

                            if (n.title === undefined || n.content === undefined) {
                                continue;
                            }

                            const not: Notification = {
                                title: n.title,
                                description: n.content,
                                type: nType,
                                category: n.category,
                                data: n.data,
                                actions: [],
                            };

                            if (n.data?.link !== undefined) {
                                not.actions?.push({
                                    label: { key: 'common.click_here' },
                                    to: n.data.link.to,
                                    external: n.data.link.external,
                                });
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
                                    } catch (e) {
                                        logger.warn(
                                            'Error while retrieving calendar/entry for calendar notification, Notification ID:',
                                            n.id,
                                            'Error:',
                                            e,
                                        );
                                    }
                                }
                            }

                            this.add(not);
                        }
                    } else if (resp.data.oneofKind === 'jobEvent') {
                        if (resp.data.jobEvent.data.oneofKind === 'jobProps') {
                            authStore.setJobProps(resp.data.jobEvent.data.jobProps);
                        } else {
                            logger.warn('Unknown job event data received - Kind: ', resp.data.oneofKind, resp.data);
                        }
                    } else if (resp.data.oneofKind === 'jobGradeEvent') {
                        if (resp.data.jobGradeEvent.data.oneofKind === 'refreshToken') {
                            await authStore.chooseCharacter(undefined);
                        }
                    } else if (resp.data.oneofKind === 'systemEvent') {
                        logger.warn('No systemEvent handlers available.', resp.data);
                    } else if (resp.data.oneofKind === 'mailerEvent') {
                        if (can('MailerService.ListEmails').value) {
                            useMailerStore().handleEvent(resp.data.mailerEvent);
                        }
                    }

                    if (resp.restart) {
                        logger.debug('Server requested stream to be restarted');
                        this.reconnectBackoffTime = 0;
                        this.stopStream();
                        useGRPCWebsocketTransport().close();
                        this.restartStream();
                        return;
                    }
                }
            } catch (e) {
                const error = e as RpcError;
                if (error.code !== 'CANCELLED' && error.code !== 'ABORTED') {
                    logger.debug('Stream failed', error.code, error.message, error.cause);

                    if (error.message.includes('ErrCharLock')) {
                        await handleGRPCError(error);
                    }
                    if (this.abort !== undefined && !this.abort?.signal.aborted) {
                        this.restartStream();
                    }
                } else {
                    await this.restartStream();
                }
            }

            logger.debug('Stream ended');
        },

        async stopStream(): Promise<void> {
            if (this.abort === undefined) {
                return;
            }

            this.abort.abort();
            this.abort = undefined;
            logger.debug('Stopping Stream');
        },

        async restartStream(): Promise<void> {
            this.reconnecting = true;

            // Reset back off time when over the max back off time
            if (this.reconnectBackoffTime > maxBackoffTime) {
                this.reconnectBackoffTime = initialReconnectBackoffTime;
            } else {
                this.reconnectBackoffTime += initialReconnectBackoffTime;
            }

            logger.debug('Restart back off time in', this.reconnectBackoffTime, 'seconds');
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
                this.notificationsCount -= req.ids.length;
            }
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useNotificatorStore, import.meta.hot));
}
