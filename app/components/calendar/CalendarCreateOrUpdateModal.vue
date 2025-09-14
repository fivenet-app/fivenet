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
import { AccessLevel } from '~~/gen/ts/resources/calendar/access';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { CreateCalendarResponse, UpdateCalendarResponse } from '~~/gen/ts/services/calendar/calendar';
import { jobAccessEntry, userAccessEntry } from '~~/shared/types/validation';

const props = defineProps<{
    calendarId?: number;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { attr, activeChar } = useAuth();

const calendarStore = useCalendarStore();
const { hasPrivateCalendar } = storeToRefs(calendarStore);

const notifications = useNotificationsStore();

const { maxAccessEntries } = useAppConfig();

const canDo = computed(() => ({
    privateCalendar: attr('calendar.CalendarService/CreateCalendar', 'Fields', 'Job').value,
    publicCalendar: attr('calendar.CalendarService/CreateCalendar', 'Fields', 'Public').value,
}));

const schema = z.object({
    name: z.coerce.string().min(3).max(255),
    description: z.coerce.string().max(512).optional(),
    private: z.coerce.boolean(),
    public: z.coerce.boolean(),
    closed: z.coerce.boolean(),
    color: z.coerce.string().max(12),
    access: z.object({
        jobs: jobAccessEntry.array().max(maxAccessEntries).default([]),
        users: userAccessEntry.array().max(maxAccessEntries).default([]),
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
            name: values.name,
            job: values.private ? undefined : activeChar.value?.job,
            public: values.public,
            closed: values.closed,
            color: values.color,
            creatorJob: '',
            access: values.access,
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
        state.access = calendar.access;
    }
}

watch(data, () => setFromProps());
watch(props, async () => refresh());

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateCalendar(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal
        :title="
            calendarId
                ? $t('components.calendar.CalendarCreateOrUpdateModal.update.title')
                : $t('components.calendar.CalendarCreateOrUpdateModal.create.title')
        "
    >
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" class="flex flex-col gap-2" @submit="onSubmitThrottle">
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
                    <UFormField class="flex-1" name="title" :label="$t('common.name')" required>
                        <UInput v-model="state.name" name="name" type="text" :placeholder="$t('common.name')" class="w-full" />
                    </UFormField>

                    <UFormField class="flex-1" name="color" :label="$t('common.color')">
                        <ColorPickerTW v-model="state.color" class="w-full" />
                    </UFormField>

                    <UFormField class="flex-1" name="description" :label="$t('common.description')">
                        <UTextarea
                            v-model="state.description"
                            name="description"
                            :placeholder="$t('common.description')"
                            class="w-full"
                        />
                    </UFormField>

                    <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
                        <UFormField
                            class="flex-1"
                            name="private"
                            :label="$t('components.calendar.CalendarCreateOrUpdateModal.private')"
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
            </UForm>
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
