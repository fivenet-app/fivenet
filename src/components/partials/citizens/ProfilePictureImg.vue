<script lang="ts" setup>
import { Float } from '@headlessui-float/vue';
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { type imageSizes } from '~/components/partials/helpers';
import SquareImg from '~/components/partials/elements/SquareImg.vue';

const { t } = useI18n();

const props = withDefaults(
    defineProps<{
        url?: string;
        name: string;
        size?: imageSizes;
        rounded?: boolean;
        enablePopup?: boolean;
        noBlur?: boolean;
        altText?: string;
    }>(),
    {
        url: undefined,
        size: 'lg',
        rounded: false,
        enablePopup: false,
        noBlur: undefined,
        altText: undefined,
    },
);

const altText = computed(() => (props.altText !== undefined ? props.altText : t('common.avatar')));
</script>

<template>
    <template v-if="!url || !enablePopup">
        <SquareImg :url="url" :text="altText" :size="size" :rounded="rounded" :no-blur="noBlur">
            <template #initials>
                {{ getInitials(name) }}
            </template>
        </SquareImg>
    </template>
    <Popover v-else class="relative">
        <Float
            portal
            auto-placement
            :offset="15"
            enter="transition duration-150 ease-out"
            enter-from="scale-95 opacity-0"
            enter-to="scale-100 opacity-100"
            leave="transition duration-100 ease-in"
            leave-from="scale-100 opacity-100"
            leave-to="scale-95 opacity-0"
        >
            <PopoverButton class="inline-flex items-center">
                <SquareImg :url="url" :text="altText" :size="size" :rounded="rounded" :no-blur="noBlur">
                    <template #initials>
                        {{ getInitials(name) }}
                    </template>
                </SquareImg>
            </PopoverButton>

            <PopoverPanel
                class="absolute z-5 w-96 min-w-fit max-w-[18rem] rounded-lg border border-gray-600 bg-gray-800 text-sm text-gray-400 shadow-sm transition-opacity"
            >
                <div class="p-3">
                    <img class="rounded-md w-96 max-w-full" :src="url" :alt="altText" />
                </div>
            </PopoverPanel>
        </Float>
    </Popover>
</template>
