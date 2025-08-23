<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
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
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        updateCallback?: () => Record<string, any>;
        badge?: boolean;
        size?: BadgeProps['size'];
    }>(),
    {
        type: 'short',
        tooltipType: undefined,
        ago: false,
        updateInterval: 1000,
        updateCallback: undefined,
        badge: false,
        size: 'xs',
    },
);

const tooltipTypeMap: Partial<{ [key in dateTimeFormats]: dateTimeFormats }> = {
    date: 'longDate',
    shortDate: 'longDate',
    short: 'long',
    compact: 'long',
};

const date = computed<Date>(() => (props.value instanceof Date ? props.value : (toDate(props.value) ?? new Date())));

const timeAttrs = ref({});
if (props.updateCallback) {
    timeAttrs.value = props.updateCallback();

    useIntervalFn(() => (timeAttrs.value = props.updateCallback!()), props.updateInterval);
}
</script>

<template>
    <UTooltip :text="$d(date, tooltipType ?? tooltipTypeMap[type] ?? type)">
        <UBadge v-if="badge" v-bind="{ ...$attrs, ...timeAttrs }" :size="size">
            <time :datetime="date.toLocaleTimeString()" v-bind="$attrs">
                {{ !ago ? $d(date, type) : useLocaleTimeAgo(date).value }}
            </time>
        </UBadge>
        <time v-else :datetime="date.toLocaleTimeString()" v-bind="{ ...$attrs, ...timeAttrs }">
            {{ !ago ? $d(date, type) : useLocaleTimeAgo(date).value }}
        </time>
    </UTooltip>
</template>
