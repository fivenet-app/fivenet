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
    <time :datetime="date.toLocaleTimeString()" :title="!ago ? $d(date, type) : useLocaleTimeAgo(date).value">
        {{ !ago ? $d(date, type) : useLocaleTimeAgo(date).value }}
    </time>
</template>
