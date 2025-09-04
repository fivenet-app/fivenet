<script lang="ts" setup>
defineProps<{
    disabled?: boolean;
    block?: boolean;
    hideLabel?: boolean;
}>();

defineEmits<{
    (e: 'close'): void;
}>();

defineOptions({
    inheritAttrs: false,
});

const modelValue = defineModel<string>();
</script>

<template>
    <UPopover>
        <UButton
            :label="!hideLabel ? $t('common.choose_color') : undefined"
            color="neutral"
            variant="outline"
            :disabled="disabled"
            :block="block"
        >
            <template #leading>
                <span :style="{ backgroundColor: modelValue ?? 'black' }" class="size-3 rounded-sm" />
            </template>
        </UButton>

        <template #content>
            <UColorPicker v-model="modelValue" class="p-2" />
        </template>
    </UPopover>
</template>
