<script lang="ts" setup>
import { Float } from '@headlessui-float/vue';
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { type imageSizes } from '~/components/partials/helpers';
import SquareImg from '~/components/partials/SquareImg.vue';

withDefaults(
    defineProps<{
        url?: string;
        name: string;
        size?: imageSizes;
        rounded?: boolean;
        enablePopup?: boolean;
    }>(),
    {
        url: undefined,
        size: 'lg',
        rounded: false,
        enablePopup: false,
    },
);
</script>

<template>
    <template v-if="!enablePopup">
        <SquareImg :url="url" :text="$t('common.avatar')" :size="size" :rounded="rounded">
            <template #initials>
                {{ getInitials(name) }}
            </template>
        </SquareImg>
    </template>
    <Popover v-else-if="url" class="relative">
        <Float portal auto-placement :offset="16">
            <PopoverButton class="inline-flex items-center">
                <SquareImg :url="url" :text="$t('common.avatar')" :size="size" :rounded="rounded">
                    <template #initials>
                        {{ getInitials(name) }}
                    </template>
                </SquareImg>
            </PopoverButton>

            <PopoverPanel
                class="absolute z-5 w-96 min-w-fit max-w-[18rem] rounded-lg border border-gray-600 bg-gray-800 text-sm text-gray-400 shadow-sm transition-opacity"
            >
                <div class="p-3">
                    <img class="rounded-md max-w-96" :src="url" />
                </div>
            </PopoverPanel>
        </Float>
    </Popover>
</template>
