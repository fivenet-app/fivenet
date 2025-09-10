<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3';

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
        <UTooltip :text="$t('components.partials.TiptapEditor.extensions.template_block.title')">
            <UButton color="neutral" variant="ghost" icon="i-mdi-application-variable" :disabled="disabled" />
        </UTooltip>
        <template #content>
            <div class="flex flex-1 flex-col gap-1 p-4">
                <h3 class="block font-medium">
                    {{ $t('components.partials.TiptapEditor.extensions.template_block.title') }}
                </h3>

                <UFormField>
                    <USelectMenu v-model="selected" class="w-full" :items="options" value-key="value" />
                </UFormField>

                <div class="flex flex-row gap-2">
                    <UFormField
                        class="justify-center"
                        :label="$t('components.partials.TiptapEditor.extensions.template_var.trim_left')"
                    >
                        <UCheckbox v-model="leftTrim" />
                    </UFormField>

                    <UFormField
                        class="justify-center"
                        :label="$t('components.partials.TiptapEditor.extensions.template_var.trim_right')"
                    >
                        <UCheckbox v-model="rightTrim" />
                    </UFormField>
                </div>

                <UFormField>
                    <UInput
                        v-model="expression"
                        :placeholder="
                            selected
                                ? $t('components.partials.TiptapEditor.extensions.template_block.block_placeholder.select')
                                : $t('components.partials.TiptapEditor.extensions.template_block.block_placeholder.empty')
                        "
                        :disabled="!selected"
                    />
                </UFormField>

                <UFormField>
                    <UButton block :disabled="!canInsert" @click="insertBlock">
                        {{ $t('components.partials.TiptapEditor.extensions.template_block.insert_block') }}
                    </UButton>
                </UFormField>
            </div>
        </template>
    </UPopover>
</template>
