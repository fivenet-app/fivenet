<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { addHours, addMinutes, isSameDay, isSameHour, isSameMinute } from 'date-fns';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DatePickerPopoverClient from '~/components/partials/DatePickerPopover.client.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { useCalendarStore } from '~/store/calendar';
import { useCompletorStore } from '~/store/completor';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access';
import type { CalendarShort } from '~~/gen/ts/resources/calendar/calendar';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { CreateOrUpdateCalendarEntryResponse } from '~~/gen/ts/services/calendar/calendar';

const props = defineProps<{
    calendarId?: string;
    entryId?: string;
}>();

const { isOpen } = useModal();

const calendarStore = useCalendarStore();

const completorStore = useCompletorStore();

const usersLoading = ref(false);

const schema = z.object({
    calendar: z.custom<CalendarShort>().optional(),
    title: z.string().min(3).max(512),
    startTime: z.date(),
    endTime: z.date(),
    content: z.string().min(3).max(1000000),
    closed: z.boolean(),
    rsvpOpen: z.boolean(),
    users: z.custom<UserShort>().array().max(20),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    calendar: undefined,
    title: '',
    startTime: addHours(new Date(), 1),
    endTime: addHours(new Date(), 2),
    content: '',
    closed: false,
    rsvpOpen: true,
    users: [],
});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`calendar-entry:${props.entryId}`, () => calendarStore.getCalendarEntry({ entryId: props.entryId! }), {
    immediate: !!props.calendarId && !!props.entryId,
});

