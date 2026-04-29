<script lang="ts" setup>
import { clampDuration, durationToSeconds, secondsToDuration } from '~/utils/duration';
import type { Duration } from '~~/gen/ts/google/protobuf/duration';

type DurationUnit = 'second' | 'minute' | 'hour' | 'day';

interface Props {
    modelValue: Duration | undefined;
    mode?: 'single' | 'composite';
    units?: DurationUnit[];
    min?: Duration;
    max?: Duration;
    clearable?: boolean;
    disabled?: boolean;
    readonly?: boolean;
    step?: number;
}

const props = withDefaults(defineProps<Props>(), {
    mode: 'single',
    units: () => ['second', 'minute', 'hour', 'day'],
    clearable: false,
    disabled: false,
    readonly: false,
    step: 0.25,
    min: undefined,
    max: undefined,
});

const emit = defineEmits<{
    (e: 'update:modelValue', value: Duration | undefined): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const unitSecondsMap: Record<DurationUnit, number> = {
    second: 1,
    minute: 60,
    hour: 60 * 60,
    day: 24 * 60 * 60,
};

const unitOrder: DurationUnit[] = ['day', 'hour', 'minute', 'second'];

const activeUnits = computed<DurationUnit[]>(() => {
    const deduped = [...new Set(props.units)].filter((u): u is DurationUnit => u in unitSecondsMap);
    return deduped.length > 0 ? deduped : ['minute'];
});

const orderedUnits = computed(() => unitOrder.filter((unit) => activeUnits.value.includes(unit)));

const unitOptions = computed(() => activeUnits.value.map((unit) => ({ value: unit })));

const selectedUnit = ref<DurationUnit>(activeUnits.value[0] ?? 'minute');
const isInternalSingleValueUpdate = ref(false);
watch(
    activeUnits,
    (units) => {
        if (!units.includes(selectedUnit.value)) {
            selectedUnit.value = units[0] ?? 'minute';
        }
    },
    { immediate: true },
);

function areDurationsEqual(a: Duration | undefined, b: Duration | undefined): boolean {
    if (!a && !b) return true;
    if (!a || !b) return false;
    return a.seconds === b.seconds && a.nanos === b.nanos;
}

function clampModelValue(value: Duration | undefined): Duration | undefined {
    if (!value) return undefined;
    return clampDuration(value, props.min, props.max);
}

function emitClampedSeconds(seconds: number): void {
    const clamped = clampDuration(secondsToDuration(seconds), props.min, props.max);
    emit('update:modelValue', clamped);
}

watch(
    () => [props.modelValue, props.min, props.max] as const,
    () => {
        if (!props.modelValue) return;

        const clamped = clampModelValue(props.modelValue);
        if (!areDurationsEqual(clamped, props.modelValue)) {
            emit('update:modelValue', clamped);
        }
    },
    { immediate: true },
);
watch(
    () => props.mode,
    () => {
        if (!props.modelValue) return;

        const clamped = clampModelValue(props.modelValue);
        if (!areDurationsEqual(clamped, props.modelValue)) {
            emit('update:modelValue', clamped);
        }
    },
);

const singleUnitSeconds = computed(() => unitSecondsMap[selectedUnit.value]);

function getStepPrecision(step: number): number {
    if (!Number.isFinite(step) || step <= 0) return 6;
    const [, decimals = ''] = step.toString().split('.');
    return Math.min(Math.max(decimals.length, 0) + 3, 9);
}

function snapSingleValueToStep(value: number): number {
    if (!Number.isFinite(value)) return 0;
    if (!Number.isFinite(props.step) || props.step <= 0) return value;

    const snapped = Math.round(value / props.step) * props.step;
    return Number(snapped.toFixed(getStepPrecision(props.step)));
}

function isStepAligned(value: number): boolean {
    if (!Number.isFinite(value)) return false;
    if (!Number.isFinite(props.step) || props.step <= 0) return Number.isInteger(value);

    const quotient = value / props.step;
    return Math.abs(quotient - Math.round(quotient)) < 1e-9;
}

function resolvePreferredSingleUnit(seconds: number): DurationUnit {
    const units = orderedUnits.value;
    if (units.length === 0) return 'minute';

    const alignedUnit = units.find((unit) => isStepAligned(seconds / unitSecondsMap[unit]));
    if (alignedUnit) return alignedUnit;

    const unitAtOrAboveOne = units.find((unit) => seconds / unitSecondsMap[unit] >= 1);
    if (unitAtOrAboveOne) return unitAtOrAboveOne;

    return units[units.length - 1] ?? 'minute';
}

const singleValue = computed(() => {
    if (!props.modelValue) return 0;
    const normalized = durationToSeconds(clampModelValue(props.modelValue)) / singleUnitSeconds.value;
    return snapSingleValueToStep(normalized);
});

const singleMin = computed(() =>
    props.min ? durationToSeconds(clampModelValue(props.min) ?? props.min) / singleUnitSeconds.value : undefined,
);
const singleMax = computed(() =>
    props.max ? durationToSeconds(clampModelValue(props.max) ?? props.max) / singleUnitSeconds.value : undefined,
);

watch(
    () => [props.modelValue, orderedUnits.value, props.step] as const,
    () => {
        if (isInternalSingleValueUpdate.value) {
            isInternalSingleValueUpdate.value = false;
            return;
        }

        if (!props.modelValue) return;

        const seconds = durationToSeconds(clampModelValue(props.modelValue));
        const preferredUnit = resolvePreferredSingleUnit(seconds);
        if (preferredUnit !== selectedUnit.value) {
            selectedUnit.value = preferredUnit;
        }
    },
    { immediate: true },
);

function onSingleValueUpdate(value: number | undefined): void {
    if (value === undefined || value === null) return;
    isInternalSingleValueUpdate.value = true;
    emitClampedSeconds(snapSingleValueToStep(value) * singleUnitSeconds.value);
}

function onSingleUnitUpdate(unit: DurationUnit): void {
    selectedUnit.value = unit;
    if (!props.modelValue) return;

    const clamped = clampModelValue(props.modelValue);
    if (!areDurationsEqual(clamped, props.modelValue)) {
        emit('update:modelValue', clamped);
    }
}

type CompositeState = Record<DurationUnit, number>;

const compositeState = ref<CompositeState>({
    day: 0,
    hour: 0,
    minute: 0,
    second: 0,
});

function decomposeSecondsForComposite(totalSeconds: number, units: DurationUnit[]): CompositeState {
    const state: CompositeState = {
        day: 0,
        hour: 0,
        minute: 0,
        second: 0,
    };

    if (units.length === 0) return state;

    const smallestUnit = units[units.length - 1]!;
    const smallestSeconds = unitSecondsMap[smallestUnit];
    const roundedTotalSeconds = Math.max(0, Math.floor(totalSeconds / smallestSeconds) * smallestSeconds);

    let remaining = roundedTotalSeconds;
    units.forEach((unit, index) => {
        const seconds = unitSecondsMap[unit];
        if (index === units.length - 1) {
            state[unit] = Math.round(remaining / seconds);
            return;
        }

        const count = Math.floor(remaining / seconds);
        state[unit] = count;
        remaining -= count * seconds;
    });

    return state;
}

watch(
    () => [props.modelValue, activeUnits.value, props.min, props.max] as const,
    () => {
        const clamped = clampModelValue(props.modelValue);
        const seconds = clamped ? durationToSeconds(clamped) : 0;
        compositeState.value = decomposeSecondsForComposite(seconds, orderedUnits.value);
    },
    { immediate: true },
);

function onCompositeFieldUpdate(unit: DurationUnit, value: number | undefined): void {
    const normalized = Math.max(0, Math.trunc(value ?? 0));
    compositeState.value[unit] = normalized;

    const seconds = orderedUnits.value.reduce(
        (acc, currentUnit) => acc + compositeState.value[currentUnit] * unitSecondsMap[currentUnit],
        0,
    );
    emitClampedSeconds(seconds);
}

function clearValue(): void {
    emit('update:modelValue', undefined);
}
</script>

<template>
    <div class="flex flex-col gap-2" v-bind="$attrs">
        <div v-if="mode === 'single'" class="flex items-center gap-2">
            <UFieldGroup class="w-full">
                <UInputNumber
                    class="w-full"
                    :model-value="singleValue"
                    :step="step"
                    :min="singleMin"
                    :max="singleMax"
                    :disabled="disabled"
                    :readonly="readonly"
                    @update:model-value="onSingleValueUpdate"
                />

                <ClientOnly>
                    <USelectMenu
                        class="min-w-36"
                        :model-value="selectedUnit"
                        value-key="value"
                        :items="unitOptions"
                        :disabled="disabled || readonly || activeUnits.length <= 1"
                        @update:model-value="onSingleUnitUpdate"
                    >
                        <template #default>
                            {{ $t(`common.time_ago.${selectedUnit}`, 2) }}
                        </template>

                        <template #item-label="{ item }">
                            {{ $t(`common.time_ago.${item.value}`, 2) }}
                        </template>
                    </USelectMenu>
                </ClientOnly>
            </UFieldGroup>

            <UTooltip v-if="clearable" :text="$t('common.clear')">
                <UButton
                    color="error"
                    variant="outline"
                    icon="i-mdi-clear"
                    :disabled="disabled || readonly"
                    @click="clearValue"
                />
            </UTooltip>
        </div>

        <div
            v-else
            class="grid gap-2"
            :class="orderedUnits.length >= 3 ? 'grid-cols-3' : orderedUnits.length === 2 ? 'grid-cols-2' : 'grid-cols-1'"
        >
            <UFormField v-for="unit in orderedUnits" :key="unit" :label="$t(`common.time_ago.${unit}`, 2)">
                <UInputNumber
                    class="w-full"
                    :model-value="compositeState[unit]"
                    :step="1"
                    :min="0"
                    :disabled="disabled"
                    :readonly="readonly"
                    @update:model-value="(value) => onCompositeFieldUpdate(unit, value)"
                />
            </UFormField>

            <UTooltip v-if="clearable" :text="$t('common.clear')">
                <UButton
                    color="error"
                    variant="outline"
                    icon="i-mdi-clear"
                    :disabled="disabled || readonly"
                    @click="clearValue"
                />
            </UTooltip>
        </div>
    </div>
</template>
