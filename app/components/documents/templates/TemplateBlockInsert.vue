<script setup lang="ts">
import type { Editor } from '@tiptap/core';

const props = defineProps<{
    editor: Editor;
    disabled?: boolean;
}>();

const { t } = useI18n();

const options = [
    { label: t('components.partials.TiptapEditor.extensions.template_block.options.range'), value: 'range' },
    { label: t('components.partials.TiptapEditor.extensions.template_block.options.if'), value: 'if' },
    { label: t('components.partials.TiptapEditor.extensions.template_block.options.with'), value: 'with' },
];

const selected = ref('');
const expression = ref('');
const leftTrim = ref(false);
const rightTrim = ref(false);

const canInsert = computed(() => selected.value && expression.value);

const insertBlock = () => {
    if (!canInsert.value) return;
    const val = `${selected.value} ${expression.value}`.trim();
    props.editor?.commands.insertTemplateBlock({
        value: val,
        leftTrim: leftTrim.value,
        rightTrim: rightTrim.value,
    });
    selected.value = '';
    expression.value = '';
};
</script>

<template>
    <UPopover>
        <UTooltip :text="$t('components.partials.TiptapEditor.extensions.template_block.title')" :popper="{ placement: 'top' }">
            <UButton color="white" variant="ghost" icon="i-mdi-application-variable" :disabled="disabled" />
        </UTooltip>
        <template #panel>
            <div class="flex flex-1 flex-col gap-1 p-4">
                <h3 class="block font-medium">
                    {{ $t('components.partials.TiptapEditor.extensions.template_block.title') }}
                </h3>

                <UFormGroup>
                    <USelectMenu v-model="selected" class="w-full" :options="options" value-attribute="value" />
                </UFormGroup>

                <div class="flex flex-row gap-2">
                    <UFormGroup
                        class="justify-center"
                        :label="$t('components.partials.TiptapEditor.extensions.template_var.trim_left')"
                    >
                        <UCheckbox v-model="leftTrim" />
                    </UFormGroup>

                    <UFormGroup
                        class="justify-center"
                        :label="$t('components.partials.TiptapEditor.extensions.template_var.trim_right')"
                    >
                        <UCheckbox v-model="rightTrim" />
                    </UFormGroup>
                </div>

                <UFormGroup>
                    <UInput
                        v-model="expression"
                        :placeholder="
                            selected
                                ? $t('components.partials.TiptapEditor.extensions.template_block.block_placeholder.select')
                                : $t('components.partials.TiptapEditor.extensions.template_block.block_placeholder.empty')
                        "
                        :disabled="!selected"
                    />
                </UFormGroup>

                <UFormGroup>
                    <UButton block :disabled="!canInsert" @click="insertBlock">
                        {{ $t('components.partials.TiptapEditor.extensions.template_block.insert_block') }}
                    </UButton>
                </UFormGroup>
            </div>
        </template>
    </UPopover>
</template>
