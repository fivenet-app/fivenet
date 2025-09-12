<script lang="ts" setup>
import TemplatePreviewModal from '~/components/documents/templates/TemplatePreviewModal.vue';
import TemplateRequirementsList from '~/components/documents/templates/TemplateRequirementsList.vue';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums, type AccessType } from '~/components/partials/access/helpers';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import CategoryBadge from '~/components/partials/documents/CategoryBadge.vue';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { Template, TemplateRequirements } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    templateId: number;
}>();

const { t } = useI18n();

const { can } = useAuth();

const overlay = useOverlay();

const notifications = useNotificationsStore();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const reqs = ref<undefined | TemplateRequirements>();

const {
    data: template,
    status,
    refresh,
    error,
} = useLazyAsyncData(`documents-template-${props.templateId}`, () => getTemplate());

async function getTemplate(): Promise<Template | undefined> {
    try {
        const call = documentsDocumentsClient.getTemplate({
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

useHead({
    title: () =>
        template.value?.title
            ? `${template.value?.title} - ${t('pages.documents.templates.view.title')}`
            : t('pages.documents.templates.view.title'),
});

async function deleteTemplate(id: number): Promise<void> {
    try {
        await documentsDocumentsClient.deleteTemplate({ id });

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

const templateAccessTypes: AccessType[] = [{ label: t('common.job', 2), value: 'job' }];
const contentAccessTypes: AccessType[] = [
    { label: t('common.citizen', 2), value: 'user' },
    { label: t('common.job', 2), value: 'job' },
];

const confirmModal = overlay.create(ConfirmModal);
const templatePreviewModal = overlay.create(TemplatePreviewModal, { props: { templateId: props.templateId } });
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('pages.documents.templates.view.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton to="/documents/templates" />

                    <UButtonGroup v-if="template" class="inline-flex">
                        <UButton
                            v-if="can('documents.DocumentsService/CreateTemplate').value"
                            class="flex-1"
                            block
                            color="neutral"
                            trailing-icon="i-mdi-print-preview"
                            @click="
                                templatePreviewModal.open({
                                    templateId: templateId,
                                })
                            "
                        >
                            {{ $t('common.preview') }}
                        </UButton>

                        <UButton
                            v-if="can('documents.DocumentsService/CreateTemplate').value"
                            class="flex-1"
                            block
                            trailing-icon="i-mdi-pencil"
                            :to="{ name: 'documents-templates-edit-id', params: { id: templateId } }"
                        >
                            {{ $t('common.edit') }}
                        </UButton>

                        <UButton
                            v-if="can('documents.DocumentsService/DeleteTemplate').value"
                            class="flex-1"
                            block
                            trailing-icon="i-mdi-delete"
                            color="error"
                            @click="
                                confirmModal.open({
                                    confirm: async () => deleteTemplate(templateId),
                                })
                            "
                        >
                            {{ $t('common.delete') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <template #default>
                    <div v-if="template" class="mx-auto my-2 w-full max-w-(--breakpoint-xl)">
                        <div class="mb-4">
                            <h1 class="inline-flex items-center gap-2 px-0.5 py-1 text-4xl font-bold break-words sm:pl-1">
                                <UIcon
                                    class="shrink-0"
                                    :class="`text-${template.color}-500 dark:text-${template.color}-400`"
                                    :name="
                                        template.icon ? convertComponentIconNameToDynamic(template.icon) : 'i-mdi-file-outline'
                                    "
                                />

                                <span>{{ template.title }}</span>
                            </h1>

                            <p class="text-base">
                                <span class="font-semibold">{{ $t('common.description') }}:</span> {{ template.description }}
                            </p>
                        </div>
                    </div>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <UContainer class="w-full">
                <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.template', 2)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.template', 2)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!template" :type="$t('common.template', 2)" />

                <template v-else>
                    <div class="mx-0 -my-2 flex flex-col gap-y-2 overflow-x-auto">
                        <div>
                            <h3 class="block text-base leading-6 font-medium text-toned">
                                {{ $t('common.template', 2) }} {{ $t('common.weight') }}
                            </h3>
                            <div class="my-2">
                                <UInputNumber type="text" name="weight" disabled :value="template.weight" />
                            </div>
                        </div>

                        <div v-if="template.jobAccess">
                            <h3 class="block text-base leading-6 font-medium text-toned">
                                {{ $t('common.template', 2) }} {{ $t('common.access') }}
                            </h3>

                            <div class="my-2">
                                <AccessManager
                                    v-model:jobs="template.jobAccess"
                                    :target-id="templateId ?? 0"
                                    :access-roles="
                                        enumToAccessLevelEnums(AccessLevel, 'enums.documents.AccessLevel').filter(
                                            (e) => e.value === AccessLevel.VIEW || e.value === AccessLevel.EDIT,
                                        )
                                    "
                                    :access-types="templateAccessTypes"
                                    disabled
                                    name="jobAccess"
                                    full-name
                                />
                            </div>
                        </div>

                        <div>
                            <h3 class="block text-base leading-6 font-medium text-toned">
                                {{ $t('common.content') }} {{ $t('common.title') }}
                            </h3>
                            <div class="my-2">
                                <UTextarea
                                    class="w-full whitespace-pre-wrap"
                                    name="contentTitle"
                                    disabled
                                    resize
                                    :rows="3"
                                    :value="template.contentTitle"
                                />
                            </div>
                        </div>

                        <div v-if="template.state">
                            <h3 class="block text-base leading-6 font-medium text-toned">
                                {{ $t('common.content') }} {{ $t('common.state') }}
                            </h3>
                            <div class="my-2">
                                <UInput
                                    class="focus:ring-base-300 block w-full rounded-md border-0 bg-neutral-900 py-1.5 whitespace-pre-wrap focus:ring-1 focus:ring-inset sm:text-sm sm:leading-6"
                                    type="text"
                                    name="state"
                                    disabled
                                    :value="template.state"
                                />
                            </div>
                        </div>

                        <div v-if="template.category">
                            <h3 class="block text-base leading-6 font-medium text-toned">
                                {{ $t('common.category') }}
                            </h3>
                            <div class="my-2">
                                <CategoryBadge :category="template.category" />
                            </div>
                        </div>

                        <div>
                            <h3 class="block text-base leading-6 font-medium text-toned">
                                {{ $t('common.content') }}
                            </h3>
                            <div class="my-2">
                                <UTextarea
                                    class="w-full whitespace-pre-wrap"
                                    name="content"
                                    disabled
                                    resize
                                    :rows="4"
                                    :value="template.content"
                                />
                            </div>
                        </div>

                        <div v-if="reqs">
                            <h3 class="block text-base leading-6 font-medium text-toned">
                                {{ $t('common.schema') }}
                            </h3>

                            <div class="my-2">
                                <ul class="mb-2 max-w-md list-inside space-y-1 text-sm font-medium text-toned">
                                    <li v-if="reqs.users">
                                        <TemplateRequirementsList name="User" :specs="reqs.users!" />
                                    </li>
                                    <li v-if="reqs.vehicles">
                                        <TemplateRequirementsList name="Vehicle" :specs="reqs.vehicles!" />
                                    </li>
                                    <li v-if="reqs.documents">
                                        <TemplateRequirementsList name="Document" :specs="reqs.documents!" />
                                    </li>
                                </ul>
                            </div>
                        </div>

                        <div v-if="template.contentAccess">
                            <h3 class="block text-base leading-6 font-medium text-toned">
                                {{ $t('common.access') }}
                            </h3>

                            <div class="my-2">
                                <AccessManager
                                    v-model:jobs="template.contentAccess.jobs"
                                    :target-id="templateId ?? 0"
                                    :access-types="contentAccessTypes"
                                    :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.documents.AccessLevel')"
                                    disabled
                                    show-required
                                />
                            </div>
                        </div>

                        <div v-if="!template.workflow">
                            {{ $t('common.none', [$t('common.workflow')]) }}
                        </div>
                        <div v-else>
                            <h3 class="block text-base leading-6 font-medium text-toned">
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

                            <h3 class="block text-base leading-6 font-medium text-toned">
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
                </template>
            </UContainer>
        </template>
    </UDashboardPanel>
</template>
