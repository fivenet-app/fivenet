<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ColorPickerTW from '~/components/partials/ColorPickerTW.vue';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useCalendarStore } from '~/stores/calendar';
import { isSystemManagedCalendar } from '~/components/calendar/helpers';
import { AccessLevel } from '~~/gen/ts/resources/calendar/access/access';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { CreateCalendarResponse, UpdateCalendarResponse } from '~~/gen/ts/services/calendar/calendar';
import TiptapEditor from '../../partials/editor/TiptapEditor.vue';

const props = defineProps<{
    calendarId?: number;
    systemManaged?: boolean;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { t } = useI18n();

const { attr, activeChar } = useAuth();

const calendarStore = useCalendarStore();
const { hasPrivateCalendar } = storeToRefs(calendarStore);

const notifications = useNotificationsStore();

const { maxAccessEntries } = useAppConfig();

const canDo = computed(() => ({
    privateCalendar: attr('calendar.CalendarService/CreateCalendar', 'Fields', 'Job').value,
    publicCalendar: attr('calendar.CalendarService/CreateCalendar', 'Fields', 'Public').value,
}));

const isSystemManaged = computed(() => props.systemManaged || isSystemManagedCalendar(data.value?.calendar));

const schema = z.object({
    name: z.coerce.string().min(3).max(255),
    description: z.string().optional(),
    private: z.coerce.boolean(),
    public: z.coerce.boolean(),
    closed: z.coerce.boolean(),
    color: z.coerce.string().max(12),
    access: z.object({
        jobs: jobsAccessEntries(t).max(maxAccessEntries).default([]),
        users: userAccessEntries(t).max(maxAccessEntries).default([]),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    description: '',
    private: !hasPrivateCalendar,
    public: false,
    closed: false,
    color: 'blue',
    access: {
        jobs: [],
        users: [],
    },
});

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(state);

const {
    data: data,
    status,
    refresh,
    error,
} = useLazyAsyncData(
    `calendar-calendar:${props.calendarId}`,
    () => calendarStore.getCalendar({ calendarId: props.calendarId! }),
    {
        immediate: !!props.calendarId,
    },
);

async function createOrUpdateCalendar(values: Schema): Promise<CreateCalendarResponse | UpdateCalendarResponse> {
    values.access.users.forEach((user) => {
        if (user.id < 0) user.id = 0;
        user.user = undefined; // Clear user object to avoid sending unnecessary data
    });
    values.access.jobs.forEach((job) => job.id < 0 && (job.id = 0));

    try {
        const response = await calendarStore.createOrUpdateCalendar({
            id: data.value?.calendar?.id ?? 0,
            job: isSystemManaged.value
                ? (data.value?.calendar?.job ?? activeChar.value?.job)
                : values.private
                  ? undefined
                  : activeChar.value?.job,
            name: values.name,
            description: values.description,
            public: values.public,
            closed: values.closed,
            color: values.color,
            access: values.access,
            creatorJob: '',
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('close', false);

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function setFromProps(): void {
    if (!data.value?.calendar) return;

    const calendar = data.value?.calendar;
    state.name = calendar.name;
    state.description = calendar.description;
    state.private = calendar.job === undefined;
    state.public = calendar.public;
    state.closed = calendar.closed;
    state.color = calendar.color ?? 'primary';
    if (calendar.access) {
        state.access = calendar.access;
    }

    syncSnapshot();
}

watch(data, () => setFromProps());
watch(props, async () => refresh());

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateCalendar(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
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
            calendarId
                ? $t('components.calendar.calendar.CreateOrUpdateModal.update.title')
                : $t('components.calendar.calendar.CreateOrUpdateModal.create.title')
        "
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-2">
                <h3 class="font-semibold text-highlighted">
                    {{
                        calendarId
                            ? $t('components.calendar.calendar.CreateOrUpdateModal.update.title')
                            : $t('components.calendar.calendar.CreateOrUpdateModal.create.title')
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
                    v-if="props.calendarId && isRequestPending(status)"
                    :message="$t('common.loading', [$t('common.calendar')])"
                />
                <DataErrorBlock
                    v-else-if="props.calendarId && error"
                    :title="$t('common.unable_to_load', [$t('common.calendar')])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-if="props.calendarId && (!data || !data.calendar)"
                    :type="$t('common.calendar')"
                    icon="i-mdi-calendar"
                />

                <template v-else>
                    <p v-if="isSystemManaged" class="text-sm text-neutral-500 dark:text-neutral-400">
                        {{ $t('common.read_only') }}
                    </p>

                    <UFormField class="flex-1" name="color" :label="$t('common.color')">
                        <ColorPickerTW v-model="state.color" class="w-full" />
                    </UFormField>

                    <template v-if="!isSystemManaged">
                        <UFormField class="flex-1" name="description" :label="$t('common.description')">
                            <TiptapEditor
                                v-model="state.description"
                                class="w-full"
                                name="content"
                                wrapper-class="min-h-80"
                                :placeholder="$t('common.description')"
                                :limit="1_000"
                            />
                        </UFormField>

                        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
                            <UFormField
                                class="flex-1"
                                name="private"
                                :label="$t('components.calendar.calendar.CreateOrUpdateModal.private')"
                            >
                                <USwitch
                                    v-model="state.private"
                                    :disabled="
                                        !canDo.privateCalendar ||
                                        calendarId !== undefined ||
                                        (!props.calendarId && hasPrivateCalendar)
                                    "
                                />
                            </UFormField>

                            <UFormField v-if="canDo.publicCalendar" class="flex-1" name="public" :label="$t('common.public')">
                                <USwitch v-model="state.public" />
                            </UFormField>

                            <UFormField class="flex-1" name="closed" :label="`${$t('common.close', 2)}?`">
                                <USwitch v-model="state.closed" />
                            </UFormField>
                        </div>

                        <UFormField class="flex-1" name="access" :label="$t('common.access')">
                            <AccessManager
                                v-model:jobs="state.access.jobs"
                                v-model:users="state.access.users"
                                :target-id="calendarId ?? 0"
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.calendar.AccessLevel')"
                            />
                        </UFormField>
                    </template>
                </template>
            </UForm>
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
