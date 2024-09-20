<script lang="ts" setup>
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        value: Date | Timestamp | undefined;
        type?: 'date' | 'short' | 'long' | 'compact' | 'time';
        ago?: boolean;
        updateInterval?: number;
        updateCallback?: () => string;
    }>(),
    {
        type: 'short',
        ago: false,
        updateInterval: 1000,
    },
);

const date = computed<Date>(() => (props.value instanceof Date ? props.value : (toDate(props.value) ?? new Date())));

const timeClass = ref('');
if (props.updateCallback) {
    timeClass.value = props.updateCallback();

    useIntervalFn(() => (timeClass.value = props.updateCallback!()), props.updateInterval);
}
</script>

<template>
    <UTooltip :text="$d(date, 'long')">
        <time v-bind="$attrs" :datetime="date.toLocaleTimeString()" :class="timeClass">
            {{ !ago ? $d(date, type) : useLocaleTimeAgo(date).value }}
        </time>
    </UTooltip>
</template>
