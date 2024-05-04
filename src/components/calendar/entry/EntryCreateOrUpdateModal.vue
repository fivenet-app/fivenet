<script lang="ts" setup>
import { addHours, format } from 'date-fns';
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import DatePickerClient from '~/components/partials/DatePicker.client.vue';
import DocEditor from '~/components/partials/DocEditor.vue';
import type { Calendar, CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';
import type { CreateOrUpdateCalendarEntryResponse } from '~~/gen/ts/services/calendar/calendar';
import { useCalendarStore } from '~/store/calendar';

const props = defineProps<{
    calendar?: Calendar;
    entry?: CalendarEntry;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const calendarStore = useCalendarStore();

const schema = z.object({
    calendar: z.custom<Calendar>().optional(),
    title: z.string().min(3).max(512),
    startTime: z.date(),
    endTime: z.date(),
    content: z.string().max(10240),
    public: z.boolean(),
    rsvpOpen: z.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    calendar: props.calendar,
    title: '',
    startTime: new Date(),
    endTime: addHours(new Date(), 1),
    content: '',
    rsvpOpen: false,
    public: false,
});

async function createOrUpdateCalendarEntry(values: Schema): Promise<CreateOrUpdateCalendarEntryResponse> {
    if (!values.calendar) {
        throw 'No Calendar selected';
    }

    try {
        const response = await calendarStore.createOrUpdateCalendarEntry({
            id: '0',
            calendarId: values.calendar.id,
            title: values.title,
            startTime: toTimestamp(values.startTime),
            endTime: toTimestamp(values.endTime),
            content: values.content,
            public: values.public,
            rsvpOpen: values.rsvpOpen,
            creatorJob: '',
        });

        isOpen.value = false;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function setFromProps(): void {
    if (props.calendar) {
        state.calendar = props.calendar;
    }

    if (!props.entry) {
        return;
    }

    state.title = props.entry.title;
    state.startTime = toDate(props.entry.startTime);
    state.endTime = toDate(props.entry.endTime);
    state.content = props.entry.content;
    state.public = props.entry.public;
}

setFromProps();

watch(props, () => setFromProps());

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
                                entry
                                    ? $t('components.calendar.EntryCreateOrUpdateModal.update.title')
                                    : $t('components.calendar.EntryCreateOrUpdateModal.create.title')
                            }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup v-if="!calendar" name="calendar" :label="$t('common.calendar')" class="flex-1" required>
                        <USelectMenu
                            v-model="state.calendar"
                            :searchable="
                                async (query) => (await calendarStore.listCalendars({ pagination: { offset: 0 } })).calendars
                            "
                            :search-attributes="['name']"
                            option-attribute="name"
                            by="id"
                            :placeholder="$t('common.calendar')"
                            :searchable-placeholder="$t('common.search_field')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
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
                                {{ $t('common.not_found', [$t('common.calendar', 1)]) }}
                            </template>
                        </USelectMenu>
                    </UFormGroup>

                    <UFormGroup name="title" :label="$t('common.title')" class="flex-1" required>
                        <UInput
                            v-model="state.title"
                            name="title"
                            type="text"
                            :placeholder="$t('common.title')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </UFormGroup>

                    <UFormGroup name="startTime" :label="$t('common.begins_at')" class="flex-1" required>
                        <UPopover :popper="{ placement: 'bottom-start' }">
                            <UButton
                                variant="outline"
                                color="gray"
                                block
                                icon="i-mdi-calendar-month"
                                :label="state.startTime ? format(state.startTime, 'dd.MM.yyyy HH:mm') : 'dd.mm.yyyy HH:mm'"
                            />

                            <template #panel="{ close }">
                                <DatePickerClient v-model="state.startTime" mode="dateTime" is24hr @close="close" />
                            </template>
                        </UPopover>
                    </UFormGroup>

                    <UFormGroup name="endTime" :label="$t('common.ends_at')" class="flex-1" required>
                        <UPopover :popper="{ placement: 'bottom-start' }">
                            <UButton
                                variant="outline"
                                color="gray"
                                block
                                icon="i-mdi-calendar-month"
                                :label="state.endTime ? format(state.endTime, 'dd.MM.yyyy HH:mm') : 'dd.mm.yyyy HH:mm'"
                            />

                            <template #panel="{ close }">
                                <DatePickerClient v-model="state.endTime" mode="dateTime" is24hr @close="close" />
                            </template>
                        </UPopover>
                    </UFormGroup>

                    <UFormGroup name="content" :label="$t('common.content')" class="flex-1" required>
                        <ClientOnly>
                            <DocEditor v-model="state.content" :min-height="250" />
                        </ClientOnly>
                    </UFormGroup>

                    <UFormGroup name="rsvpOpen" :label="$t('common.rsvp')" class="flex-1" required>
                        <UToggle v-model="state.rsvpOpen" />
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ entry ? $t('common.save') : $t('common.create') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
