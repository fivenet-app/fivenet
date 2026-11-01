<script lang="ts" setup>
import type { Form } from '@nuxt/ui';
import { z } from 'zod';
import type { PropertyGroup, PropertyScope } from '~~/gen/ts/resources/properties/properties';
import ResponsiveDrawerModal from '~/components/partials/ResponsiveDrawerModal.vue';

const props = defineProps<{
    group?: PropertyGroup;
    scope: PropertyScope;
    onSave: (group: PropertyGroup) => Promise<void>;
}>();

const emit = defineEmits<{ (e: 'close', value: boolean): void }>();
const { t } = useI18n();
const isEditing = computed(() => props.group !== undefined);

const schema = z.object({
    key: z
        .string()
        .min(1)
        .max(64)
        .regex(/^[a-z][a-z0-9_.-]*$/),
    label: z.string().trim().min(1).max(128),
});
const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const editor = reactive<PropertyGroup>({ scope: props.scope, key: '', label: '', displayKey: '', sortOrder: 0 });

const { submit, isSubmitting, canSubmit } = useSubmitGuard(async () => props.onSave({ ...editor, displayKey: '' }));
const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(() => ({ editor: { ...editor } }));

function resetEditor(): void {
    Object.assign(editor, {
        scope: props.scope,
        key: props.group?.key ?? '',
        label: props.group?.label ?? '',
        displayKey: '',
        sortOrder: props.group?.sortOrder ?? 0,
        deletedAt: undefined,
    });
}

watch([() => props.group, () => props.scope], resetEditor, { immediate: true });
watch([() => props.group, () => props.scope], () => syncSnapshot(), { immediate: true });

async function close(): Promise<void> {
    if (!canSubmit.value) return;
    if (hasUnsavedChanges.value && !(await confirmLeave())) return;
    emit('close', false);
}
</script>

<template>
    <ResponsiveDrawerModal
        :title="t('components.settings.property_definitions.group_editor_title')"
        :close="false"
        :dismissible="!hasUnsavedChanges"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-1.5">
                <h3 class="font-semibold text-highlighted">
                    {{ t('components.settings.property_definitions.group_editor_title') }}
                </h3>
                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="close"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" class="flex flex-col gap-4" :schema="schema" :state="editor" @submit="submit">
                <UFormField
                    name="key"
                    :label="t('components.settings.property_definitions.group_key')"
                    :description="t('components.settings.property_definitions.group_key_description')"
                    required
                >
                    <UInput v-model="editor.key" class="w-full" :disabled="isEditing" />
                </UFormField>
                <UFormField
                    name="label"
                    :label="t('components.settings.property_definitions.group_label')"
                    :description="t('components.settings.property_definitions.group_label_description')"
                    required
                >
                    <UInput v-model="editor.label" class="w-full" />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="close"
                />
                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="isSubmitting"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </ResponsiveDrawerModal>
</template>
