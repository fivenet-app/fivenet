<script lang="ts" setup>
const props = withDefaults(
    defineProps<{
        value: string | number | (() => string | number);
        showText?: boolean;
    }>(),
    {
        showText: false,
    },
);

const { notifyCopied } = useClipboardNotifications();

async function addToClipboard(): Promise<void> {
    await copyToClipboardWrapper(typeof props.value === 'function' ? props.value().toString() : props.value.toString());

    notifyCopied();
}

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UTooltip :text="$t('components.clipboard.clipboard_button.add')">
        <UButton
            class="px-1 py-1"
            trailing-icon="i-mdi-clipboard-plus"
            variant="outline"
            color="neutral"
            size="xs"
            :label="props.showText ? $t('common.copy') : undefined"
            @click="addToClipboard()"
        />
    </UTooltip>
</template>
