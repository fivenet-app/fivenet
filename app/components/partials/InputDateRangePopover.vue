<script setup lang="ts" generic="R extends boolean, M extends boolean">
import { CalendarDate, DateFormatter, getLocalTimeZone } from '@internationalized/date';
import type { CalendarProps } from '@nuxt/ui';
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';

export type DateRange = { start: Date; end: Date };

export type TimeSplit = { hours: number; minutes: number };

export interface Props<R extends boolean = false, M extends boolean = false> extends /* @vue-ignore */ Omit<
    CalendarProps<R, M>,
    'modelValue' | 'range'
> {
    modelValue: DateRange | undefined;
    clearable?: boolean;
    time?: boolean;
}

const props = withDefaults(defineProps<Props<R, M>>(), {
    clearable: false,
    time: false,
    class: '',
});

const emits = defineEmits<{
    (e: 'update:modelValue', date: DateRange | undefined): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const { locale } = useI18n();

const df = new DateFormatter(locale.value, {
    dateStyle: 'medium',
});

// State derived from modelValue
const state = computed({
    get() {
        if (!props.modelValue) {
            return {
                start: { hours: 0, minutes: 0 },
                end: { hours: 23, minutes: 59 },
            };
        }

        return {
            start: {
                hours: props.modelValue.start.getHours(),
                minutes: props.modelValue.start.getMinutes(),
            },
            end: {
                hours: props.modelValue.end.getHours(),
                minutes: props.modelValue.end.getMinutes(),
            },
        };
    },
    set(value) {
        if (props.modelValue && value) {
            const updatedValue: DateRange = {
                start: new Date(props.modelValue.start),
                end: new Date(props.modelValue.end),
            };

            updatedValue.start.setHours(value.start.hours, value.start.minutes);
            updatedValue.end.setHours(value.end.hours, value.end.minutes);

            emits('update:modelValue', updatedValue);
        }
    },
});

// Internal modelValue using CalendarDate
const internalModelValue = computed({
    get() {
        if (!props.modelValue) return undefined;

        return {
            start: new CalendarDate(
                props.modelValue.start.getFullYear(),
                props.modelValue.start.getMonth() + 1,
                props.modelValue.start.getDate(),
            ),
            end: new CalendarDate(
                props.modelValue.end.getFullYear(),
                props.modelValue.end.getMonth() + 1,
                props.modelValue.end.getDate(),
            ),
        };
    },
    set(value) {
        if (value) {
            const startDate = value.start.toDate(getLocalTimeZone());
            const endDate = value.end.toDate(getLocalTimeZone());

            // Apply the time state to the dates
            startDate.setHours(state.value.start.hours, state.value.start.minutes);
            endDate.setHours(state.value.end.hours, state.value.end.minutes);

            emits('update:modelValue', { start: startDate, end: endDate });
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
            <template v-if="modelValue?.start">
                <template v-if="modelValue.end">
                    {{ df.format(modelValue.start) }} -
                    {{ df.format(modelValue.end) }}
                </template>

                <template v-else>
                    {{ df.format(modelValue.start) }}
                </template>
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
                    range
                />

                <UForm
                    v-if="time"
                    :schema="{}"
                    :state="state"
                    class="flex w-full flex-col items-center justify-center gap-2 pb-2 md:flex-row"
                >
                    <div class="flex flex-1 items-center justify-center">
                        <UFormField :label="$t('common.from')">
                            <div class="inline-flex flex-row gap-1">
                                <UInputNumber v-model="state.start.hours" class="max-w-24" :min="0" :max="23" />
                                <span class="font-bold">:</span>
                                <UInputNumber v-model="state.start.minutes" class="max-w-24" :min="0" :max="59" />
                            </div>
                        </UFormField>
                    </div>

                    <div class="flex flex-1 items-center justify-center">
                        <UFormField :label="$t('common.to')">
                            <div class="inline-flex flex-row gap-1">
                                <UInputNumber v-model="state.end.hours" class="max-w-24" :min="0" :max="23" />
                                <span class="font-bold">:</span>
                                <UInputNumber v-model="state.end.minutes" class="max-w-24" :min="0" :max="59" />
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
