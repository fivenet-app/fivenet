<script lang="ts" setup>
import {
    AccountAlertIcon,
    AccountCancelIcon,
    AccountCheckIcon,
    AccountPlusIcon,
    AccountRemoveIcon,
    ArchiveIcon,
    CancelIcon,
    CarIcon,
    CheckIcon,
    HelpIcon,
    MapMarkerIcon,
    NewBoxIcon,
} from 'mdi-vue3';
import DispatchStatusInfoPopover from '~/components/centrum/dispatches/DispatchStatusInfoPopover.vue';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { DispatchStatus, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

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

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
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
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <NewBoxIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.NEW') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNASSIGNED">
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <AccountAlertIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNASSIGNED') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNIT_ASSIGNED">
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <AccountPlusIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNIT_ASSIGNED') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNIT_UNASSIGNED">
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <AccountRemoveIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNIT_UNASSIGNED') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNIT_ACCEPTED">
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <AccountCheckIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNIT_ACCEPTED') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNIT_DECLINED">
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <AccountCancelIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNIT_DECLINED') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.EN_ROUTE">
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <CarIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.EN_ROUTE') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.ON_SCENE">
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <MapMarkerIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.ON_SCENE') }}
                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.NEED_ASSISTANCE">
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <HelpIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.NEED_ASSISTANCE') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.COMPLETED">
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <CheckIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.COMPLETED') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.CANCELLED">
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <CancelIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.CANCELLED') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.ARCHIVED">
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <ArchiveIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.ARCHIVED') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else>
            <div class="relative flex h-5 w-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <NewBoxIcon class="h-5 w-5 text-primary-600" aria-hidden="true" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between py-0.5 text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNSPECIFIED') }}

                    <DispatchStatusInfoPopover v-if="showId" text-class="font-medium text-gray-400 pl-1" :status="item" />
                    <UnitInfoPopover
                        v-if="item.unit && item.unitId"
                        text-class="font-medium text-gray-400 pl-1"
                        :unit="item.unit"
                        :initials-only="true"
                        :badge="true"
                    />
                </span>

                <span class="inline-flex items-center">
                    <button v-if="item.x && item.y" type="button" @click="$emit('goto', { x: item.x, y: item.y })">
                        <MapMarkerIcon class="h-5 w-5 text-primary-400 hover:text-primary-600" />
                    </button>
                    <CitizenInfoPopover v-if="item.user" text-class="font-medium text-gray-400 pl-1" :user="item.user" />
                </span>
            </p>
            <span class="flex-none py-0.5 text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
    </li>
</template>
