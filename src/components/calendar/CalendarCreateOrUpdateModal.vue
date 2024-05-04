<script lang="ts" setup>
import { z } from 'zod';
import type { Calendar } from '~~/gen/ts/resources/calendar/calendar';

defineProps<{
    calendar?: Calendar;
}>();

const { isOpen } = useModal();

const schema = z.object({
    name: z.string().min(3).max(255),
    description: z.string().max(512),
    public: z.boolean(),
    closed: z.boolean(),
    color: z.string().max(12),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    description: '',
    public: false,
    closed: false,
    color: 'primary',
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
