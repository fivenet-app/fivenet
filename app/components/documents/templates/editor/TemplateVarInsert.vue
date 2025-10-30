<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3';

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
        <UTooltip :text="$t('components.partials.tiptap_editor.extensions.template_var.title')">
            <UButton color="neutral" variant="ghost" icon="i-mdi-variable" :disabled="disabled" />
        </UTooltip>

        <template #content>
            <div class="flex flex-col gap-2 p-4">
                <h3 class="block font-medium">
                    {{ $t('components.partials.tiptap_editor.extensions.template_var.title') }}
                </h3>

                <div class="flex flex-col gap-2">
                    <UFormField :label="$t('common.category', 1)">
                        <USelectMenu v-model="selectedCategory" class="w-full" :items="categories" />
                    </UFormField>

                    <UFormField :label="$t('common.property', 1)">
                        <USelectMenu
                            v-model="selectedProperty"
                            class="w-full"
                            :items="templateVars[selectedCategory?.key ?? '']"
                            value-key="value"
                        />
                    </UFormField>

                    <div class="flex flex-row gap-2">
                        <UFormField
                            class="justify-center"
                            :label="$t('components.partials.tiptap_editor.extensions.template_var.trim_left')"
                        >
                            <USwitch v-model="leftTrim" />
                        </UFormField>

                        <UFormField
                            class="justify-center"
                            :label="$t('components.partials.tiptap_editor.extensions.template_var.trim_right')"
                        >
                            <USwitch v-model="rightTrim" />
                        </UFormField>
                    </div>

                    <UFormField>
                        <UButton
                            block
                            :label="$t('common.insert')"
                            :disabled="!selectedCategory || !selectedProperty"
                            @click="insert"
                        />
                    </UFormField>

                    <UFormField :label="$t('components.partials.tiptap_editor.extensions.template_var.custom_template')">
                        <UInput
                            v-model="customInput"
                            class="w-full"
                            :placeholder="$t('components.partials.tiptap_editor.extensions.template_var.custom_placeholder')"
                        />
                    </UFormField>

                    <UFormField>
                        <UButton
                            block
                            :label="$t('components.partials.tiptap_editor.extensions.template_var.insert_custom')"
                            :disabled="!customInput"
                            @click="insertCustom"
                        />
                    </UFormField>
                </div>
            </div>
        </template>
    </UPopover>
</template>
