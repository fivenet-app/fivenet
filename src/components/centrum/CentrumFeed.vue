<script lang="ts" setup>
import DispatchFeedItem from '~/components/centrum/dispatches/DispatchFeedItem.vue';
import UnitFeedItem from '~/components/centrum/units/UnitFeedItem.vue';
import { DispatchStatus } from '~~/gen/ts/resources/centrum/dispatches';
import { UnitStatus } from '~~/gen/ts/resources/centrum/units';

defineProps<{
    items: (DispatchStatus | UnitStatus)[];
}>();
</script>

<template>
    <div class="flex size-full grow flex-col overflow-y-auto px-1">
        <div class="flex justify-between">
            <h2 class="text-base font-semibold leading-6 text-gray-100">
                {{ $t('common.activity', 2) }}
            </h2>
        </div>
        <div class="flex-1">
            <ul role="list" class="space-y-2">
                <template v-for="(activityItem, activityItemIdx) in items">
                    <DispatchFeedItem
                        v-if="'dispatchId' in activityItem"
                        :key="'dsp' + activityItem.id"
                        :activity-length="items?.length ?? 0"
                        :item="activityItem"
                        :activity-item-idx="activityItemIdx"
                        :show-id="true"
                    />
                    <UnitFeedItem
                        v-else
                        :key="'unit' + activityItem.id"
                        :activity-length="items?.length ?? 0"
                        :item="activityItem"
                        :activity-item-idx="activityItemIdx"
                    />
                </template>
            </ul>
        </div>
    </div>
</template>
