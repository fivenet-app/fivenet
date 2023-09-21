<script lang="ts" setup>
import {
    AccountAlertIcon,
    AccountCancelIcon,
    AccountPlusIcon,
    AccountRemoveIcon,
    ArchiveIcon,
    CarIcon,
    CheckIcon,
    HelpIcon,
    MapMarkerIcon,
    NewBoxIcon,
} from 'mdi-vue3';
import Time from '~/components/partials/elements/Time.vue';
import { DispatchStatus, StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';

defineProps<{
    activityLength: number;
    item: DispatchStatus;
    activityItemIdx: number;
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
        <template v-if="item.status === StatusDispatch.NEW">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <NewBoxIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Dispatch created
                <span class="font-medium text-gray-400 pl-1" v-if="item.user">
                    {{ item.user?.firstname }}, {{ item.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNASSIGNED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <AccountAlertIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                No units assigned
                <span class="font-medium text-gray-400 pl-1" v-if="item.unit">
                    {{ item.unit?.initials }}
                </span>
                <span class="font-medium text-gray-400 pl-1" v-if="item.user">
                    {{ item.user?.firstname }}, {{ item.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNIT_ASSIGNED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <AccountPlusIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Unit assigned to Dispatch
                <span class="font-medium text-gray-400 pl-1" v-if="item.unit">
                    {{ item.unit?.initials }}
                </span>
                <span class="font-medium text-gray-400 pl-1" v-if="item.user">
                    {{ item.user?.firstname }}, {{ item.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNIT_UNASSIGNED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <AccountRemoveIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Unit unassigned from Dispatch
                <span class="font-medium text-gray-400 pl-1" v-if="item.unit">
                    {{ item.unit?.initials }}
                </span>
                <span class="font-medium text-gray-400 pl-1" v-if="item.user">
                    {{ item.user?.firstname }}, {{ item.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.EN_ROUTE">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <CarIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                En Route to Dispatch
                <span class="font-medium text-gray-400 pl-1" v-if="item.unit">
                    {{ item.unit?.initials }}
                </span>
                <span class="font-medium text-gray-400 pl-1" v-if="item.user">
                    {{ item.user?.firstname }}, {{ item.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.ON_SCENE">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <MapMarkerIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Arrived on scene
                <span class="font-medium text-gray-400 pl-1" v-if="item.unit">
                    {{ item.unit?.initials }}
                </span>
                <span class="font-medium text-gray-400 pl-1" v-if="item.user">
                    {{ item.user?.firstname }}, {{ item.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.NEED_ASSISTANCE">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <HelpIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Need Assistance
                <span class="font-medium text-gray-400 pl-1" v-if="item.unit">
                    {{ item.unit?.initials }}
                </span>
                <span class="font-medium text-gray-400 pl-1" v-if="item.user">
                    {{ item.user?.firstname }}, {{ item.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.COMPLETED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <CheckIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Dispatch completed
                <span class="font-medium text-gray-400 pl-1" v-if="item.unit">
                    {{ item.unit?.initials }}
                </span>
                <span class="font-medium text-gray-400 pl-1" v-if="item.user">
                    {{ item.user?.firstname }}, {{ item.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.CANCELLED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <AccountCancelIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Dispatch cancelled
                <span class="font-medium text-gray-400 pl-1" v-if="item.unit">
                    {{ item.unit?.initials }}
                </span>
                <span class="font-medium text-gray-400 pl-1" v-if="item.user">
                    {{ item.user?.firstname }}, {{ item.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.ARCHIVED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <ArchiveIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                Dispatch archived
                <span class="font-medium text-gray-400 pl-1" v-if="item.unit">
                    {{ item.unit?.initials }}
                </span>
                <span class="font-medium text-gray-400 pl-1" v-if="item.user">
                    {{ item.user?.firstname }}, {{ item.user?.lastname }}
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <Time :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
    </li>
</template>
