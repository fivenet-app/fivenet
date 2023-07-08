<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import Time from '~/components/partials/elements/Time.vue';
import { DispatchStatus, DISPATCH_STATUS } from '~~/gen/ts/resources/dispatch/dispatch';
import {
    mdiAccountCancel,
    mdiAccountPlus,
    mdiAccountRemove,
    mdiCar,
    mdiCheck,
    mdiHelp,
    mdiMapMarker,
    mdiNewBox,
} from '@mdi/js';

defineProps<{
    activityLength: number;
    activityItem: DispatchStatus;
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
        <template v-if="activityItem.status === DISPATCH_STATUS.NEW">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiNewBox" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">Dispatch created</p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === DISPATCH_STATUS.DECLINED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiAccountCancel" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Dispatch declined by
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === DISPATCH_STATUS.UNASSIGNED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiAccountRemove" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Dispatch unassigned by
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === DISPATCH_STATUS.UNIT_ASSIGNED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiAccountPlus" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Dispatch accepted
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === DISPATCH_STATUS.EN_ROUTE">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiCar" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                En Route to Dispatch
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === DISPATCH_STATUS.AT_SCENE">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiMapMarker" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Arrived on scene at Dispatch
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === DISPATCH_STATUS.NEED_ASSISTANCE">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiHelp" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Need Assistance
                <span class="font-medium text-gray-400">
                    {{ activityItem.user?.firstname }}, {{ activityItem.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="activityItem.createdAt" />
            </span>
        </template>
        <template v-else-if="activityItem.status === DISPATCH_STATUS.COMPLETED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <SvgIcon type="mdi" :path="mdiCheck" class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Dispatch completed
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
