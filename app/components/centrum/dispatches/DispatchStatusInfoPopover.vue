<script lang="ts" setup>
import DispatchDetailsByIDSlideover from '~/components/centrum/dispatches//DispatchDetailsByIDSlideover.vue';
import { useCentrumStore } from '~/store/centrum';
import type { DispatchStatus } from '~~/gen/ts/resources/centrum/dispatches';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { dispatchStatusToBGColor } from '../helpers';

const props = defineProps<{
    status: DispatchStatus | undefined;
}>();

const modal = useModal();

const centrumStore = useCentrumStore();

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
        <UButton
            variant="outline"
            :padded="false"
            size="xs"
            class="inline-flex items-center p-0.5"
            trailing-icon="i-mdi-chevron-down"
        >
            <slot name="before" />
            <span> DSP-{{ status.dispatchId }} </span>
            <slot name="after" />
        </UButton>

        <template #panel>
            <div class="inline-flex min-w-48 flex-col gap-1 p-4">
                <div class="flex items-center gap-2">
                    <UButton
                        variant="link"
                        icon="i-mdi-car-emergency"
                        :title="$t('common.detail', 2)"
                        @click="
                            modal.open(DispatchDetailsByIDSlideover, {
                                dispatchId: status.dispatchId,
                            })
                        "
                    >
                        {{ $t('common.detail', 2) }}
                    </UButton>
                </div>

                <p class="text-base font-semibold leading-none text-gray-900 dark:text-white">DSP-{{ status.dispatchId }}</p>

                <UBadge class="rounded font-semibold" :class="dispatchStatusColor" size="xs">
                    {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[status.status ?? 0]}`) }}
                </UBadge>

                <div v-if="dispatch" class="text-gray-900 dark:text-white">
                    <p class="text-sm font-medium leading-none">
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
