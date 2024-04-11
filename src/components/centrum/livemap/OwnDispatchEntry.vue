<script lang="ts" setup>
import { CarEmergencyIcon } from 'mdi-vue3';
import DispatchDetailsSlideover from '~/components/centrum/dispatches/DispatchDetailsSlideover.vue';
import { dispatchStatusToBGColor } from '~/components/centrum/helpers';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/store/livemap';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

defineProps<{
    dispatch: Dispatch;
    selectedDispatch: string | undefined;
}>();

defineEmits<{
    (e: 'update:selectedDispatch', dsp: string | undefined): void;
}>();

const { goto } = useLivemapStore();

const slideover = useSlideover();
</script>

<template>
    <li class="flex flex-row items-center">
        <div class="flex flex-col items-center gap-2">
            <input
                :value="dispatch.id"
                name="active"
                type="radio"
                class="form-radio focus-visible:ring-primary-500 dark:focus-visible:ring-primary-400 text-primary-500 dark:text-primary-400 h-4 w-4 border border-gray-300 bg-white focus:ring-0 focus:ring-transparent focus:ring-offset-transparent focus-visible:ring-1 focus-visible:ring-offset-2 focus-visible:ring-offset-white disabled:cursor-not-allowed disabled:opacity-50 dark:border-gray-700 dark:bg-gray-900 dark:checked:border-transparent dark:checked:bg-current dark:focus-visible:ring-offset-gray-900"
                :checked="selectedDispatch === dispatch.id"
                @change="$emit('update:selectedDispatch', dispatch.id)"
            />

            <UButton variant="link" icon="i-mdi-map-marker" @click="goto({ x: dispatch.x, y: dispatch.y })" />
        </div>
        <UButton
            color="red"
            :padded="false"
            class="my-0.5 flex w-full max-w-full shrink flex-col items-center p-2 text-xs"
            @click="
                slideover.open(DispatchDetailsSlideover, {
                    dispatch: dispatch,
                })
            "
        >
            <span class="mb-0.5 inline-flex w-full flex-col place-content-between items-center sm:flex-row sm:gap-1">
                <span class="inline-flex items-center font-bold md:gap-1">
                    <CarEmergencyIcon class="hidden h-3 w-auto md:block" />
                    DSP-{{ dispatch.id }}
                </span>
                <span>
                    <span class="font-semibold">{{ $t('common.postal') }}:</span> <span>{{ dispatch.postal }}</span>
                </span>
            </span>
            <span class="mb-0.5 inline-flex flex-col place-content-between items-center sm:flex-row sm:gap-1">
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
                        <span class="truncate"> {{ dispatch.creator.firstname }} {{ dispatch.creator.lastname }} </span>
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
