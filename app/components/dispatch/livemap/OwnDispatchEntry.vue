<script lang="ts" setup>
import DispatchDetailsSlideover from '~/components/dispatch/dispatches/DispatchDetailsSlideover.vue';
import { dispatchTimeToTextColorSidebar } from '~/components/dispatch/helpers';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useCentrumStore } from '~/stores/centrum';
import { useLivemapStore } from '~/stores/livemap';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches/dispatches';
import DispatchStatusBadge from '~/components/dispatch/partials/DispatchStatusBadge.vue';

const props = defineProps<{
    dispatch: Dispatch;
    now: Date;
}>();

const modelValue = defineModel<number | undefined>({ required: true });

const centrumStore = useCentrumStore();
const { settings } = storeToRefs(centrumStore);

const { gotoCoords } = useLivemapStore();

const overlay = useOverlay();

const dispatchDetailsSlideover = overlay.create(DispatchDetailsSlideover, {
    props: {
        dispatchId: props.dispatch.id,
    },
});

const dispatchTimeStyle = computed(() =>
    dispatchTimeToTextColorSidebar(
        props.dispatch.createdAt,
        props.dispatch.status?.status,
        settings.value?.timings?.dispatchMaxWait,
        props.now.getTime(),
    ),
);
</script>

<template>
    <li class="my-1 flex items-stretch -space-x-px">
        <div class="grid grid-rows-[minmax(0,1fr)_minmax(0,1fr)] -space-y-px self-stretch">
            <UTooltip class="h-full min-h-0" :text="$t('common.select')">
                <URadioGroup
                    v-model="modelValue"
                    name="active"
                    variant="card"
                    size="sm"
                    color="primary"
                    :items="[
                        {
                            value: dispatch.id,
                        },
                    ]"
                    value-key="value"
                    :ui="{
                        root: 'h-full',
                        fieldset: 'h-full',
                        item: 'h-full w-7 cursor-pointer items-center justify-center rounded-tl-md rounded-bl-none rounded-r-none border-r-0 p-0',
                        container: 'h-full',
                        wrapper: 'm-0 flex h-full items-center justify-center',
                    }"
                />
            </UTooltip>

            <UTooltip class="h-full min-h-0" :text="$t('common.goto')">
                <UButton
                    block
                    class="h-full rounded-t-none rounded-r-none border-t-0 border-r-0 p-0"
                    color="neutral"
                    icon="i-mdi-map-marker"
                    variant="subtle"
                    @click="gotoCoords({ x: dispatch.x, y: dispatch.y })"
                />
            </UTooltip>
        </div>

        <UChip
            class="min-w-0 flex-1 self-stretch"
            :show="dispatchTimeStyle.ping"
            position="top-left"
            size="md"
            :ui="{ base: dispatchTimeStyle.class + ' ' + (dispatchTimeStyle.ping ? 'animate-pulse' : '') }"
        >
            <UButton
                class="h-full w-full flex-col items-center rounded-l-none p-2 text-xs"
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
                <div class="grid w-full min-w-0 grid-cols-2 gap-1 text-xs">
                    <div class="inline-flex w-full min-w-0 flex-col items-center">
                        <span class="font-medium">{{ $t('common.status') }}:</span>

                        <DispatchStatusBadge :status="dispatch.status?.status" class="max-w-full min-w-0 justify-center" />
                    </div>

                    <div class="inline-flex flex-col items-center">
                        <span class="font-medium">{{ $t('common.sent_by') }}:</span>
                        <span class="line-clamp-2 break-all">
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
