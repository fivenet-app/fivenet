<script setup lang="ts">
import type { CalendarDate, Time } from '@internationalized/date';
import type { ButtonProps, CalendarProps, InputDateProps } from '@nuxt/ui';
import InputTimePicker from './InputTimePicker.vue';

export type TimeSplit = { hours: number; minutes: number };

type InputDatePickerAttrs = Partial<InputDateProps<false> & CalendarProps<false, false>>;

export interface Props {
    modelValue: Date | undefined;
    clearable?: boolean;
    time?: boolean;
    numberOfMonths?: CalendarProps<false, false>['numberOfMonths'];
    isDateDisabled?: CalendarProps<false, false>['isDateDisabled'];
    dateFormat?: string | undefined;
    customDateFormat?: string | 'ago' | undefined;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    button?: ButtonProps & { style?: Record<string, any> };
    hideIcon?: boolean;
    minValue?: CalendarProps<true, false>['minValue'];
    maxValue?: CalendarProps<true, false>['maxValue'];
}

const props = withDefaults(defineProps<Props>(), {
    clearable: false,
    time: false,
    numberOfMonths: undefined,
    isDateDisabled: undefined,
    dateFormat: 'short',
    customDateFormat: undefined,
    button: undefined,
    hideIcon: false,
    minValue: undefined,
    maxValue: undefined,
});

const emits = defineEmits<{
    (e: 'update:modelValue', date: Date | undefined): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const inputDate = useTemplateRef('inputDate');

const numberOfMonths = computed(() => props.numberOfMonths ?? 2);

const attrs = useAttrs() as InputDatePickerAttrs & { class?: unknown };

const inputDateAttrs = computed(() => {
    const { class: _class, ...forwardedAttrs } = attrs;

    return {
        ...forwardedAttrs,
        class: 'min-w-0 flex-1',
        icon: props.hideIcon ? undefined : forwardedAttrs.icon,
    };
});

const internalModelValue = computed<CalendarDate | undefined>({
    get() {
        if (!props.modelValue) return undefined;

        return dateToCalendarDate(props.modelValue);
    },
    set(value) {
        if (value) {
            const date = calendarDateToDate(value);
            const time = dateToTime(props.modelValue);

            if (props.time) {
                date.setHours(time?.hour ?? 0, time?.minute ?? 0, time?.second ?? 0, time?.millisecond ?? 0);
            }

            emits('update:modelValue', date);
        } else {
            emits('update:modelValue', undefined);
        }
    },
});

const internalTimeValue = computed<Time | undefined>({
    get() {
        return dateToTime(props.modelValue);
    },
    set(value) {
        if (!value) {
            emits('update:modelValue', undefined);
            return;
        }

        const date = props.modelValue ? new Date(props.modelValue) : new Date();
        date.setHours(value.hour, value.minute, value.second, value.millisecond);

        emits('update:modelValue', date);
    },
});
</script>

<template>
    <UFieldGroup class="w-full" :class="$attrs.class">
        <UInputDate ref="inputDate" v-model="internalModelValue" v-bind="inputDateAttrs">
            <template #trailing>
                <div class="flex items-center gap-1">
                    <UButton
                        v-if="clearable && modelValue"
                        class="px-0"
                        color="error"
                        variant="link"
                        size="sm"
                        icon="i-mdi-clear"
                        :aria-label="$t('common.clear')"
                        @click.stop="emits('update:modelValue', undefined)"
                    />

                    <UPopover :reference="inputDate?.inputsRef[0]?.$el">
                        <UButton
                            class="px-0"
                            color="neutral"
                            variant="link"
                            size="sm"
                            :icon="props.hideIcon ? undefined : 'i-mdi-calendar'"
                            :aria-label="$t('common.pick_date')"
                            v-bind="props.button"
                        />

                        <template #content>
                            <UCalendar
                                v-model="internalModelValue"
                                class="p-2"
                                :min-value="minValue"
                                :max-value="maxValue"
                                :is-date-disabled="props.isDateDisabled"
                                :is-date-unavailable="attrs.isDateUnavailable"
                                :number-of-months="numberOfMonths"
                            />
                        </template>
                    </UPopover>
                </div>
            </template>
        </UInputDate>

        <InputTimePicker
            v-if="time"
            v-model="internalTimeValue"
            class="shrink-0"
            :disabled="attrs.disabled"
            :readonly="attrs.readonly"
            :required="attrs.required"
            :size="attrs.size"
            :color="attrs.color"
            :variant="attrs.variant"
            :hour-cycle="24"
        />
    </UFieldGroup>
</template>
