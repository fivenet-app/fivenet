<script lang="ts" setup>
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';
import { format } from 'date-fns';
import type { DatePicker as VCalendarDatePicker } from 'v-calendar';
import DateRangePickerClient from './DateRangePicker.client.vue';

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        modelValue?: { start: Date; end: Date } | undefined;
        clearable?: boolean;
        dateFormat?: string;
        popover?: object;
        button?: object;
        datePicker?: (typeof VCalendarDatePicker)['$props'];
    }>(),
    {
        modelValue: undefined,
        clearable: false,
        dateFormat: 'dd.MM.yyyy',
        popover: undefined,
        button: undefined,
        datePicker: undefined,
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', modelValue: { start: Date; end: Date } | undefined): void;
}>();

const date = useVModel(props, 'modelValue', emit);

const breakpoints = useBreakpoints(breakpointsTailwind);

const smallerThanSm = breakpoints.smaller('sm');

const open = ref(false);
</script>

<template>
    <ClientOnly>
        <template v-if="smallerThanSm">
            <UButton
                variant="outline"
                color="neutral"
                block
                icon="i-mdi-calendar-month"
                truncate
                :label="
                    modelValue
                        ? `${format(modelValue.start, dateFormat)} - ${format(modelValue.end, dateFormat)}`
                        : `${dateFormat} - ${dateFormat}`
                "
                v-bind="button"
                @click="open = true"
                @touchstart="open = true"
            />

            <UModal v-model:open="open">
                <UCard>
                    <template #header>
                        <div class="flex items-center justify-between">
                            <h3 class="text-2xl font-semibold leading-6">
                                {{ $t('common.date') }}
                            </h3>

                            <UButton
                                class="-my-1"
                                color="neutral"
                                variant="ghost"
                                icon="i-mdi-window-close"
                                @click="open = false"
                            />
                        </div>
                    </template>

                    <div class="flex flex-1 items-center">
                        <DateRangePickerClient v-model="date" v-bind="datePicker" @close="open = false" />
                    </div>

                    <template #footer>
                        <UButton class="flex-1" color="neutral" block @click="open = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </template>
                </UCard>
            </UModal>
        </template>

        <UPopover v-else v-model:open="open" v-bind="popover">
            <UButton
                variant="outline"
                color="neutral"
                block
                truncate
                icon="i-mdi-calendar-month"
                :label="
                    modelValue
                        ? `${format(modelValue.start, dateFormat)} - ${format(modelValue.end, dateFormat)}`
                        : `${dateFormat} - ${dateFormat}`
                "
                v-bind="button"
                @touchstart="open = true"
            />

            <template #panel="{ close }">
                <DateRangePickerClient v-model="date" v-bind="datePicker" @close="close" />
            </template>
        </UPopover>
    </ClientOnly>
</template>
