<script lang="ts" setup>
import DispatchDetailsSlideover from '~/components/centrum/dispatches/DispatchDetailsSlideover.vue';
import { dispatchStatusToBGColor, dispatchTimeToTextColorSidebar } from '~/components/centrum/helpers';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

const props = defineProps<{
    dispatch: Dispatch;
    selectedDispatch: number | undefined;
}>();

defineEmits<{
    (e: 'update:selectedDispatch', dsp: number | undefined): void;
}>();

const centrumStore = useCentrumStore();
const { settings } = storeToRefs(centrumStore);

const { goto } = useLivemapStore();

const slideover = useSlideover();

const dispatchTimeStyle = ref<{ ping: boolean; class: string }>({ ping: false, class: '' });

useIntervalFn(
    () =>
        (dispatchTimeStyle.value = dispatchTimeToTextColorSidebar(
            props.dispatch.createdAt,
            props.dispatch.status?.status,
            settings.value?.timings?.dispatchMaxWait,
        )),
    1000,
);
</script>

<template>
    <li class="flex flex-row items-center gap-1">
        <div class="flex flex-col items-center gap-2">
            <URadio
                :value="dispatch.id"
                name="active"
                :checked="selectedDispatch === dispatch.id"
                @change="$emit('update:selectedDispatch', dispatch.id)"
            />

            <UButton variant="link" :padded="false" icon="i-mdi-map-marker" @click="goto({ x: dispatch.x, y: dispatch.y })" />
        </div>

        <UChip
            :show="dispatchTimeStyle.ping"
            position="top-left"
            size="md"
            class="flex w-full max-w-full shrink flex-col items-center"
            :ui="{ base: dispatchTimeStyle.ping ? 'animate-pulse' : '', background: dispatchTimeStyle.class }"
        >
            <UButton
                color="red"
                :padded="false"
                class="my-0.5 inline-flex w-full max-w-full shrink flex-col items-center p-2 text-xs"
                @click="
                    slideover.open(DispatchDetailsSlideover, {
                        dispatchId: dispatch.id,
                    })
                "
            >
                <span class="mb-0.5 inline-flex w-full flex-col place-content-between items-center md:flex-row md:gap-1">
                    <span class="inline-flex items-center font-bold md:gap-1">
                        <UIcon name="i-mdi-car-emergency" class="hidden h-3 w-auto md:block" />
                        DSP-{{ dispatch.id }}
                    </span>
                    <span>
                        <span class="font-semibold">{{ $t('common.postal') }}:</span> {{ dispatch.postal }}
                    </span>
                </span>

                <span class="mb-0.5 inline-flex flex-col place-content-between items-center md:flex-row md:gap-1">
                    <span class="font-semibold">{{ $t('common.status') }}:</span>
                    <span class="line-clamp-2 break-words" :class="dispatchStatusToBGColor(dispatch.status?.status)">{{
                        $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`)
                    }}</span>
                </span>

                <span class="line-clamp-2 inline-flex flex-col md:flex-row md:gap-1">
                    <span class="font-semibold">{{ $t('common.sent_by') }}:</span>
                    <span class="truncate">
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

                <span class="inline-flex flex-col items-center md:flex-row md:gap-1">
                    <span class="font-semibold">{{ $t('common.sent_at') }}:</span>
                    <GenericTime :value="dispatch.createdAt" type="compact" />
                </span>
            </UButton>
        </UChip>
    </li>
</template>
