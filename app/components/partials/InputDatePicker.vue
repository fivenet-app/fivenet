<script setup lang="ts" generic="R extends boolean, M extends boolean">
import { CalendarDate, getLocalTimeZone } from '@internationalized/date';
import type { CalendarProps } from '@nuxt/ui';
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';
import { format } from 'date-fns';

export type TimeSplit = { hours: number; minutes: number };

export interface Props<R extends boolean = false, M extends boolean = false>
    extends /* @vue-ignore */ Omit<CalendarProps<R, M>, 'modelValue' | 'range'> {
    modelValue: Date | undefined;
    clearable?: boolean;
    time?: boolean;
    dateFormat?: string | undefined;
    customDateFormat?: string | 'ago' | undefined;
}

const props = withDefaults(defineProps<Props<R, M>>(), {
    clearable: false,
    time: false,
    class: '',
    dateFormat: 'short',
    customDateFormat: undefined,
});

const emits = defineEmits<{
    (e: 'update:modelValue', date: Date | undefined): void;
}>();

defineOptions({
    inheritAttrs: false,
});

// State derived from modelValue
const timeState = computed({
    get() {
        if (!props.modelValue) {
            return {
                hours: 0,
                minutes: 0,
            };
        }

        return {
            hours: props.modelValue.getHours(),
            minutes: props.modelValue.getMinutes(),
        };
    },
    set(value) {
        if (value) {
            const updatedValue: Date = props.modelValue ? new Date(props.modelValue) : new Date();
            updatedValue.setHours(value.hours, value.minutes);

            emits('update:modelValue', updatedValue);
        }
    },
});

// Internal modelValue using CalendarDate
const internalModelValue = computed({
    get() {
        if (!props.modelValue) return undefined;

        return new CalendarDate(props.modelValue.getFullYear(), props.modelValue.getMonth() + 1, props.modelValue.getDate());
    },
    set(value) {
        if (value) {
            const date = value.toDate(getLocalTimeZone());

            // Apply the time state to the dates
            date.setHours(timeState.value.hours, timeState.value.minutes);

            emits('update:modelValue', date);
        } else {
            emits('update:modelValue', undefined);
        }
    },
});

const breakpoints = useBreakpoints(breakpointsTailwind);

const smallerThanSm = breakpoints.smaller('sm');
</script>

<template>
    <UPopover>
        <UButton
            color="neutral"
            variant="subtle"
            :icon="smallerThanSm ? undefined : 'i-mdi-calendar'"
            class="inline-flex w-full gap-2"
            block
        >
            <template v-if="modelValue">
                {{
                    customDateFormat
                        ? customDateFormat === 'ago'
                            ? useLocaleTimeAgo(modelValue)
                            : format(modelValue, customDateFormat)
                        : $d(modelValue, dateFormat)
                }}
            </template>
            <template v-else> {{ $t('common.pick_date') }} </template>
        </UButton>

        <template #content>
            <div class="flex flex-col items-center pb-2">
                <UCalendar
                    v-model="internalModelValue"
                    class="p-1"
                    :number-of-months="smallerThanSm ? 1 : 2"
                    v-bind="$attrs"
                    :range="false"
                />

                <UForm
                    v-if="time"
                    :schema="{}"
                    :state="timeState"
                    class="flex w-full flex-col items-center justify-center gap-2 pb-2 md:flex-row"
                >
                    <div class="flex flex-1 items-center justify-center">
                        <UFormField :label="$t('common.time')">
                            <div class="inline-flex flex-row items-center gap-1">
                                <UInputNumber
                                    :model-value="timeState.hours"
                                    class="max-w-24"
                                    :min="0"
                                    :max="23"
                                    @update:model-value="
                                        ($event) => (timeState = { hours: $event, minutes: timeState.minutes })
                                    "
                                />
                                <span class="font-bold">:</span>
                                <UInputNumber
                                    :model-value="timeState.minutes"
                                    class="max-w-24"
                                    :min="0"
                                    :max="59"
                                    @update:model-value="($event) => (timeState = { hours: timeState.hours, minutes: $event })"
                                />
                            </div>
                        </UFormField>
                    </div>
                </UForm>

                <div v-if="clearable" class="w-full px-2">
                    <UButton
                        variant="outline"
                        color="error"
                        block
                        :label="$t('common.clear')"
                        trailing-icon="i-mdi-clear"
                        @click="internalModelValue = undefined"
                    />
                </div>
            </div>
        </template>
    </UPopover>
</template>