async function createOrUpdateCalendarEntry(values: Schema): Promise<CreateOrUpdateCalendarEntryResponse> {
    if (!values.calendar) {
        throw 'No Calendar selected';
    }

    try {
        const response = await calendarStore.createOrUpdateCalendarEntry(
            {
                id: data.value?.entry?.id ?? '0',
                calendarId: values.calendar.id,
                title: values.title,
                startTime: toTimestamp(values.startTime),
                endTime: toTimestamp(values.endTime),
                content: {
                    rawContent: values.content,
                },
                closed: values.closed,
                rsvpOpen: values.rsvpOpen,
                creatorJob: '',
            },
            state.users,
        );

        isOpen.value = false;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function setFromProps(): void {
    if (!data.value?.entry) {
        return;
    }

    const entry = data.value?.entry;
    if (entry.calendar) {
        state.calendar = entry.calendar;
    }

    state.title = entry.title;
    state.startTime = toDate(entry.startTime);
    state.endTime = toDate(entry.endTime);
    state.content = entry.content?.rawContent ?? '';
    state.closed = entry.closed;
    state.rsvpOpen = entry.rsvpOpen !== undefined;
}

watch(data, () => setFromProps());
watch(props, async () => refresh());

watch(
    () => state.startTime,
    () => {
        const endTime = state.endTime;

        if (state.startTime && !isSameDay(state.startTime, state.endTime)) {
            endTime.setFullYear(state.startTime.getFullYear());
            endTime.setMonth(state.startTime.getMonth());
            endTime.setDate(state.startTime.getDate());

            if (isSameHour(state.startTime, endTime) && isSameMinute(state.startTime, endTime)) {
                endTime.setHours(addHours(state.startTime, 1).getHours());
                endTime.setMinutes(addMinutes(state.startTime, 30).getMinutes());
            } else if (isSameHour(state.startTime, endTime)) {
                endTime.setHours(addHours(state.startTime, 1).getHours());
            }

            state.endTime = endTime;
        }
    },
);

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateCalendarEntry(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{
                                entryId
                                    ? $t('components.calendar.EntryCreateOrUpdateModal.update.title')
                                    : $t('components.calendar.EntryCreateOrUpdateModal.create.title')
                            }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <DataPendingBlock
                        v-if="props.entryId && loading"
                        :message="$t('common.loading', [$t('common.entry', 1)])"
                    />
                    <DataErrorBlock
                        v-else-if="props.entryId && error"
                        :title="$t('common.unable_to_load', [$t('common.entry', 1)])"
                        :error="error"
                        :retry="refresh"
                    />
                    <DataNoDataBlock
                        v-if="props.entryId && (!data || !data.entry)"
                        :type="$t('common.entry', 1)"
                        icon="i-mdi-calendar"
                    />

                    <template v-else>
                        <UFormGroup name="calendar" :label="$t('common.calendar')" class="flex-1" required>
                            <ClientOnly>
                                <USelectMenu
                                    v-model="state.calendar"
                                    :disabled="!!entryId"
                                    :searchable="
                                        async () =>
                                            (
                                                await calendarStore.listCalendars({
                                                    pagination: {
                                                        offset: 0,
                                                    },
                                                    onlyPublic: false,
                                                    minAccessLevel: AccessLevel.EDIT,
                                                })
                                            ).calendars?.filter((c) => !c.closed)
                                    "
                                    searchable-lazy
                                    :searchable-placeholder="$t('common.search_field')"
                                    :search-attributes="['name']"
                                    option-attribute="name"
                                    by="id"
                                    :placeholder="$t('common.calendar')"
                                >
                                    <template #label>
                                        <template v-if="state.calendar">
                                            <span
                                                class="size-2 rounded-full"
                                                :class="`bg-${state.calendar?.color ?? 'primary'}-500 dark:bg-${state.calendar?.color ?? 'primary'}-400`"
                                            />
                                            <span class="truncate">{{ state.calendar?.name }}</span>
                                        </template>
                                    </template>

                                    <template #option="{ option }">
                                        <span
                                            class="size-2 rounded-full"
                                            :class="`bg-${option.color ?? 'primary'}-500 dark:bg-${option.color ?? 'primary'}-400`"
                                        />
                                        <span class="truncate">{{ option.name }}</span>
                                    </template>

                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.calendar')]) }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormGroup>

                        <UFormGroup name="title" :label="$t('common.title')" class="flex-1" required>
                            <UInput v-model="state.title" name="title" type="text" :placeholder="$t('common.title')" />
                        </UFormGroup>

                        <UFormGroup name="startTime" :label="$t('common.begins_at')" class="flex-1" required>
                            <DatePickerPopoverClient
                                v-model="state.startTime"
                                :popover="{ popper: { placement: 'bottom-start' } }"
                                :date-picker="{ mode: 'dateTime', is24hr: true, clearable: true }"
                            />
                        </UFormGroup>

                        <UFormGroup name="endTime" :label="$t('common.ends_at')" class="flex-1" required>
                            <DatePickerPopoverClient
                                v-model="state.endTime"
                                :popover="{ popper: { placement: 'bottom-start' } }"
                                :date-picker="{ mode: 'dateTime', is24hr: true, clearable: true }"
                            />
                        </UFormGroup>

                        <UFormGroup name="content" :label="$t('common.content')" class="flex-1" required>
                            <ClientOnly>
                                <TiptapEditor v-model="state.content" wrapper-class="min-h-80" />
                            </ClientOnly>
                        </UFormGroup>

                        <UFormGroup name="closed" :label="`${$t('common.close', 2)}?`" class="flex-1">
                            <UToggle v-model="state.closed" />
                        </UFormGroup>

                        <UFormGroup name="rsvpOpen" :label="$t('common.rsvp')" class="flex-1">
                            <UToggle v-model="state.rsvpOpen" />
                        </UFormGroup>

                        <UFormGroup name="users" :label="$t('common.guest', 2)" class="flex-1">
                            <ClientOnly>
                                <USelectMenu
                                    v-model="state.users"
                                    multiple
                                    :searchable="
                                        async (query: string) => {
                                            usersLoading = true;
                                            const users = await completorStore.completeCitizens({
                                                search: query,
                                            });
                                            usersLoading = false;
                                            return users;
                                        }
                                    "
                                    searchable-lazy
                                    :searchable-placeholder="$t('common.search_field')"
                                    :search-attributes="['firstname', 'lastname']"
                                    block
                                    :placeholder="$t('common.citizen', 2)"
                                    trailing
                                    by="userId"
                                >
                                    <template #label>
                                        {{ $t('common.selected', state.users.length) }}
                                    </template>

                                    <template #option="{ option: user }">
                                        {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                                    </template>

                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>

                                    <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                                </USelectMenu>
                            </ClientOnly>
                        </UFormGroup>

                        <div
                            v-if="state.users.length > 0"
                            class="mt-2 overflow-hidden rounded-md bg-neutral-100 dark:bg-base-900"
                        >
                            <ul role="list" class="grid grid-cols-2 text-sm font-medium text-gray-100 lg:grid-cols-3">
                                <li
                                    v-for="user in state.users"
                                    :key="user.userId"
                                    class="flex items-center border-b border-gray-100 px-4 py-2 dark:border-gray-800"
                                >
                                    <CitizenInfoPopover :user="user" show-avatar show-avatar-in-name />
                                </li>
                            </ul>
                        </div>
                    </template>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ data ? $t('common.save') : $t('common.create') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
