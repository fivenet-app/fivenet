<script setup lang="ts">
import { Time } from '@internationalized/date';
import type { ButtonProps, InputTimeProps } from '@nuxt/ui';

type TimeValue = Pick<Time, 'hour' | 'minute' | 'second' | 'millisecond'>;
type TimeRangeValue = { start: TimeValue | undefined; end: TimeValue | undefined };
type ModelValue = TimeValue | TimeRangeValue | null | undefined;
type InputTimeModelValue = InputTimeProps<boolean>['modelValue'];
type IsTimeUnavailable = (value: TimeValue) => boolean;
type InputTimePickerAttrs = Partial<
    Omit<InputTimeProps<boolean>, 'modelValue' | 'defaultValue' | 'minValue' | 'maxValue' | 'isTimeUnavailable'> & {
        minValue?: TimeValue;
        maxValue?: TimeValue;
        isTimeUnavailable?: IsTimeUnavailable;
    }
>;

export interface Props {
    modelValue: ModelValue;
    clearable?: boolean;
    picker?: boolean;
    stepMinutes?: number;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    button?: ButtonProps & { style?: Record<string, any> };
    hideIcon?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
    clearable: false,
    picker: true,
    stepMinutes: 15,
    button: undefined,
    hideIcon: false,
});

const emits = defineEmits<{
    (e: 'update:modelValue', value: ModelValue): void;
    (e: 'change', event: Event): void;
    (e: 'blur', event: FocusEvent): void;
    (e: 'focus', event: FocusEvent): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const { t } = useI18n();

const inputTime = useTemplateRef('inputTime');

const normalizedStepMinutes = computed(() => Math.max(1, Math.min(720, Math.trunc(props.stepMinutes))));

const attrs = useAttrs() as InputTimePickerAttrs;

const isDisabled = computed(() => !!attrs.disabled || !!attrs.readonly);

const isTimeUnavailable = computed(() => attrs.isTimeUnavailable as IsTimeUnavailable | undefined);

const localModelValue = ref<ModelValue>(props.modelValue);

watch(
    () => props.modelValue,
    (value) => {
        localModelValue.value = value;
    },
);

const internalModelValue = computed<InputTimeModelValue>({
    get() {
        return localModelValue.value as InputTimeModelValue;
    },
    set(value) {
        updateModelValue(value as ModelValue);
    },
});

const inputTimeAttrs = computed<Partial<InputTimeProps<boolean>>>(
    () =>
        ({
            ...attrs,
            hourCycle: attrs.hourCycle ?? 24,
        }) as Partial<InputTimeProps<boolean>>,
);

const pickerButtonProps = computed<ButtonProps>(() => ({
    class: 'px-0',
    color: 'neutral',
    variant: 'link',
    size: 'sm',
    icon: props.hideIcon ? undefined : 'i-mdi-clock-outline',
    disabled: isDisabled.value,
    'aria-label': t('common.pick_time'),
    ...props.button,
}));

const timeOptions = computed(() => {
    const options: { label: string; value: Time; disabled: boolean }[] = [];

    for (let minuteOfDay = 0; minuteOfDay < 24 * 60; minuteOfDay += normalizedStepMinutes.value) {
        const value = new Time(Math.floor(minuteOfDay / 60), minuteOfDay % 60);

        options.push({
            label: formatTime(value),
            value,
            disabled: isTimeDisabled(value),
        });
    }

    return options;
});

function formatTime(value: TimeValue): string {
    return `${value.hour.toString().padStart(2, '0')}:${value.minute.toString().padStart(2, '0')}`;
}

function getTimeMinutes(value: TimeValue | null | undefined): number | undefined {
    if (!value) return undefined;
    return value.hour * 60 + value.minute;
}

function isTimeRange(value: ModelValue): value is TimeRangeValue {
    return typeof value === 'object' && value !== null && ('start' in value || 'end' in value);
}

function isTimeDisabled(value: Time): boolean {
    const minutes = getTimeMinutes(value);
    const minMinutes = getTimeMinutes(attrs.minValue);
    const maxMinutes = getTimeMinutes(attrs.maxValue);

    return (
        (minMinutes !== undefined && minutes !== undefined && minutes < minMinutes) ||
        (maxMinutes !== undefined && minutes !== undefined && minutes > maxMinutes) ||
        !!isTimeUnavailable.value?.(value)
    );
}

function isSelectedTime(value: Time, part?: keyof TimeRangeValue): boolean {
    if (attrs.range) {
        if (!part || !isTimeRange(localModelValue.value)) return false;
        return getTimeMinutes(localModelValue.value[part]) === getTimeMinutes(value);
    }

    return !isTimeRange(localModelValue.value) && getTimeMinutes(localModelValue.value) === getTimeMinutes(value);
}

function updateModelValue(value: ModelValue): void {
    localModelValue.value = value;
    emits('update:modelValue', value);
}

function selectTime(value: Time, part?: keyof TimeRangeValue): void {
    if (isTimeDisabled(value)) return;

    if (attrs.range) {
        const current = isTimeRange(localModelValue.value) ? localModelValue.value : { start: undefined, end: undefined };
        updateModelValue({
            start: part === 'start' ? value : current.start,
            end: part === 'end' ? value : current.end,
        });
        return;
    }

    updateModelValue(value);
}

function clearValue(): void {
    updateModelValue(undefined);
}

function onChange(event: Event): void {
    emits('change', event);
}

function onBlur(event: FocusEvent): void {
    emits('blur', event);
}

function onFocus(event: FocusEvent): void {
    emits('focus', event);
}
</script>

<template>
    <UInputTime
        ref="inputTime"
        v-model="internalModelValue"
        v-bind="inputTimeAttrs"
        @change="onChange"
        @blur="onBlur"
        @focus="onFocus"
    >
        <template v-if="$slots.leading" #leading="{ ui }">
            <slot name="leading" :ui="ui" />
        </template>

        <template v-if="$slots.separator" #separator="{ ui }">
            <slot name="separator" :ui="ui" />
        </template>

        <template #default="{ ui }">
            <slot :ui="ui" />
        </template>

        <template #trailing="{ ui }">
            <div class="flex items-center gap-1">
                <slot name="trailing" :ui="ui" />

                <UButton
                    v-if="clearable && localModelValue"
                    class="px-0"
                    color="error"
                    variant="link"
                    size="sm"
                    icon="i-mdi-clear"
                    :disabled="isDisabled"
                    aria-label="$t('common.clear')"
                    @click.stop="clearValue"
                />

                <UPopover v-if="picker" :reference="inputTime?.inputsRef[0]?.$el">
                    <UButton v-bind="pickerButtonProps" />

                    <template #content>
                        <div v-if="$attrs.range" class="grid max-w-96 grid-cols-2 gap-3 p-2">
                            <div class="min-w-0">
                                <div class="px-1 pb-2 text-sm font-medium text-highlighted">
                                    {{ $t('common.begins_at') }}
                                </div>

                                <div class="grid max-h-64 grid-cols-3 gap-1 overflow-y-auto pr-1">
                                    <UButton
                                        v-for="option in timeOptions"
                                        :key="`start-${option.label}`"
                                        :label="option.label"
                                        :disabled="option.disabled"
                                        variant="ghost"
                                        :color="isSelectedTime(option.value, 'start') ? 'primary' : 'neutral'"
                                        size="sm"
                                        block
                                        @click="selectTime(option.value, 'start')"
                                    />
                                </div>
                            </div>

                            <div class="min-w-0">
                                <div class="px-1 pb-2 text-sm font-medium text-highlighted">
                                    {{ $t('common.ends_at') }}
                                </div>

                                <div class="grid max-h-64 grid-cols-3 gap-1 overflow-y-auto pr-1">
                                    <UButton
                                        v-for="option in timeOptions"
                                        :key="`end-${option.label}`"
                                        :label="option.label"
                                        :disabled="option.disabled"
                                        variant="ghost"
                                        :color="isSelectedTime(option.value, 'end') ? 'primary' : 'neutral'"
                                        size="sm"
                                        block
                                        @click="selectTime(option.value, 'end')"
                                    />
                                </div>
                            </div>
                        </div>

                        <div v-else class="grid max-h-64 w-72 grid-cols-4 gap-1 overflow-y-auto p-2">
                            <UButton
                                v-for="option in timeOptions"
                                :key="option.label"
                                :label="option.label"
                                :disabled="option.disabled"
                                :variant="isSelectedTime(option.value) ? 'solid' : 'ghost'"
                                color="neutral"
                                size="sm"
                                block
                                @click="selectTime(option.value)"
                            />
                        </div>
                    </template>
                </UPopover>
            </div>
        </template>
    </UInputTime>
</template>
