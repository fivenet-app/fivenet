<script setup lang="ts">
import type { DateValue, Time } from '@internationalized/date';
import type { CalendarProps, InputDateProps } from '@nuxt/ui';
import InputTimePicker from './InputTimePicker.vue';

export type DateRange = { start: Date; end: Date };

export type TimeSplit = { hours: number; minutes: number };

type InputDateRangePickerAttrs = Partial<InputDateProps<true> & CalendarProps<true, false>>;
type CalendarDateRangeValue = { start: DateValue | undefined; end: DateValue | undefined } | null | undefined;

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
        class: props.time ? 'date-range-input w-full' : 'w-full',
        granularity: props.time ? 'minute' : forwardedAttrs.granularity,
        hourCycle: props.time ? (forwardedAttrs.hourCycle ?? 24) : forwardedAttrs.hourCycle,
        range: true,
    };
});

const internalModelValue = computed<CalendarDateRangeValue>({
    get() {
        if (!props.modelValue) return undefined;

        return {
            start: props.time ? dateToCalendarDateTime(props.modelValue.start) : dateToCalendarDate(props.modelValue.start),
            end: props.time ? dateToCalendarDateTime(props.modelValue.end) : dateToCalendarDate(props.modelValue.end),
        };
    },
    set(value) {
        if (value?.start && value.end) {
            const startDate = calendarDateToDate(value.start)!;
            const endDate = calendarDateToDate(value.end)!;
            const startTime = props.time && 'hour' in value.start ? dateToTime(startDate)! : timeState.value.start;
            const endTime = props.time && 'hour' in value.end ? dateToTime(endDate)! : timeState.value.end;

            // Calendar selection returns date-only values; typed time segments return date-time values.
            startDate.setHours(startTime.hour, startTime.minute);
            endDate.setHours(endTime.hour, endTime.minute);

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

<style scoped>
:deep(.date-range-input) {
    container-type: inline-size;
}

:deep(
    .date-range-input
        :is(
            [data-segment='hour'],
            [data-segment='minute'],
            [data-segment='dayPeriod'],
            [data-segment='year'] + [data-segment='literal'],
            [data-segment='hour'] + [data-segment='literal'],
            [data-segment='minute'] + [data-segment='literal']
        )
) {
    display: var(--date-range-time-display, none);
}

@container (min-width: 22rem) {
    :deep(.date-range-input [data-slot='segment']) {
        --date-range-time-display: inline;
    }
}
</style>
