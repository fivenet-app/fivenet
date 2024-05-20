<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { StatusUnit, UnitStatus } from '~~/gen/ts/resources/centrum/units';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import { useLivemapStore } from '~/store/livemap';

defineProps<{
    activityLength: number;
    item: UnitStatus;
    activityItemIdx: number;
}>();

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
        <template v-if="item.status === StatusUnit.USER_ADDED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-account-plus" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.USER_ADDED') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" :initials-only="true" :badge="true" />
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
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>
        <template v-else-if="item.status === StatusUnit.USER_REMOVED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-account-remove" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.USER_REMOVED') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" :initials-only="true" :badge="true" />
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
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>
        <template v-else-if="item.status === StatusUnit.UNAVAILABLE">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-stop" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.UNAVAILABLE') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" :initials-only="true" :badge="true" />
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
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>
        <template v-else-if="item.status === StatusUnit.AVAILABLE">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-play" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.AVAILABLE') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" :initials-only="true" :badge="true" />
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
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>
        <template v-else-if="item.status === StatusUnit.ON_BREAK">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-coffee" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.ON_BREAK') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" :initials-only="true" :badge="true" />
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
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>
        <template v-else-if="item.status === StatusUnit.BUSY">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-briefcase" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.BUSY') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" :initials-only="true" :badge="true" />
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
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>
        <template v-else>
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-gray-300">
                <UIcon name="i-mdi-help" class="text-primary-500 size-5" />
            </div>
            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.UNKNOWN') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" :initials-only="true" :badge="true" />
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
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>
    </li>
</template>
