<script lang="ts" setup>
import { DispatchStatus } from '~~/gen/ts/resources/dispatch/dispatch';
import { UnitStatus } from '~~/gen/ts/resources/dispatch/units';
import { default as DispatchFeedItem } from './dispatches/FeedItem.vue';
import { default as UnitFeedItem } from './units/FeedItem.vue';

defineProps<{
    items: (DispatchStatus | UnitStatus)[];
}>();
</script>

<template>
    <div class="px-4 sm:px-6 lg:px-8 h-full">
        <div class="sm:flex sm:items-center">
            <div class="sm:flex-auto">
                <h1 class="text-base font-semibold leading-6 text-gray-100">Feed</h1>
            </div>
        </div>
        <div class="mt-2 flow-root">
            <div class="-mx-2 -my-2 overflow-y-scroll sm:-mx-6 lg:-mx-8">
                <div class="inline-block min-w-full py-2 align-middle sm:px-2 lg:px-2">
                    <ul role="list" class="space-y-2">
                        <template v-for="(activityItem, activityItemIdx) in items">
                            <template v-if="'dispatchId' in activityItem">
                                <DispatchFeedItem
                                    :activityLength="items?.length ?? 0"
                                    :activityItem="activityItem"
                                    :activityItemIdx="activityItemIdx"
                                    :showName="true"
                                />
                            </template>
                            <template v-else>
                                <UnitFeedItem
                                    :activityLength="items?.length ?? 0"
                                    :activityItem="activityItem"
                                    :activityItemIdx="activityItemIdx"
                                    :showName="true"
                                />
                            </template>
                        </template>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
