<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { FileOutlineIcon } from 'mdi-vue3';
import { z } from 'zod';
import SingleHint from '~/components/SingleHint.vue';
import TemplateSchemaEditor from '~/components/documents/templates/TemplateSchemaEditor.vue';
import { zWorkflow } from '~/components/documents/templates/types';
import ColorPickerTW from '~/components/partials/ColorPickerTW.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import SelectMenu from '~/components/partials/SelectMenu.vue';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums, type AccessType } from '~/components/partials/access/helpers';
import CategoryBadge from '~/components/partials/documents/CategoryBadge.vue';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { TemplateBlock } from '~/composables/tiptap/extensions/TemplateBlock';
import { TemplateVar } from '~/composables/tiptap/extensions/TemplateVar';
import { useAuthStore } from '~/stores/auth';
import { useCompletorStore } from '~/stores/completor';
import { jobAccessEntry, userAccessEntry } from '~/utils/validation';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { ApprovalAssigneeKind, ApprovalRuleKind, OnEditBehavior } from '~~/gen/ts/resources/documents/approval';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { Template, TemplateRequirements } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { CreateTemplateRequest, UpdateTemplateRequest } from '~~/gen/ts/services/documents/documents';
import PolicyEditor from '../approval/PolicyEditor.vue';
import ApprovalTasksEditor from './ApprovalTasksEditor.vue';
import EditorButtons from './EditorButtons.vue';
import TemplateWorkflowEditor from './TemplateWorkflowEditor.vue';

const props = defineProps<{
    templateId?: number;
}>();

const { t } = useI18n();

const { game } = useAppConfig();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const notifications = useNotificationsStore();

const completorStore = useCompletorStore();

