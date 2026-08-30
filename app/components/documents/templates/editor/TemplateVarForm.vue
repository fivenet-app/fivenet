<script setup lang="ts">
import TemplateTrimControls from './TemplateTrimControls.vue';

const props = withDefaults(
    defineProps<{
        editing?: boolean;
        value?: string;
        leftTrim?: boolean;
        rightTrim?: boolean;
    }>(),
    {
        editing: false,
        value: '',
        leftTrim: false,
        rightTrim: false,
    },
);

const emit = defineEmits<{
    (e: 'submit', value: { value: string; leftTrim: boolean; rightTrim: boolean }): void;
    (e: 'delete'): void;
}>();

const { t } = useI18n();

type Category = { label: string; value: string; key: string };

const categories: Category[] = [
    { label: t('common.date'), value: 'now', key: 'date' },
    { label: t('common.active_user'), value: '.ActiveChar', key: 'activeChar' },
    { label: t('common.first_citizen'), value: '(first .Users)', key: 'user' },
];

const baseUserProperties: { label: string; value: string }[] = [
    { label: t('common.firstname'), value: '.Firstname' },
    { label: t('common.lastname'), value: '.Lastname' },
    { label: t('common.date_of_birth'), value: '.Dateofbirth' },
    { label: t('common.sex'), value: '.Sex' },
    { label: t('common.height'), value: '.Height' },
];

const templateVars = computed<Record<string, { label: string; value: string }[]>>(() => ({
    date: [
        { label: `${t('common.date')} "02.01.2006 15:04"`, value: ' | date "02.01.2006 15:04"' },
        { label: `${t('common.date')} "02.01.2006"`, value: ' | date "02.01.2006"' },
        { label: `${t('common.time')} "15:04"`, value: ' | date "15:04"' },
    ],
    activeChar: [
        ...baseUserProperties,
        { label: t('common.phone'), value: '.PhoneNumber' },
        { label: t('common.prefix'), value: '.Props.NamePrefix' },
        { label: t('common.suffix'), value: '.Props.NameSuffix' },
        { label: t('common.mail'), value: '.Email' },
    ],
    user: [
        ...baseUserProperties,
        { label: t('common.wanted'), value: '.Props.Wanted' },
        { label: t('common.phone'), value: '.PhoneNumber' },
    ],
}));

const selectedCategory = ref<Category>();
const selectedProperty = ref<string>();
const draftValue = ref(props.value);
const customInput = ref('');
const draftLeftTrim = ref(props.leftTrim);
const draftRightTrim = ref(props.rightTrim);

watch(
    () => [props.value, props.leftTrim, props.rightTrim] as const,
    ([value, leftTrim, rightTrim]) => {
        draftValue.value = value;
        draftLeftTrim.value = leftTrim;
        draftRightTrim.value = rightTrim;
    },
    { immediate: true },
);

function submit(value: string): void {
    if (!value.trim()) return;

    emit('submit', {
        value: value.trim(),
        leftTrim: draftLeftTrim.value,
        rightTrim: draftRightTrim.value,
    });
}

function insertSelected(): void {
    if (!selectedCategory.value || !selectedProperty.value) return;
    submit(selectedCategory.value.value + selectedProperty.value);
    selectedCategory.value = undefined;
    selectedProperty.value = undefined;
}

function insertCustom(): void {
    submit(customInput.value);
    customInput.value = '';
}
</script>

<template>
    <div class="flex flex-col gap-2">
        <template v-if="!editing">
            <UFormField name="category" :label="$t('common.category', 1)">
                <USelectMenu v-model="selectedCategory" class="w-full" :items="categories" />
            </UFormField>

            <UFormField name="property" :label="$t('common.property', 1)">
                <USelectMenu
                    v-model="selectedProperty"
                    class="w-full"
                    :items="templateVars[selectedCategory?.key ?? '']"
                    :disabled="!templateVars[selectedCategory?.key ?? '']"
                    value-key="value"
                />
            </UFormField>
        </template>

        <UFormField
            v-if="editing"
            name="templateValue"
            :label="$t('components.partials.tiptap_editor.extensions.template_var.custom_template')"
        >
            <UTextarea v-model="draftValue" class="w-full" :rows="3" autofocus />
        </UFormField>

        <UFieldGroup v-if="editing">
            <UButton block :label="$t('common.save')" :disabled="!draftValue.trim()" @click="submit(draftValue)" />
            <UTooltip :text="$t('common.delete')">
                <UButton color="error" icon="i-mdi-trash-can" @click="$emit('delete')" />
            </UTooltip>
        </UFieldGroup>

        <template v-else>
            <UButton
                block
                :label="$t('components.partials.tiptap_editor.insert')"
                :disabled="!selectedCategory || !selectedProperty"
                @click="insertSelected"
            />

            <USeparator class="my-2" />

            <UFormField
                name="customInput"
                :label="$t('components.partials.tiptap_editor.extensions.template_var.custom_template')"
            >
                <UInput
                    v-model="customInput"
                    class="w-full"
                    :placeholder="$t('components.partials.tiptap_editor.extensions.template_var.custom_placeholder')"
                />
            </UFormField>

            <UButton
                block
                :label="$t('components.partials.tiptap_editor.extensions.template_var.insert_custom')"
                :disabled="!customInput"
                @click="insertCustom"
            />
        </template>

        <USeparator class="my-2" />

        <TemplateTrimControls v-model:left-trim="draftLeftTrim" v-model:right-trim="draftRightTrim" />
    </div>
</template>
