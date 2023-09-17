<script lang="ts" setup>
import { CarEmergencyIcon } from 'mdi-vue3';
import { default as DispatchDetails } from '~/components/centrum/dispatches/Details.vue';
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';

defineProps<{
    dispatch: Dispatch;
    selectedDispatch: bigint | undefined;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
    (e: 'update:selectedDispatch', dsp: bigint | undefined): void;
}>();

const openDetails = ref(false);
</script>

<template>
    <li class="flex flex-row items-center">
        <DispatchDetails :dispatch="dispatch" :open="openDetails" @close="openDetails = false" @goto="$emit('goto', $event)" />

        <div class="mr-1.5">
            <input
                name="active"
                type="radio"
                class="h-4 w-4 border-gray-300 text-primary-600 focus:ring-primary-600"
                v-bind:value="dispatch.id"
                @change="$emit('update:selectedDispatch', dispatch.id)"
            />
        </div>
        <button
            type="button"
            class="text-white bg-error-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
            @click="openDetails = true"
        >
            <span class="font-bold truncate inline-flex items-center">
                <CarEmergencyIcon class="h-4 w-4 mr-0.5" /> DSP-{{ dispatch.id }}</span
            >
            <span class="mt-2 truncate">
                {{ $t(`enums.centrum.DISPATCH_STATUS.${DISPATCH_STATUS[dispatch.status?.status ?? (0 as number)]}`) }}
            </span>
        </button>
    </li>
</template>
