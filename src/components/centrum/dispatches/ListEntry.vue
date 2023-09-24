<script lang="ts" setup>
import { AccountMultiplePlusIcon, DotsVerticalIcon, MapMarkerIcon } from 'mdi-vue3';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import Time from '~/components/partials/elements/Time.vue';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { dispatchStatusAnimate, dispatchStatusToBGColor } from '../helpers';
import AssignDispatchModal from './AssignDispatchModal.vue';
import Details from './Details.vue';

const props = defineProps<{
    dispatch: Dispatch;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
    (e: 'details', dsp: Dispatch): void;
    (e: 'assignUnit', dsp: Dispatch): void;
}>();

const openDetails = ref(false);
const openAssign = ref(false);
</script>

<template>
    <tr>
        <Details :dispatch="dispatch" :open="openDetails" @close="openDetails = false" @goto="$emit('goto', $event)" />
        <AssignDispatchModal v-if="openAssign" :open="openAssign" :dispatch="dispatch" @close="openAssign = false" />

        <td
            class="relative whitespace-nowrap pl-0 py-1 pr-0 text-left text-sm font-medium sm:pr-0.5 flex flex-row justify-start"
        >
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.assign')"
                @click="openAssign = true"
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
                @click="openDetails = true"
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
        <td
            class="whitespace-nowrap px-1 py-1 text-sm text-gray-100"
            :class="dispatchStatusToBGColor(props.dispatch.status?.status ?? 0)"
        >
            <span :class="dispatchStatusAnimate(props.dispatch.status?.status ?? 0) ? 'animate-pulse' : ''">
                {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[props.dispatch.status?.status ?? 0]}`) }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            {{ dispatch.postal ?? 'N/A' }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <span v-if="dispatch.units.length === 0" class="italic">{{ $t('enums.centrum.StatusDispatch.UNASSIGNED') }}</span>
            <span v-else class="mr-1">
                {{ dispatch.units.map((unit) => unit.unit?.initials ?? 'N/A').join(', ') }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <span v-if="dispatch.anon">
                {{ $t('common.anon') }}
            </span>
            <span v-else-if="dispatch.user">
                <CitizenInfoPopover :user="dispatch.user" />
            </span>
            <span v-else>
                {{ $t('common.unknown') }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300 truncate">{{ dispatch.message }}</td>
    </tr>
</template>
