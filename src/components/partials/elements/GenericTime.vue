<script lang="ts" setup>
import { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

const props = withDefaults(
    defineProps<{
        value: Date | Timestamp | undefined;
        type?: 'short' | 'long' | 'compact' | 'date';
        ago?: boolean;
    }>(),
    {
        type: 'short',
        ago: false,
    },
);

const date: Date = props.value instanceof Date ? props.value : toDate(props.value) ?? new Date();
</script>

<template>
    <UTooltip :text="$d(date, 'long')">
        <time :datetime="date.toLocaleTimeString()">
            {{ !ago ? $d(date, type) : useLocaleTimeAgo(date).value }}
        </time>
    </UTooltip>
</template>
