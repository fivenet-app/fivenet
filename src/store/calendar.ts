import { defineStore, type StoreDefinition } from 'pinia';
import { CalendarEntry, type Calendar } from '~~/gen/ts/resources/calendar/calendar';
import type {
    CreateOrUpdateCalendarEntryResponse,
    CreateOrUpdateCalendarResponse,
    GetCalendarEntryRequest,
    GetCalendarEntryResponse,
    GetCalendarRequest,
    GetCalendarResponse,
    ListCalendarEntriesRequest,
    ListCalendarEntriesResponse,
    ListCalendarsRequest,
    ListCalendarsResponse,
} from '~~/gen/ts/services/calendar/calendar';
import { useAuthStore } from './auth';

export interface CalendarState {
    calendars: Calendar[];
    entries: CalendarEntry[];
}

export const useCalendarStore = defineStore('calendar', {
    state: () =>
        ({
            calendars: [],
            entries: [],
        }) as CalendarState,
    persist: {
        key(id) {
            return `state-${useAuthStore().activeChar?.userId}-${id}`;
        },
    },
    actions: {
        // Calendars
        async getCalendar(req: GetCalendarRequest): Promise<GetCalendarResponse> {
            const { $grpc } = useNuxtApp();

            try {
                const call = $grpc.getCalendarClient().getCalendar(req);
                const { response } = await call;

                return response;
            } catch (e) {
                throw e;
            }
        },
        async listCalendars(req: ListCalendarsRequest): Promise<ListCalendarsResponse> {
            const { $grpc } = useNuxtApp();

            try {
                const call = $grpc.getCalendarClient().listCalendars(req);
                const { response } = await call;

                this.calendars = response.calendars;

                return response;
            } catch (e) {
                $grpc.handleError(e as RpcError);
                throw e;
            }
        },
        async createOrUpdateCalendar(calendar: Calendar): Promise<CreateOrUpdateCalendarResponse> {
            const { $grpc } = useNuxtApp();

            try {
                const call = $grpc.getCalendarClient().createOrUpdateCalendar({
                    calendar: calendar,
                });
                const { response } = await call;

                return response;
            } catch (e) {
                throw e;
            }
        },

        // Entries
        async getCalendarEntry(req: GetCalendarEntryRequest): Promise<GetCalendarEntryResponse> {
            const { $grpc } = useNuxtApp();

            try {
                const call = $grpc.getCalendarClient().getCalendarEntry(req);
                const { response } = await call;

                return response;
            } catch (e) {
                throw e;
            }
        },
        async listCalendarEntries(req: ListCalendarEntriesRequest): Promise<ListCalendarEntriesResponse> {
            const { $grpc } = useNuxtApp();

            try {
                const call = $grpc.getCalendarClient().listCalendarEntries(req);
                const { response } = await call;

                return response;
            } catch (e) {
                $grpc.handleError(e as RpcError);
                throw e;
            }
        },
        async createOrUpdateCalendarEntry(entry: CalendarEntry): Promise<CreateOrUpdateCalendarEntryResponse> {
            const { $grpc } = useNuxtApp();

            try {
                const call = $grpc.getCalendarClient().createOrUpdateCalendarEntry({
                    entry: entry,
                });
                const { response } = await call;

                return response;
            } catch (e) {
                throw e;
            }
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useCalendarStore as unknown as StoreDefinition, import.meta.hot));
}
