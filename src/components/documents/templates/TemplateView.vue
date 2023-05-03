<script lang="ts" setup>
import { DocumentTemplate, TemplateRequirements } from '@fivenet/gen/resources/documents/templates_pb';
import { DeleteTemplateRequest, GetTemplateRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { RpcError } from 'grpc-web';
import TemplateRequirementsList from './TemplateRequirementsList.vue';
import { useNotificationsStore } from '~/store/notifications';

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const props = defineProps({
    templateId: {
        type: Number,
        required: true,
    }
});

const { data: template, pending, refresh, error } = useLazyAsyncData(`documents-template-${props.templateId}`, () => getTemplate());
const reqs = ref<undefined | TemplateRequirements>();

async function getTemplate(): Promise<DocumentTemplate | undefined> {
    return new Promise(async (res, rej) => {
        const req = new GetTemplateRequest();
        req.setTemplateId(props.templateId);

        try {
            const resp = await $grpc.getDocStoreClient().
                getTemplate(req, null);

            if (resp.getTemplate()?.hasSchema()) {
                reqs.value = resp.getTemplate()?.getSchema()?.getRequirements();
            }

            return res(resp.getTemplate()!);
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function deleteTemplate(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new DeleteTemplateRequest();
        req.setId(props.templateId);

        try {
            await $grpc.getDocStoreClient().
                deleteTemplate(req, null);

            notifications.dispatchNotification({
                title: 'Template: Deleted',
                content: 'Template deleted successfully.',
                type: 'success',
            });

            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function editTemplate(): Promise<void> {
    await navigateTo({ name: 'documents-templates-edit-id', params: { id: props.templateId } });
}
</script>

<template>
    <div v-if="template" class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <button type="submit" v-can="'DocStoreService.CreateTemplate'" @click="editTemplate()"
                        class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                        {{ $t('common.edit') }}
                    </button>
                    <button type="submit" v-can="'DocStoreService.DeleteTemplate'" @click="deleteTemplate()"
                        class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-error-600 text-neutral hover:bg-error-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                        {{ $t('common.delete') }}
                    </button>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <label for="content" class="block text-sm font-medium leading-6 text-gray-100">
                        {{ $t('common.content') }} {{ $t('common.title') }}</label>
                    <div class="mt-2">
                        <textarea rows="4" name="content" id="content"
                            class="block w-full rounded-md border-0 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:py-1.5 sm:text-sm sm:leading-6"
                            disabled :value="template.getContentTitle()" />
                    </div>
                    <label for="content" class="block text-sm font-medium leading-6 text-gray-100">
                        {{ $t('common.content') }}</label>
                    <div class="mt-2">
                        <textarea rows="4" name="content" id="content"
                            class="block w-full rounded-md border-0 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:py-1.5 sm:text-sm sm:leading-6"
                            disabled :value="template.getContent()" />
                    </div>
                    <div v-if="reqs">
                        <label for="content" class="block text-sm font-medium leading-6 text-gray-100">
                            {{ $t('common.schema') }}</label>
                        <div class="mt-2">
                            <ul
                                class="text-sm font-medium max-w-md space-y-1 text-gray-100 list-disc list-inside dark:text-gray-300">
                                <li v-if="reqs.hasUsers()">
                                    <TemplateRequirementsList name="User" :required="reqs.getUsers()?.getRequired()!"
                                        :min="reqs.getUsers()?.getMin()!" :max="reqs.getUsers()?.getMax()!" />
                                </li>
                                <li v-if="reqs.hasVehicles()">
                                    <TemplateRequirementsList name="Vehicle" :required="reqs.getVehicles()?.getRequired()!"
                                        :min="reqs.getVehicles()?.getMin()!" :max="reqs.getVehicles()?.getMax()!" />
                                </li>
                                <li v-if="reqs.hasDocuments()">
                                    <TemplateRequirementsList name="User" :required="reqs.getDocuments()?.getRequired()!"
                                        :min="reqs.getDocuments()?.getMin()!" :max="reqs.getDocuments()?.getMax()!" />
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
