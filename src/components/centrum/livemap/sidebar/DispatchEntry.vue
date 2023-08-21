<script lang="ts" setup>
import { CarEmergencyIcon } from 'mdi-vue3';
import { DISPATCH_STATUS, Dispatch } from '~~/gen/ts/resources/dispatch/dispatches';

defineProps<{
    dispatch: Dispatch;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
    (e: 'details', dsp: Dispatch): void;
}>();
</script>

<template>
    <li>
        <button
            type="button"
            class="text-white bg-error-700 hover:bg-primary-100/10 hover:text-neutral font-medium hover:transition-all group flex w-full flex-col items-center rounded-md p-2 text-xs my-0.5"
            @click="$emit('details', dispatch)"
        >
            <CarEmergencyIcon class="h-5 w-5" aria-hidden="true" />
            <span class="mt-2 truncate">DSP-{{ dispatch.id }}</span>
            <span class="mt-2 truncate">
                {{ $t(`enums.centrum.DISPATCH_STATUS.${DISPATCH_STATUS[dispatch.status?.status ?? (0 as number)]}`) }}
            </span>
        </button>
    </li>
</template>
