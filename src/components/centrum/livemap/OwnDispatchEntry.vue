<script lang="ts" setup>
import { CarEmergencyIcon, MapMarkerIcon } from 'mdi-vue3';
import DispatchDetailsSlideover from '~/components/centrum/dispatches/DispatchDetailsSlideover.vue';
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
        <DispatchDetailsSlideover
            :dispatch="dispatch"
            :open="openDetails"
            @close="openDetails = false"
            @goto="$emit('goto', $event)"
        />

        <div class="mr-1.5 flex flex-col items-center gap-2">
            <UInput
                :value="dispatch.id"
                name="active"
                type="radio"
                class="size-4 border-gray-300 text-primary-600 focus:ring-primary-600"
                :checked="selectedDispatch === dispatch.id"
                @change="$emit('update:selectedDispatch', dispatch.id)"
            />

            <UButton
                class="inline-flex items-center text-primary-400 hover:text-primary-600"
                @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
            >
                <MapMarkerIcon class="size-5" />
            </UButton>
        </div>
        <UButton
            class="group my-0.5 flex w-full max-w-full flex-col items-center rounded-md bg-error-700 p-2 text-xs font-medium hover:bg-primary-100/10 hover:transition-all"
            @click="openDetails = true"
        >
            <span class="mb-0.5 flex w-full flex-col place-content-between items-center sm:flex-row sm:gap-1">
                <span class="inline-flex items-center font-bold md:gap-1">
                    <CarEmergencyIcon class="hidden h-3 w-auto md:block" />
                    DSP-{{ dispatch.id }}
                </span>
                <span>
                    <span class="font-semibold">{{ $t('common.postal') }}:</span> <span>{{ dispatch.postal }}</span>
                </span>
            </span>
            <span class="mb-0.5 flex flex-col place-content-between items-center sm:flex-row sm:gap-1">
                <span class="font-semibold">{{ $t('common.status') }}:</span>
                <span class="line-clamp-2 break-words" :class="dispatchStatusToBGColor(dispatch.status?.status)">{{
                    $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`)
                }}</span>
            </span>
            <span class="line-clamp-2 inline-flex flex-col sm:flex-row sm:gap-1">
                <span class="font-semibold">{{ $t('common.sent_by') }}:</span>
                <span>
                    <template v-if="dispatch.anon">
                        {{ $t('common.anon') }}
                    </template>
                    <template v-else-if="dispatch.creator">
                        {{ dispatch.creator.firstname }} {{ dispatch.creator.lastname }}
                    </template>
                    <template v-else>
                        {{ $t('common.unknown') }}
                    </template>
                </span>
            </span>
            <span class="inline-flex flex-col sm:flex-row sm:gap-1">
                <span class="font-semibold">{{ $t('common.sent_at') }}:</span>
                <GenericTime :value="dispatch.createdAt" type="compact" />
            </span>
        </UButton>
    </li>
</template>