const { maxAccessEntries } = useAppConfig();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const schema = z.object({
    weight: z.coerce.number().min(0).max(999_999),
    title: z.coerce.string().min(3).max(255),
    description: z.coerce.string().min(3).max(255),
    color: z.coerce.string().max(7),
    icon: z.coerce.string().max(128).optional(),
    contentTitle: z.coerce.string().min(3).max(2048),
    content: z.coerce.string().min(3).max(1500000),
    contentState: z.union([z.coerce.string().min(1).max(512), z.coerce.string().length(0)]),
    category: z.custom<Category>().optional(),
    jobAccess: jobAccessEntry.array().max(maxAccessEntries).default([]),
    contentAccess: z.object({
        jobs: jobAccessEntry.array().max(maxAccessEntries).default([]),
        users: userAccessEntry.array().max(maxAccessEntries).default([]),
    }),
    workflow: zWorkflow,
    approval: z
        .object({
            enabled: z.boolean().default(false),

            policy: z.object({
                ruleKind: z.enum(ApprovalRuleKind).default(ApprovalRuleKind.REQUIRE_ALL),
                onEditBehavior: z.enum(OnEditBehavior).default(OnEditBehavior.KEEP_PROGRESS),
                requiredCount: z.number().min(1).max(10).default(2),
                signatureRequired: z.boolean().default(false),
            }),

            tasks: z
                .union([
                    z.object({
                        ruleKind: z.enum(ApprovalAssigneeKind).default(ApprovalAssigneeKind.JOB_GRADE),
                        userId: z.coerce.number(),
                        job: z.coerce.string().optional(),
                        minimumGrade: z.coerce.number().min(game.startJobGrade).optional(),
                        label: z.string().max(120).optional(),
                        signatureRequired: z.coerce.boolean().default(false),
                        slots: z.coerce.number().min(1).max(10).optional().default(1),
                        dueInDays: z.coerce.number().min(1).optional(),
                        comment: z.coerce.string().max(255).optional(),
                    }),
                    z.object({
                        ruleKind: z.enum(ApprovalAssigneeKind).default(ApprovalAssigneeKind.JOB_GRADE),
                        userId: z.coerce.number().optional().default(0),
                        job: z.coerce.string(),
                        minimumGrade: z.coerce.number().min(game.startJobGrade),
                        label: z.string().max(120).optional(),
                        signatureRequired: z.coerce.boolean().default(false),
                        slots: z.coerce.number().min(1).max(10).optional().default(1),
                        dueInDays: z.coerce.number().min(1).optional(),
                        comment: z.coerce.string().max(255).optional(),
                    }),
                ])
                .array()
                .max(20)
                .default([]),
        })
        .default({
            enabled: false,
            policy: {
                ruleKind: ApprovalRuleKind.REQUIRE_ALL,
                onEditBehavior: OnEditBehavior.KEEP_PROGRESS,
                requiredCount: 2,
                signatureRequired: false,
            },
            tasks: [],
        }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    weight: 0,
    title: '',
    description: '',
    color: 'primary',
    contentTitle: '',
    content: '',
    contentState: '',
    category: undefined,
    jobAccess: [],
    contentAccess: {
        jobs: [],
        users: [],
    },
    workflow: {
        autoClose: {
            autoClose: false,
            autoCloseSettings: {
                message: '',
                duration: 7,
            },
        },

        reminders: {
            reminder: false,
            reminderSettings: {
                reminders: [],
                maxReminderCount: 10,
            },
        },
    },
    approval: {
        enabled: false,
        policy: {
            ruleKind: ApprovalRuleKind.REQUIRE_ALL,
            onEditBehavior: OnEditBehavior.KEEP_PROGRESS,
            requiredCount: 2,
            signatureRequired: false,
        },

        tasks: [],
    },
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (event.submitter?.getAttribute('role') === 'tab') return;

    canSubmit.value = false;
    await createOrUpdateTemplate(event.data, props.templateId).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const schemaEditor = ref<TemplateRequirements>({
    users: {
        required: false,
        min: 0,
        max: 0,
    },

    documents: {
        required: false,
        min: 0,
        max: 0,
    },

    vehicles: {
        required: false,
        min: 0,
        max: 0,
    },
});

const accessTypes: AccessType[] = [{ label: t('common.job', 2), value: 'job' }];
const contentAccessTypes: AccessType[] = [
    { label: t('common.citizen', 2), value: 'user' },
    { label: t('common.job', 2), value: 'job' },
];

async function createOrUpdateTemplate(values: Schema, templateId?: number): Promise<void> {
    values.contentAccess.users.forEach((user) => user.id < 0 && (user.id = 0));
    values.contentAccess.jobs.forEach((job) => job.id < 0 && (job.id = 0));

    values.jobAccess.forEach((job) => job.id < 0 && (job.id = 0));

    const tRequirements: TemplateRequirements = {
        users: schemaEditor.value.users,
        documents: schemaEditor.value.documents,
        vehicles: schemaEditor.value.vehicles,
    };

    const req: CreateTemplateRequest | UpdateTemplateRequest = {
        template: {
            id: templateId ?? 0,
            weight: values.weight,
            title: values.title,
            description: values.description,
            color: values.color,
            icon: values.icon,
            contentTitle: values.contentTitle,
            content: values.content,
            state: values.contentState,
            schema: {
                requirements: tRequirements,
            },
            contentAccess: values.contentAccess,
            jobAccess: values.jobAccess,
            category: values.category,
            creatorJob: '',
            workflow: {
                reminder: values.workflow.reminders.reminder,
                reminderSettings: {
                    reminders: values.workflow.reminders.reminderSettings.reminders
                        .filter((r) => r.duration !== undefined)
                        .map((r) => ({
                            duration: toDuration((r.duration ? r.duration : 0) * 24 * 60 * 60),
                            message: r.message ?? '',
                        })),
                    maxReminderCount: values.workflow?.reminders.reminderSettings?.maxReminderCount ?? 10,
                },

                autoClose: values.workflow.autoClose.autoClose,
                autoCloseSettings: {
                    duration: toDuration(
                        (values.workflow.autoClose.autoCloseSettings.duration > 0
                            ? values.workflow.autoClose.autoCloseSettings.duration
                            : 1) *
                            24 *
                            60 *
                            60,
                    ),
                    message: values.workflow.autoClose.autoCloseSettings.message ?? '',
                },
            },
            approval: {
                enabled: values.approval.enabled,
                policy: {
                    ruleKind: values.approval.policy.ruleKind,
                    onEditBehavior: values.approval.policy.onEditBehavior,
                    requiredCount: values.approval.policy.requiredCount,
                    signatureRequired: values.approval.policy.signatureRequired,
                },
                tasks: values.approval.tasks.map((task) => ({
                    ruleKind: task.ruleKind,
                    userId: task.userId,
                    job: task.job ?? '',
                    minimumGrade: task.minimumGrade ?? 0,
                    label: task.label,
                    signatureRequired: task.signatureRequired,
                    slots: task.slots,
                    dueInDays: task.dueInDays,
                    comment: task.comment,
                })),
            },
        },
    };

    try {
        if (templateId === undefined) {
            const call = documentsDocumentsClient.createTemplate(req);
            const { response } = await call;

            notifications.add({
                title: { key: 'notifications.templates.created.title', parameters: {} },
                description: { key: 'notifications.templates.created.title', parameters: {} },
                type: NotificationType.SUCCESS,
            });

            await navigateTo({
                name: 'documents-templates-id',
                params: { id: response.id },
            });
        } else {
            const call = documentsDocumentsClient.updateTemplate(req);
            const { response } = await call;
            if (response.template) {
                setValuesFromTemplate(response.template);
            }

            notifications.add({
                title: { key: 'notifications.templates.updated.title', parameters: {} },
                description: { key: 'notifications.templates.updated.content', parameters: {} },
                type: NotificationType.SUCCESS,
            });
        }
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const entriesCategories = ref<Category[]>([]);
const queryCategories = ref('');

watchDebounced(queryCategories, () => findCategories(), {
    debounce: 200,
    maxWait: 1250,
});

async function findCategories(): Promise<void> {
    entriesCategories.value = await completorStore.completeDocumentCategories(queryCategories.value);
}

function setValuesFromTemplate(tpl: Template): void {
    state.weight = tpl.weight;
    state.title = tpl.title;
    state.description = tpl.description;
    state.color = tpl.color ?? 'primary';
    state.icon = tpl.icon;
    state.contentTitle = tpl.contentTitle;
    state.content = tpl.content;
    state.contentState = tpl.state;
    state.category = tpl.category;
    state.jobAccess = tpl.jobAccess;
    state.contentAccess = tpl.contentAccess ?? {
        jobs: [],
        users: [],
    };

    const autoCloseDuration = fromDuration(tpl.workflow?.autoCloseSettings?.duration);
    state.workflow = {
        reminders: {
            reminder: tpl.workflow?.reminder ?? false,
            reminderSettings: {
                reminders:
                    tpl.workflow?.reminderSettings?.reminders.map((r) => {
                        const dur = fromDuration(r.duration);
                        return {
                            duration: dur > 0 ? dur / 24 / 60 / 60 : 7,
                            message: r.message ?? '',
                        };
                    }) ?? [],
                maxReminderCount:
                    (tpl.workflow?.reminderSettings?.maxReminderCount ?? 10) <= 0
                        ? 10
                        : (tpl.workflow?.reminderSettings?.maxReminderCount ?? 10),
            },
        },

        autoClose: {
            autoClose: tpl.workflow?.autoClose ?? false,
            autoCloseSettings: {
                message: tpl.workflow?.autoCloseSettings?.message ?? '',
                duration: autoCloseDuration > 0 ? autoCloseDuration / 24 / 60 / 60 : 7,
            },
        },
    };

    state.approval = {
        enabled: tpl.approval?.enabled ?? false,
        policy: {
            ruleKind: tpl.approval?.policy?.ruleKind ?? ApprovalRuleKind.REQUIRE_ALL,
            onEditBehavior: tpl.approval?.policy?.onEditBehavior ?? OnEditBehavior.KEEP_PROGRESS,
            requiredCount: tpl.approval?.policy?.requiredCount ?? 2,
            signatureRequired: tpl.approval?.policy?.signatureRequired ?? false,
        },
        tasks:
            tpl.approval?.tasks.map((task) => ({
                ruleKind: task.userId == 0 ? ApprovalAssigneeKind.JOB_GRADE : ApprovalAssigneeKind.USER,
                userId: task.userId,
                job: task.job,
                minimumGrade: task.minimumGrade,
                label: task.label,
                signatureRequired: task.signatureRequired,
                slots: task.slots,
                dueInDays: task.dueInDays,
                comment: task.comment,
            })) ?? [],
    };

    schemaEditor.value.users = tpl.schema?.requirements?.users;

    schemaEditor.value.documents = tpl.schema?.requirements?.documents;

    schemaEditor.value.vehicles = tpl.schema?.requirements?.vehicles;
}

const extensions = [TemplateVar.configure(), TemplateBlock.configure()];

onBeforeMount(async () => {
    if (props.templateId) {
        try {
            const call = documentsDocumentsClient.getTemplate({
                templateId: props.templateId,
                render: false,
            });
            const { response } = await call;

            const tpl = response.template;
            if (!tpl) return;

            setValuesFromTemplate(tpl);

            useHead({
                title: () =>
                    tpl?.title
                        ? `${tpl.title} - ${t('pages.documents.templates.edit.title')}`
                        : t('pages.documents.templates.edit.title'),
            });
        } catch (e) {
            handleGRPCError(e as RpcError);
        }
    } else {
        state.weight = 0;

        state.jobAccess.push({
            id: 0,
            targetId: props.templateId ?? 0,
            job: activeChar.value!.job,
            minimumGrade: game.startJobGrade,
            access: AccessLevel.VIEW,
        });
    }

    findCategories();
});

const items = [
    {
        slot: 'details' as const,
        label: t('common.detail', 2),
        icon: 'i-mdi-details',
        value: 'details',
    },
    {
        slot: 'workflow' as const,
        label: t('common.workflow'),
        icon: 'i-mdi-workflow',
        value: 'workflow',
    },
    {
        slot: 'approval' as const,
        label: t('common.approvals', 2),
        icon: 'i-mdi-approval',
        value: 'approval',
    },
    {
        slot: 'content' as const,
        label: t('common.content'),
        icon: 'i-mdi-file-edit',
        value: 'content',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        return (route.query.tab as string) || 'details';
    },
    set(tab) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.push({ query: { tab: tab }, hash: '#control-active-item' });
    },
});

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0', body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.documents.templates.edit.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <UButton
                        color="neutral"
                        icon="i-mdi-arrow-left"
                        :to="
                            templateId ? { name: 'documents-templates-id', params: { id: templateId } } : `/documents/templates`
                        "
                        :label="$t('common.back')"
                    />

                    <UButton
                        trailing-icon="i-mdi-content-save"
                        :disabled="!canSubmit"
                        :loading="!canSubmit"
                        @click="formRef?.submit()"
                    >
                        <span class="hidden truncate sm:block">
                            {{ templateId ? $t('common.save') : $t('common.create') }}
                        </span>
                    </UButton>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <UForm
                ref="formRef"
                class="mb-4 flex w-full max-w-full flex-1 flex-col overflow-y-hidden"
                :schema="schema"
                :state="state"
                @submit="onSubmitThrottle"
            >
                <UTabs
                    v-model="selectedTab"
                    class="flex-1 flex-col overflow-y-hidden"
                    :items="items"
                    variant="link"
                    :unmount-on-hide="false"
                    :ui="{
                        content:
                            'p-4 flex flex-1 min-h-0 max-h-full overflow-y-auto flex-col gap-4 max-w-(--ui-container) mx-auto',
                    }"
                >
                    <template #details>
                        <UPageCard :title="$t('common.detail', 2)">
                            <UFormField
                                name="weight"
                                :label="`${$t('common.template', 1)} ${$t('common.weight')}`"
                                class="grid grid-cols-2 items-center gap-2"
                            >
                                <UInputNumber
                                    v-model="state.weight"
                                    name="weight"
                                    :min="0"
                                    :max="999999"
                                    :step="1"
                                    :placeholder="$t('common.weight')"
                                />
                            </UFormField>

                            <UFormField
                                name="title"
                                :label="`${$t('common.template')} ${$t('common.title')}`"
                                class="grid grid-cols-2 items-center gap-2"
                                required
                            >
                                <UInput
                                    v-model="state.title"
                                    name="title"
                                    :placeholder="$t('common.title')"
                                    size="lg"
                                    class="w-full"
                                />
                            </UFormField>

                            <UFormField
                                name="description"
                                :label="`${$t('common.template')} ${$t('common.description')}`"
                                class="grid grid-cols-2 items-center gap-2"
                                required
                            >
                                <UTextarea
                                    v-model="state.description"
                                    name="description"
                                    :rows="4"
                                    :label="$t('common.description')"
                                    class="w-full"
                                />
                            </UFormField>

                            <UFormField
                                name="color"
                                :label="$t('common.color')"
                                class="grid grid-cols-2 items-center gap-2"
                                required
                            >
                                <div class="flex flex-1 gap-1">
                                    <ColorPickerTW v-model="state.color" class="flex-1" />
                                </div>
                            </UFormField>

                            <UFormField name="icon" :label="$t('common.icon')" class="grid grid-cols-2 items-center gap-2">
                                <div class="flex flex-1 gap-1">
                                    <IconSelectMenu
                                        v-model="state.icon"
                                        class="flex-1"
                                        :color="state.color"
                                        :fallback-icon="FileOutlineIcon"
                                    />

                                    <UButton icon="i-mdi-backspace" @click="state.icon = undefined" />
                                </div>
                            </UFormField>
                        </UPageCard>

                        <UPageCard :title="`${$t('common.template')} ${$t('common.access')}`">
                            <AccessManager
                                v-model:jobs="state.jobAccess"
                                :target-id="templateId ?? 0"
                                :access-types="accessTypes"
                                :access-roles="
                                    enumToAccessLevelEnums(AccessLevel, 'enums.documents.AccessLevel').filter(
                                        (e) => e.value === AccessLevel.VIEW || e.value === AccessLevel.EDIT,
                                    )
                                "
                                name="jobAccess"
                                full-name
                            />
                        </UPageCard>

                        <UPageCard :title="$t('common.requirements', 2)">
                            <TemplateSchemaEditor v-model="schemaEditor" />
                        </UPageCard>
                    </template>

                    <template #workflow>
                        <TemplateWorkflowEditor v-model="state.workflow" />
                    </template>

                    <template #approval>
                        <UPageCard :title="$t('components.documents.approval.policy_form.title', 2)">
                            <UFormField name="approval.enabled" :label="$t('common.enabled')">
                                <USwitch v-model="state.approval.enabled" />
                            </UFormField>

                            <PolicyEditor v-model="state.approval.policy" :disabled="!state.approval.enabled" />
                        </UPageCard>

                        <UPageCard :title="$t('components.documents.approval.tasks', 2)">
                            <ApprovalTasksEditor
                                v-model="state.approval.tasks"
                                :disabled="!state.approval.enabled"
                                :signature-required="state.approval.policy.signatureRequired"
                            />
                        </UPageCard>
                    </template>

                    <template #content>
                        <UPageCard>
                            <UFormField name="contentTitle" :label="`${$t('common.content')} ${$t('common.title')}`" required>
                                <UTextarea v-model="state.contentTitle" name="contentTitle" :rows="2" class="w-full" />
                            </UFormField>

                            <UFormField name="category" :label="$t('common.category', 1)">
                                <SelectMenu
                                    v-model="state.category"
                                    :filter-fields="['name']"
                                    block
                                    nullable
                                    class="w-full"
                                    :searchable="
                                        async (q: string) => {
                                            try {
                                                const categories = await completorStore.completeDocumentCategories(q);
                                                return categories;
                                            } catch (e) {
                                                handleGRPCError(e as RpcError);
                                                throw e;
                                            }
                                        }
                                    "
                                    searchable-key="completor-document-categories"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                >
                                    <template v-if="state.category" #default>
                                        <CategoryBadge :category="state.category" />
                                    </template>

                                    <template #item-label="{ item }">
                                        <CategoryBadge :category="item" />
                                    </template>

                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.category', 2)]) }}
                                    </template>
                                </SelectMenu>
                            </UFormField>

                            <UFormField name="contentState" :label="`${$t('common.content')} ${$t('common.state')}`">
                                <UTextarea v-model="state.contentState" name="contentState" :rows="2" class="w-full" />
                            </UFormField>

                            <SingleHint
                                hint-id="template_editor_templating"
                                to="https://fivenet.app/user-guides/documents/templates"
                                external
                                link-target="_blank"
                            />

                            <UFormField
                                class="flex flex-1 flex-col overflow-y-hidden"
                                name="content"
                                :label="`${$t('common.content')} ${$t('common.template')}`"
                                required
                                :ui="{ container: 'flex flex-1 overflow-y-hidden flex-col', error: 'hidden' }"
                            >
                                <ClientOnly>
                                    <TiptapEditor
                                        v-model="state.content"
                                        name="content"
                                        class="mx-auto min-h-120 w-full max-w-(--breakpoint-xl) flex-1 overflow-y-hidden"
                                        :extensions="extensions"
                                    >
                                        <template #toolbar="{ editor }">
                                            <EditorButtons :editor="editor" />
                                        </template>
                                    </TiptapEditor>
                                </ClientOnly>
                            </UFormField>
                        </UPageCard>

                        <UPageCard :title="$t('common.access')">
                            <AccessManager
                                v-model:jobs="state.contentAccess.jobs"
                                v-model:users="state.contentAccess.users"
                                :target-id="templateId ?? 0"
                                :access-types="contentAccessTypes"
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.documents.AccessLevel')"
                                show-required
                                name="contentAccess"
                            />
                        </UPageCard>
                    </template>
                </UTabs>
            </UForm>
        </template>
    </UDashboardPanel>
</template>
