import type { DuplexStreamingCall } from '@protobuf-ts/runtime-rpc';
import { defineStore } from 'pinia';
import { useGRPCWebsocketTransport } from '~/composables/grpc/grpcws';
import { notificationsEvents } from '~/composables/useClientUpdate';
import { useAuthStore } from '~/stores/auth';
import type { Notification } from '~/types/notifications';
import { getNotificationsNotificationsClient } from '~~/gen/ts/clients';
import type { ObjectEvent, ObjectType } from '~~/gen/ts/resources/notifications/client_view';
import { NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { MarkNotificationsRequest, StreamRequest, StreamResponse } from '~~/gen/ts/services/notifications/notifications';
import { useCalendarStore } from './calendar';
import { useMailerStore } from './mailer';

const logger = useLogger('📣 Notificator');

// In seconds
const maxBackoffTime = 20;
const initialReconnectBackoffTime = 2;

export const useNotificationsStore = defineStore(
    'notifications',
    () => {
        // State
        const doNotDisturb = ref<boolean>(false);
        const notifications = ref<Notification[]>([]);
        const notificationsCount = ref<number>(0);

        const abort = ref<AbortController | undefined>(undefined);
        const ready = ref<boolean>(false);
        const reconnecting = ref<boolean>(false);
        const reconnectBackoffTime = ref<number>(0);

        const dismissedBannerMessageID = ref<string | undefined>();

        // Actions
        const remove = (notId: number): void => {
            notifications.value = notifications.value.filter((n) => n.id !== notId);
        };

        const add = (notification: Notification): void => {
            if (notification.duration === undefined) {
                notification.duration = useAppConfig().timeouts.notification;
            }
            notifications.value.push(notification);
        };

        const resetData = (): void => {
            notificationsCount.value = 0;
            notifications.value = [];
        };

        // Stream
        let currentStream: DuplexStreamingCall<StreamRequest, StreamResponse> | undefined = undefined;

        const startStream = async (): Promise<void> => {
            if (abort.value !== undefined) return;

            logger.debug('Starting Stream');
            abort.value = new AbortController();
            reconnecting.value = false;

            const notificationsNotificationsClient = await getNotificationsNotificationsClient();

            const authStore = useAuthStore();
            const { activeChar, can } = useAuth();

            const mailerStore = useMailerStore();

            const notificationSound = useSounds('/sounds/notification.mp3');

            try {
                currentStream = notificationsNotificationsClient.stream({ abort: abort.value.signal });
                ready.value = true;
                notificationsEvents.emit('ready', true);

                for await (const resp of currentStream.responses) {
                    notificationsCount.value = resp.notificationCount;

                    if (!resp || !resp.data || resp.data.oneofKind === undefined) continue;

                    if (resp.data.oneofKind === 'userEvent') {
                        if (resp.data.userEvent.data.oneofKind === 'refreshToken') {
                            logger.info('Refreshing token...');
                            await authStore.chooseCharacter(undefined);
                        } else if (resp.data.userEvent.data.oneofKind === 'notification') {
                            if (doNotDisturb.value) continue;

                            const n = resp.data.userEvent.data.notification;
                            const nType = n.type !== NotificationType.UNSPECIFIED ? n.type : NotificationType.INFO;

                            if (!n.title || !n.content) continue;

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
                        } else if (resp.data.userEvent.data.oneofKind === 'notificationsReadCount') {
                            notificationsCount.value = resp.data.userEvent.data.notificationsReadCount;
                        } else if (resp.data.userEvent.data.oneofKind === 'userInfoChanged') {
                            // Update the user job info in the auth store
                            if (resp.data.userEvent.data.userInfoChanged.newJob)
                                activeChar.value!.job = resp.data.userEvent.data.userInfoChanged.newJob;

                            if (resp.data.userEvent.data.userInfoChanged.newJobLabel)
                                activeChar.value!.jobLabel = resp.data.userEvent.data.userInfoChanged.newJobLabel;

                            if (resp.data.userEvent.data.userInfoChanged.newJobGrade)
                                activeChar.value!.jobGrade = resp.data.userEvent.data.userInfoChanged.newJobGrade;

                            if (resp.data.userEvent.data.userInfoChanged.newJobGradeLabel)
                                activeChar.value!.jobGradeLabel = resp.data.userEvent.data.userInfoChanged.newJobGradeLabel;
                        } else {
                            logger.warn('Unknown userEvent data received - oneofKind:', resp.data.oneofKind, resp.data);
                        }
                    } else if (resp.data.oneofKind === 'jobEvent') {
                        if (resp.data.jobEvent.data.oneofKind === 'jobProps') {
                            // Check that the job props are for the active character's job
                            if (
                                !resp.data.jobEvent.data.jobProps ||
                                (activeChar.value?.job && resp.data.jobEvent.data.jobProps.job !== activeChar.value?.job)
                            )
                                continue;

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
                        if (resp.data.systemEvent.data.oneofKind === 'clientConfig') {
                            logger.info('Client config update received');
                            updateAppConfig({ ...resp.data.systemEvent.data.clientConfig });
                        } else {
                            logger.warn('Unknown systemEvent event data received - oneofKind:', resp.data.oneofKind, resp.data);
                        }
                    } else if (resp.data.oneofKind === 'objectEvent') {
                        const event = resp.data.objectEvent;

                        notificationsEvents.emit(event.type, event);
                    } else if (resp.data.oneofKind === 'mailerEvent') {
                        if (can('mailer.MailerService/ListEmails').value) mailerStore.handleEvent(resp.data.mailerEvent);
                    }

                    if (resp.restart) {
                        logger.debug('Server requested stream to be restarted');
                        reconnectBackoffTime.value = 0;
                        await stopStream();
                        useGRPCWebsocketTransport().close();
                        await restartStream();
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

            ready.value = false;
            notificationsEvents.emit('ready', false);
            logger.debug('Stream ended');
        };

        const stopStream = async (): Promise<void> => {
            if (!abort.value) return;

            ready.value = false;
            notificationsEvents.emit('ready', false);
            abort.value.abort();
            logger.debug('Stopping Stream');
            abort.value = undefined;
        };

        const restartStream = async (): Promise<void> => {
            if (!abort.value || abort.value.signal.aborted) return;

            reconnecting.value = true;

            if (reconnectBackoffTime.value > maxBackoffTime) {
                reconnectBackoffTime.value = initialReconnectBackoffTime;
            } else {
                reconnectBackoffTime.value += initialReconnectBackoffTime;
            }

            logger.debug('Restart back off time in', reconnectBackoffTime.value, 'seconds');
            await stopStream();

            useTimeoutFn(async () => {
                if (reconnecting.value) startStream();
            }, reconnectBackoffTime.value * 1000);
        };

        const markNotifications = async (req: MarkNotificationsRequest): Promise<void> => {
            const notificationsNotificationsClient = await getNotificationsNotificationsClient();

            try {
                await notificationsNotificationsClient.markNotifications(req);
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

        const sendClientView = async (viewType: ObjectType, id?: number): Promise<void> => {
            logger.debug('Sending client view', viewType, id, currentStream !== undefined);
            try {
                await currentStream?.requests.send({
                    data: {
                        oneofKind: 'clientView',
                        clientView: {
                            type: viewType,
                            id: id,
                        },
                    },
                });
            } catch (e) {
                logger.error('Failed to send client view', e);
                // Ignore any errors here, as this is just a client-side view
            }
        };

        const onClientUpdate = (objType: ObjectType, callback: (event: ObjectEvent) => void): void => {
            notificationsEvents.on(objType, callback);
        };

        const offClientUpdate = (objType: ObjectType, callback: (event: ObjectEvent) => void): void => {
            notificationsEvents.off(objType, callback);
        };

        return {
            // State
            doNotDisturb,
            notifications,
            notificationsCount,
            abort,
            ready,
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
            sendClientView,
            onClientUpdate,
            offClientUpdate,
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
