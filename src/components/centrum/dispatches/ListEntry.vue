<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiAccountMultiplePlus, mdiDetails, mdiMapMarker } from '@mdi/js';
import Time from '~/components/partials/elements/Time.vue';
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatch';
import { Unit } from '~~/gen/ts/resources/dispatch/units';
import AssignDispatchModal from './AssignDispatchModal.vue';
import Details from './Details.vue';

defineProps<{
    dispatch: Dispatch;
    units: Unit[] | null;
}>();

defineEmits<{
    (e: 'goto', location: { x: number; y: number }): void;
}>();

const detailsOpen = ref(false);
const assignOpen = ref(false);
</script>

<template>
    <tr>
        <Details @close="detailsOpen = false" :dispatch="dispatch" :open="detailsOpen" />
        <AssignDispatchModal @close="assignOpen = false" :dispatch="dispatch" :units="units" :open="assignOpen" />
        <td
            class="relative whitespace-nowrap py-2 pl-2 text-right text-sm font-medium sm:pr-0 max-w-[42px] flex flex-row justify-start"
        >
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.detail', 2)"
                @click="detailsOpen = true"
            >
                <SvgIcon type="mdi" :path="mdiDetails" class="w-6 h-auto ml-auto mr-2.5" />
            </button>
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.assign')"
                @click="assignOpen = true"
            >
                <SvgIcon type="mdi" :path="mdiAccountMultiplePlus" class="w-6 h-auto ml-auto mr-2.5" aria-hidden="true" />
            </button>
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
            {{ dispatch.id }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm font-medium text-gray-100">
            {{ DISPATCH_STATUS[dispatch.status?.status as number] }}
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
            <Time :value="dispatch.createdAt" />
        </td>
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm text-gray-300 sm:pl-0">
            <span v-if="dispatch.units.length === 0"> No Units assigned. </span>
            <span v-else v-for="unit in dispatch.units" class="mr-1">
                {{ (units ?? []).find((u) => u.id === unit.unitId)?.initials }}
            </span>
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
            <button
                type="button"
                class="inline-flex items-center rounded bg-white/10 px-2 py-1 text-xs font-semibold text-white shadow-sm hover:bg-white/20"
                @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
            >
                Go to
                <SvgIcon type="mdi" :path="mdiMapMarker" class="-mr-0.5 h-4 w-4" aria-hidden="true" />
            </button>
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-100">{{ dispatch.message }}</td>
    </tr>
</template>
