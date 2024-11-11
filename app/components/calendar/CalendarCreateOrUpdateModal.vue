<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { useCalendarStore } from '~/store/calendar';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel, type CalendarJobAccess, type CalendarUserAccess } from '~~/gen/ts/resources/calendar/access';
import type { CreateOrUpdateCalendarResponse } from '~~/gen/ts/services/calendar/calendar';
import ColorPickerTW from '../partials/ColorPickerTW.vue';
import AccessManager from '../partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '../partials/access/helpers';

const props = defineProps<{
    calendarId?: string;
}>();

const { isOpen } = useModal();

const { attr, activeChar } = useAuth();

const calendarStore = useCalendarStore();

const completorStore = useCompletorStore();

const notifications = useNotificatorStore();

const { maxAccessEntries } = useAppConfig();

const canDo = computed(() => ({
    privatecalendar: attr('CalendarService.CreateOrUpdateCalendar', 'Fields', 'Job').value,
    publicCalendar: attr('CalendarService.CreateOrUpdateCalendar', 'Fields', 'Public').value,
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
    private: true,
    public: false,
    closed: false,
    color: 'primary',
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

async function createOrUpdateCalendar(values: Schema): Promise<CreateOrUpdateCalendarResponse> {
    try {
        const response = await calendarStore.createOrUpdateCalendar({
            id: data.value?.calendar?.id ?? '0',
            name: values.name,
            job: values.private ? undefined : activeChar.value?.job,
            public: values.public,
            closed: values.closed,
            color: values.color,
            creatorJob: '',
            access: values.access,
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
                            <UInput v-model="state.name" name="name" type="text" :placeholder="$t('common.name')" />
                        </UFormGroup>

                        <UFormGroup name="color" :label="$t('common.color')" class="flex-1">
                            <ColorPickerTW v-model="state.color" />
                        </UFormGroup>

                        <UFormGroup name="description" :label="$t('common.description')" class="flex-1">
                            <UTextarea v-model="state.description" name="description" :placeholder="$t('common.description')" />
                        </UFormGroup>

                        <UFormGroup
                            name="private"
                            :label="$t('components.calendar.CalendarCreateOrUpdateModal.private')"
                            class="flex-1"
                        >
                            <UToggle v-model="state.private" :disabled="!canDo.privatecalendar || calendarId !== undefined" />
                        </UFormGroup>

                        <UFormGroup v-if="canDo.publicCalendar" name="public" :label="$t('common.public')" class="flex-1">
                            <UToggle v-model="state.public" />
                        </UFormGroup>

                        <UFormGroup name="closed" :label="`${$t('common.close', 2)}?`" class="flex-1">
                            <UToggle v-model="state.closed" />
                        </UFormGroup>

                        <UFormGroup name="access" :label="$t('common.access')" class="flex-1">
                            <AccessManager
                                v-model:jobs="state.access.jobs"
                                :target-id="calendarId ?? '0'"
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.calendar.AccessLevel')"
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
