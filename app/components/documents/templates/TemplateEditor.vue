<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { FileOutlineIcon } from 'mdi-vue3';
import { z } from 'zod';
import SingleHint from '~/components/SingleHint.vue';
import TemplateSchemaEditor, { type SchemaEditorValue } from '~/components/documents/templates/TemplateSchemaEditor.vue';
import { zWorkflow, type ObjectSpecsValue } from '~/components/documents/templates/types';
import ColorPickerTW from '~/components/partials/ColorPickerTW.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums, type AccessType } from '~/components/partials/access/helpers';
import TiptapEditor from '~/components/partials/editor/TiptapEditor.vue';
import { TemplateBlock } from '~/composables/tiptap/extensions/TemplateBlock';
import { TemplateVar } from '~/composables/tiptap/extensions/TemplateVar';
import { useAuthStore } from '~/stores/auth';
import { useCompletorStore } from '~/stores/completor';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { AccessLevel, type DocumentJobAccess, type DocumentUserAccess } from '~~/gen/ts/resources/documents/access';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { ObjectSpecs, Template, TemplateJobAccess, TemplateRequirements } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { CreateTemplateRequest, UpdateTemplateRequest } from '~~/gen/ts/services/documents/documents';
import TemplateEditorButtons from './TemplateEditorButtons.vue';
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
    title: z.string().min(3).max(255),
    description: z.string().min(3).max(255),
    color: z.string().max(7),
    icon: z.string().max(128).optional(),
    contentTitle: z.string().min(3).max(2048),
    content: z.string().min(3).max(1500000),
    contentState: z.union([z.string().min(1).max(512), z.string().length(0)]),
    category: z.custom<Category>().optional(),
    jobAccess: z.custom<TemplateJobAccess>().array().max(maxAccessEntries).default([]),
    contentAccess: z.object({
        jobs: z.custom<DocumentJobAccess>().array().max(maxAccessEntries).default([]),
        users: z.custom<DocumentUserAccess>().array().max(maxAccessEntries).default([]),
    }),
    workflow: zWorkflow,
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
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    if (event.submitter?.getAttribute('role') === 'tab') {
        return;
    }

    canSubmit.value = false;
    await createOrUpdateTemplate(event.data, props.templateId).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const schemaEditor = ref<SchemaEditorValue>({
    users: {
        req: false,
        min: 0,
        max: 0,
    },

    documents: {
        req: false,
        min: 0,
        max: 0,
    },

    vehicles: {
        req: false,
        min: 0,
        max: 0,
    },
});

const accessTypes: AccessType[] = [{ type: 'job', name: t('common.job', 2) }];
const contentAccessTypes: AccessType[] = [
    { type: 'user', name: t('common.citizen', 2) },
    { type: 'job', name: t('common.job', 2) },
];

function createObjectSpec(v: ObjectSpecsValue): ObjectSpecs {
    const o: ObjectSpecs = {
        required: v.req,
        min: v.min,
        max: v.max,
    };
    return o;
}

