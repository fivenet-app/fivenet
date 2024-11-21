<script lang="ts" setup>
import { backgroundColors, primaryColors } from '~/composables/colors';

const props = withDefaults(
    defineProps<{
        modelValue: string | undefined;
    }>(),
    {
        modelValue: 'primary',
    },
);

const emit = defineEmits<{
    (e: 'update:modelValue', value: string | undefined): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const color = useVModel(props, 'modelValue', emit);

const availableColorOptions = [...primaryColors, ...backgroundColors];
</script>

<template>
    <ClientOnly>
        <USelectMenu
            v-bind="$attrs"
            v-model="color"
            :options="availableColorOptions"
            option-attribute="label"
            value-attribute="label"
            :searchable-placeholder="$t('common.color')"
        >
            <template #label>
                <span class="size-2 rounded-full" :class="availableColorOptions.find((o) => o.label === color)?.class" />
                <span class="truncate">{{ color }}</span>
            </template>

            <template #option="{ option }">
                <span class="size-2 rounded-full" :class="option.class" />
                <span class="truncate">{{ option.label }}</span>
            </template>
        </USelectMenu>
    </ClientOnly>
</template>
