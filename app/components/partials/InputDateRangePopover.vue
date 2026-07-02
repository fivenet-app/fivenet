<script setup lang="ts">
import type { CalendarDate, Time } from '@internationalized/date';
import type { CalendarProps, InputDateProps } from '@nuxt/ui';
import InputTimePicker from './InputTimePicker.vue';

export type DateRange = { start: Date; end: Date };

export type TimeSplit = { hours: number; minutes: number };

type InputDateRangePickerAttrs = Partial<InputDateProps<true> & CalendarProps<true, false>>;
type CalendarDateRangeValue = { start: CalendarDate | undefined; end: CalendarDate | undefined } | null | undefined;

export interface Props {
    modelValue: DateRange | undefined;
    clearable?: boolean;
    time?: boolean;
    numberOfMonths?: CalendarProps<true, false>['numberOfMonths'];
    isDateDisabled?: CalendarProps<true, false>['isDateDisabled'];
    isDateUnavailable?: CalendarProps<true, false>['isDateUnavailable'];
    isDateHighlightable?: CalendarProps<true, false>['isDateHighlightable'];
    allowNonContiguousRanges?: CalendarProps<true, false>['allowNonContiguousRanges'];
    fixedDate?: CalendarProps<true, false>['fixedDate'];
    maximumDays?: CalendarProps<true, false>['maximumDays'];
    minValue?: CalendarProps<true, false>['minValue'];
    maxValue?: CalendarProps<true, false>['maxValue'];
}

const props = withDefaults(defineProps<Props>(), {
    clearable: false,
    time: false,
    isDateDisabled: undefined,
    isDateUnavailable: undefined,
    isDateHighlightable: undefined,
    allowNonContiguousRanges: undefined,
    fixedDate: undefined,
    maximumDays: undefined,
    numberOfMonths: undefined,
    minValue: undefined,
    maxValue: undefined,
});

const emits = defineEmits<{
    (e: 'update:modelValue', date: DateRange | undefined): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const dateFormatter = useDateFormatterWithOptions('medium');

const inputDate = useTemplateRef('inputDate');

const numberOfMonths = computed(() => props.numberOfMonths ?? 2);

const attrs = useAttrs() as InputDateRangePickerAttrs & { class?: unknown };

const inputDateAttrs = computed(() => {
    const { class: _class, ...forwardedAttrs } = attrs;

    return {
        ...forwardedAttrs,
        class: 'w-full',
        range: true,
    };
});

const internalModelValue = computed<CalendarDateRangeValue>({
    get() {
        if (!props.modelValue) return undefined;

        return {
            start: dateToCalendarDate(props.modelValue.start),
            end: dateToCalendarDate(props.modelValue.end),
        };
    },
    set(value) {
        if (value?.start && value.end) {
            const startDate = calendarDateToDate(value.start)!;
            const endDate = calendarDateToDate(value.end)!;

            // Apply the time state to the dates
            startDate.setHours(timeState.value.start.hour, timeState.value.start.minute);
            endDate.setHours(timeState.value.end.hour, timeState.value.end.minute);

            emits('update:modelValue', { start: startDate, end: endDate });
        } else {
            emits('update:modelValue', undefined);
        }
    },
});

const timeState = computed<{ start: Time; end: Time }>({
    get() {
        return {
            start: dateToTime(props.modelValue?.start) ?? dateToTime(new Date(new Date().setHours(0, 0, 0, 0))),
            end: dateToTime(props.modelValue?.end) ?? dateToTime(new Date(new Date().setHours(23, 59, 0, 0))),
        };
    },
    set(value) {
        if (!props.modelValue) return;

        const startDate = new Date(props.modelValue.start);
        const endDate = new Date(props.modelValue.end);

        startDate.setHours(value.start.hour, value.start.minute, value.start.second, value.start.millisecond);
        endDate.setHours(value.end.hour, value.end.minute, value.end.second, value.end.millisecond);

        emits('update:modelValue', { start: startDate, end: endDate });
    },
});
</script>

<template>
    <div class="flex w-full flex-col gap-2" :class="$attrs.class">
        <UInputDate ref="inputDate" v-model="internalModelValue" v-bind="inputDateAttrs">
            <template #trailing>
                <div class="flex items-center gap-1">
                    <UTooltip v-if="clearable && modelValue" :text="$t('common.clear')">
                        <UButton
                            class="px-0"
                            color="error"
                            variant="link"
                            size="sm"
                            icon="i-mdi-clear"
                            :aria-label="$t('common.clear')"
                            @click.stop="emits('update:modelValue', undefined)"
                        />
                    </UTooltip>

                    <UPopover :reference="inputDate?.inputsRef[0]?.$el">
                        <UTooltip :text="$t('common.pick_date')">
                            <UButton
                                class="px-0"
                                color="neutral"
                                variant="link"
                                size="sm"
                                icon="i-mdi-calendar"
                                :aria-label="$t('common.pick_date')"
                            />
                        </UTooltip>

                        <template #content>
                            <div>
                                <UCalendar
                                    v-model="internalModelValue"
                                    class="p-2"
                                    :range="true"
                                    :min-value="minValue"
                                    :max-value="maxValue"
                                    :is-date-disabled="isDateDisabled"
                                    :is-date-unavailable="isDateUnavailable"
                                    :is-date-highlightable="isDateHighlightable"
                                    :allow-non-contiguous-ranges="allowNonContiguousRanges"
                                    :fixed-date="fixedDate"
                                    :maximum-days="maximumDays"
                                    :number-of-months="numberOfMonths"
                                />

                                <div class="px-2 py-1 pb-2">
                                    <InputTimePicker
                                        v-if="time"
                                        v-model="timeState"
                                        class="w-full border-t border-default p-2"
                                        :range="true"
                                        :disabled="$attrs.disabled"
                                        :readonly="$attrs.readonly"
                                        :required="$attrs.required"
                                        :size="$attrs.size"
                                        :color="$attrs.color"
                                        :variant="$attrs.variant"
                                        :hour-cycle="24"
                                    />
                                </div>
                            </div>
                        </template>
                    </UPopover>
                </div>
            </template>
        </UInputDate>

        <span v-if="modelValue?.start" class="sr-only">
            {{ dateFormatter.format(modelValue.start) }}
            <template v-if="modelValue.end"> - {{ dateFormatter.format(modelValue.end) }} </template>
        </span>
    </div>
</template>
