<script setup lang="ts">
import type { Editor } from '@tiptap/core';

const props = defineProps<{
    editor: Editor;
    disabled?: boolean;
}>();

const { t } = useI18n();

type Category = { label: string; value: string; key: string };

const categories: Category[] = [
    { label: t('common.date'), value: '', key: 'date' },
    { label: t('common.active_user'), value: '.ActiveChar', key: 'user' },
    { label: t('common.first_citizen'), value: '(first .Users)', key: 'user' },
];

const templateVars: Record<string, { label: string; value: string }[]> = {
    user: [
        { label: t('common.firstname'), value: '.Firstname' },
        { label: t('common.lastname'), value: '.Lastname' },
        { label: t('common.date_of_birth'), value: '.Dateofbirth' },
    ],
    date: [
        { label: `${t('common.date')} "02.01.2006 15:04"`, value: 'now | date "02.01.2006 15:04"' },
        { label: `${t('common.time')} "15:04"`, value: 'now | date "15:04"' },
    ],
};

const selectedCategory = ref<(typeof categories)[0] | undefined>(undefined);
const selectedProperty = ref<string | undefined>(undefined);
const customInput = ref('');
const leftTrim = ref(false);
const rightTrim = ref(false);

const insert = () => {
    if (!selectedCategory.value || !selectedProperty.value) return;
    props.editor?.commands.insertTemplateVar({
        value: selectedCategory.value.value + selectedProperty.value,
        leftTrim: leftTrim.value,
        rightTrim: rightTrim.value,
    });
    selectedCategory.value = undefined;
    selectedProperty.value = undefined;
};

const insertCustom = () => {
    if (!customInput.value) return;
    props.editor?.commands.insertTemplateVar({
        value: customInput.value,
        leftTrim: leftTrim.value,
        rightTrim: rightTrim.value,
    });
    customInput.value = '';
};
</script>

<template>
    <UPopover>
        <UTooltip :text="$t('components.partials.TiptapEditor.extensions.template_var.title')" :popper="{ placement: 'top' }">
            <UButton color="white" variant="ghost" icon="i-mdi-variable" :disabled="disabled" />
        </UTooltip>

        <template #panel>
            <div class="flex flex-1 flex-col gap-1 p-4">
                <h3 class="block font-medium">
                    {{ $t('components.partials.TiptapEditor.extensions.template_var.title') }}
                </h3>

                <UFormGroup :label="$t('common.category', 1)">
                    <USelectMenu v-model="selectedCategory" class="w-full" :options="categories" />
                </UFormGroup>

                <UFormGroup :label="$t('common.property', 1)">
                    <USelectMenu
                        v-model="selectedProperty"
                        class="w-full"
                        :options="templateVars[selectedCategory?.key ?? '']"
                        value-attribute="value"
                    />
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
                    <UButton
                        block
                        :label="$t('common.insert')"
                        :disabled="!selectedCategory || !selectedProperty"
                        @click="insert"
                    />
                </UFormGroup>

                <UFormGroup :label="$t('components.partials.TiptapEditor.extensions.template_var.custom_template')">
                    <UInput
                        v-model="customInput"
                        class="w-full"
                        :placeholder="$t('components.partials.TiptapEditor.extensions.template_var.custom_placeholder')"
                    />
                </UFormGroup>

                <UFormGroup>
                    <UButton
                        block
                        :label="$t('components.partials.TiptapEditor.extensions.template_var.insert_custom')"
                        :disabled="!customInput"
                        @click="insertCustom"
                    />
                </UFormGroup>
            </div>
        </template>
    </UPopover>
</template>
