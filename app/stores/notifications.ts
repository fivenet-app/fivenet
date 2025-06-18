import { defineStore } from 'pinia';
import { useGRPCWebsocketTransport } from '~/composables/grpc/grpcws';
import { useAuthStore } from '~/stores/auth';
import type { Notification } from '~/utils/notifications';
import { NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { MarkNotificationsRequest } from '~~/gen/ts/services/notifications/notifications';
import { useCalendarStore } from './calendar';
import { useMailerStore } from './mailer';

const logger = useLogger('📣 Notificator');

// In seconds
const maxBackoffTime = 20;
const initialReconnectBackoffTime = 2;

export const useNotificationsStore = defineStore(
    'notifications',
    () => {
        const { $grpc } = useNuxtApp();

        const notificationSound = useSounds('/sounds/notification.mp3');

        // State
        const doNotDisturb = ref<boolean>(false);
        const notifications = ref<Notification[]>([]);
        const notificationsCount = ref<number>(0);

        const abort = ref<AbortController | undefined>(undefined);
        const reconnecting = ref<boolean>(false);
        const reconnectBackoffTime = ref<number>(0);

        const dismissedBannerMessageID = ref<string | undefined>();

        // Actions
        const remove = (notId: number): void => {
            notifications.value = notifications.value.filter((n) => n.id !== notId);
        };

        const add = (notification: Notification): void => {
            if (notification.timeout === undefined) {
                notification.timeout = useAppConfig().timeouts.notification;
            }
            notifications.value.push(notification);
        };

        const resetData = (): void => {
            notificationsCount.value = 0;
            notifications.value = [];
        };

        const startStream = async (): Promise<void> => {
            if (abort.value !== undefined) {
                return;
            }

            logger.debug('Starting Stream');
            abort.value = new AbortController();
            reconnecting.value = false;

            const authStore = useAuthStore();
            const { can } = useAuth();

            try {
                const call = $grpc.notifications.notifications.stream({ abort: abort.value.signal });

                for await (const resp of call.responses) {
                    notificationsCount.value = resp.notificationCount;

                    if (!resp || !resp.data || resp.data.oneofKind === undefined) {
                        continue;
                    }

                    if (resp.data.oneofKind === 'userEvent') {
                        if (resp.data.userEvent.data.oneofKind === 'refreshToken') {
                            logger.info('Refreshing token...');
                            await authStore.chooseCharacter(undefined);
                        } else if (resp.data.userEvent.data.oneofKind === 'notification') {
                            if (doNotDisturb.value) {
                                continue;
                            }

                            const n = resp.data.userEvent.data.notification;
                            const nType = n.type !== NotificationType.UNSPECIFIED ? n.type : NotificationType.INFO;

                            if (!n.title || !n.content) {
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
                                notificationSound.play();

                                if (n.data?.calendar !== undefined) {
                                    const calendarStore = useCalendarStore();
                                    try {
                                        if (n.data.calendar.calendarId !== undefined) {
                                            calendarStore.getCalendar({ calendarId: n.data.calendar.calendarId });
                                        } else if (n.data.calendar.calendarEntryId !== undefined) {
                                            calendarStore.getCalendarEntry({ entryId: n.data.calendar.calendarEntryId });
                                        }
                                    } catch (e) {
                                        logger.warn('Error retrieving calendar notification data:', e);
                                    }
                                }
                            }

                            add(not);
                        } else {
                            logger.warn('Unknown userEvent data received - oneofKind:', resp.data.oneofKind, resp.data);
                        }
                    } else if (resp.data.oneofKind === 'jobEvent') {
                        if (resp.data.jobEvent.data.oneofKind === 'jobProps') {
                            authStore.setJobProps(resp.data.jobEvent.data.jobProps);
                        } else {
                            logger.warn('Unknown jobEvent data received - oneofKind:', resp.data.oneofKind, resp.data);
                        }
                    } else if (resp.data.oneofKind === 'jobGradeEvent') {
                        if (resp.data.jobGradeEvent.data.oneofKind === 'refreshToken') {
                            await authStore.chooseCharacter(undefined);
                        } else {
                            logger.warn(
                                'Unknown jobGradeEvent event data received - oneofKind:',
                                resp.data.oneofKind,
                                resp.data,
                            );
                        }
                    } else if (resp.data.oneofKind === 'systemEvent') {
                        if (resp.data.systemEvent.data.oneofKind === 'ping') {
                            // Pong!
                        } else if (resp.data.systemEvent.data.oneofKind === 'bannerMessage') {
                            const { system } = useAppConfig();
                            resp.data.systemEvent.data.bannerMessage.bannerMessageEnabled =
                                resp.data.systemEvent.data.bannerMessage.bannerMessageEnabled ?? false;
                            if (resp.data.systemEvent.data.bannerMessage.bannerMessage === undefined) {
                                system.bannerMessage = undefined;
                                continue;
                            }

                            if (system.bannerMessage?.id === resp.data.systemEvent.data.bannerMessage.bannerMessage.id) {
                                continue;
                            }

                            system.bannerMessage = resp.data.systemEvent.data.bannerMessage.bannerMessage;
                        } else {
                            logger.warn('Unknown systemEvent event data received - oneofKind:', resp.data.oneofKind, resp.data);
                        }
                    } else if (resp.data.oneofKind === 'mailerEvent') {
                        if (can('mailer.MailerService/ListEmails').value) {
                            useMailerStore().handleEvent(resp.data.mailerEvent);
                        }
                    }

                    if (resp.restart) {
                        logger.debug('Server requested stream to be restarted');
                        reconnectBackoffTime.value = 0;
                        await stopStream();
                        useGRPCWebsocketTransport().close();
                        restartStream();
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
                }

                if (!abort.value?.signal.aborted) {
                    await restartStream();
                }
            }

            logger.debug('Stream ended');
        };

        const stopStream = async (): Promise<void> => {
            if (!abort.value) {
                return;
            }
            abort.value.abort();
            logger.debug('Stopping Stream');
            abort.value = undefined;
        };

        const restartStream = async (): Promise<void> => {
            if (!abort.value || abort.value.signal.aborted) {
                return;
            }
            reconnecting.value = true;

            if (reconnectBackoffTime.value > maxBackoffTime) {
                reconnectBackoffTime.value = initialReconnectBackoffTime;
            } else {
                reconnectBackoffTime.value += initialReconnectBackoffTime;
            }

            logger.debug('Restart back off time in', reconnectBackoffTime.value, 'seconds');
            await stopStream();

            setTimeout(async () => {
                if (reconnecting.value) {
                    startStream();
                }
            }, reconnectBackoffTime.value * 1000);
        };

        const markNotifications = async (req: MarkNotificationsRequest): Promise<void> => {
            try {
                await $grpc.notifications.notifications.markNotifications(req);
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }

            if (req.all === true || req.ids.length >= notificationsCount.value) {
                notificationsCount.value = 0;
            } else {
                notificationsCount.value -= req.ids.length;
            }
        };

        return {
            // State
            doNotDisturb,
            notifications,
            notificationsCount,
            abort,
            reconnecting,
            reconnectBackoffTime,

            dismissedBannerMessageID,

            // Actions
            remove,
            add,
            resetData,
            startStream,
            stopStream,
            restartStream,
            markNotifications,
        };
    },
    {
        persist: {
            pick: ['doNotDisturb', 'dismissedBannerMessageID'],
        },
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useNotificationsStore, import.meta.hot));
}
