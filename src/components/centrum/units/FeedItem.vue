<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccountPlus, mdiAccountRemove, mdiBriefcase, mdiCoffee, mdiHelp, mdiPlay, mdiStop } from '@mdi/js';
import Time from '~/components/partials/elements/Time.vue';
import { UNIT_STATUS, UnitStatus } from '~~/gen/ts/resources/dispatch/units';

defineProps<{
    activityLength: number;
    activityItem: UnitStatus;
    activityItemIdx: number;
    showName?: boolean; // TODO
}>();
</script>

<template>
    <li class="relative flex gap-x-2">
        <div
            :class="[
                activityItemIdx === activityLength - 1 ? 'h-6' : '-bottom-6',
                'absolute left-0 top-0 flex w-6 justify-center',
            ]"
        >
            <div class="w-px bg-gray-200" />
        </div>
        <template v-if="activityItem.status === UNIT_STATUS.UNKNOWN">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiHelp" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">Unit created</p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === UNIT_STATUS.USER_ADDED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiAccountPlus" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Colleague added to Unit
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === UNIT_STATUS.USER_REMOVED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiAccountRemove" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Colleague removed from Unit
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === UNIT_STATUS.UNAVAILABLE">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiStop" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Unit unavailable
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === UNIT_STATUS.AVAILABLE">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiPlay" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Unit available
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === UNIT_STATUS.ON_BREAK">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiCoffee" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Unit on break
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === UNIT_STATUS.BUSY">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiBriefcase" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Unit busy
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
    </li>
</template>
