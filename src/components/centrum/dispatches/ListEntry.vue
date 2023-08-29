<script lang="ts" setup>
import { AccountMultiplePlusIcon, DotsVerticalIcon, MapMarkerIcon } from 'mdi-vue3';
import Time from '~/components/partials/elements/Time.vue';
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { dispatchStatusAnimate, dispatchStatusToBGColor } from '../helpers';

const props = defineProps<{
    dispatch: Dispatch;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
    (e: 'details', dsp: Dispatch): void;
    (e: 'assignUnit', dsp: Dispatch): void;
    (e: 'status', dsp: Dispatch): void;
}>();

const status = computed(() => props.dispatch.status?.status ?? 0);
</script>

<template>
    <tr>
        <td
            class="relative whitespace-nowrap pl-0 py-1 pr-0 text-left text-sm font-medium sm:pr-0.5 flex flex-row justify-start"
        >
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.assign')"
                @click="$emit('assignUnit', dispatch)"
            >
                <AccountMultiplePlusIcon class="w-6 h-auto ml-auto mr-1.5" aria-hidden="true" />
            </button>
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
            >
                <MapMarkerIcon class="w-6 h-auto ml-auto mr-1.5" aria-hidden="true" />
            </button>
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.detail', 2)"
                @click="$emit('details', dispatch)"
            >
                <DotsVerticalIcon class="w-6 h-auto ml-auto mr-1.5" />
            </button>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            {{ dispatch.id }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <Time :value="dispatch.createdAt" type="compact" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-100" :class="dispatchStatusToBGColor(status)">
            <span :class="dispatchStatusAnimate(status) ? 'animate-pulse' : ''">
                {{ $t(`enums.centrum.DISPATCH_STATUS.${DISPATCH_STATUS[status]}`) }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <span v-if="dispatch.units.length === 0" class="italic">{{ $t('enums.centrum.DISPATCH_STATUS.UNASSIGNED') }}</span>
            <span v-else class="mr-1">
                {{ dispatch.units.map((unit) => unit.unit?.initials ?? 'N/A').join(', ') }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <span v-if="!dispatch.anon && dispatch.user"> {{ dispatch.user.firstname }}, {{ dispatch.user.lastname }} </span>
            <span v-else-if="dispatch.anon">
                {{ $t('common.anon') }}
            </span>
            <span v-else>
                {{ $t('common.unknown') }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300 truncate">{{ dispatch.message }}</td>
    </tr>
</template>
