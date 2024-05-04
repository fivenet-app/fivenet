<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import type { CreateOrUpdateCalendarResponse } from '~~/gen/ts/services/calendar/calendar';
import { useCalendarStore } from '~/store/calendar';
import { primaryColors } from '~/components/auth/account/settings';
import { useAuthStore } from '~/store/auth';

const props = defineProps<{
    calendarId?: string;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const calendarStore = useCalendarStore();

const schema = z.object({
    name: z.string().min(3).max(255),
    description: z.string().max(512).optional(),
    private: z.boolean(),
    public: z.boolean(),
    closed: z.boolean(),
    color: z.string().max(12),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    description: '',
    private: false,
    public: false,
    closed: false,
    color: 'primary',
});

const {
    data: calendar,
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
// TODO show data loading blocks and error

async function createOrUpdateCalendar(values: Schema): Promise<CreateOrUpdateCalendarResponse> {
    try {
        const response = await calendarStore.createOrUpdateCalendar({
            id: '0',
            name: values.name,
            job: values.private ? undefined : activeChar.value?.job,
            public: values.public,
            closed: values.closed,
            color: values.color,
            creatorJob: '',
        });

        isOpen.value = false;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const availableColorOptions = primaryColors.map((color) => ({
    label: color,
    chip: color,
}));

function setFromProps(): void {
    if (!calendar.value?.calendar) {
        return;
    }

    state.name = calendar.value?.calendar?.name;
    state.description = calendar.value?.calendar?.description;
    state.private = calendar.value?.calendar?.job === undefined;
    state.public = calendar.value?.calendar?.public;
    state.closed = calendar.value?.calendar?.closed;
    state.color = calendar.value?.calendar?.color ?? 'primary';
}

watch(calendar, () => setFromProps());
watch(props, () => refresh());

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
                                calendar
                                    ? $t('components.calendar.CalendarCreateOrUpdateModal.update.title')
                                    : $t('components.calendar.CalendarCreateOrUpdateModal.create.title')
                            }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
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
                                <span class="size-2 rounded-full" :class="`bg-${state.color}-500 dark:bg-${state.color}-400`" />
                                <span class="truncate">{{ state.color }}</span>
                            </template>

                            <template #option="{ option }">
                                <span class="size-2 rounded-full" :class="`bg-${option.chip}-500 dark:bg-${option.chip}-400`" />
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
                        v-if="attr('CalendarService.CreateOrUpdateCalendar', 'Fields', 'Job')"
                        name="private"
                        :label="$t('components.calendar.CalendarCreateOrUpdateModal.private')"
                        class="flex-1"
                        required
                    >
                        <UToggle v-model="state.private" :disabled="calendar !== undefined" />
                    </UFormGroup>

                    <UFormGroup name="access" :label="$t('common.access')" class="flex-1">
                        <!-- TODO -->
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ calendar ? $t('common.save') : $t('common.create') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
