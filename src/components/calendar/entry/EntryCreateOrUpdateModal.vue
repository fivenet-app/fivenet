<script lang="ts" setup>
import { addHours, format } from 'date-fns';
import { z } from 'zod';
import DatePickerClient from '~/components/partials/DatePicker.client.vue';
import DocEditor from '~/components/partials/DocEditor.vue';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';

defineProps<{
    entry?: CalendarEntry;
}>();

const { isOpen } = useModal();

const schema = z.object({
    title: z.string(),
    startTime: z.date(),
    endTime: z.date(),
    content: z.string(),
    public: z.boolean(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    title: '',
    startTime: new Date(),
    endTime: addHours(new Date(), 1),
    content: '',
    public: false,
});

const canSubmit = ref(true);

// TODO
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state">
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
                    <UFormGroup name="name" :label="$t('common.name')" class="flex-1">
                        <UInput
                            v-model="state.title"
                            name="name"
                            type="text"
                            :placeholder="$t('common.name')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </UFormGroup>

                    <UFormGroup name="startTime" :label="$t('common.begins_at')" class="flex-1">
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

                    <UFormGroup name="endTime" :label="$t('common.ends_at')" class="flex-1">
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

                    <UFormGroup name="content" :label="$t('common.content')" class="flex-1">
                        <ClientOnly>
                            <DocEditor v-model="state.content" :min-height="250" />
                        </ClientOnly>
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
