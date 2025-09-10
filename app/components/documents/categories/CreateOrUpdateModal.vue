<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { ShapeIcon } from 'mdi-vue3';
import { z } from 'zod';
import ColorPickerTW from '~/components/partials/ColorPickerTW.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import type { Category } from '~~/gen/ts/resources/documents/category';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    category?: Category;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'update'): void;
}>();

const { can } = useAuth();

const { fallbackColor } = useAppConfig();

const notifications = useNotificationsStore();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const canEdit = can('documents.DocumentsService/CreateOrUpdateCategory');

const schema = z.object({
    name: z.string().min(3).max(128),
    description: z.union([z.string().min(0).max(255), z.string().optional()]),
    color: z.string().max(7),
    icon: z.string().max(128).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    description: '',
    color: 'primary',
});

async function createOrUpdateCategory(values: Schema): Promise<void> {
    try {
        await documentsDocumentsClient.createOrUpdateCategory({
            category: {
                id: props.category?.id ?? 0,
                name: values.name,
                description: values.description,
                color: values.color,
                icon: values.icon,
            },
        });

        if (!props.category?.id) {
            notifications.add({
                title: { key: 'notifications.category_created.title', parameters: {} },
                description: { key: 'notifications.category_created.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        } else {
            notifications.add({
                title: { key: 'notifications.category_updated.title', parameters: {} },
                description: { key: 'notifications.category_updated.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        }

        emit('update');
        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteCategory(): Promise<void> {
    if (props.category === undefined) return;

    try {
        await documentsDocumentsClient.deleteCategory({
            id: props.category.id!,
        });

        if (!props.category.deletedAt) {
            notifications.add({
                title: { key: 'notifications.category_deleted.title', parameters: {} },
                description: { key: 'notifications.category_deleted.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        } else {
            notifications.add({
                title: { key: 'notifications.category_restored.title', parameters: {} },
                description: { key: 'notifications.category_restored.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        }

        emit('update');
        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateCategory(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

function setFromProps(): void {
    if (!props.category) {
        return;
    }

    state.name = props.category.name;
    state.description = props.category.description;
    state.color = props.category.color ?? fallbackColor;
    state.icon = props.category.icon;
}

setFromProps();
watch(props, () => setFromProps());

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal
        :title="
            canEdit
                ? category
                    ? $t('components.documents.categories.modal.update_category')
                    : $t('components.documents.categories.modal.create_category')
                : $t('common.category') + (category ? `: ${category.name}` : '')
        "
    >
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField class="flex-1" name="name" :label="$t('common.name')">
                    <UInput
                        v-model="state.name"
                        type="text"
                        name="name"
                        :disabled="!canEdit"
                        :placeholder="$t('common.name', 1)"
                        :label="$t('common.name', 1)"
                        class="w-full"
                    />
                </UFormField>

                <UFormField class="flex-1" name="description" :label="$t('common.description')">
                    <UTextarea
                        v-model="state.description"
                        name="description"
                        :disabled="!canEdit"
                        :placeholder="$t('common.description')"
                        class="w-full"
                    />
                </UFormField>

                <UFormField class="flex-1 flex-row" name="color" :label="$t('common.color')">
                    <div class="flex flex-1 gap-1">
                        <ColorPickerTW v-model="state.color" class="flex-1" :disabled="!canEdit" />
                    </div>
                </UFormField>

                <UFormField class="flex-1" name="icon" :label="$t('common.icon')">
                    <div class="flex flex-1 gap-1">
                        <IconSelectMenu
                            v-model="state.icon"
                            class="flex-1"
                            :color="state.color"
                            :disabled="!canEdit"
                            :fallback-icon="ShapeIcon"
                        />

                        <UButton v-if="canEdit" color="error" icon="i-mdi-backspace" @click="state.icon = undefined" />
                    </div>
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    v-if="category !== undefined && canEdit && can('documents.DocumentsService/DeleteCategory').value"
                    class="flex-1"
                    block
                    :color="!category.deletedAt ? 'error' : 'success'"
                    :icon="!category.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore'"
                    :label="!category.deletedAt ? $t('common.delete') : $t('common.restore')"
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    @click="() => deleteCategory()"
                />

                <UButton
                    v-if="canEdit"
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="category === undefined ? $t('common.create') : $t('common.update')"
                    @click="() => formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>
