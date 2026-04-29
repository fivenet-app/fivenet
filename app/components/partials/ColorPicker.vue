<script lang="ts" setup>
import { vMaska } from 'maska/vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

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

const notifications = useNotificationsStore();

watch(modelValue, (val) => (manual.value = val ?? ''));

const manual = ref(modelValue.value ?? '');
watch(manual, (val) => {
    if (/^#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$/.test(val)) {
        modelValue.value = val;
    }
});

async function copyColor() {
    copyToClipboardWrapper(manual.value);

    notifications.add({
        title: { key: 'notifications.clipboard.copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.copied.content', parameters: {} },
        duration: 1500,
        type: NotificationType.INFO,
    });
}
</script>

<template>
    <UPopover>
        <UButton
            :label="!hideLabel ? (modelValue ? modelValue : $t('common.choose_color')) : undefined"
            color="neutral"
            variant="outline"
            :disabled="disabled"
            :block="block"
            v-bind="$attrs"
        >
            <template #leading>
                <span class="size-5 rounded-sm" :style="{ backgroundColor: modelValue ?? 'black' }" />
            </template>
        </UButton>

        <template #content>
            <UColorPicker v-model="modelValue" class="p-2" format="hex" />

            <UInput
                v-model="manual"
                v-maska
                type="text"
                :ui="{ trailing: 'pr-0.5' }"
                data-maska="!#HHHHHH"
                data-maska-tokens="H:[0-9a-fA-F]"
            >
                <template v-if="manual?.length" #trailing>
                    <UTooltip :text="$t('common.copy')" :content="{ side: 'right' }">
                        <UButton
                            color="neutral"
                            variant="link"
                            size="sm"
                            icon="i-mdi-content-copy"
                            :aria-label="$t('common.copy')"
                            @click="copyColor"
                        />
                    </UTooltip>
                </template>
            </UInput>
        </template>
    </UPopover>
</template>
