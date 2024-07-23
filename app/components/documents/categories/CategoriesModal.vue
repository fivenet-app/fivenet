<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useNotificatorStore } from '~/store/notificator';
import { Category } from '~~/gen/ts/resources/documents/category';
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
    description: z.string().min(0).max(255),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    description: '',
});

interface FormData {
    name: string;
    description: string;
}

async function createCategory(values: FormData): Promise<void> {
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

async function updateCategory(values: FormData): Promise<void> {
    props.category!.name = values.name;
    props.category!.description = values.description;

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
                        <UFormGroup name="name" :label="$t('common.category', 1)" class="flex-1">
                            <UInput
                                v-model="state.name"
                                type="text"
                                name="name"
                                :placeholder="$t('common.category', 1)"
                                :label="$t('common.category', 1)"
                                :value="category?.name"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <UFormGroup name="description" :label="$t('common.description')" class="flex-1">
                            <UTextarea
                                v-model="state.description"
                                name="description"
                                :placeholder="$t('common.description')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
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
