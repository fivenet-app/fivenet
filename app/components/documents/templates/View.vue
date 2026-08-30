<script lang="ts" setup>
import PreviewModal from '~/components/documents/templates/PreviewModal.vue';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums, type AccessType } from '~/components/partials/access/helpers';
import ColorPickerTW from '~/components/partials/ColorPickerTW.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import CategoryBadge from '~/components/partials/documents/CategoryBadge.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { TemplateBlock, TemplateBlockEnd } from '~/composables/tiptap/extensions/TemplateBlock';
import { TemplateVar } from '~/composables/tiptap/extensions/TemplateVar';
import type { ResponsiveActionEntry } from '~/components/partials/ResponsiveActions.types';
import { getDocumentsTemplatesClient } from '~~/gen/ts/clients';
import { AccessLevel } from '~~/gen/ts/resources/documents/access/access';
import { ApprovalAssigneeKind } from '~~/gen/ts/resources/documents/approval/approval';
import type { Template, TemplateRequirements } from '~~/gen/ts/resources/documents/templates/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import PolicyEditor from '../approval/PolicyEditor.vue';
import ApprovalTasksEditor from './editor/ApprovalTasksEditor.vue';
import SchemaEditor from './editor/SchemaEditor.vue';
import ResponsiveActions from '~/components/partials/ResponsiveActions.vue';

const props = defineProps<{
    templateId: number;
}>();

const { t } = useI18n();

const { can } = useAuth();

const overlay = useOverlay();

const notifications = useNotificationsStore();

const documentsTemplatesClient = await getDocumentsTemplatesClient();

const reqs = ref<undefined | TemplateRequirements>();

const {
    data: template,
    status,
    refresh,
    error,
} = useLazyAsyncData(`documents-template-${props.templateId}`, () => getTemplate());

async function getTemplate(): Promise<Template | undefined> {
    try {
        const call = documentsTemplatesClient.getTemplate({
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
        await documentsTemplatesClient.deleteTemplate({ id });

        notifications.add({
            title: { key: 'notifications.templates.deleted.title', parameters: {} },
            description: { key: 'notifications.templates.deleted.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        await navigateTo('/documents/templates');
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

const actionItems = computed<ResponsiveActionEntry[]>(() => {
    const items: ResponsiveActionEntry[] = [];

    if (can('documents.TemplatesService/CreateTemplate').value) {
        items.push({
            label: t('common.preview'),
            tooltip: t('common.preview'),
            icon: 'i-mdi-print-preview',
            color: 'neutral',
            variant: 'ghost',
            onClick: () => {
                templatePreviewModal.open({
                    templateId: props.templateId,
                });
            },
        });

        items.push({
            label: t('common.edit'),
            tooltip: t('common.edit'),
            icon: 'i-mdi-pencil',
            color: 'neutral',
            variant: 'ghost',
            to: `/documents/templates/${props.templateId}/edit`,
        });
    }

    if (can('documents.TemplatesService/DeleteTemplate').value) {
        if (items.length > 0) {
            items.push({ kind: 'separator' });
        }

        items.push({
            label: t('common.delete'),
            tooltip: t('common.delete'),
            icon: 'i-mdi-delete',
            color: 'error',
            variant: 'ghost',
            onClick: () => {
                confirmModal.open({
                    confirm: async () => deleteTemplate(props.templateId),
                });
            },
        });
    }

    return items;
});

const confirmModal = overlay.create(ConfirmModal);
const templatePreviewModal = overlay.create(PreviewModal, { props: { templateId: props.templateId } });
</script>

<template>
    <UDashboardPanel :ui="{ root: 'pb-(--page-content-bottom-offset)' }">
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
                <ResponsiveActions :items="actionItems" :label="$t('common.action', 2)" />
            </UDashboardToolbar>

            <UDashboardToolbar v-if="template">
                <template #default>
                    <div class="mx-auto my-2 w-full max-w-(--breakpoint-xl)">
                        <div class="mb-2">
                            <h1 class="inline-flex items-center gap-2 px-0.5 py-1 text-4xl font-bold break-words sm:pl-1">
                                <UIcon
                                    class="shrink-0"
                                    :class="`text-${template.color ?? 'primary'}`"
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
                            <UFormField name="color" :label="$t('common.color')">
                                <div class="flex flex-1 gap-1">
                                    <ColorPickerTW v-model="template.color" class="flex-1" disabled />
                                </div>
                            </UFormField>

                            <UFormField :label="`${$t('common.template', 2)} ${$t('common.access')}`">
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
                                <UInput class="w-full" type="text" name="state" disabled :value="template.state" />
                            </UFormField>

                            <UFormField v-if="template.category" :label="$t('common.category')">
                                <CategoryBadge :category="template.category" />
                            </UFormField>

                            <UCollapsible>
                                <UButton
                                    class="group"
                                    block
                                    color="neutral"
                                    variant="subtle"
                                    trailing-icon="i-mdi-chevron-down"
                                    :label="$t('common.content')"
                                    :ui="{
                                        trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                                    }"
                                />

                                <template #content>
                                    <TiptapEditor
                                        :model-value="template.content"
                                        class="mt-2 min-h-64"
                                        content-type="html"
                                        disabled
                                        hide-toolbar
                                        :extensions="[TemplateVar, TemplateBlock, TemplateBlockEnd]"
                                    />
                                </template>
                            </UCollapsible>

                            <UFormField v-if="template.contentAccess" :label="$t('common.access')">
                                <AccessManager
                                    v-model:jobs="template.contentAccess.jobs"
                                    :target-id="templateId ?? 0"
                                    :access-types="contentAccessTypes"
                                    :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.documents.AccessLevel')"
                                    disabled
                                    required-mode="badge"
                                />
                            </UFormField>
                        </UPageCard>

                        <UPageCard v-if="reqs" :title="$t('common.requirements')">
                            <SchemaEditor :model-value="reqs" disabled />
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

                                <UFormField class="flex-1" :label="$t('common.message')">
                                    <UInput
                                        class="w-full"
                                        :model-value="template.workflow?.autoCloseSettings?.message ?? $t('common.na')"
                                        disabled
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

                        <UPageCard :title="$t('components.documents.approval.policy_form.title', 2)">
                            <UFormField name="approval.enabled" :label="$t('common.enabled')">
                                <USwitch :model-value="template?.approval?.enabled" disabled />
                            </UFormField>

                            <PolicyEditor
                                v-if="template?.approval?.policy"
                                :model-value="template?.approval?.policy"
                                disabled
                            />
                        </UPageCard>

                        <UPageCard :title="$t('components.documents.approval.tasks', 2)">
                            <ApprovalTasksEditor
                                :model-value="
                                    (template?.approval?.tasks ?? [])?.map((task) => ({
                                        ...task,
                                        ruleKind: task.userId ? ApprovalAssigneeKind.USER : ApprovalAssigneeKind.JOB_GRADE,
                                    }))
                                "
                                disabled
                                :signature-required="template?.approval?.policy?.signatureRequired"
                            />
                        </UPageCard>
                    </div>
                </template>
            </UContainer>
        </template>
    </UDashboardPanel>
</template>
