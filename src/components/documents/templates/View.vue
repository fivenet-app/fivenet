<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { useConfirmDialog } from '@vueuse/core';
import ConfirmDialog from '~/components/partials/ConfirmDialog.vue';
import { useNotificationsStore } from '~/store/notifications';
import { Template, TemplateRequirements } from '~~/gen/ts/resources/documents/templates';
import PreviewModal from './PreviewModal.vue';
import RequirementsList from './RequirementsList.vue';

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const props = defineProps<{
    templateId: bigint;
}>();

const {
    data: template,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`documents-template-${props.templateId}`, () => getTemplate());
const reqs = ref<undefined | TemplateRequirements>();

async function getTemplate(): Promise<Template | undefined> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getDocStoreClient().getTemplate({
                templateId: props.templateId,
                render: false,
            });
            const { response } = await call;

            if (response.template?.schema) {
                reqs.value = response.template?.schema?.requirements;
            }

            return res(response.template!);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function deleteTemplate(id: bigint): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getDocStoreClient().deleteTemplate({
                id: id,
            });

            notifications.dispatchNotification({
                title: { key: 'notifications.templates.deleted.title', parameters: [] },
                content: { key: 'notifications.templates.deleted.content', parameters: [] },
                type: 'success',
            });

            await navigateTo({ name: 'documents-templates' });

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function editTemplate(): Promise<void> {
    await navigateTo({
        name: 'documents-templates-edit-id',
        params: { id: props.templateId.toString() },
    });
}

const openPreview = ref(false);

const { isRevealed, reveal, confirm, cancel, onConfirm } = useConfirmDialog();

onConfirm(async (id) => deleteTemplate(id));
</script>

<template>
    <ConfirmDialog :open="isRevealed" :cancel="cancel" :confirm="() => confirm(templateId)" />

    <PreviewModal :id="templateId" :open="openPreview" @close="openPreview = false" v-if="openPreview" />

    <div v-if="template" class="py-2">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto inline-flex">
                    <button
                        v-if="can('DocStoreService.CreateTemplate')"
                        type="submit"
                        @click="editTemplate()"
                        class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
                    >
                        {{ $t('common.edit') }}
                    </button>
                    <button
                        v-if="can('DocStoreService.CreateTemplate')"
                        type="button"
                        @click="openPreview = true"
                        class="flex justify-center w-full px-3 py-2 ml-4 text-sm font-semibold transition-colors rounded-md bg-accent-600 text-neutral hover:bg-accent-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
                    >
                        {{ $t('common.preview') }}
                    </button>
                </div>
            </div>
            <div class="sm:flex sm:items-center">
                <div>
                    <h2 class="text-white text-2xl">
                        {{ template.title }}
                    </h2>
                    <p class="text-white text-sm">{{ $t('common.description') }}: {{ template.description }}</p>
                </div>
            </div>
            <div class="flow-root mt-4 mb-6">
                <div class="mx-0 -my-2 overflow-x-auto">
                    <label for="weight" class="block text-sm font-medium leading-6 text-gray-100">
                        {{ $t('common.template', 2) }} {{ $t('common.weight') }}
                    </label>
                    <div class="mt-2">
                        <input
                            type="text"
                            name="weight"
                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            disabled
                            :value="template.weight"
                        />
                    </div>
                    <label for="contentTitle" class="block text-sm font-medium leading-6 text-gray-100">
                        {{ $t('common.content') }} {{ $t('common.title') }}
                    </label>
                    <div class="mt-2">
                        <textarea
                            rows="4"
                            name="contentTitle"
                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            disabled
                            :value="template.contentTitle"
                        />
                    </div>
                    <label for="content" class="block text-sm font-medium leading-6 text-gray-100">
                        {{ $t('common.content') }}
                    </label>
                    <div class="mt-2">
                        <textarea
                            rows="4"
                            name="content"
                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            disabled
                            :value="template.content"
                        />
                    </div>
                    <div v-if="template.category">
                        <label for="category" class="block text-sm font-medium leading-6 text-gray-100">
                            {{ $t('common.category') }}
                        </label>
                        <div class="mt-2">
                            <p class="text-sm font-medium leading-6 text-gray-100">
                                {{ template.category?.name }} ({{ $t('common.description') }}:
                                {{ template.category?.description }})
                            </p>
                        </div>
                    </div>
                    <div v-if="reqs">
                        <label for="reqs" class="block text-sm font-medium leading-6 text-gray-100">
                            {{ $t('common.schema') }}
                        </label>
                        <div class="mt-2">
                            <ul
                                class="text-sm font-medium max-w-md space-y-1 text-gray-100 list-disc list-inside dark:text-gray-300"
                            >
                                <li v-if="reqs.users">
                                    <RequirementsList name="User" :specs="reqs.users!" />
                                </li>
                                <li v-if="reqs.vehicles">
                                    <RequirementsList name="Vehicle" :specs="reqs.vehicles!" />
                                </li>
                                <li v-if="reqs.documents">
                                    <RequirementsList name="User" :specs="reqs.documents!" />
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
            <div class="flow-root mt-4">
                <button
                    v-if="can('DocStoreService.DeleteTemplate')"
                    type="submit"
                    @click="reveal"
                    class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-error-600 text-neutral hover:bg-error-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300"
                >
                    {{ $t('common.delete') }}
                </button>
            </div>
        </div>
    </div>
</template>
