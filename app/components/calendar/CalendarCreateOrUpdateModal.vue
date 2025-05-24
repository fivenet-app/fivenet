<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ColorPickerTW from '~/components/partials/ColorPickerTW.vue';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useCalendarStore } from '~/stores/calendar';
import { useNotificatorStore } from '~/stores/notificator';
import { AccessLevel, type CalendarJobAccess, type CalendarUserAccess } from '~~/gen/ts/resources/calendar/access';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { CreateCalendarResponse, UpdateCalendarResponse } from '~~/gen/ts/services/calendar/calendar';

const props = defineProps<{
    calendarId?: number;
}>();

const { isOpen } = useModal();

const { attr, activeChar } = useAuth();

const calendarStore = useCalendarStore();
const { hasPrivateCalendar } = storeToRefs(calendarStore);

const notifications = useNotificatorStore();

const { maxAccessEntries } = useAppConfig();

const canDo = computed(() => ({
    privateCalendar: attr('calendar.CalendarService.CreateCalendar', 'Fields', 'Job').value,
    publicCalendar: attr('calendar.CalendarService.CreateCalendar', 'Fields', 'Public').value,
}));

const schema = z.object({
    name: z.string().min(3).max(255),
    description: z.string().max(512).optional(),
    private: z.boolean(),
    public: z.boolean(),
    closed: z.boolean(),
    color: z.string().max(12),
    access: z.object({
        jobs: z.custom<CalendarJobAccess>().array().max(maxAccessEntries),
        users: z.custom<CalendarUserAccess>().array().max(maxAccessEntries),
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

async function createOrUpdateCalendar(values: Schema): Promise<CreateCalendarResponse | UpdateCalendarResponse> {
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
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        isOpen.value = false;

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

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
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
                        :error="error"
                        :retry="refresh"
                    />
                    <DataNoDataBlock
                        v-if="props.calendarId && (!data || !data.calendar)"
                        :type="$t('common.calendar')"
                        icon="i-mdi-calendar"
                    />

                    <template v-else>
                        <UFormGroup class="flex-1" name="title" :label="$t('common.name')" required>
                            <UInput v-model="state.name" name="name" type="text" :placeholder="$t('common.name')" />
                        </UFormGroup>

                        <UFormGroup class="flex-1" name="color" :label="$t('common.color')">
                            <ColorPickerTW v-model="state.color" />
                        </UFormGroup>

                        <UFormGroup class="flex-1" name="description" :label="$t('common.description')">
                            <UTextarea v-model="state.description" name="description" :placeholder="$t('common.description')" />
                        </UFormGroup>

                        <UFormGroup
                            class="flex-1"
                            name="private"
                            :label="$t('components.calendar.CalendarCreateOrUpdateModal.private')"
                        >
                            <UToggle
                                v-model="state.private"
                                :disabled="
                                    !canDo.privateCalendar ||
                                    calendarId !== undefined ||
                                    (!props.calendarId && hasPrivateCalendar)
                                "
                            />
                        </UFormGroup>

                        <UFormGroup v-if="canDo.publicCalendar" class="flex-1" name="public" :label="$t('common.public')">
                            <UToggle v-model="state.public" />
                        </UFormGroup>

                        <UFormGroup class="flex-1" name="closed" :label="`${$t('common.close', 2)}?`">
                            <UToggle v-model="state.closed" />
                        </UFormGroup>

                        <UFormGroup class="flex-1" name="access" :label="$t('common.access')">
                            <AccessManager
                                v-model:jobs="state.access.jobs"
                                :target-id="calendarId ?? 0"
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.calendar.AccessLevel')"
                            />
                        </UFormGroup>
                    </template>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="black" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ data ? $t('common.save') : $t('common.create') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
