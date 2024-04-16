<script lang="ts" setup>
import { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        value: Date | Timestamp | undefined;
        type?: 'short' | 'long' | 'compact' | 'date';
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

const date: Date = props.value instanceof Date ? props.value : toDate(props.value) ?? new Date();

const timeClass = ref('');
if (props.updateCallback) {
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
