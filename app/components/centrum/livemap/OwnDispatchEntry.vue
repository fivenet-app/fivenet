<script lang="ts" setup>
import DispatchDetailsSlideover from '~/components/centrum/dispatches/DispatchDetailsSlideover.vue';
import { dispatchStatusToBadgeColor, dispatchTimeToTextColorSidebar } from '~/components/centrum/helpers';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useCentrumStore } from '~/stores/centrum';
import { useLivemapStore } from '~/stores/livemap';
import { type Dispatch, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

const props = defineProps<{
    dispatch: Dispatch;
}>();

const modelValue = defineModel<number | undefined>({ required: true });

const centrumStore = useCentrumStore();
const { settings } = storeToRefs(centrumStore);

const { goto } = useLivemapStore();

const overlay = useOverlay();

const dispatchDetailsSlideover = overlay.create(DispatchDetailsSlideover, {
    props: {
        dispatchId: props.dispatch.id,
    },
});

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
    <li class="my-1 flex flex-row items-center gap-1">
        <div class="flex flex-col items-center gap-2">
            <URadioGroup
                v-model="modelValue"
                name="active"
                :items="[dispatch.id]"
                :ui="{ label: 'hidden', item: 'items-end', wrapper: 'ms-0' }"
            />

            <UButton variant="link" icon="i-mdi-map-marker" @click="goto({ x: dispatch.x, y: dispatch.y })" />
        </div>

        <UChip
            class="flex w-full max-w-full shrink flex-col items-center"
            :show="dispatchTimeStyle.ping"
            position="top-left"
            size="md"
            :ui="{ base: dispatchTimeStyle.class + ' ' + (dispatchTimeStyle.ping ? 'animate-pulse' : '') }"
        >
            <UButton
                class="my-0.5 inline-flex w-full max-w-full shrink flex-col items-center p-2 text-xs"
                block
                color="error"
                @click="
                    dispatchDetailsSlideover.open({
                        dispatchId: dispatch.id,
                    })
                "
            >
                <!-- Row 1: ID + Postal -->
                <div class="flex w-full items-center justify-between">
                    <div class="flex items-center space-x-2 text-sm font-bold">
                        <Icon class="h-4 w-4" name="mdi-car-emergency" />
                        <span>DSP-{{ dispatch.id }}</span>
                    </div>
                    <div class="text-sm">
                        <span class="font-medium">{{ $t('common.postal') }}:</span>
                        <span>{{ dispatch.postal }}</span>
                    </div>
                </div>

                <!-- Row 2: Grid of Status & Sent By, plus full-width Sent At -->
                <div class="grid w-full grid-cols-2 gap-1 text-xs">
                    <div class="inline-flex flex-col items-center">
                        <span class="font-medium">{{ $t('common.status') }}:</span>
                        <UBadge
                            class="line-clamp-2 px-px py-0.5 break-words"
                            variant="solid"
                            :color="dispatchStatusToBadgeColor(dispatch.status?.status)"
                        >
                            {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`) }}
                        </UBadge>
                    </div>

                    <div class="inline-flex flex-col items-center">
                        <span class="font-medium">{{ $t('common.sent_by') }}:</span>
                        <span class="line-clamp-2 break-words">
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
                    </div>

                    <div class="col-span-2">
                        <span class="font-medium">{{ $t('common.sent_at') }}:</span>
                        <GenericTime class="ml-1" :value="dispatch.createdAt" type="compact" />
                    </div>
                </div>
            </UButton>
        </UChip>
    </li>
</template>
