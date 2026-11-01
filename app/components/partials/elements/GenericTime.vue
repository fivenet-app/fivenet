<script lang="ts" setup>
import { isSameDay } from 'date-fns';
import type { BadgeProps } from '@nuxt/ui';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import type { I18nDateTimeFormats } from '~/utils/time';

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        value: Date | Timestamp | undefined;
        type?: I18nDateTimeFormats;
        compactToday?: boolean;
        tooltipType?: I18nDateTimeFormats;
        ago?: boolean;
        updateInterval?: number;
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        updateCallback?: () => Record<string, any>;
        badge?: boolean;
        size?: BadgeProps['size'];
    }>(),
    {
        type: 'short',
        compactToday: false,
        tooltipType: undefined,
        ago: false,
        updateInterval: 1000,
        updateCallback: undefined,
        badge: false,
        size: 'xs',
    },
);

const tooltipTypeMap: Partial<{ [key in I18nDateTimeFormats]: I18nDateTimeFormats }> = {
    date: 'longDate',
    shortDate: 'longDate',
    short: 'long',
    compact: 'long',
};

const date = computed<Date>(() => (props.value instanceof Date ? props.value : (toDate(props.value) ?? new Date())));
const minuteClock = useMinuteClock();
const displayType = computed<I18nDateTimeFormats>(() =>
    props.compactToday && props.type === 'compact' && !isSameDay(date.value, minuteClock.value) ? 'short' : props.type,
);

const timeAttrs = ref({});
if (props.updateCallback) {
    timeAttrs.value = props.updateCallback();

    useIntervalFn(() => (timeAttrs.value = props.updateCallback!()), props.updateInterval);
}
</script>

<template>
    <UTooltip :text="$d(date, tooltipType ?? tooltipTypeMap[displayType] ?? displayType)">
        <UBadge v-if="badge" :size="size" v-bind="{ ...$attrs, ...timeAttrs }">
            <time :datetime="date.toLocaleTimeString()" v-bind="$attrs">
                {{ !ago ? $d(date, displayType) : useLocaleTimeAgo(date).value }}
            </time>
        </UBadge>
        <time v-else :datetime="date.toLocaleTimeString()" v-bind="{ ...$attrs, ...timeAttrs }">
            {{ !ago ? $d(date, displayType) : useLocaleTimeAgo(date).value }}
        </time>
    </UTooltip>
</template>
