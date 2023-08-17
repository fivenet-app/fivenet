<script lang="ts" setup>
import { AccountMultiplePlusIcon, DotsVerticalIcon, MapMarkerIcon } from 'mdi-vue3';
import Time from '~/components/partials/elements/Time.vue';
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import AssignDispatchModal from './AssignDispatchModal.vue';
import Details from './Details.vue';

defineProps<{
    dispatch: Dispatch;
    units: Unit[];
}>();

defineEmits<{
    (e: 'goto', loc: { x: number; y: number }): void;
}>();

const detailsOpen = ref(false);
const assignOpen = ref(false);
</script>

<template>
    <tr>
        <Details
            @close="detailsOpen = false"
            :dispatch="dispatch"
            :units="units"
            :open="detailsOpen"
            @goto="$emit('goto', $event)"
        />
        <AssignDispatchModal @close="assignOpen = false" :dispatch="dispatch" :units="units" :open="assignOpen" />
        <td class="relative whitespace-nowrap py-2 pl-2 text-right text-sm font-medium sm:pr-0 flex flex-row justify-start">
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.assign')"
                @click="assignOpen = true"
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
                @click="detailsOpen = true"
            >
                <DotsVerticalIcon class="w-6 h-auto ml-auto mr-1.5" />
            </button>
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
            {{ dispatch.id }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
            <Time :value="dispatch.createdAt" type="compact" />
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm font-medium text-gray-100">
            {{ $t(`enums.centrum.DISPATCH_STATUS.${DISPATCH_STATUS[dispatch.status?.status ?? (0 as number)]}`) }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
            <span v-if="dispatch.units.length === 0" class="italic">{{ $t('enums.centrum.DISPATCH_STATUS.UNASSIGNED') }}</span>
            <span v-else class="mr-1">
                {{ dispatch.units.map((unit) => unit.unit?.initials ?? 'N/A').join(', ') }}
            </span>
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
            <span v-if="!dispatch.anon && dispatch.user"> {{ dispatch.user.firstname }}, {{ dispatch.user.lastname }} </span>
            <span v-else-if="dispatch.anon">
                {{ $t('common.anon') }}
            </span>
            <span v-else>
                {{ $t('common.unknown') }}
            </span>
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-100 truncate">{{ dispatch.message }}</td>
    </tr>
</template>