async function createOrUpdateTemplate(values: Schema, templateId?: number): Promise<void> {
    values.contentAccess.users.forEach((user) => user.id < 0 && (user.id = 0));
    values.contentAccess.jobs.forEach((job) => job.id < 0 && (job.id = 0));

    values.jobAccess.forEach((job) => job.id < 0 && (job.id = 0));

    const tRequirements: TemplateRequirements = {
        users: createObjectSpec(schemaEditor.value.users),
        documents: createObjectSpec(schemaEditor.value.documents),
        vehicles: createObjectSpec(schemaEditor.value.vehicles),
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
                maxReminderCount: tpl.workflow?.reminderSettings?.maxReminderCount ?? 10,
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

    schemaEditor.value.users.req = tpl.schema?.requirements?.users?.required ?? false;
    schemaEditor.value.users.min = tpl.schema?.requirements?.users?.min ?? 0;
    schemaEditor.value.users.max = tpl.schema?.requirements?.users?.max ?? 0;

    schemaEditor.value.documents.req = tpl.schema?.requirements?.documents?.required ?? false;
    schemaEditor.value.documents.min = tpl.schema?.requirements?.documents?.min ?? 0;
    schemaEditor.value.documents.max = tpl.schema?.requirements?.documents?.max ?? 0;

    schemaEditor.value.vehicles.req = tpl.schema?.requirements?.vehicles?.required ?? false;
    schemaEditor.value.vehicles.min = tpl.schema?.requirements?.vehicles?.min ?? 0;
    schemaEditor.value.vehicles.max = tpl.schema?.requirements?.vehicles?.max ?? 0;
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
            if (!tpl) {
                return;
            }

            setValuesFromTemplate(tpl);
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

const categoriesLoading = ref(false);
</script>

<template>
    <UForm
        class="flex min-h-dvh w-full max-w-full flex-1 flex-col overflow-y-auto"
        :schema="schema"
        :state="state"
        @submit="onSubmitThrottle"
    >
        <UDashboardNavbar :title="$t('pages.documents.templates.edit.title')">
            <template #right>
                <UButton
                    color="neutral"
                    icon="i-mdi-arrow-left"
                    :to="templateId ? { name: 'documents-templates-id', params: { id: templateId } } : `/documents/templates`"
                >
                    {{ $t('common.back') }}
                </UButton>

                <UButton type="submit" trailing-icon="i-mdi-content-save" :disabled="!canSubmit" :loading="!canSubmit">
                    <span class="hidden truncate sm:block">
                        {{ templateId ? $t('common.save') : $t('common.create') }}
                    </span>
                </UButton>
            </template>
        </UDashboardNavbar>

        <UDashboardPanelContent class="p-0 sm:pb-0">
            <UTabs v-model="selectedTab" class="flex flex-1 flex-col" :items="items">
                <template #details>
                    <UContainer class="mt-2 w-full overflow-y-scroll">
                        <div>
                            <UFormField name="weight" :label="`${$t('common.template', 1)} ${$t('common.weight')}`">
                                <UInput
                                    v-model="state.weight"
                                    type="number"
                                    name="weight"
                                    :min="0"
                                    :max="999999"
                                    :placeholder="$t('common.weight')"
                                />
                            </UFormField>

                            <UFormField name="title" :label="`${$t('common.template')} ${$t('common.title')}`" required>
                                <UTextarea v-model="state.title" name="title" :rows="1" :placeholder="$t('common.title')" />
                            </UFormField>

                            <UFormField
                                name="description"
                                :label="`${$t('common.template')} ${$t('common.description')}`"
                                required
                            >
                                <UTextarea
                                    v-model="state.description"
                                    name="description"
                                    :rows="4"
                                    :label="$t('common.description')"
                                />
                            </UFormField>

                            <UFormField class="flex-1 flex-row" name="color" :label="$t('common.color')" required>
                                <div class="flex flex-1 gap-1">
                                    <ColorPickerTW v-model="state.color" class="flex-1" />
                                </div>
                            </UFormField>

                            <UFormField class="flex-1" name="icon" :label="$t('common.icon')">
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
                        </div>

                        <div class="my-2">
                            <h2 class="text-sm">{{ $t('common.template') }} {{ $t('common.access') }}</h2>

                            <AccessManager
                                v-model:jobs="state.jobAccess"
                                :target-id="templateId ?? 0"
                                :access-types="accessTypes"
                                :access-roles="
                                    enumToAccessLevelEnums(AccessLevel, 'enums.documents.AccessLevel').filter(
                                        (e) => e.value === AccessLevel.VIEW || e.value === AccessLevel.EDIT,
                                    )
                                "
                            />
                        </div>

                        <div class="my-2">
                            <UAccordion
                                :items="[
                                    { slot: 'schema' as const, label: $t('common.requirements', 2), icon: 'i-mdi-asterisk' },
                                    {
                                        slot: 'workflow' as const,
                                        label: $t('common.workflow'),
                                        icon: 'i-mdi-reminder',
                                    },
                                ]"
                            >
                                <template #schema>
                                    <TemplateSchemaEditor v-model="schemaEditor" />
                                </template>

                                <template #workflow>
                                    <TemplateWorkflowEditor v-model="state.workflow" />
                                </template>
                            </UAccordion>
                        </div>

                        <div class="my-2">
                            <h2 class="text-sm">{{ $t('common.content') }} {{ $t('common.access') }}</h2>

                            <AccessManager
                                v-model:jobs="state.contentAccess.jobs"
                                v-model:users="state.contentAccess.users"
                                :target-id="templateId ?? 0"
                                :access-types="contentAccessTypes"
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.documents.AccessLevel')"
                                :show-required="true"
                            />
                        </div>
                    </UContainer>
                </template>

                <template #content>
                    <UContainer class="flex w-full flex-1 flex-col overflow-y-hidden">
                        <SingleHint
                            class="my-2"
                            hint-id="template_editor_templating"
                            to="https://fivenet.app/user-guides/documents/templates"
                            external
                            link-target="_blank"
                        />

                        <UFormField name="contentTitle" :label="`${$t('common.content')} ${$t('common.title')}`" required>
                            <UTextarea v-model="state.contentTitle" name="contentTitle" :rows="2" />
                        </UFormField>

                        <UFormField name="category" :label="$t('common.category', 1)">
                            <ClientOnly>
                                <UInputMenu
                                    v-model="state.category"
                                    option-attribute="name"
                                    :search-attributes="['name']"
                                    nullable
                                    :search="
                                        async (search: string) => {
                                            try {
                                                categoriesLoading = true;
                                                const categories = await completorStore.completeDocumentCategories(search);
                                                categoriesLoading = false;
                                                return categories;
                                            } catch (e) {
                                                handleGRPCError(e as RpcError);
                                                throw e;
                                            } finally {
                                                categoriesLoading = false;
                                            }
                                        }
                                    "
                                >
                                    <template #empty> {{ $t('common.not_found', [$t('common.category', 2)]) }} </template>
                                </UInputMenu>
                            </ClientOnly>
                        </UFormField>

                        <UFormField name="contentState" :label="`${$t('common.content')} ${$t('common.state')}`">
                            <UTextarea v-model="state.contentState" name="contentState" :rows="2" />
                        </UFormField>

                        <UFormField
                            class="flex flex-1 flex-col overflow-y-hidden"
                            name="content"
                            :label="`${$t('common.content')} ${$t('common.template')}`"
                            required
                            :ui="{ container: 'flex flex-1 overflow-y-hidden flex-col' }"
                        >
                            <ClientOnly>
                                <TiptapEditor
                                    v-model="state.content"
                                    class="mx-auto w-full max-w-(--breakpoint-xl) flex-1 overflow-y-hidden"
                                    :extensions="extensions"
                                >
                                    <template #toolbar="{ editor }">
                                        <TemplateEditorButtons :editor="editor" />
                                    </template>
                                </TiptapEditor>
                            </ClientOnly>
                        </UFormField>
                    </UContainer>
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
