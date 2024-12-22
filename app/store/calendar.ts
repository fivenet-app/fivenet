import { format } from 'date-fns';
import Dexie, { type Table } from 'dexie';
import 'dexie-observable';
import 'dexie-syncable';
import { defineStore } from 'pinia';
import type { Calendar, CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';
import { RsvpResponses } from '~~/gen/ts/resources/calendar/calendar';
import { NotificationCategory, NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type {
    CreateOrUpdateCalendarEntryResponse,
    CreateOrUpdateCalendarResponse,
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
} from '~~/gen/ts/services/calendar/calendar';
import { useNotificatorStore } from './notificator';
import { useSettingsStore } from './settings';

const logger = useLogger('📅 Calendar');

export interface CalendarState {
    activeCalendarIds: string[];
    view: 'month' | 'week' | 'summary';
    currentDate: {
        year: number;
        month: number;
    };
    calendars: Calendar[];
    entries: CalendarEntry[];
    eventReminders: Map<string, number>;
}

export const useCalendarStore = defineStore('calendar', {
    state: () =>
        ({
            activeCalendarIds: [],
            view: 'month',
            currentDate: {
                year: new Date().getFullYear(),
                month: new Date().getMonth() + 1,
            },
            calendars: [],
            entries: [],
            eventReminders: new Map<string, number>(),
        }) as CalendarState,
    persist: {
        pick: ['activeCalendarIds', 'view'],
    },
    actions: {
        async checkAppointments(): Promise<void> {
            try {
                const reminderTimes = useSettingsStore().calendar.reminderTimes;
                const highestReminder = Math.max(...reminderTimes);

                const response = await this.getUpcomingEntries({
                    seconds: highestReminder + 10,
                });

                const now = new Date();
                response.entries.forEach((entry) => {
                    calendarDB.entries.add(entry);
                    calendarDB.entries.delete(entry);

                    const startTime = toDate(entry.startTime);
                    const time = startTime.getTime() - now.getTime();

                    const closestTime = reminderTimes.reduce((prev, curr) =>
                        Math.abs(curr - time) < Math.abs(prev - time) ? curr : prev,
                    );

                    if (this.eventReminders.get(entry.id) === closestTime) {
                        return;
                    }

                    if (closestTime > time) {
                        return;
                    }

                    if (time <= 0) {
                        this.eventReminders.delete(entry.id);
                    } else {
                        this.eventReminders.set(entry.id, closestTime);
                    }

                    useNotificatorStore().add({
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

                    useSound().play({ name: 'notification' });
                });
            } catch (e) {
                logger.error('error while getting upcoming events', e);
            }
        },
        // Calendars
        async getCalendar(req: GetCalendarRequest): Promise<GetCalendarResponse> {
            const call = getGRPCCalendarClient().getCalendar(req);
            const { response } = await call;

            if (response.calendar) {
                const idx = this.calendars.findIndex((c) => c.id === response.calendar!.id);
                if (idx > -1) {
                    this.calendars[idx] = response.calendar;
                } else {
                    this.calendars.push(response.calendar);
                }
            }

            return response;
        },
        async listCalendars(req: ListCalendarsRequest): Promise<ListCalendarsResponse> {
            try {
                const call = getGRPCCalendarClient().listCalendars(req);
                const { response } = await call;

                // Only "register" calendars in list when they are accessible by the user
                if (!req.onlyPublic) {
                    if (response.calendars.length > 0) {
                        const foundCalendars: string[] = [];
                        response.calendars.forEach((calendar) => {
                            const idx = this.calendars.findIndex((c) => c.id === calendar!.id);
                            if (idx > -1) {
                                this.calendars[idx] = calendar;
                            } else {
                                this.calendars.push(calendar);
                            }
                            foundCalendars.push(calendar.id);
                        });

                        // Remove non-accessible calendars (ignore public ones) and their entries from our list
                        this.calendars = this.calendars.filter((calendar): boolean => {
                            if (!calendar.public) {
                                return true;
                            }

                            if (foundCalendars.find((c) => c === calendar.id)) {
                                return true;
                            }

                            this.entries = this.entries.filter((entry) => entry.calendarId === calendar.id);

                            return false;
                        });
                    } else {
                        this.calendars.length = 0;
                    }
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
        async createOrUpdateCalendar(calendar: Calendar): Promise<CreateOrUpdateCalendarResponse> {
            const call = getGRPCCalendarClient().createOrUpdateCalendar({
                calendar: calendar,
            });
            const { response } = await call;

            if (response.calendar) {
                const idx = this.calendars.findIndex((c) => c.id === response.calendar!.id);
                if (idx > -1) {
                    this.calendars[idx] = response.calendar;
                } else {
                    this.calendars.push(response.calendar);
                }

                this.activeCalendarIds.push(response.calendar.id);
            }

            return response;
        },
        async deleteCalendar(id: string): Promise<void> {
            try {
                const call = getGRPCCalendarClient().deleteCalendar({
                    calendarId: id,
                });
                await call;

                const idx = this.calendars.findIndex((c) => c.id === id);
                if (idx > -1) {
                    this.calendars.splice(idx, 1);
                }
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // Entries
        async getCalendarEntry(req: GetCalendarEntryRequest): Promise<GetCalendarEntryResponse> {
            const call = getGRPCCalendarClient().getCalendarEntry(req);
            const { response } = await call;

            if (response.entry) {
                const idx = this.entries.findIndex((c) => c.id === response.entry!.id);
                if (idx > -1) {
                    this.entries[idx] = response.entry;
                } else {
                    this.entries.push(response.entry);
                }
            }

            return response;
        },
        async listCalendarEntries(req?: ListCalendarEntriesRequest): Promise<ListCalendarEntriesResponse> {
            if (req === undefined) {
                req = {
                    calendarIds: this.activeCalendarIds,
                    year: this.currentDate.year,
                    month: this.currentDate.month,
                };
            }

            try {
                const call = getGRPCCalendarClient().listCalendarEntries(req);
                const { response } = await call;

                if (response.entries.length > 0) {
                    response.entries.forEach((entry) => {
                        // Make sure that we have the calendar in our list before adding it
                        if (!this.calendars.find((c) => c.id === entry.calendarId)) {
                            return;
                        }

                        const idx = this.entries.findIndex((c) => c.id === entry!.id);
                        if (idx > -1) {
                            this.entries[idx] = entry;
                        } else {
                            this.entries.push(entry);
                        }
                    });
                } else {
                    this.entries.length = 0;
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
        async getUpcomingEntries(req: GetUpcomingEntriesRequest): Promise<GetUpcomingEntriesResponse> {
            try {
                const call = getGRPCCalendarClient().getUpcomingEntries(req);
                const { response } = await call;

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
        async createOrUpdateCalendarEntry(
            entry: CalendarEntry,
            users?: UserShort[],
        ): Promise<CreateOrUpdateCalendarEntryResponse> {
            const call = getGRPCCalendarClient().createOrUpdateCalendarEntry({
                entry: entry,
                userIds: users?.map((u) => u.userId) ?? [],
            });
            const { response } = await call;

            if (response.entry) {
                const idx = this.entries.findIndex((e) => e.id === response.entry?.id);
                if (idx > -1) {
                    this.entries[idx] = response.entry;
                } else {
                    this.entries.push(response.entry);
                }
            }

            return response;
        },

        async deleteCalendarEntry(entryId: string): Promise<void> {
            try {
                const call = getGRPCCalendarClient().deleteCalendarEntry({
                    entryId: entryId,
                });
                await call;

                const idx = this.entries.findIndex((c) => c.id === entryId);
                if (idx > -1) {
                    this.entries.splice(idx, 1);
                }
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // RSVP
        async listCalendarEntryRSVP(req: ListCalendarEntryRSVPRequest): Promise<ListCalendarEntryRSVPResponse> {
            try {
                const call = getGRPCCalendarClient().listCalendarEntryRSVP(req);
                const { response } = await call;

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
        async rsvpCalendarEntry(req: RSVPCalendarEntryRequest): Promise<RSVPCalendarEntryResponse> {
            try {
                const call = getGRPCCalendarClient().rSVPCalendarEntry(req);
                const { response } = await call;

                // Retrieve calendar entry if a "should be visible" response and it is not in our list yet
                if (req.entry?.entryId && response.entry?.response && response.entry.response > RsvpResponses.HIDDEN) {
                    const entry = this.entries.find((e) => e.id === response.entry?.entryId);
                    if (!entry) {
                        await this.getCalendarEntry({ entryId: req.entry?.entryId });
                    } else {
                        entry.rsvp = response.entry;
                    }
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
    },
    getters: {
        hasPrivateCalendar: (state) => {
            const { activeChar } = useAuth();
            return !!state.calendars.find((c) => c.job === undefined && c.creatorId === activeChar.value?.userId);
        },
    },
});

class CalendarDexie extends Dexie {
    calendars!: Table<Calendar>;
    entries!: Table<CalendarEntry>;

    constructor() {
        super('calendar');
        this.version(1).stores({
            calendars: 'id',
            entries: 'id, calendarId',
            rsvps: 'entry_id, user_id',
        });
    }
}

export const calendarDB = new CalendarDexie();

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCalendarStore, import.meta.hot));
}
