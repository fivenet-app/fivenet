<script lang="ts" setup>
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';
import { format } from 'date-fns';
import type { DatePicker as VCalendarDatePicker } from 'v-calendar';
import DatePickerClient from './DatePicker.client.vue';

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        modelValue?: Date | undefined;
        popover?: object;
        button?: object;
        datePicker?: (typeof VCalendarDatePicker)['$props'];
        dateFormat?: string;
    }>(),
    {
        modelValue: undefined,
        popover: undefined,
        button: undefined,
        datePicker: undefined,
        dateFormat: 'dd.MM.yyyy',
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', modelValue: Date | undefined): void;
}>();

const date = useVModel(props, 'modelValue', emit);

const breakpoints = useBreakpoints(breakpointsTailwind);

const smallerBreakpoint = breakpoints.smaller('sm');

const open = ref(false);
</script>

<template>
    <ClientOnly>
        <template v-if="smallerBreakpoint">
            <UButton
                variant="outline"
                color="black"
                block
                icon="i-mdi-calendar-month"
                :label="modelValue ? format(modelValue, dateFormat) : dateFormat"
                v-bind="button"
                @click="open = true"
                @touchstart="open = true"
            />

            <UModal v-model="open">
                <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                    <template #header>
                        <div class="flex items-center justify-between">
                            <h3 class="text-2xl font-semibold leading-6">
                                {{ $t('common.date') }}
                            </h3>

                            <UButton
                                class="-my-1"
                                color="gray"
                                variant="ghost"
                                icon="i-mdi-window-close"
                                @click="open = false"
                            />
                        </div>
                    </template>

                    <div class="flex flex-1 items-center">
                        <DatePickerClient v-model="date" class="mx-auto" v-bind="datePicker" @close="open = false" />
                    </div>

                    <template #footer>
                        <UButton class="flex-1" color="black" block @click="open = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </template>
                </UCard>
            </UModal>
        </template>

        <UPopover v-else v-model:open="open" v-bind="popover">
            <UButton
                variant="outline"
                color="black"
                block
                icon="i-mdi-calendar-month"
                :label="modelValue ? format(modelValue, dateFormat) : dateFormat"
                v-bind="button"
                @touchstart="open = true"
            />

            <template #panel="{ close }">
                <DatePickerClient v-model="date" v-bind="datePicker" @close="close" />
            </template>
        </UPopover>
    </ClientOnly>
</template>
