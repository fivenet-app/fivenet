import { defineStore } from 'pinia';
import type { Calendar, CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';
import { RsvpResponses } from '~~/gen/ts/resources/calendar/calendar';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type {
    CreateOrUpdateCalendarEntryResponse,
    CreateOrUpdateCalendarResponse,
    GetCalendarEntryRequest,
    GetCalendarEntryResponse,
    GetCalendarRequest,
    GetCalendarResponse,
    ListCalendarEntriesRequest,
    ListCalendarEntriesResponse,
    ListCalendarEntryRSVPRequest,
    ListCalendarEntryRSVPResponse,
    ListCalendarsRequest,
    ListCalendarsResponse,
    RSVPCalendarEntryRequest,
    RSVPCalendarEntryResponse,
} from '~~/gen/ts/services/calendar/calendar';

export interface CalendarState {
    activeCalendarIds: string[];
    weeklyView: boolean;
    currentDate: {
        year: number;
        month: number;
    };
    calendars: Calendar[];
    entries: CalendarEntry[];
}

export const useCalendarStore = defineStore('calendar', {
    state: () =>
        ({
            activeCalendarIds: [],
            weeklyView: false,
            currentDate: {
                year: new Date().getFullYear(),
                month: new Date().getMonth() + 1,
            },
            calendars: [],
            entries: [],
        }) as CalendarState,
    persist: {
        pick: ['activeCalendarIds', 'weeklyView'],
    },
    actions: {
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
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCalendarStore, import.meta.hot));
}
