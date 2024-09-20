<script lang="ts" setup>
import DispatchDetailsByIDSlideover from '~/components/centrum/dispatches//DispatchDetailsByIDSlideover.vue';
import type { DispatchStatus} from '~~/gen/ts/resources/centrum/dispatches';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

defineProps<{
    status: DispatchStatus | undefined;
}>();

const modal = useModal();
</script>

<template>
    <template v-if="!status">
        <span class="inline-flex items-center">
            <slot name="before" />
            <span>{{ $t('common.na') }}</span>
            <slot name="after" />
        </span>
    </template>
    <template v-else>
        <UPopover>
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
                <div class="p-4">
                    <div class="mb-2 flex items-center gap-2">
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
                    <div class="text-gray-900 dark:text-white">
                        <p class="text-base font-semibold leading-none">DSP-{{ status.dispatchId }}</p>
                        <p class="inline-flex items-center justify-center text-sm font-normal">
                            <span class="font-semibold"> {{ $t('common.status') }} </span>:
                            {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[status.status ?? 0]}`) }}
                        </p>
                    </div>
                </div>
            </template>
        </UPopover>
    </template>
</template>
