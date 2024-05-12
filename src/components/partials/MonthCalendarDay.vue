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
                truncate
                class="vc-day-entry mb-1 mt-0 w-full rounded-sm p-1 text-left text-xs leading-tight"
                :color="attr.customData.color"
                @click="$emit('selected', attr.customData)"
            >
                {{ attr.customData.title }}
                <template
                    v-if="
                        attr.customData.time &&
                        (isSameDay(day.date, toDate(attr.customData.startTime)) ||
                            isSameDay(day.date, toDate(attr.customData.endTime)))
                    "
                >
                    <br />
                    {{ attr.customData.time }}
                </template>
            </UButton>

            <UDivider
                v-if="day.isToday && (attributes.past.length > 0 || attributes.upcoming.length > 0)"
                size="sm"
                :ui="{ border: { base: 'border-red-300 dark:border-red-500' } }"
                class="mb-1"
            />

            <UButton
                v-for="attr in attributes.upcoming"
                :key="attr.key"
                truncate
                class="vc-day-entry mb-1 mt-0 w-full rounded-sm p-1 text-left text-xs leading-tight"
                :color="attr.customData.color"
                @click="$emit('selected', attr.customData)"
            >
                {{ attr.customData.title }}
                <template
                    v-if="
                        attr.customData.time &&
                        (isSameDay(day.date, toDate(attr.customData.startTime)) ||
                            isSameDay(day.date, toDate(attr.customData.endTime)))
                    "
                >
                    <br />
                    {{ attr.customData.time }}
                </template>
            </UButton>
        </div>
    </div>
</template>
