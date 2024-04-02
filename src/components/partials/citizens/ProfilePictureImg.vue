<script lang="ts" setup>
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
    <UPopover v-else>
        <UButton variant="link" :padded="false" class="inline-flex items-center">
            <SquareImg :url="url" :text="altText" :size="size" :rounded="rounded" :no-blur="noBlur">
                <template #initials>
                    {{ getInitials(name) }}
                </template>
            </SquareImg>
        </UButton>

        <template #panel>
            <div class="p-3">
                <img class="w-96 max-w-full rounded-md" :src="url" :alt="altText" />
            </div>
        </template>
    </UPopover>
</template>
