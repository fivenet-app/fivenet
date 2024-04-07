<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import { Category } from '~~/gen/ts/resources/documents/category';

const props = defineProps<{
    category?: Category;
}>();

const emit = defineEmits<{
    (e: 'update'): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const { isOpen } = useModal();

interface FormData {
    name: string;
    description: string;
}

async function createCategory(values: FormData): Promise<void> {
    try {
        await $grpc.getDocStoreClient().createCategory({
            category: {
                id: '0',
                name: values.name,
                description: values.description,
            },
        });

        notifications.add({
            title: { key: 'notifications.category_created.title', parameters: {} },
            description: { key: 'notifications.category_created.content', parameters: {} },
            type: 'success',
        });

        emit('update');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function updateCategory(values: FormData): Promise<void> {
    props.category!.name = values.name;
    props.category!.description = values.description;

    try {
        await $grpc.getDocStoreClient().updateCategory({
            category: props.category!,
        });

        notifications.add({
            title: { key: 'notifications.category_updated.title', parameters: {} },
            description: { key: 'notifications.category_updated.content', parameters: {} },
            type: 'success',
        });

        emit('update');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function deleteCategory(): Promise<void> {
    if (props.category === undefined) {
        return;
    }

    try {
        await $grpc.getDocStoreClient().deleteCategory({
            ids: [props.category.id!],
        });

        notifications.add({
            title: { key: 'notifications.category_deleted.title', parameters: {} },
            description: { key: 'notifications.category_deleted.content', parameters: {} },
            type: 'success',
        });

        emit('update');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 128 },
        description: { required: true, min: 0, max: 255 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await (props.category === undefined ? createCategory(values) : updateCategory(values)).finally(() =>
            useTimeoutFn(() => (canSubmit.value = true), 400),
        ),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :state="{}" @submit="onSubmitThrottle">
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
                        <div class="text-sm text-gray-100">
                            <div class="flex-1">
                                <label for="name" class="block text-sm font-medium leading-6">
                                    {{ $t('common.category', 1) }}
                                </label>
                                <VeeField
                                    type="text"
                                    name="name"
                                    :placeholder="$t('common.category', 1)"
                                    :label="$t('common.category', 1)"
                                    :value="category?.name"
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 focus:ring-1 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                                <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
                            </div>
                            <div class="flex-1">
                                <label for="description" class="block text-sm font-medium leading-6">
                                    {{ $t('common.description') }}
                                </label>
                                <VeeField
                                    as="textarea"
                                    name="description"
                                    :placeholder="$t('common.description')"
                                    :label="$t('common.description')"
                                    :value="category?.description"
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 pr-14 focus:ring-1 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                />
                                <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
                            </div>
                        </div>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                        <UButton
                            v-if="category !== undefined && can('DocStoreService.DeleteCategory')"
                            block
                            class="flex-1"
                            :disabled="!meta.valid || !canSubmit"
                            :loading="!canSubmit"
                            @click="deleteCategory()"
                        >
                            {{ $t('common.delete') }}
                        </UButton>
                        <UButton block class="flex-1" :disabled="!meta.valid || !canSubmit" :loading="!canSubmit">
                            {{ category === undefined ? $t('common.create') : $t('common.update') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
