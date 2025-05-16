<script lang="ts" setup>
import { isSameDay } from 'date-fns';
import type { Attribute } from 'v-calendar/dist/types/src/utils/attribute.js';
import type { CalendarDay } from 'v-calendar/dist/types/src/utils/page.js';
import type { CalendarEntry } from '~~/gen/ts/resources/calendar/calendar';

const props = defineProps<{
    day: CalendarDay;
    attributes: Attribute[];
}>();

defineEmits<{
    (e: 'selected', entry: CalendarEntry): void;
}>();

const attributes = computed(() => ({
    past: props.attributes.filter((a: Attribute) => a.customData.isPast),
    upcoming: props.attributes.filter((a: Attribute) => !a.customData.isPast),
}));
</script>

<template>
    <div class="z-10 flex h-full flex-col overflow-hidden">
        <div class="day-label my-px inline-flex justify-between text-sm text-gray-900 dark:text-white">
            {{ day.day }}
            <UBadge v-if="day.isToday" size="xs" color="amber" :label="$t('common.today')" />
        </div>

        <div class="flex-grow overflow-x-auto overflow-y-auto">
            <UButton
                v-for="attr in attributes.past"
                :key="attr.key"
                class="vc-day-entry mb-1 mt-0 flex w-full flex-col !items-start justify-start rounded-sm p-1 text-left text-xs leading-tight"
                truncate
                :color="attr.customData.color"
                @click="$emit('selected', attr.customData)"
            >
                <span class="inline-flex items-center gap-0.5">
                    {{ attr.customData.title }}
                </span>

                <span v-if="attr.customData.time">
                    <template v-if="attr.customData.timeEnd && isSameDay(day.date, toDate(attr.customData.endTime))">
                        {{ attr.customData.timeEnd }}
                    </template>
                    <template v-else-if="isSameDay(day.date, toDate(attr.customData.startTime))">
                        {{ attr.customData.time }}
                    </template>
                </span>
            </UButton>

            <UDivider
                v-if="day.isToday && (attributes.past.length > 0 || attributes.upcoming.length > 0)"
                class="mb-1"
                size="sm"
                :ui="{ border: { base: 'border-red-300 dark:border-red-600' } }"
            />

            <UButton
                v-for="attr in attributes.upcoming"
                :key="attr.key"
                class="vc-day-entry mb-1 mt-0 flex w-full flex-col !items-start justify-start rounded-sm p-1 text-left text-xs leading-tight"
                truncate
                :color="attr.customData.color"
                @click="$emit('selected', attr.customData)"
            >
                <span class="inline-flex items-center gap-0.5">
                    <UIcon v-if="attr.customData.ongoing" class="size-3 text-amber-800" name="i-mdi-timer-sand" />
                    {{ attr.customData.title }}
                </span>

                <span v-if="attr.customData.time">
                    <template v-if="attr.customData.timeEnd && isSameDay(day.date, toDate(attr.customData.endTime))">
                        {{ attr.customData.timeEnd }}
                    </template>
                    <template v-else-if="isSameDay(day.date, toDate(attr.customData.startTime))">
                        {{ attr.customData.time }}
                    </template>
                </span>
            </UButton>
        </div>
    </div>
</template>
