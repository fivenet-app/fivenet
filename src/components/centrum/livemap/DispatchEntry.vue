<script lang="ts" setup>
import { CarEmergencyIcon, MapMarkerIcon } from 'mdi-vue3';
import DispatchDetails from '~/components/centrum/dispatches/DispatchDetails.vue';
import { dispatchStatusToBGColor } from '~/components/centrum/helpers';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

defineProps<{
    dispatch: Dispatch;
    selectedDispatch: string | undefined;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
    (e: 'update:selectedDispatch', dsp: string | undefined): void;
}>();

const openDetails = ref(false);
</script>

<template>
    <li class="flex flex-row items-center">
        <DispatchDetails :dispatch="dispatch" :open="openDetails" @close="openDetails = false" @goto="$emit('goto', $event)" />

        <div class="mr-1.5 flex flex-col items-center gap-2">
            <input
                :value="dispatch.id"
                name="active"
                type="radio"
                class="h-4 w-4 border-gray-300 text-primary-600 focus:ring-primary-600"
                :checked="selectedDispatch === dispatch.id"
                @change="$emit('update:selectedDispatch', dispatch.id)"
            />

            <button
                type="button"
                class="inline-flex items-center text-primary-400 hover:text-primary-600"
                @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
            >
                <MapMarkerIcon class="h-5 w-5" aria-hidden="true" />
            </button>
        </div>
        <button
            type="button"
            class="group my-0.5 flex w-full flex-col items-center rounded-md bg-error-700 p-2 text-xs font-medium text-neutral hover:bg-primary-100/10 hover:text-neutral hover:transition-all"
            @click="openDetails = true"
        >
            <span class="mb-0.5 inline-flex w-full place-content-between items-center flex-col sm:flex-row">
                <span class="truncate font-bold"> DSP-{{ dispatch.id }} </span>
                <span>
                    <CarEmergencyIcon class="h-4 w-4 hidden sm:block" />
                </span>
                <span>
                    <span class="font-semibold">{{ $t('common.postal') }}:</span> {{ dispatch.postal }}
                </span>
            </span>
            <span class="truncate inline-flex flex-col sm:flex-row sm:gap-1">
                <span class="font-semibold">{{ $t('common.status') }}:</span>
                <span :class="dispatchStatusToBGColor(dispatch.status?.status)">{{
                    $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`)
                }}</span>
            </span>
            <span class="mt-1 truncate inline-flex flex-col sm:flex-row sm:gap-1">
                <span class="font-semibold">{{ $t('common.sent_by') }}:</span>
                <span v-if="dispatch.anon">
                    {{ $t('common.anon') }}
                </span>
                <template v-else-if="dispatch.creator">
                    {{ dispatch.creator.firstname }}, {{ dispatch.creator.lastname }}
                </template>
                <span v-else>
                    {{ $t('common.unknown') }}
                </span>
            </span>
            <span class="">
                <span class="font-semibold">{{ $t('common.sent_at') }}:</span>
                <GenericTime :value="dispatch.createdAt" type="compact" />
            </span>
        </button>
    </li>
</template>
