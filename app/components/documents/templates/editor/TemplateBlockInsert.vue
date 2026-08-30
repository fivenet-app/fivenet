<script setup lang="ts">
import type { Editor } from '@tiptap/core';

const props = defineProps<{
    editor: Editor;
    disabled?: boolean;
}>();

const { t } = useI18n();

const options = [
    { label: t('components.partials.tiptap_editor.extensions.template_block.options.range'), value: 'range' },
    { label: t('components.partials.tiptap_editor.extensions.template_block.options.if'), value: 'if' },
    { label: t('components.partials.tiptap_editor.extensions.template_block.options.with'), value: 'with' },
];

const selected = ref('');
const expression = ref('');
const leftTrim = ref<boolean>(false);
const rightTrim = ref<boolean>(false);

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
        <UTooltip :text="$t('components.partials.tiptap_editor.extensions.template_block.title')">
            <UButton color="neutral" variant="ghost" icon="i-mdi-application-variable" :disabled="disabled" />
        </UTooltip>

        <template #content>
            <div class="flex w-full max-w-86 flex-col gap-2 p-4">
                <h3 class="block font-medium">
                    {{ $t('components.partials.tiptap_editor.extensions.template_block.title') }}
                </h3>

                <UAlert icon="i-mdi-information-outline" variant="subtle">
                    <template #description>
                        <I18nT keypath="components.partials.tiptap_editor.extensions.template_block.close_info">
                            <template #end>
                                <code class="font-mono" v-text="'{{ end }}'" />
                            </template>
                        </I18nT>
                    </template>
                </UAlert>

                <div class="flex flex-col gap-2">
                    <UFormField name="selected">
                        <USelectMenu v-model="selected" class="w-full" :items="options" value-key="value" />
                    </UFormField>

                    <div class="flex flex-row gap-2">
                        <UFormField
                            class="justify-center"
                            name="leftTrim"
                            :label="$t('components.partials.tiptap_editor.extensions.template_var.trim_left')"
                        >
                            <USwitch v-model="leftTrim" />
                        </UFormField>

                        <UFormField
                            class="justify-center"
                            name="rightTrim"
                            :label="$t('components.partials.tiptap_editor.extensions.template_var.trim_right')"
                        >
                            <USwitch v-model="rightTrim" />
                        </UFormField>
                    </div>

                    <UFormField name="expression">
                        <UInput
                            v-model="expression"
                            class="w-full"
                            :placeholder="
                                selected
                                    ? $t('components.partials.tiptap_editor.extensions.template_block.block_placeholder.select')
                                    : $t('components.partials.tiptap_editor.extensions.template_block.block_placeholder.empty')
                            "
                            :disabled="!selected"
                        />
                    </UFormField>

                    <UFormField>
                        <UButton
                            block
                            :disabled="!canInsert"
                            :label="$t('components.partials.tiptap_editor.extensions.template_block.insert_block')"
                            @click="insertBlock"
                        />
                    </UFormField>

                    <UButton block variant="outline" @click="props.editor?.commands.insertTemplateBlockEnd()">
                        <I18nT keypath="components.partials.tiptap_editor.extensions.template_block.insert_end">
                            <template #end>
                                <code class="font-mono" v-text="'{{ end }}'" />
                            </template>
                        </I18nT>
                    </UButton>
                </div>
            </div>
        </template>
    </UPopover>
</template>
