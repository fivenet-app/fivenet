<script lang="ts" setup>
import DispatchStatusInfoPopover from '~/components/centrum/dispatches/DispatchStatusInfoPopover.vue';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/store/livemap';
import type { DispatchStatus} from '~~/gen/ts/resources/centrum/dispatches';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

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

const { goto } = useLivemapStore();
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
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-new-box" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.NEW') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNASSIGNED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-account-alert" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNASSIGNED') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNIT_ASSIGNED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-account-plus" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNIT_ASSIGNED') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNIT_UNASSIGNED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-account-remove" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNIT_UNASSIGNED') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNIT_ACCEPTED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-account-check" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNIT_ACCEPTED') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.UNIT_DECLINED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-account-cancel" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNIT_DECLINED') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.EN_ROUTE">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-car" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.EN_ROUTE') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.ON_SCENE">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-map-marker" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.ON_SCENE') }}
                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.NEED_ASSISTANCE">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-help" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.NEED_ASSISTANCE') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.COMPLETED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-check" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.COMPLETED') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.CANCELLED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-cancel" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.CANCELLED') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else-if="item.status === StatusDispatch.ARCHIVED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-archive" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.ARCHIVED') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
        <template v-else>
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-new-box" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.dispatches.feed.item.UNSPECIFIED') }}

                    <DispatchStatusInfoPopover v-if="showId" :status="item" />
                    <UnitInfoPopover v-if="item.unit && item.unitId" :unit="item.unit" :initials-only="true" :badge="true" />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="goto({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>
            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" :type="'compact'" />
            </span>
        </template>
    </li>
</template>
