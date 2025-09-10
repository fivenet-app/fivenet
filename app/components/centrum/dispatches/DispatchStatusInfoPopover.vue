<script lang="ts" setup>
import DispatchDetailsByIDSlideover from '~/components/centrum/dispatches//DispatchDetailsByIDSlideover.vue';
import { useCentrumStore } from '~/stores/centrum';
import { type DispatchStatus, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { dispatchStatusToBGColor } from '../helpers';

const props = defineProps<{
    status: DispatchStatus | undefined;
}>();

const overlay = useOverlay();

const centrumStore = useCentrumStore();

const dispatchDetailsByIDSlideover = overlay.create(DispatchDetailsByIDSlideover, {
    props: { dispatchId: props.status?.dispatchId ?? 0 },
});

const dispatch = props.status?.dispatchId ? centrumStore.dispatches.get(props.status.dispatchId) : undefined;
const dispatchStatusColor = computed(() => dispatchStatusToBGColor(props.status?.status));
</script>

<template>
    <template v-if="!status">
        <span class="inline-flex items-center">
            <slot name="before" />
            <span>{{ $t('common.na') }}</span>
            <slot name="after" />
        </span>
    </template>
    <UPopover v-else>
        <UButton class="inline-flex items-center p-0.5" variant="outline" size="xs" trailing-icon="i-mdi-chevron-down">
            <slot name="before" />
            <span> DSP-{{ status.dispatchId }} </span>
            <slot name="after" />
        </UButton>

        <template #content>
            <div class="inline-flex min-w-48 flex-col gap-1 p-4">
                <div class="flex items-center gap-2">
                    <UTooltip :text="$t('common.detail', 2)">
                        <UButton
                            variant="link"
                            icon="i-mdi-car-emergency"
                            @click="
                                dispatchDetailsByIDSlideover.open({
                                    dispatchId: status.dispatchId,
                                })
                            "
                        >
                            {{ $t('common.detail', 2) }}
                        </UButton>
                    </UTooltip>
                </div>

                <p class="text-base leading-none font-semibold text-highlighted">DSP-{{ status.dispatchId }}</p>

                <UBadge class="rounded-sm font-semibold" :class="dispatchStatusColor" size="xs">
                    {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[status.status ?? 0]}`) }}
                </UBadge>

                <div v-if="dispatch" class="text-highlighted">
                    <p class="text-sm leading-none font-medium">
                        {{ $t('common.unit', 2) }}
                    </p>

                    <template v-if="dispatch.units.length === 0">
                        <p class="text-xs font-normal">
                            {{ $t('common.units', 0) }}
                        </p>
                    </template>
                    <ul v-else class="inline-flex flex-col text-xs font-normal">
                        <li v-for="unit in dispatch.units" :key="unit.unitId" class="inline-flex items-center gap-1">
                            <span>{{ unit.unit?.initials }}: {{ unit.unit?.name }}</span>
                        </li>
                    </ul>
                </div>
            </div>
        </template>
    </UPopover>
</template>
