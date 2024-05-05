<script lang="ts" setup>
import { addHours, format } from 'date-fns';
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import DatePickerClient from '~/components/partials/DatePicker.client.vue';
import DocEditor from '~/components/partials/DocEditor.vue';
import type { CalendarShort } from '~~/gen/ts/resources/calendar/calendar';
import type { CreateOrUpdateCalendarEntryResponse } from '~~/gen/ts/services/calendar/calendar';
import { useCalendarStore } from '~/store/calendar';
import type { AccessLevel, CalendarAccess } from '~~/gen/ts/resources/calendar/access';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';

const props = defineProps<{
    calendarId?: string;
    entryId?: string;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const calendarStore = useCalendarStore();

const completorStore = useCompletorStore();

const notifications = useNotificatorStore();

const maxAccessEntries = 10;

const schema = z.object({
    calendar: z.custom<CalendarShort>().optional(),
    title: z.string().min(3).max(512),
    startTime: z.date(),
    endTime: z.date(),
    content: z.string().max(1000000),
    rsvpOpen: z.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    calendar: undefined,
    title: '',
    startTime: new Date(),
    endTime: addHours(new Date(), 1),
    content: '',
    rsvpOpen: false,
});

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `calendar-entry:${props.entryId}`,
    () => calendarStore.getCalendarEntry({ calendarId: props.calendarId!, entryId: props.entryId! }),
    {
        immediate: !!props.calendarId && !!props.entryId,
    },
);

async function createOrUpdateCalendarEntry(values: Schema): Promise<CreateOrUpdateCalendarEntryResponse> {
    if (!values.calendar) {
        throw 'No Calendar selected';
    }

    const reqAccess: CalendarAccess = {
        jobs: [],
        users: [],
    };
    access.value.forEach((entry) => {
        if (entry.values.accessRole === undefined) {
            return;
        }

        if (entry.type === 0) {
            if (!entry.values.userId) {
                return;
            }

            reqAccess.users.push({
                id: '0',
                calendarId: values.calendar!.id,
                entryId: data.value?.entry?.id ?? '0',
                userId: entry.values.userId,
                access: entry.values.accessRole,
            });
        } else if (entry.type === 1) {
            if (!entry.values.job) {
                return;
            }

            reqAccess.jobs.push({
                id: '0',
                calendarId: values.calendar!.id,
                entryId: data.value?.entry?.id ?? '0',
                job: entry.values.job,
                minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                access: entry.values.accessRole,
            });
        }
    });

    try {
        const response = await calendarStore.createOrUpdateCalendarEntry({
            id: data.value?.entry?.id ?? '0',
            calendarId: values.calendar.id,
            title: values.title,
            startTime: toTimestamp(values.startTime),
            endTime: toTimestamp(values.endTime),
            content: values.content,
            rsvpOpen: values.rsvpOpen,
            creatorJob: '',
            access: reqAccess,
        });

        isOpen.value = false;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
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
    state.content = entry.content;
    state.rsvpOpen = entry.rsvpOpen !== undefined;

    if (entry.access) {
        access.value.clear();

        let accessId = 0;
        entry.access?.users.forEach((user) => {
            const id = accessId.toString();
            access.value.set(id, {
                id,
                type: 0,
                values: { userId: user.userId, accessRole: user.access },
            });
            accessId++;
        });

        entry.access?.jobs.forEach((job) => {
            const id = accessId.toString();
            access.value.set(id, {
                id,
                type: 1,
                values: {
                    job: job.job,
                    accessRole: job.access,
                    minimumGrade: job.minimumGrade,
                },
            });
            accessId++;
        });
    }
}

watch(data, () => setFromProps());
watch(props, () => refresh());

const access = ref<
    Map<
        string,
        {
            id: string;
            type: number;
            values: {
                job?: string;
                userId?: number;
                accessRole?: AccessLevel;
                minimumGrade?: number;
            };
        }
    >
>(new Map());

function addAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.add({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            description: {
                key: 'notifications.max_access_entry.content',
                parameters: { max: maxAccessEntries.toString() },
            } as TranslateItem,
            type: 'error',
        });
        return;
    }

    const id = access.value.size > 0 ? parseInt([...access.value.keys()]?.pop() ?? '1', 10) + 1 : 0;
    access.value.set(id.toString(), {
        id: id.toString(),
        type: 1,
        values: {},
    });
}

function removeAccessEntry(event: { id: string }): void {
    access.value.delete(event.id);
}

function updateAccessEntryType(event: { id: string; type: number }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.type = event.type;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryName(event: { id: string; job?: Job; char?: UserShort }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    if (event.job) {
        accessEntry.values.job = event.job.name;
        accessEntry.values.userId = undefined;
    } else if (event.char) {
        accessEntry.values.job = undefined;
        accessEntry.values.userId = event.char.userId;
    }

    access.value.set(event.id, accessEntry);
}

function updateAccessEntryRank(event: { id: string; rank: JobGrade }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.minimumGrade = event.rank.grade;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryAccess(event: { id: string; access: AccessLevel }): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) {
        return;
    }

    accessEntry.values.accessRole = event.access;
    access.value.set(event.id, accessEntry);
}

const { data: jobs } = useAsyncData('completor-jobs', () => completorStore.listJobs());

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
                        :retry="refresh"
                    />
                    <DataNoDataBlock
                        v-if="props.entryId && (!data || !data.entry)"
                        :type="$t('common.entry', 1)"
                        icon="i-mdi-calendar"
                    />

                    <template v-else>
                        <UFormGroup name="calendar" :label="$t('common.calendar')" class="flex-1" required>
                            <USelectMenu
                                v-model="state.calendar"
                                :disabled="!entryId"
                                :searchable="
                                    async (query) =>
                                        (await calendarStore.listCalendars({ pagination: { offset: 0 } })).calendars
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

                        <UFormGroup name="access" :label="$t('common.access')" class="flex-1">
                            <CalendarAccessEntry
                                v-for="entry in access.values()"
                                :key="entry.id"
                                :init="entry"
                                :jobs="jobs"
                                @type-change="updateAccessEntryType($event)"
                                @name-change="updateAccessEntryName($event)"
                                @rank-change="updateAccessEntryRank($event)"
                                @access-change="updateAccessEntryAccess($event)"
                                @delete-request="removeAccessEntry($event)"
                            />

                            <UButton
                                :ui="{ rounded: 'rounded-full' }"
                                icon="i-mdi-plus"
                                :title="$t('components.documents.document_editor.add_permission')"
                                @click="addAccessEntry()"
                            />
                        </UFormGroup>
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
