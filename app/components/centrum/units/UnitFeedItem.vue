<script lang="ts" setup>
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/stores/livemap';
import { type UnitStatus, StatusUnit } from '~~/gen/ts/resources/centrum/units';

defineProps<{
    activityLength: number;
    item: UnitStatus;
    activityItemIdx: number;
}>();

const { gotoCoords } = useLivemapStore();
</script>

<template>
    <li class="relative flex gap-x-2">
        <div
            :class="[
                activityItemIdx === activityLength - 1 ? 'h-6' : '-bottom-6',
                'absolute top-0 left-0 flex w-6 justify-center',
            ]"
        >
            <div class="w-px bg-neutral-200" />
        </div>
        <template v-if="item.status === StatusUnit.USER_ADDED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-neutral-300">
                <UIcon class="size-5 text-primary-500" name="i-mdi-account-plus" />
            </div>

            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.USER_ADDED') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" initials-only badge />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="gotoCoords({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>

            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>

        <template v-else-if="item.status === StatusUnit.USER_REMOVED">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-neutral-300">
                <UIcon class="size-5 text-primary-500" name="i-mdi-account-remove" />
            </div>

            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.USER_REMOVED') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" initials-only badge />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="gotoCoords({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>

            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>

        <template v-else-if="item.status === StatusUnit.UNAVAILABLE">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-neutral-300">
                <UIcon class="size-5 text-primary-500" name="i-mdi-stop" />
            </div>

            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.UNAVAILABLE') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" initials-only badge />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="gotoCoords({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>

            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>

        <template v-else-if="item.status === StatusUnit.AVAILABLE">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-neutral-300">
                <UIcon class="size-5 text-primary-500" name="i-mdi-play" />
            </div>

            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.AVAILABLE') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" initials-only badge />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="gotoCoords({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>

            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>

        <template v-else-if="item.status === StatusUnit.ON_BREAK">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-neutral-300">
                <UIcon class="size-5 text-primary-500" name="i-mdi-coffee" />
            </div>

            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.ON_BREAK') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" initials-only badge />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="gotoCoords({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>

            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>

        <template v-else-if="item.status === StatusUnit.BUSY">
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-neutral-300">
                <UIcon class="size-5 text-primary-500" name="i-mdi-briefcase" />
            </div>

            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.BUSY') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" initials-only badge />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="gotoCoords({ x: item.x, y: item.y })"
                    />
                    <CitizenInfoPopover v-if="item.user" :user="item.user" :trailing="false" text-class="text-xs" />
                </span>
            </p>

            <span class="flex-none text-xs leading-5 text-gray-200">
                <GenericTime :value="item.createdAt" type="compact" />
            </span>
        </template>

        <template v-else>
            <div class="relative flex size-5 flex-none items-center justify-center rounded-lg bg-neutral-300">
                <UIcon class="size-5 text-primary-500" name="i-mdi-help" />
            </div>

            <p class="inline-flex flex-auto flex-row justify-between text-xs leading-5 text-gray-200">
                <span class="inline-flex items-center gap-1">
                    {{ $t('components.centrum.units.feed.item.UNKNOWN') }}

                    <UnitInfoPopover v-if="item.unit" :unit="item.unit" initials-only badge />
                </span>

                <span class="inline-flex items-center text-xs">
                    <UButton
                        v-if="item.x !== undefined && item.y !== undefined"
                        variant="link"
                        size="xs"
                        icon="i-mdi-map-marker"
                        @click="gotoCoords({ x: item.x, y: item.y })"
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
