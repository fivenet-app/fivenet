<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { primaryColors } from '~/components/auth/account/settings';
import { useAuthStore } from '~/store/auth';
import { useCalendarStore } from '~/store/calendar';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import type { AccessLevel, CalendarAccess } from '~~/gen/ts/resources/calendar/access';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Job, JobGrade } from '~~/gen/ts/resources/users/jobs';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type { CreateOrUpdateCalendarResponse } from '~~/gen/ts/services/calendar/calendar';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';

const props = defineProps<{
    calendarId?: string;
}>();

const { isOpen } = useModal();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const calendarStore = useCalendarStore();

const completorStore = useCompletorStore();

const notifications = useNotificatorStore();

const maxAccessEntries = 10;

const canCreateNonPrivateCalendar = attr('CalendarService.CreateOrUpdateCalendar', 'Fields', 'Job').value;

const schema = z.object({
    name: z.string().min(3).max(255),
    description: z.string().max(512).optional(),
    private: z.boolean(),
    public: z.boolean(),
    closed: z.boolean(),
    color: z.string().max(12),
    access: z.custom<CalendarAccess>().optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    description: '',
    private: true,
    public: false,
    closed: false,
    color: 'primary',
});

const {
    data: data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(
    `calendar-calendar:${props.calendarId}`,
    () => calendarStore.getCalendar({ calendarId: props.calendarId! }),
    {
        immediate: !!props.calendarId,
    },
);

async function createOrUpdateCalendar(values: Schema): Promise<CreateOrUpdateCalendarResponse> {
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
                calendarId: data.value?.calendar?.id ?? '0',
                userId: entry.values.userId,
                access: entry.values.accessRole,
            });
        } else if (entry.type === 1) {
            if (!entry.values.job) {
                return;
            }

            reqAccess.jobs.push({
                id: '0',
                calendarId: data.value?.calendar?.id ?? '0',
                job: entry.values.job,
                minimumGrade: entry.values.minimumGrade ? entry.values.minimumGrade : 0,
                access: entry.values.accessRole,
            });
        }
    });

    try {
        const response = await calendarStore.createOrUpdateCalendar({
            id: data.value?.calendar?.id ?? '0',
            name: values.name,
            job: values.private ? undefined : activeChar.value?.job,
            public: values.public,
            closed: values.closed,
            color: values.color,
            creatorJob: '',
            access: reqAccess,
        });

        isOpen.value = false;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const availableColorOptions = primaryColors.map((color) => ({
    label: color,
    chip: color,
}));

function setFromProps(): void {
    if (!data.value?.calendar) {
        return;
    }

    const calendar = data.value?.calendar;
    state.name = calendar.name;
    state.description = calendar.description;
    state.private = calendar.job === undefined;
    state.public = calendar.public;
    state.closed = calendar.closed;
    state.color = calendar.color ?? 'primary';

    if (calendar.access) {
        access.value.clear();

        let accessId = 0;
        calendar.access?.users.forEach((user) => {
            const id = accessId.toString();
            access.value.set(id, {
                id,
                type: 0,
                values: { userId: user.userId, accessRole: user.access },
            });
            accessId++;
        });

        calendar.access?.jobs.forEach((job) => {
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
watch(props, async () => refresh());

const access = ref(
    new Map<
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
    >(),
);

function addAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.add({
            title: { key: 'notifications.max_access_entry.title', parameters: {} },
            description: {
                key: 'notifications.max_access_entry.content',
                parameters: { max: maxAccessEntries.toString() },
            },
            type: NotificationType.ERROR,
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
    await createOrUpdateCalendar(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
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
                                calendarId
                                    ? $t('components.calendar.CalendarCreateOrUpdateModal.update.title')
                                    : $t('components.calendar.CalendarCreateOrUpdateModal.create.title')
                            }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <DataPendingBlock
                        v-if="props.calendarId && loading"
                        :message="$t('common.loading', [$t('common.calendar')])"
                    />
                    <DataErrorBlock
                        v-else-if="props.calendarId && error"
                        :title="$t('common.unable_to_load', [$t('common.calendar')])"
                        :retry="refresh"
                    />
                    <DataNoDataBlock
                        v-if="props.calendarId && (!data || !data.calendar)"
                        :type="$t('common.calendar')"
                        icon="i-mdi-calendar"
                    />

                    <template v-else>
                        <UFormGroup name="title" :label="$t('common.name')" class="flex-1" required>
                            <UInput
                                v-model="state.name"
                                name="name"
                                type="text"
                                :placeholder="$t('common.name')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup name="color" :label="$t('common.color')" class="flex-1">
                            <USelectMenu
                                v-model="state.color"
                                name="color"
                                :options="availableColorOptions"
                                option-attribute="label"
                                value-attribute="chip"
                                :searchable-placeholder="$t('common.search_field')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            >
                                <template #label>
                                    <span
                                        class="size-2 rounded-full"
                                        :class="`bg-${state.color}-500 dark:bg-${state.color}-400`"
                                    />
                                    <span class="truncate">{{ state.color }}</span>
                                </template>

                                <template #option="{ option }">
                                    <span
                                        class="size-2 rounded-full"
                                        :class="`bg-${option.chip}-500 dark:bg-${option.chip}-400`"
                                    />
                                    <span class="truncate">{{ option.label }}</span>
                                </template>
                            </USelectMenu>
                        </UFormGroup>

                        <UFormGroup name="description" :label="$t('common.description')" class="flex-1">
                            <UTextarea
                                v-model="state.description"
                                name="description"
                                :placeholder="$t('common.description')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup
                            name="private"
                            :label="$t('components.calendar.CalendarCreateOrUpdateModal.private')"
                            class="flex-1"
                        >
                            <UToggle
                                v-model="state.private"
                                :disabled="!canCreateNonPrivateCalendar || calendarId !== undefined"
                            />
                        </UFormGroup>

                        <UFormGroup
                            v-if="attr('CalendarService.CreateOrUpdateCalendar', 'Fields', 'Public').value"
                            name="public"
                            :label="$t('common.public')"
                            class="flex-1"
                        >
                            <UToggle v-model="state.public" />
                        </UFormGroup>

                        <UFormGroup name="closed" :label="`${$t('common.close', 2)}?`" class="flex-1">
                            <UToggle v-model="state.closed" />
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
