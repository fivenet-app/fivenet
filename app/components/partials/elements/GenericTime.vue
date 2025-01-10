<script lang="ts" setup>
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

defineOptions({
    inheritAttrs: false,
});

type dateTimeFormats = 'date' | 'shortDate' | 'longDate' | 'short' | 'long' | 'compact' | 'time';

const props = withDefaults(
    defineProps<{
        value: Date | Timestamp | undefined;
        type?: dateTimeFormats;
        tooltipType?: dateTimeFormats;
        ago?: boolean;
        updateInterval?: number;
        updateCallback?: () => string;
    }>(),
    {
        type: 'short',
        tooltipType: undefined,
        ago: false,
        updateInterval: 1000,
        updateCallback: undefined,
    },
);

const tooltipTypeMap: Partial<{ [key in dateTimeFormats]: dateTimeFormats }> = {
    date: 'longDate',
    shortDate: 'longDate',
    short: 'long',
};

const date = computed<Date>(() => (props.value instanceof Date ? props.value : (toDate(props.value) ?? new Date())));

const timeClass = ref('');
if (props.updateCallback) {
    timeClass.value = props.updateCallback();

    useIntervalFn(() => (timeClass.value = props.updateCallback!()), props.updateInterval);
}
</script>

<template>
    <UTooltip :text="$d(date, tooltipType ?? tooltipTypeMap[type] ?? type)">
        <time v-bind="$attrs" :datetime="date.toLocaleTimeString()" :class="timeClass">
            {{ !ago ? $d(date, type) : useLocaleTimeAgo(date).value }}
        </time>
    </UTooltip>
</template>
