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
                color="neutral"
                block
                icon="i-mdi-calendar-month"
                :label="modelValue ? format(modelValue, dateFormat) : dateFormat"
                v-bind="button"
                @click="open = true"
                @touchstart="open = true"
            />

            <UModal v-model:open="open" :title="$t('common.date')">
                <template #body>
                    <div class="flex flex-1 items-center">
                        <DatePickerClient v-model="date" class="mx-auto" v-bind="datePicker" @close="open = false" />
                    </div>
                </template>

                <template #footer>
                    <UButton class="flex-1" color="neutral" block @click="open = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </template>
            </UModal>
        </template>

        <UPopover v-else v-model:open="open" v-bind="popover">
            <UButton
                variant="outline"
                color="neutral"
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
