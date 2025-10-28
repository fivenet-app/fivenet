import type { DuplexStreamingCall } from '@protobuf-ts/runtime-rpc';
import { defineStore } from 'pinia';
import { useGRPCWebsocketTransport } from '~/composables/grpc/grpcws';
import { notificationsEvents } from '~/composables/useClientUpdate';
import type { Notification } from '~/types/notifications';
import { getNotificationsNotificationsClient } from '~~/gen/ts/clients';
import type { ObjectEvent, ObjectType } from '~~/gen/ts/resources/notifications/client_view';
import {
    NotificationCategory,
    NotificationType,
    type Notification as ProtoNotification,
} from '~~/gen/ts/resources/notifications/notifications';
import type { UserInfoChanged } from '~~/gen/ts/resources/userinfo/user_info';
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
        /**
         * Indicates whether the user has enabled Do Not Disturb mode.
         */
        const doNotDisturb = ref<boolean>(false);

        /**
         * List of notifications currently displayed to the user.
         */
        const notifications = ref<Notification[]>([]);

        /**
         * Count of unread notifications.
         */
        const notificationsCount = ref<number>(0);

        /**
         * Controller to manage the abort signal for the notification stream.
         */
        const abort = ref<AbortController | undefined>(undefined);

        /**
         * Indicates whether the notification stream is ready.
         */
        const ready = ref<boolean>(false);

        /**
         * Indicates whether the notification stream is reconnecting.
         */
        const reconnecting = ref<boolean>(false);

        /**
         * Time in seconds for the next reconnect attempt.
         */
        const reconnectBackoffTime = ref<number>(0);

        /**
         * ID of the dismissed banner message, if any.
         */
        const dismissedBannerMessageID = ref<string | undefined>();

        const notificationSound = useSounds('/sounds/notification.mp3');

        // Actions
        /**
         * Removes a notification by its ID.
         * @param notId - The ID of the notification to remove.
         */
        const remove = (notId: number): void => {
            notifications.value = notifications.value.filter((n) => n.id !== notId);
        };

        /**
         * Adds a new notification to the list.
         * @param notification - The notification to add.
         */
        const add = (notification: Notification): void => {
            if (notification.duration === undefined) {
                notification.duration = useAppConfig().timeouts.notification;
            }
            if (notification.id === undefined) {
                notification.id = Date.now() + Math.floor(Math.random() * 1000);
            }
            notifications.value.push(notification);
        };

        /**
         * Resets the notification data, clearing the list and count.
         */
        const resetData = (): void => {
            notificationsCount.value = 0;
            notifications.value = [];
        };

        // Stream
        let currentStream: DuplexStreamingCall<StreamRequest, StreamResponse> | undefined = undefined;

        // Handles userEvent logic
        /**
         * Handles the refresh token event by refreshing the user's token.
         * @param authStore - The authentication store instance.
         */
        const handleRefreshTokenEvent = async (authStore: ReturnType<typeof useAuthStore>): Promise<void> => {
            logger.info('Refreshing token...');
            await authStore.chooseCharacter(undefined);
        };

        /**
         * Handles a notification event by processing and adding it to the notifications list.
         * @param n - The notification data from the server.
         * @param calendarStore - The calendar store instance for handling calendar-related notifications.
         */
        const handleNotificationEvent = (n: ProtoNotification, calendarStore: ReturnType<typeof useCalendarStore>): void => {
            if (doNotDisturb.value) return;

            const nType = n.type !== NotificationType.UNSPECIFIED ? n.type : NotificationType.INFO;

            if (!n.title || !n.content) return;

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
                    label: { key: 'common.click_here', parameters: {} },
                    to: n.data.link.to,
                    external: n.data.link.external,
                });
            }

            if (n.category === NotificationCategory.CALENDAR) {
                notificationSound.play();

                if (n.data?.calendar !== undefined) {
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
        };

        /**
         * Handles the user info changed event by updating the active character's job information.
         * @param userInfoChanged - The user info change data from the server.
         */
        const handleUserInfoChangedEvent = (userInfoChanged: UserInfoChanged): void => {
            const { activeChar } = useAuth();

            if (userInfoChanged.newJob) activeChar.value!.job = userInfoChanged.newJob;
            if (userInfoChanged.newJobLabel) activeChar.value!.jobLabel = userInfoChanged.newJobLabel;
            if (userInfoChanged.newJobGrade) activeChar.value!.jobGrade = userInfoChanged.newJobGrade;
            if (userInfoChanged.newJobGradeLabel) activeChar.value!.jobGradeLabel = userInfoChanged.newJobGradeLabel;
        };

        /**
         * Handles a user event by delegating to the appropriate handler based on the event type.
         * @param resp - The stream response containing the user event data.
         * @param authStore - The authentication store instance.
         */
        const handleUserEvent = async (resp: StreamResponse, authStore: ReturnType<typeof useAuthStore>): Promise<void> => {
            if (resp.data.oneofKind !== 'userEvent') return;

            const userEvent = resp.data.userEvent;

            if (userEvent.data.oneofKind === 'refreshToken') {
                await handleRefreshTokenEvent(authStore);
            } else if (userEvent.data.oneofKind === 'notification') {
                const calendarStore = useCalendarStore();
                handleNotificationEvent(userEvent.data.notification, calendarStore);
            } else if (userEvent.data.oneofKind === 'notificationsReadCount') {
                notificationsCount.value = userEvent.data.notificationsReadCount;
            } else if (userEvent.data.oneofKind === 'userInfoChanged') {
                handleUserInfoChangedEvent(userEvent.data.userInfoChanged);
            } else {
                logger.warn('Unknown userEvent data received - oneofKind:', userEvent.data.oneofKind, userEvent.data);
            }
        };

        /**
         * Sets the ready state of the notification stream.
         * @param isReady - Whether the stream is ready.
         */
        const setStreamReadyState = (isReady: boolean): void => {
            ready.value = isReady;
            notificationsEvents.emit('ready', isReady);
        };

        /**
         * Starts the notification stream.
         */
        const startStream = async (): Promise<void> => {
            if (abort.value !== undefined) return;

            logger.debug('Starting Stream');
            abort.value = new AbortController();
            reconnecting.value = true;

            const notificationsNotificationsClient = await getNotificationsNotificationsClient();

            const authStore = useAuthStore();
            const { activeChar, can } = useAuth();

            const mailerStore = useMailerStore();

            try {
                currentStream = notificationsNotificationsClient.stream({ abort: abort.value.signal });
                setStreamReadyState(true);

                for await (const resp of currentStream.responses) {
                    notificationsCount.value = resp.notificationCount;

                    if (!resp || !resp.data || resp.data.oneofKind === undefined) continue;

                    if (resp.data.oneofKind === 'userEvent') {
                        await handleUserEvent(resp, authStore);
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
            } finally {
                setStreamReadyState(false);
                logger.debug('Stream ended');
            }
        };

        /**
         * Stops the notification stream.
         * @param end - Whether to end the stream completely.
         */
        const stopStream = async (end?: boolean): Promise<void> => {
            if (!abort.value) return;
            logger.debug('Stopping Stream');

            if (end) {
                reconnecting.value = false;
                reconnectBackoffTime.value = 0;
            } else {
                reconnecting.value = true;
            }

            setStreamReadyState(false);
            abort.value?.abort();
            abort.value = undefined;
        };

        /**
         * Restarts the notification stream.
         */
        const restartStream = async (): Promise<void> => {
            if (!reconnecting.value) {
                logger.debug('Reconnect is disabled, not restarting stream');
                return;
            }

            reconnecting.value = true;

            if (!abort.value || abort.value?.signal?.aborted) return;

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

        /**
         * Marks notifications as read.
         * @param req - The request containing notification IDs to mark as read.
         */
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

        /**
         * Sends a client view event.
         * @param viewType - The type of the object view.
         * @param id - The ID of the object (optional).
         */
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

        /**
         * Registers a callback for client update events.
         * @param objType - The type of the object.
         * @param callback - The callback function to execute on update.
         */
        const onClientUpdate = (objType: ObjectType, callback: (event: ObjectEvent) => void): void => {
            notificationsEvents.on(objType, callback);
        };

        /**
         * Unregisters a callback for client update events.
         * @param objType - The type of the object.
         * @param callback - The callback function to remove.
         */
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
