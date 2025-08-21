<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { ShapeIcon } from 'mdi-vue3';
import { z } from 'zod';
import ColorPickerTW from '~/components/partials/ColorPickerTW.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import type { Category } from '~~/gen/ts/resources/documents/category';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    category?: Category;
}>();

const emit = defineEmits<{
    (e: 'update'): void;
}>();

const { $grpc } = useNuxtApp();

const { can } = useAuth();

const { fallbackColor } = useAppConfig();

const notifications = useNotificationsStore();

const { isOpen } = useModal();

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
        await $grpc.documents.documents.createOrUpdateCategory({
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
        isOpen.value = false;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteCategory(): Promise<void> {
    if (props.category === undefined) {
        return;
    }

    try {
        await $grpc.documents.documents.deleteCategory({
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
        isOpen.value = false;
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
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            <template v-if="!canEdit">
                                {{ $t('common.category') }}:
                                {{ category?.name }}
                            </template>
                            <template v-else-if="category">
                                {{ $t('components.documents.categories.modal.update_category') }}:
                                {{ category?.name }}
                            </template>
                            <template v-else>
                                {{ $t('components.documents.categories.modal.create_category') }}
                            </template>
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <div>
                        <UFormGroup class="flex-1" name="name" :label="$t('common.name')">
                            <UInput
                                v-model="state.name"
                                type="text"
                                name="name"
                                :disabled="!canEdit"
                                :placeholder="$t('common.name', 1)"
                                :label="$t('common.name', 1)"
                            />
                        </UFormGroup>

                        <UFormGroup class="flex-1" name="description" :label="$t('common.description')">
                            <UTextarea
                                v-model="state.description"
                                name="description"
                                :disabled="!canEdit"
                                :placeholder="$t('common.description')"
                            />
                        </UFormGroup>

                        <UFormGroup class="flex-1 flex-row" name="color" :label="$t('common.color')">
                            <div class="flex flex-1 gap-1">
                                <ColorPickerTW v-model="state.color" class="flex-1" :disabled="!canEdit" />
                            </div>
                        </UFormGroup>

                        <UFormGroup class="flex-1" name="icon" :label="$t('common.icon')">
                            <div class="flex flex-1 gap-1">
                                <IconSelectMenu
                                    v-model="state.icon"
                                    class="flex-1"
                                    :color="state.color"
                                    :disabled="!canEdit"
                                    :fallback-icon="ShapeIcon"
                                />

                                <UButton v-if="canEdit" icon="i-mdi-backspace" @click="state.icon = undefined" />
                            </div>
                        </UFormGroup>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="black" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton
                            v-if="category !== undefined && canEdit && can('documents.DocumentsService/DeleteCategory').value"
                            class="flex-1"
                            block
                            :color="!category.deletedAt ? 'error' : 'success'"
                            :icon="!category.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore'"
                            :label="!category.deletedAt ? $t('common.delete') : $t('common.restore')"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                            @click="deleteCategory()"
                        />

                        <UButton v-if="canEdit" class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ category === undefined ? $t('common.create') : $t('common.update') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
