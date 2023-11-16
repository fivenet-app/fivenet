<script lang="ts" setup>
import { DispatchStatus } from '~~/gen/ts/resources/dispatch/dispatches';
import { UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import DispatchFeedItem from '~/components/centrum/dispatches/DispatchFeedItem.vue';
import UnitFeedItem from '~/components/centrum/units/UnitFeedItem.vue';

defineProps<{
    items: (DispatchStatus | UnitStatus)[];
}>();
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full overflow-y-auto">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto inline-flex items-center">
                <h2 class="text-base font-semibold leading-6 text-gray-100">
                    {{ $t('common.activity', 2) }}
                </h2>
            </div>
        </div>
        <div class="mt-0.5 flow-root">
            <div class="-mx-2 -my-2 sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
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
        </div>
    </div>
</template>
