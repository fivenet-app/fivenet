<script lang="ts" setup>
import TemplatePreviewModal from '~/components/documents/templates/TemplatePreviewModal.vue';
import TemplateRequirementsList from '~/components/documents/templates/TemplateRequirementsList.vue';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums, type AccessType } from '~/components/partials/access/helpers';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import DocumentCategoryBadge from '~/components/partials/documents/DocumentCategoryBadge.vue';
import { useNotificatorStore } from '~/stores/notificator';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { Template, TemplateRequirements } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    templateId: number;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { can } = useAuth();

const modal = useModal();

const notifications = useNotificatorStore();

const reqs = ref<undefined | TemplateRequirements>();

const {
    data: template,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`documents-template-${props.templateId}`, () => getTemplate());

async function getTemplate(): Promise<Template | undefined> {
    try {
        const call = $grpc.docstore.docStore.getTemplate({
            templateId: props.templateId,
            render: false,
        });
        const { response } = await call;

        if (response.template?.schema) {
            reqs.value = response.template?.schema?.requirements;
        }

        return response.template!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteTemplate(id: number): Promise<void> {
    try {
        await $grpc.docstore.docStore.deleteTemplate({ id });

        notifications.add({
            title: { key: 'notifications.templates.deleted.title', parameters: {} },
            description: { key: 'notifications.templates.deleted.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await navigateTo({ name: 'documents-templates' });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const templateAccessTypes: AccessType[] = [{ type: 'job', name: t('common.job', 2) }];
const contentAccessTypes: AccessType[] = [
    { type: 'user', name: t('common.citizen', 2) },
    { type: 'job', name: t('common.job', 2) },
];
</script>

<template>
    <UDashboardNavbar :title="$t('pages.documents.templates.view.title')">
        <template #right>
            <PartialsBackButton to="/documents/templates" />

            <UButtonGroup v-if="template" class="inline-flex">
                <UButton
                    v-if="can('DocStoreService.CreateTemplate').value"
                    block
                    class="flex-1"
                    color="white"
                    trailing-icon="i-mdi-print-preview"
                    @click="
                        modal.open(TemplatePreviewModal, {
                            templateId: templateId,
                        })
                    "
                >
                    {{ $t('common.preview') }}
                </UButton>

                <UButton
                    v-if="can('DocStoreService.CreateTemplate').value"
                    block
                    class="flex-1"
                    trailing-icon="i-mdi-pencil"
                    :to="{ name: 'documents-templates-edit-id', params: { id: templateId } }"
                >
                    {{ $t('common.edit') }}
                </UButton>

                <UButton
                    v-if="can('DocStoreService.DeleteTemplate').value"
                    block
                    class="flex-1"
                    trailing-icon="i-mdi-trash-can"
                    color="red"
                    @click="
                        modal.open(ConfirmModal, {
                            confirm: async () => deleteTemplate(templateId),
                        })
                    "
                >
                    {{ $t('common.delete') }}
                </UButton>
            </UButtonGroup>
        </template>
    </UDashboardNavbar>

    <UDashboardPanelContent class="p-0 sm:pb-0">
        <UContainer class="w-full">
            <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.template', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.template', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock v-else-if="!template" :type="$t('common.template', 2)" />

            <template v-else>
                <div class="mt-2 sm:flex sm:items-center">
                    <div>
                        <h2 class="inline-flex items-center gap-2 text-2xl">
                            <UIcon
                                :name="template.icon ?? 'i-mdi-file-outline'"
                                :class="`text-${template.color}-500 dark:text-${template.color}-400`"
                            />

                            <span>
                                {{ template.title }}
                            </span>
                        </h2>

                        <p class="text-base">
                            <span class="font-semibold">{{ $t('common.description') }}:</span> {{ template.description }}
                        </p>
                    </div>
                </div>

                <div class="mb-6 mt-4 flow-root">
                    <div class="-my-2 mx-0 overflow-x-auto">
                        <div class="my-2">
                            <h3 class="block text-base font-medium leading-6 text-gray-100">
                                {{ $t('common.template', 2) }} {{ $t('common.weight') }}
                            </h3>
                            <div class="my-2">
                                <UInput type="text" name="weight" disabled :value="template.weight" />
                            </div>
                        </div>

                        <div v-if="template.jobAccess" class="my-2">
                            <h3 class="block text-base font-medium leading-6 text-gray-100">
                                {{ $t('common.template', 2) }} {{ $t('common.access') }}
                            </h3>
                            <div class="my-2">
                                <AccessManager
                                    v-model:jobs="template.jobAccess"
                                    :target-id="templateId ?? 0"
                                    :access-roles="
                                        enumToAccessLevelEnums(AccessLevel, 'enums.docstore.AccessLevel').filter(
                                            (e) => e.value === AccessLevel.VIEW || e.value === AccessLevel.EDIT,
                                        )
                                    "
                                    :access-types="templateAccessTypes"
                                    :disabled="true"
                                />
                            </div>
                        </div>

                        <div class="my-2">
                            <h3 class="block text-base font-medium leading-6 text-gray-100">
                                {{ $t('common.content') }} {{ $t('common.title') }}
                            </h3>
                            <div class="my-2">
                                <UTextarea
                                    name="contentTitle"
                                    class="w-full whitespace-pre-wrap"
                                    disabled
                                    resize
                                    :rows="3"
                                    :value="template.contentTitle"
                                />
                            </div>
                        </div>

                        <div v-if="template.state">
                            <h3 class="block text-base font-medium leading-6 text-gray-100">
                                {{ $t('common.content') }} {{ $t('common.state') }}
                            </h3>
                            <div class="my-2">
                                <UInput
                                    type="text"
                                    name="state"
                                    class="block w-full whitespace-pre-wrap rounded-md border-0 bg-base-900 py-1.5 focus:ring-1 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    disabled
                                    :value="template.state"
                                />
                            </div>
                        </div>

                        <div v-if="template.category">
                            <h3 class="block text-base font-medium leading-6 text-gray-100">
                                {{ $t('common.category') }}
                            </h3>
                            <div class="my-2">
                                <DocumentCategoryBadge :category="template.category" />
                            </div>
                        </div>

                        <div class="my-2">
                            <h3 class="block text-base font-medium leading-6 text-gray-100">
                                {{ $t('common.content') }}
                            </h3>
                            <div class="my-2">
                                <UTextarea
                                    name="content"
                                    class="w-full whitespace-pre-wrap"
                                    disabled
                                    resize
                                    :rows="4"
                                    :value="template.content"
                                />
                            </div>
                        </div>

                        <div v-if="reqs">
                            <h3 class="block text-base font-medium leading-6 text-gray-100">
                                {{ $t('common.schema') }}
                            </h3>
                            <div class="my-2">
                                <ul
                                    class="mb-2 max-w-md list-inside list-disc space-y-1 text-sm font-medium text-gray-100 dark:text-gray-300"
                                >
                                    <li v-if="reqs.users">
                                        <TemplateRequirementsList name="User" :specs="reqs.users!" />
                                    </li>
                                    <li v-if="reqs.vehicles">
                                        <TemplateRequirementsList name="Vehicle" :specs="reqs.vehicles!" />
                                    </li>
                                    <li v-if="reqs.documents">
                                        <TemplateRequirementsList name="User" :specs="reqs.documents!" />
                                    </li>
                                </ul>
                            </div>
                        </div>

                        <div v-if="template.contentAccess" class="my-2">
                            <h3 class="block text-base font-medium leading-6 text-gray-100">
                                {{ $t('common.access') }}
                            </h3>
                            <div class="my-2">
                                <AccessManager
                                    v-model:jobs="template.contentAccess.jobs"
                                    :target-id="templateId ?? 0"
                                    :access-types="contentAccessTypes"
                                    :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.docstore.AccessLevel')"
                                    :disabled="true"
                                    :show-required="true"
                                />
                            </div>
                        </div>

                        <div v-if="!template.workflow">
                            {{ $t('common.none', [$t('common.workflow')]) }}
                        </div>
                        <div v-else>
                            <h3 class="block text-base font-medium leading-6 text-gray-100">
                                {{ $t('common.auto_close') }}
                            </h3>

                            <div class="my-2">
                                <div class="flex gap-2">
                                    <span
                                        ><span class="font-semibold">{{ $t('common.enabled') }}:</span>
                                        {{ $t(template.workflow?.autoClose ? 'common.yes' : 'common.no') }}</span
                                    >

                                    <span v-if="template.workflow?.autoCloseSettings?.duration">
                                        <span class="font-semibold">{{ $t('common.time_ago.day', 2) }}:</span>
                                        {{
                                            (template.workflow.autoCloseSettings.duration.seconds / 24 / 60 / 60).toFixed(0)
                                        }}</span
                                    >
                                </div>

                                <span v-if="template.workflow?.autoCloseSettings?.message">
                                    <span class="font-semibold">{{ $t('common.message') }}:</span>
                                    "{{ template.workflow.autoCloseSettings.message }}"</span
                                >
                            </div>

                            <h3 class="block text-base font-medium leading-6 text-gray-100">
                                {{ $t('common.reminder', 2) }}
                            </h3>

                            <div class="my-2">
                                <span
                                    >{{ $t('common.enabled') }}:
                                    {{ $t(template.workflow?.reminder ? 'common.yes' : 'common.no') }}</span
                                >

                                <ol class="list-inside list-decimal">
                                    <li
                                        v-for="(reminder, idx) in template.workflow?.reminderSettings?.reminders"
                                        :key="idx"
                                        class="gap-2"
                                    >
                                        <div class="inline-flex gap-2">
                                            <span>
                                                <span class="font-semibold">{{ $t('common.time_ago.day', 2) }}:</span>
                                                {{ ((reminder?.duration?.seconds ?? 0) / 24 / 60 / 60).toFixed(0) }}
                                            </span>

                                            <span>
                                                <span class="font-semibold">{{ $t('common.message') }}:</span>
                                                "{{ reminder.message }}"</span
                                            >
                                        </div>
                                    </li>
                                </ol>
                            </div>
                        </div>
                    </div>
                </div>
            </template>
        </UContainer>
    </UDashboardPanelContent>
</template>
