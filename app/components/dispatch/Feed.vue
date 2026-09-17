<script lang="ts" setup>
import DispatchFeedItem from '~/components/dispatch/dispatches/DispatchFeedItem.vue';
import UnitFeedItem from '~/components/dispatch/units/UnitFeedItem.vue';
import { dispatchStatuses, dispatchStatusToBGColor, unitStatuses, unitStatusToBGColor } from '~/components/dispatch/helpers';
import type { DispatchStatus } from '~~/gen/ts/resources/centrum/dispatches/dispatches';
import type { UnitStatus } from '~~/gen/ts/resources/centrum/units/units';

defineProps<{
    items: (DispatchStatus | UnitStatus)[];
}>();

function timelineIcon(item: DispatchStatus | UnitStatus): string {
    if ('dispatchId' in item) {
        return dispatchStatuses.find((status) => status.status === item.status)?.icon ?? 'i-mdi-info-circle';
    }

    return unitStatuses.find((status) => status.status === item.status)?.icon ?? 'i-mdi-info-circle';
}

function timelineColor(item: DispatchStatus | UnitStatus): string {
    return `text-highlighted ${'dispatchId' in item ? dispatchStatusToBGColor(item.status) : unitStatusToBGColor(item.status)}`;
}
</script>

<template>
    <div class="flex size-full grow flex-col overflow-y-auto px-1">
        <div class="flex justify-between">
            <h2 class="text-base leading-6 font-semibold text-toned">
                {{ $t('common.activity', 2) }}
            </h2>
        </div>
        <div class="flex-1">
            <UTimeline
                :items="items.map((item) => ({ ...item, icon: timelineIcon(item), ui: { indicator: timelineColor(item) } }))"
                size="xs"
                :ui="{ wrapper: '!mt-0 !pb-2' }"
            >
                <template #wrapper="{ item }">
                    <DispatchFeedItem v-if="'dispatchId' in item" :item="item" show-id />
                    <UnitFeedItem v-else :item="item" />
                </template>
            </UTimeline>
        </div>
    </div>
</template>
