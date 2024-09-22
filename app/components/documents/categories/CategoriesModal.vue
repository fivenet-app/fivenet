<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { ShapeIcon } from 'mdi-vue3';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import { useNotificatorStore } from '~/store/notificator';
import type { Category } from '~~/gen/ts/resources/documents/category';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    category?: Category;
}>();

const emit = defineEmits<{
    (e: 'update'): void;
}>();

const notifications = useNotificatorStore();

const { isOpen } = useModal();

const schema = z.object({
    name: z.string().min(3).max(128),
    description: z.union([z.string().min(0).max(255), z.string().optional()]),
    color: z.string().max(12).optional(),
    icon: z.string().max(64).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    description: '',
});

async function createCategory(values: Schema): Promise<void> {
    try {
        await getGRPCDocStoreClient().createCategory({
            category: {
                id: '0',
                name: values.name,
                description: values.description,
            },
        });

        notifications.add({
            title: { key: 'notifications.category_created.title', parameters: {} },
            description: { key: 'notifications.category_created.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('update');
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function updateCategory(values: Schema): Promise<void> {
    props.category!.name = values.name;
    props.category!.description = values.description;
    props.category!.color = values.color;
    props.category!.icon = values.icon;

    try {
        await getGRPCDocStoreClient().updateCategory({
            category: props.category!,
        });

        notifications.add({
            title: { key: 'notifications.category_updated.title', parameters: {} },
            description: { key: 'notifications.category_updated.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('update');
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
        await getGRPCDocStoreClient().deleteCategory({
            ids: [props.category.id!],
        });

        notifications.add({
            title: { key: 'notifications.category_deleted.title', parameters: {} },
            description: { key: 'notifications.category_deleted.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('update');
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await (props.category === undefined ? createCategory(event.data) : updateCategory(event.data)).finally(() =>
        useTimeoutFn(() => (canSubmit.value = true), 400),
    );
}, 1000);

function setFromProps(): void {
    if (!props.category) {
        return;
    }

    state.name = props.category.name;
    state.description = props.category.description;
    state.color = props.category.color;
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
                            <template v-if="category">
                                {{ $t('components.documents.categories.modal.update_category') }}:
                                {{ category?.name }}
                            </template>
                            <template v-else>
                                {{ $t('components.documents.categories.modal.create_category') }}
                            </template>
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <div>
                        <UFormGroup name="name" :label="$t('common.name')" class="flex-1">
                            <UInput
                                v-model="state.name"
                                type="text"
                                name="name"
                                :placeholder="$t('common.name', 1)"
                                :label="$t('common.name', 1)"
                            />
                        </UFormGroup>

                        <UFormGroup name="description" :label="$t('common.description')" class="flex-1">
                            <UTextarea v-model="state.description" name="description" :placeholder="$t('common.description')" />
                        </UFormGroup>

                        <UFormGroup name="color" :label="$t('common.color')" class="flex-1 flex-row">
                            <div class="flex flex-1 gap-1">
                                <ColorPicker v-model="state.color" :block="true" class="flex-1" />

                                <UButton icon="i-mdi-backspace" @click="state.color = undefined" />
                            </div>
                        </UFormGroup>

                        <UFormGroup name="icon" :label="$t('common.icon')" class="flex-1">
                            <div class="flex flex-1 gap-1">
                                <IconSelectMenu v-model="state.icon" class="flex-1" :fallback-icon="ShapeIcon" />

                                <UButton icon="i-mdi-backspace" @click="state.icon = undefined" />
                            </div>
                        </UFormGroup>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ category === undefined ? $t('common.create') : $t('common.update') }}
                        </UButton>

                        <UButton
                            v-if="category !== undefined && can('DocStoreService.DeleteCategory').value"
                            block
                            class="flex-1"
                            icon="i-mdi-trash-can"
                            color="red"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                            @click="deleteCategory()"
                        >
                            {{ $t('common.delete') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
