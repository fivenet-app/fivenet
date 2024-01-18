<script lang="ts" setup>
import { Float } from '@headlessui-float/vue';
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { DispatchStatus, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

defineProps<{
    status: DispatchStatus | undefined;
    textClass?: unknown;
    buttonClass?: unknown;
}>();
</script>

<template>
    <template v-if="!status">
        <span class="inline-flex items-center">
            <slot name="before" />
            <span>N/A</span>
            <slot name="after" />
        </span>
    </template>
    <Popover v-else class="relative">
        <Float portal auto-placement :offset="16">
            <PopoverButton class="inline-flex items-center" :class="buttonClass">
                <slot name="before" />
                <span :class="textClass"> DSP-{{ status.dispatchId }} </span>
                <slot name="after" />
            </PopoverButton>

            <PopoverPanel
                class="absolute z-5 w-64 min-w-fit max-w-[18rem] rounded-lg border border-gray-600 bg-gray-800 text-sm text-gray-400 shadow-sm transition-opacity"
            >
                <div class="p-3">
                    <p class="text-base font-semibold leading-none text-gray-900 dark:text-neutral">
                        DSP-{{ status.dispatchId }}
                    </p>
                    <p class="inline-flex items-center justify-center text-sm font-normal">
                        <span class="font-semibold"> {{ $t('common.status') }} </span>:
                        {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[status.status ?? 0]}`) }}
                    </p>
                </div>
            </PopoverPanel>
        </Float>
    </Popover>
</template>
