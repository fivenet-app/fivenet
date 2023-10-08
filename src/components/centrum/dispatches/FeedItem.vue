<script lang="ts" setup>
import {
    AccountAlertIcon,
    AccountCancelIcon,
    AccountCheckIcon,
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

withDefaults(
    defineProps<{
        activityLength: number;
        item: DispatchStatus;
        activityItemIdx: number;
        showId?: boolean;
    }>(),
    {
        showId: false,
    },
);
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
            <p class="flex-auto flex items-center py-0.5 text-xs leading-5 text-gray-200">
                {{ $t('components.centrum.dispatches.feed.item.NEW') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
                {{ $t('components.centrum.dispatches.feed.item.UNASSIGNED') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
                {{ $t('components.centrum.dispatches.feed.item.UNIT_ASSIGNED') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
                {{ $t('components.centrum.dispatches.feed.item.UNIT_UNASSIGNED') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
        <template v-else-if="item.status === StatusDispatch.UNIT_ACCEPTED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <AccountCheckIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                {{ $t('components.centrum.dispatches.feed.item.UNIT_ACCEPTED') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
        <template v-else-if="item.status === StatusDispatch.UNIT_DECLINED">
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <AccountCancelIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                {{ $t('components.centrum.dispatches.feed.item.UNIT_DECLINED') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
                {{ $t('components.centrum.dispatches.feed.item.EN_ROUTE') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
                {{ $t('components.centrum.dispatches.feed.item.ON_SCENE') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
                {{ $t('components.centrum.dispatches.feed.item.NEED_ASSISTANCE') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
                {{ $t('components.centrum.dispatches.feed.item.COMPLETED') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
                {{ $t('components.centrum.dispatches.feed.item.CANCELLED') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
                {{ $t('components.centrum.dispatches.feed.item.ARCHIVED') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
        <template v-else>
            <div class="relative flex h-6 w-6 flex-none items-center justify-center bg-gray-300 rounded-lg">
                <NewBoxIcon class="h-6 w-6 text-primary-600" aria-hidden="true" />
            </div>
            <p class="flex-auto py-0.5 text-xs leading-5 text-gray-200">
                {{ $t('components.centrum.dispatches.feed.item.UNSPECIFIED') }}

                <span v-if="showId" class="font-medium text-gray-400 pl-1">(DSP-{{ item.dispatchId }})</span>
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
