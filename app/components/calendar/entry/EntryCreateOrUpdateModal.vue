<script lang="ts" setup>
import type { ChipProps, FormSubmitEvent } from '@nuxt/ui';
import type { JSONContent } from '@tiptap/core';
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
import { Struct } from '~~/gen/ts/google/protobuf/struct';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access/access';
import type { CalendarShort } from '~~/gen/ts/resources/calendar/calendar';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { UserShort } from '~~/gen/ts/resources/users/short/user';
import type { CreateOrUpdateCalendarEntryResponse } from '~~/gen/ts/services/calendar/entries';

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
    title: z.coerce.string().min(3).max(512),
    startTime: z.date(),
    endTime: z.date().optional(),
    allDay: z.coerce.boolean().default(false),
    content: z.custom<JSONContent | string>().optional(),
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
    allDay: false,
    content: '',
    closed: false,
    rsvpOpen: true,
    users: [],
});

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(state, {
    serializer: (value) =>
        JSON.stringify({
            calendarId: value.calendar?.id ?? 0,
            title: value.title,
            startTime: value.startTime ? value.startTime.toISOString() : '',
            endTime: value.endTime ? value.endTime.toISOString() : '',
            allDay: value.allDay,
            content: value.content ? JSON.stringify(value.content) : '',
            closed: value.closed,
            rsvpOpen: value.rsvpOpen,
            users: [...value.users.map((user) => user.userId)].sort((a, b) => a - b),
        }),
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
                calendarId: values.calendar?.id,
                title: values.title,
                startTime: toTimestamp(values.startTime),
                endTime: values.allDay ? undefined : toTimestamp(values.endTime),
                content: {
                    contentType: ContentType.TIPTAP_JSON,
                    version: '',
                    tiptapJson: Struct.fromJsonString(JSON.stringify(values.content)),
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
        syncSnapshot();
        return;
    }

    if (!data.value?.entry) {
        syncSnapshot();
        return;
    }

    const entry = data.value?.entry;
    if (entry.calendar) state.calendar = entry.calendar;
    state.title = entry.title;
    state.startTime = toDate(entry.startTime);
    state.endTime = entry.endTime ? toDate(entry.endTime) : undefined;
    state.allDay = state.endTime === undefined;
    state.content = entry.content?.tiptapJson
        ? (Struct.toJson(entry.content.tiptapJson) as JSONContent)
        : (entry.content?.rawHtml ?? '');
    state.closed = entry.closed;
    state.rsvpOpen = entry.rsvpOpen !== undefined;
    syncSnapshot();
}

setFromProps();
watch(data, () => setFromProps());
watch(props, async () => refresh());

watch(
    () => state.startTime,
    () => {
        if (!state.endTime) return;

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

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateCalendarEntry(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}
</script>

<template>
    <UModal
        :title="
            entryId
                ? $t('components.calendar.EntryCreateOrUpdateModal.update.title')
                : $t('components.calendar.EntryCreateOrUpdateModal.create.title')
        "
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{
                        entryId
                            ? $t('components.calendar.EntryCreateOrUpdateModal.update.title')
                            : $t('components.calendar.EntryCreateOrUpdateModal.create.title')
                    }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" class="flex flex-col gap-2" :schema="schema" :state="state" @submit="onSubmitThrottle">
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
                            class="w-full"
                            label-key="name"
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
                                            calendarIds: [state.calendar?.id].filter((id): id is number => !!id),
                                        })
                                    ).calendars
                                        .filter((c) => !c.closed)
                                        .map((c) => ({
                                            id: c.id,
                                            name: c.name,
                                            color: c.color,
                                        }))
                            "
                            searchable-key="calendar-list"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :filter-fields="['name']"
                            :placeholder="$t('common.calendar')"
                        >
                            <template #leading="{ modelValue, ui }">
                                <UChip
                                    v-if="modelValue"
                                    :class="ui.itemLeadingChip()"
                                    :color="(modelValue?.color ?? 'primary') as ChipProps['color']"
                                    inset
                                    standalone
                                    :size="ui.itemLeadingChipSize() as ChipProps['size']"
                                />
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.calendar')]) }}
                            </template>
                        </SelectMenu>
                    </UFormField>

                    <UFormField class="flex-1" name="title" :label="$t('common.title')" required>
                        <UInput
                            v-model="state.title"
                            class="w-full"
                            name="title"
                            type="text"
                            :placeholder="$t('common.title')"
                        />
                    </UFormField>

                    <UFormField class="flex-1">
                        <UTabs
                            :items="[
                                { label: $t('common.time_range'), icon: 'i-mdi-calendar-today' },
                                { label: $t('common.all_day'), icon: 'i-mdi-calendar-today' },
                            ]"
                            @update:model-value="($event) => (state.allDay = $event === '1')"
                        >
                            <template #content>
                                <UFormField class="flex-1" name="startTime" :label="$t('common.begins_at')" required>
                                    <InputDatePicker
                                        v-model="state.startTime"
                                        class="w-full"
                                        :date-format="!state.allDay ? undefined : 'date'"
                                        :time="!state.allDay"
                                    />
                                </UFormField>

                                <UFormField class="flex-1" name="endTime" :label="$t('common.ends_at')" required>
                                    <InputDatePicker
                                        v-model="state.endTime"
                                        class="w-full"
                                        :date-format="!state.allDay ? undefined : 'date'"
                                        :time="!state.allDay"
                                    />
                                </UFormField>
                            </template>
                        </UTabs>
                    </UFormField>

                    <UFormField class="flex-1" name="content" :label="$t('common.content')" required :ui="{ error: 'hidden' }">
                        <ClientOnly>
                            <TiptapEditor
                                v-model="state.content"
                                class="w-full"
                                name="content"
                                wrapper-class="min-h-80"
                                :limit="10_000"
                            />
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
                            class="w-full"
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

                            <template #item-label="{ item: user }">
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
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="data ? $t('common.save') : $t('common.create')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
