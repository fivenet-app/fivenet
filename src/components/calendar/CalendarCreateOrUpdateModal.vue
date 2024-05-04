<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import type { Calendar } from '~~/gen/ts/resources/calendar/calendar';
import type { CreateOrUpdateCalendarResponse } from '~~/gen/ts/services/calendar/calendar';
import { useCalendarStore } from '~/store/calendar';

const props = defineProps<{
    calendar?: Calendar;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

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

async function createOrUpdateCalendar(values: Schema): Promise<CreateOrUpdateCalendarResponse> {
    try {
        const response = await calendarStore.createOrUpdateCalendar({
            id: '0',
            name: values.name,
            public: values.public,
            closed: values.closed,
            creatorJob: '',
        });

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

function setFromProps(): void {
    if (!props.calendar) {
        return;
    }

    state.name = props.calendar.name;
    state.description = props.calendar.description;
    state.private = props.calendar.job === undefined;
    state.public = props.calendar.public;
    state.closed = props.calendar.closed;
    state.color = props.calendar.color ?? 'primary';
}

setFromProps();

watch(props, () => setFromProps());

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
                        <UInput
                            v-model="state.color"
                            name="color"
                            type="text"
                            :placeholder="$t('common.color')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
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
