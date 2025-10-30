<script lang="ts" setup>
import { vMaska } from 'maska/vue';

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

watch(modelValue, (val) => {
    manual.value = val ?? '';
});

const manual = ref('');
watch(manual, (val) => {
    if (/^#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$/.test(val)) {
        modelValue.value = val;
    }
});
</script>

<template>
    <UPopover>
        <UButton
            :label="!hideLabel ? $t('common.choose_color') : undefined"
            color="neutral"
            variant="outline"
            :disabled="disabled"
            :block="block"
            v-bind="$attrs"
        >
            <template #leading>
                <span :style="{ backgroundColor: modelValue ?? 'black' }" class="size-5 rounded-sm" />
            </template>
        </UButton>

        <template #content>
            <UColorPicker v-model="modelValue" class="p-2" format="hex" />

            <UInput
                v-model="manual"
                v-maska
                type="text"
                data-maska="!#HHHHHH"
                data-maska-tokens="H:[0-9a-fA-F]"
                :ui="{ trailing: 'pr-0.5' }"
            >
                <template v-if="manual?.length" #trailing>
                    <UTooltip :text="$t('common.copy')" :content="{ side: 'right' }">
                        <UButton
                            color="neutral"
                            variant="link"
                            size="sm"
                            icon="i-lucide-copy"
                            :aria-label="$t('common.copy')"
                            @click="copyToClipboardWrapper(manual)"
                        />
                    </UTooltip>
                </template>
            </UInput>
        </template>
    </UPopover>
</template>
