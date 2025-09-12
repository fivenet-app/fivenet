<script lang="ts" setup>
import TemplatePreviewModal from '~/components/documents/templates/TemplatePreviewModal.vue';
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
import TemplateSchemaEditor from './TemplateSchemaEditor.vue';

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
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <template #default>
                    <div
                        class="mx-auto flex w-full max-w-(--breakpoint-xl) flex-1 snap-x flex-row flex-wrap justify-between gap-2 overflow-x-auto"
                    >
                        <UTooltip
                            v-if="can('documents.DocumentsService/CreateTemplate').value"
                            class="flex-1"
                            :text="$t('common.preview')"
                        >
                            <UButton
                                class="flex-1"
                                block
                                color="neutral"
                                variant="ghost"
                                icon="i-mdi-print-preview"
                                :label="$t('common.preview')"
                                @click="
                                    templatePreviewModal.open({
                                        templateId: templateId,
                                    })
                                "
                            />
                        </UTooltip>

                        <UTooltip
                            v-if="can('documents.DocumentsService/CreateTemplate').value"
                            class="flex-1"
                            :text="$t('common.edit')"
                        >
                            <UButton
                                class="flex-1"
                                block
                                color="neutral"
                                variant="ghost"
                                icon="i-mdi-pencil"
                                :to="{ name: 'documents-templates-edit-id', params: { id: templateId } }"
                                :label="$t('common.edit')"
                            />
                        </UTooltip>

                        <UTooltip
                            v-if="can('documents.DocumentsService/DeleteTemplate').value"
                            class="flex-1"
                            :text="$t('common.delete')"
                        >
                            <UButton
                                class="flex-1"
                                block
                                color="error"
                                variant="ghost"
                                icon="i-mdi-delete"
                                :label="$t('common.delete')"
                                @click="
                                    confirmModal.open({
                                        confirm: async () => deleteTemplate(templateId),
                                    })
                                "
                            />
                        </UTooltip>
                    </div>
                </template>
            </UDashboardToolbar>

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

                            <p class="line-clamp-3 text-base">
                                <span class="font-semibold">{{ $t('common.description') }}:</span> {{ template.description }}
                            </p>
                        </div>
                    </div>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <UContainer class="mx-auto max-w-(--ui-container)">
                <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.template', 2)])" />
                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.template', 2)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock v-else-if="!template" :type="$t('common.template', 2)" />

                <template v-else>
                    <div class="flex flex-col gap-4">
                        <UPageCard :title="$t('common.detail', 2)">
                            <UFormField :label="`${$t('common.template', 2)} ${$t('common.weight')}`">
                                <UInputNumber type="text" name="weight" disabled :value="template.weight" class="w-full" />
                            </UFormField>

                            <UFormField v-if="template.jobAccess" :label="`${$t('common.template', 2)} ${$t('common.access')}`">
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
                            </UFormField>
                        </UPageCard>

                        <UPageCard :title="$t('common.content')">
                            <UFormField :label="$t('common.title')">
                                <UTextarea
                                    class="w-full whitespace-pre-wrap"
                                    name="contentTitle"
                                    disabled
                                    resize
                                    :rows="3"
                                    :value="template.contentTitle"
                                />
                            </UFormField>

                            <UFormField v-if="template.state" :label="$t('common.state')">
                                <UInput type="text" name="state" disabled :value="template.state" class="w-full" />
                            </UFormField>

                            <UFormField v-if="template.category" :label="$t('common.category')">
                                <CategoryBadge :category="template.category" />
                            </UFormField>

                            <UFormField :label="$t('common.content')">
                                <UTextarea
                                    class="w-full whitespace-pre-wrap"
                                    name="content"
                                    disabled
                                    resize
                                    :rows="4"
                                    :value="template.content"
                                />
                            </UFormField>

                            <UFormField v-if="template.contentAccess" :label="$t('common.access')">
                                <AccessManager
                                    v-model:jobs="template.contentAccess.jobs"
                                    :target-id="templateId ?? 0"
                                    :access-types="contentAccessTypes"
                                    :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.documents.AccessLevel')"
                                    disabled
                                    show-required
                                />
                            </UFormField>
                        </UPageCard>

                        <UPageCard v-if="reqs" :title="$t('common.requirements')">
                            <TemplateSchemaEditor :model-value="reqs" disabled />
                        </UPageCard>

                        <UPageCard :title="`${$t('common.workflow')}: ${$t('common.auto_close')}`">
                            <span v-if="!template.workflow">
                                {{ $t('common.none', [$t('common.workflow')]) }}
                            </span>
                            <template v-else>
                                <UFormField :label="$t('common.enabled')">
                                    <USwitch
                                        :model-value="template.workflow?.autoClose"
                                        disabled
                                        :label="$t(template.workflow?.autoClose ? 'common.yes' : 'common.no')"
                                    />
                                </UFormField>

                                <UFormField :label="$t('common.duration')">
                                    <div class="inline-flex items-center gap-2">
                                        <UInputNumber
                                            :model-value="
                                                parseInt(
                                                    (
                                                        (template.workflow.autoCloseSettings?.duration?.seconds ?? 0) /
                                                        24 /
                                                        60 /
                                                        60
                                                    ).toFixed(0),
                                                )
                                            "
                                            disabled
                                        />
                                        <span>{{ $t('common.time_ago.day', 2) }}</span>
                                    </div>
                                </UFormField>

                                <UFormField :label="$t('common.message')" class="flex-1">
                                    <UInput
                                        :model-value="template.workflow?.autoCloseSettings?.message ?? $t('common.na')"
                                        disabled
                                        class="w-full"
                                    />
                                </UFormField>
                            </template>
                        </UPageCard>

                        <UPageCard :title="`${$t('common.workflow')}: ${$t('common.reminder', 2)}`">
                            <span v-if="!template.workflow">
                                {{ $t('common.none', [$t('common.workflow')]) }}
                            </span>
                            <template v-else>
                                <UFormField :label="$t('common.enabled')">
                                    <USwitch
                                        :model-value="template.workflow?.reminder"
                                        disabled
                                        :label="$t(template.workflow?.reminder ? 'common.yes' : 'common.no')"
                                    />
                                </UFormField>

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
                            </template>
                        </UPageCard>
                    </div>
                </template>
            </UContainer>
        </template>
    </UDashboardPanel>
</template>
