<script lang="ts" setup>
import { CarEmergencyIcon } from 'mdi-vue3';
import Details from '~/components/centrum/dispatches/Details.vue';
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import { Unit } from '~~/gen/ts/resources/dispatch/units';

defineProps<{
    dispatch: Dispatch;
    units: Unit[];
}>();

defineEmits<{
    (e: 'goto', location: { x: number; y: number }): void;
}>();

const detailsOpen = ref(false);
</script>

<template>
    <li>
        <button
            type="button"
            class="text-white bg-error-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
            @click="detailsOpen = true"
        >
            <CarEmergencyIcon class="h-5 w-5" aria-hidden="true" />
            <span class="mt-2 truncate">DSP-{{ dispatch.id }}</span>
            <span class="mt-2 truncate">
                {{ $t(`enums.centrum.DISPATCH_STATUS.${DISPATCH_STATUS[dispatch.status?.status ?? (0 as number)]}`) }}
            </span>
        </button>
        <Details
            @close="detailsOpen = false"
            :dispatch="dispatch"
            :open="detailsOpen"
            @goto="$emit('goto', $event)"
            :units="units"
        />
    </li>
</template>
