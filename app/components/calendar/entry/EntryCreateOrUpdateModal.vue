<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { addHours, addMinutes, isSameDay, isSameHour, isSameMinute } from 'date-fns';
import type { CalendarDay } from 'v-calendar/dist/types/src/utils/page.js';
import { z } from 'zod';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import InputDatePicker from '~/components/partials/InputDatePicker.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import { useCalendarStore } from '~/stores/calendar';
import { useCompletorStore } from '~/stores/completor';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access';
import type { CalendarShort } from '~~/gen/ts/resources/calendar/calendar';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { CreateOrUpdateCalendarEntryResponse } from '~~/gen/ts/services/calendar/calendar';

const props = defineProps<{
    calendarId?: number;
    entryId?: number;

    day?: CalendarDay;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const calendarStore = useCalendarStore();

const completorStore = useCompletorStore();

const schema = z.object({
    calendar: z.custom<CalendarShort>().optional(),
    title: z.string().min(3).max(512),
    startTime: z.date(),
    endTime: z.date(),
    content: z.string().min(3).max(1000000),
    closed: z.coerce.boolean(),
    rsvpOpen: z.coerce.boolean(),
    users: z.custom<UserShort>().array().max(20).default([]),
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

const { data, status, refresh, error } = useLazyAsyncData(
    `calendar-entry:${props.entryId}`,
    () => calendarStore.getCalendarEntry({ entryId: props.entryId! }),
    {
        immediate: !!props.calendarId && !!props.entryId,
    },
);

async function createOrUpdateCalendarEntry(values: Schema): Promise<CreateOrUpdateCalendarEntryResponse> {
    if (!values.calendar) {
        throw 'No Calendar selected';
    }

    try {
        const response = await calendarStore.createOrUpdateCalendarEntry(
            {
                id: data.value?.entry?.id ?? 0,
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
            state.users.map((u) => u.userId),
        );

        emit('close', false);

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function setFromProps(): void {
    if (props.day) {
        state.startTime = addHours(props.day.date, 1);
        state.endTime = addHours(props.day.date, 2);
        return;
    }

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

setFromProps();
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

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal
        :title="
            entryId
                ? $t('components.calendar.EntryCreateOrUpdateModal.update.title')
                : $t('components.calendar.EntryCreateOrUpdateModal.create.title')
        "
    >
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <DataPendingBlock
                    v-if="props.entryId && isRequestPending(status)"
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
                    <UFormField class="flex-1" name="calendar" :label="$t('common.calendar')" required>
                        <SelectMenu
                            v-model="state.calendar"
                            label-key="name"
                            :disabled="!!entryId"
                            :searchable="
                                async () =>
                                    (
                                        (
                                            await calendarStore.listCalendars({
                                                pagination: {
                                                    offset: 0,
                                                },
                                                onlyPublic: false,
                                                minAccessLevel: AccessLevel.EDIT,
                                            })
                                        ).calendars as CalendarShort[]
                                    )?.filter((c) => !c.closed)
                            "
                            searchable-key="calendar-list"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :filter-fields="['name']"
                            :placeholder="$t('common.calendar')"
                        >
                            <template v-if="state.calendar" #default>
                                <span
                                    class="size-2 rounded-full"
                                    :class="`bg-${state.calendar?.color ?? 'primary'}-500 dark:bg-${state.calendar?.color ?? 'primary'}-400`"
                                />

                                <span class="truncate">{{ state.calendar?.name }}</span>
                            </template>

                            <template #item="{ item }">
                                <span
                                    class="size-2 rounded-full"
                                    :class="`bg-${item.color ?? 'primary'}-500 dark:bg-${item.color ?? 'primary'}-400`"
                                />

                                <span class="truncate">{{ item.name }}</span>
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.calendar')]) }}
                            </template>
                        </SelectMenu>
                    </UFormField>

                    <UFormField class="flex-1" name="title" :label="$t('common.title')" required>
                        <UInput v-model="state.title" name="title" type="text" :placeholder="$t('common.title')" />
                    </UFormField>

                    <UFormField class="flex-1" name="startTime" :label="$t('common.begins_at')" required>
                        <InputDatePicker v-model="state.startTime" date-format="dd.MM.yyyy HH:mm" clearable time />
                    </UFormField>

                    <UFormField class="flex-1" name="endTime" :label="$t('common.ends_at')" required>
                        <InputDatePicker v-model="state.endTime" date-format="dd.MM.yyyy HH:mm" clearable time />
                    </UFormField>

                    <UFormField class="flex-1" name="content" :label="$t('common.content')" required>
                        <ClientOnly>
                            <TiptapEditor v-model="state.content" wrapper-class="min-h-80" />
                        </ClientOnly>
                    </UFormField>

                    <UFormField class="flex-1" name="closed" :label="`${$t('common.close', 2)}?`">
                        <USwitch v-model="state.closed" />
                    </UFormField>

                    <UFormField class="flex-1" name="rsvpOpen" :label="$t('common.rsvp')">
                        <USwitch v-model="state.rsvpOpen" />
                    </UFormField>

                    <UFormField class="flex-1" name="users" :label="$t('common.guest', 2)">
                        <SelectMenu
                            v-model="state.users"
                            multiple
                            :searchable="
                                async (q: string) =>
                                    await completorStore.completeCitizens({
                                        search: q,
                                        userIds: state.users.map((u) => u.userId),
                                    })
                            "
                            searchable-key="completor-citizens"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :filter-fields="['firstname', 'lastname']"
                            block
                            :placeholder="$t('common.citizen', 2)"
                            trailing
                        >
                            <template #default>
                                {{ $t('common.selected', state.users.length) }}
                            </template>

                            <template #item="{ item: user }">
                                {{ `${user?.firstname} ${user?.lastname} (${user?.dateofbirth})` }}
                            </template>

                            <template #empty> {{ $t('common.not_found', [$t('common.citizen', 2)]) }} </template>
                        </SelectMenu>
                    </UFormField>
                </template>
            </UForm>

            <div v-if="state.users.length > 0" class="mt-2 overflow-hidden rounded-md bg-neutral-100 dark:bg-neutral-900">
                <ul class="grid grid-cols-2 text-sm font-medium text-toned lg:grid-cols-3" role="list">
                    <li
                        v-for="user in state.users"
                        :key="user.userId"
                        class="flex items-center border-b border-neutral-100 px-4 py-2 dark:border-neutral-800"
                    >
                        <CitizenInfoPopover :user="user" show-avatar show-avatar-in-name />
                    </li>
                </ul>
            </div>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="data ? $t('common.save') : $t('common.create')"
                    @click="formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>
