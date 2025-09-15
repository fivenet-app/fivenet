import { format } from 'date-fns';
import { defineStore } from 'pinia';
import { checkCalendarAccess } from '~/components/calendar/helpers';
import { getCalendarCalendarClient } from '~~/gen/ts/clients';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access';
import { type Calendar, type CalendarEntry, RsvpResponses } from '~~/gen/ts/resources/calendar/calendar';
import { NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type {
    CreateCalendarResponse,
    CreateOrUpdateCalendarEntryResponse,
    GetCalendarEntryRequest,
    GetCalendarEntryResponse,
    GetCalendarRequest,
    GetCalendarResponse,
    GetUpcomingEntriesRequest,
    GetUpcomingEntriesResponse,
    ListCalendarEntriesRequest,
    ListCalendarEntriesResponse,
    ListCalendarEntryRSVPRequest,
    ListCalendarEntryRSVPResponse,
    ListCalendarsRequest,
    ListCalendarsResponse,
    RSVPCalendarEntryRequest,
    RSVPCalendarEntryResponse,
    UpdateCalendarResponse,
} from '~~/gen/ts/services/calendar/calendar';

const logger = useLogger('📅 Calendar');

export const useCalendarStore = defineStore(
    'calendar',
    () => {
        const notifications = useNotificationsStore();
        const settings = useSettingsStore();

        // State
        const activeCalendarIds = ref<number[]>([]);
        const view = ref<'month' | 'week' | 'summary'>('month');
        const currentDate = ref({
            year: new Date().getFullYear(),
            month: new Date().getMonth() + 1,
        });
        const calendars = ref<Calendar[]>([]);
        const entries = ref<CalendarEntry[]>([]);
        const eventReminders = ref<Map<number, number>>(new Map());

        const notificationSound = useSounds('/sounds/notification.mp3');

        // Actions
        const checkAppointments = async (): Promise<void> => {
            try {
                const reminderTimes = settings.calendar.reminderTimes;
                const highestReminder = Math.max(...reminderTimes);

                const response = await getUpcomingEntries({
                    seconds: highestReminder + 10,
                });

                const now = new Date();
                response.entries.forEach((entry) => {
                    const startTime = toDate(entry.startTime);
                    const time = startTime.getTime() - now.getTime();

                    const closestTime = reminderTimes.reduce((prev, curr) =>
                        Math.abs(curr - time) < Math.abs(prev - time) ? curr : prev,
                    );

                    if (eventReminders.value.get(entry.id) === closestTime) return;

                    if (closestTime > time) return;

                    if (time <= 0) {
                        eventReminders.value.delete(entry.id);
                    } else {
                        eventReminders.value.set(entry.id, closestTime);
                    }

                    notifications.add({
                        title: {
                            key: 'notifications.calendar.event_starting.title',
                            parameters: {
                                title: entry.title,
                                name: entry.creator ? `${entry.creator.firstname} ${entry.creator.lastname}` : 'N/A',
                            },
                        },
                        description: {
                            key: 'notifications.calendar.event_starting.content',
                            parameters: {
                                time: format(startTime, 'HH:mm'),
                                ago: useTimeAgo(startTime).value,
                            },
                        },
                        type: NotificationType.INFO,
                        category: NotificationCategory.CALENDAR,
                        actions: [
                            {
                                label: { key: 'common.open' },
                                icon: 'i-mdi-calendar',
                                to: `/calendar?entry_id=${entry.id}`,
                            },
                        ],
                    });

                    notificationSound.play();
                });
            } catch (e) {
                logger.error('error while getting upcoming events', e);
            }
        };

        // Calendars
        const getCalendar = async (req: GetCalendarRequest): Promise<GetCalendarResponse> => {
            const calendarCalendarClient = await getCalendarCalendarClient();

            const call = calendarCalendarClient.getCalendar(req);
            const { response } = await call;

            if (response.calendar) {
                const idx = calendars.value.findIndex((c) => c.id === response.calendar!.id);
                if (idx > -1) {
                    calendars.value[idx] = response.calendar;
                } else {
                    calendars.value.push(response.calendar);
                }
            }

            return response;
        };

        const listCalendars = async (req: ListCalendarsRequest): Promise<ListCalendarsResponse> => {
            const calendarCalendarClient = await getCalendarCalendarClient();

            try {
                const call = calendarCalendarClient.listCalendars(req);
                const { response } = await call;

                // Only "register" calendars in list when they are accessible by the user
                if (!req.onlyPublic) {
                    if (response.calendars.length === 0) {
                        calendars.value.length = 0;
                    } else {
                        const foundCalendars: number[] = [];
                        response.calendars.forEach((calendar) => {
                            const idx = calendars.value.findIndex((c) => c.id === calendar!.id);
                            if (idx > -1) {
                                calendars.value[idx] = calendar;
                            } else {
                                calendars.value.push(calendar);
                            }
                            foundCalendars.push(calendar.id);
                        });

                        // Remove non-accessible calendars (ignore public ones) and their entries from our list
                        calendars.value = calendars.value.filter((calendar): boolean => {
                            if (!calendar.public) {
                                return true;
                            }

                            if (foundCalendars.find((c) => c === calendar.id)) {
                                return true;
                            }

                            entries.value = entries.value.filter((entry) => entry.calendarId === calendar.id);

                            return false;
                        });
                    }
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        const createOrUpdateCalendar = async (
            calendarParam: Calendar,
        ): Promise<CreateCalendarResponse | UpdateCalendarResponse> => {
            const calendarCalendarClient = await getCalendarCalendarClient();

            let call;
            if (calendarParam.id === 0) {
                call = calendarCalendarClient.createCalendar({
                    calendar: calendarParam,
                });
            } else {
                call = calendarCalendarClient.updateCalendar({
                    calendar: calendarParam,
                });
            }
            const { response } = await call;

            if (response.calendar) {
                const idx = calendars.value.findIndex((c) => c.id === response.calendar!.id);
                if (idx > -1) {
                    calendars.value[idx] = response.calendar;
                } else {
                    calendars.value.push(response.calendar);
                }

                activeCalendarIds.value.push(response.calendar.id);
            }

            return response;
        };

        const deleteCalendar = async (id: number): Promise<void> => {
            const calendarCalendarClient = await getCalendarCalendarClient();

            try {
                const call = calendarCalendarClient.deleteCalendar({
                    calendarId: id,
                });
                await call;

                const idx = calendars.value.findIndex((c) => c.id === id);
                if (idx > -1) {
                    calendars.value.splice(idx, 1);
                }
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        // Entries
        const getCalendarEntry = async (req: GetCalendarEntryRequest): Promise<GetCalendarEntryResponse> => {
            const calendarCalendarClient = await getCalendarCalendarClient();

            const call = calendarCalendarClient.getCalendarEntry(req);
            const { response } = await call;

            if (response.entry) {
                const idx = entries.value.findIndex((c) => c.id === response.entry!.id);
                if (idx > -1) {
                    entries.value[idx] = response.entry;
                } else {
                    entries.value.push(response.entry);
                }
            }

            return response;
        };

        const listCalendarEntries = async (req?: ListCalendarEntriesRequest): Promise<ListCalendarEntriesResponse> => {
            if (!req) {
                req = {
                    calendarIds: activeCalendarIds.value,
                    year: currentDate.value.year,
                    month: currentDate.value.month,
                };
            }
            const calendarCalendarClient = await getCalendarCalendarClient();

            try {
                const call = calendarCalendarClient.listCalendarEntries(req);
                const { response } = await call;

                if (response.entries.length > 0) {
                    response.entries.forEach((entry) => {
                        // Make sure that we have the calendar in our list before adding it
                        if (!calendars.value.find((c) => c.id === entry.calendarId)) return;

                        const idx = entries.value.findIndex((c) => c.id === entry!.id);
                        if (idx > -1) {
                            entries.value[idx] = entry;
                        } else {
                            entries.value.push(entry);
                        }
                    });
                } else {
                    entries.value.length = 0;
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        const getUpcomingEntries = async (req: GetUpcomingEntriesRequest): Promise<GetUpcomingEntriesResponse> => {
            const calendarCalendarClient = await getCalendarCalendarClient();

            try {
                const call = calendarCalendarClient.getUpcomingEntries(req);
                const { response } = await call;

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        const createOrUpdateCalendarEntry = async (
            entryParam: CalendarEntry,
            users?: number[],
        ): Promise<CreateOrUpdateCalendarEntryResponse> => {
            const calendarCalendarClient = await getCalendarCalendarClient();

            const call = calendarCalendarClient.createOrUpdateCalendarEntry({
                entry: entryParam,
                userIds: users || [],
            });
            const { response } = await call;

            if (response.entry) {
                const idx = entries.value.findIndex((e) => e.id === response.entry?.id);
                if (idx > -1) {
                    entries.value[idx] = response.entry;
                } else {
                    entries.value.push(response.entry);
                }
            }

            return response;
        };

        const deleteCalendarEntry = async (entryId: number): Promise<void> => {
            const calendarCalendarClient = await getCalendarCalendarClient();

            try {
                const call = calendarCalendarClient.deleteCalendarEntry({
                    entryId: entryId,
                });
                await call;

                const idx = entries.value.findIndex((c) => c.id === entryId);
                if (idx > -1) {
                    entries.value.splice(idx, 1);
                }
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        // RSVP
        const listCalendarEntryRSVP = async (req: ListCalendarEntryRSVPRequest): Promise<ListCalendarEntryRSVPResponse> => {
            const calendarCalendarClient = await getCalendarCalendarClient();

            try {
                const call = calendarCalendarClient.listCalendarEntryRSVP(req);
                const { response } = await call;

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        const rsvpCalendarEntry = async (req: RSVPCalendarEntryRequest): Promise<RSVPCalendarEntryResponse> => {
            const calendarCalendarClient = await getCalendarCalendarClient();

            try {
                const call = calendarCalendarClient.rSVPCalendarEntry(req);
                const { response } = await call;

                // Retrieve calendar entry if a "should be visible" response and it is not in our list yet
                if (req.entry?.entryId && response.entry?.response && response.entry.response > RsvpResponses.HIDDEN) {
                    const foundEntry = entries.value.find((e) => e.id === response.entry?.entryId);
                    if (!foundEntry) {
                        await getCalendarEntry({ entryId: req.entry?.entryId });
                    } else {
                        foundEntry.rsvp = response.entry;
                    }
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        // Getters
        const hasPrivateCalendar = computed(() => {
            const { activeChar } = useAuth();
            return !!calendars.value.find((c) => c.job === undefined && c.creatorId === activeChar.value?.userId);
        });

        const hasEditAccessToCalendar = computed(() => {
            const { activeChar } = useAuth();
            return !!calendars.value.find((c) => {
                if (c.job === undefined && c.creatorId === activeChar.value?.userId) return true;

                return checkCalendarAccess(c.access, c.creator, AccessLevel.EDIT);
            });
        });

        return {
            activeCalendarIds,
            view,
            currentDate,
            calendars,
            entries,
            eventReminders,

            // Actions
            checkAppointments,
            getCalendar,
            listCalendars,
            createOrUpdateCalendar,
            deleteCalendar,
            getCalendarEntry,
            listCalendarEntries,
            getUpcomingEntries,
            createOrUpdateCalendarEntry,
            deleteCalendarEntry,
            listCalendarEntryRSVP,
            rsvpCalendarEntry,

            // Getters
            hasPrivateCalendar,
            hasEditAccessToCalendar,
        };
    },
    {
        persist: {
            pick: ['activeCalendarIds', 'view'],
        },
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCalendarStore, import.meta.hot));
}
