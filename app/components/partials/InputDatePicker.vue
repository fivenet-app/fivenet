<script setup lang="ts" generic="R extends boolean, M extends boolean">
import { CalendarDate, DateFormatter, getLocalTimeZone } from '@internationalized/date';
import type { CalendarProps } from '@nuxt/ui';
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';

export type TimeSplit = { hours: number; minutes: number };

export interface Props<R extends boolean = false, M extends boolean = false>
    extends /* @vue-ignore */ Omit<CalendarProps<R, M>, 'modelValue' | 'range'> {
    modelValue: Date | undefined;
    clearable?: boolean;
    time?: boolean;
}

const props = withDefaults(defineProps<Props<R, M>>(), {
    clearable: false,
    time: false,
    class: '',
});

const emits = defineEmits<{
    (e: 'update:modelValue', date: Date | undefined): void;
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
        if (props.modelValue && value) {
            const updatedValue: Date = new Date(props.modelValue);
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
            date.setHours(state.value.hours, state.value.minutes);

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
                {{ df.format(modelValue) }}
            </template>
            <template v-else> {{ $t('common.pick_date') }} </template>
        </UButton>

        <template #content>
            <div class="pb-2">
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
                    :state="state"
                    class="pb-02 flex w-full flex-col items-center justify-center gap-2 md:flex-row"
                >
                    <div class="flex flex-1 items-center justify-center">
                        <UFormField :label="$t('common.time')">
                            <div class="inline-flex flex-row gap-1">
                                <UInputNumber v-model="state.hours" class="max-w-24" :min="0" :max="23" />
                                <span class="font-bold">:</span>
                                <UInputNumber v-model="state.minutes" class="max-w-24" :min="0" :max="59" />
                            </div>
                        </UFormField>
                    </div>
                </UForm>

                <UButton
                    v-if="clearable"
                    variant="outline"
                    color="neutral"
                    block
                    class="mx-2 mt-2"
                    :label="$t('common.clear')"
                    trailing-icon="i-mdi-clear"
                    @click="internalModelValue = undefined"
                />
            </div>
        </template>
    </UPopover>
</template>
