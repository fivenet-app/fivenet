<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
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
import { useAuthStore } from '~/stores/auth';
import { useCompletorStore } from '~/stores/completor';
import { useNotificatorStore } from '~/stores/notificator';
import type { DocumentJobAccess, DocumentUserAccess } from '~~/gen/ts/resources/documents/access';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { ObjectSpecs, Template, TemplateJobAccess, TemplateRequirements } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { CreateTemplateRequest, UpdateTemplateRequest } from '~~/gen/ts/services/docstore/docstore';
import TemplateWorkflowEditor from './TemplateWorkflowEditor.vue';

const props = defineProps<{
    templateId?: number;
}>();

const { $grpc } = useNuxtApp();

const { t } = useI18n();

const { game } = useAppConfig();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const notifications = useNotificatorStore();

const completorStore = useCompletorStore();

const { maxAccessEntries } = useAppConfig();

const schema = z.object({
    weight: z.coerce.number().min(0).max(999_999),
    title: z.string().min(3).max(255),
    description: z.string().min(3).max(255),
    color: z.string().max(7),
    icon: z.string().max(64).optional(),
    contentTitle: z.string().min(3).max(2048),
    content: z.string().min(3).max(1500000),
    contentState: z.union([z.string().min(1).max(2048), z.string().length(0)]),
    category: z.custom<Category>().optional(),
    jobAccess: z.custom<TemplateJobAccess>().array().max(maxAccessEntries),
    contentAccess: z.object({
        jobs: z.custom<DocumentJobAccess>().array().max(maxAccessEntries),
        users: z.custom<DocumentUserAccess>().array().max(maxAccessEntries),
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
            const call = $grpc.docstore.docStore.createTemplate(req);
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
            const call = $grpc.docstore.docStore.updateTemplate(req);
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

onBeforeMount(async () => {
    if (props.templateId) {
        try {
            const call = $grpc.docstore.docStore.getTemplate({
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
        slot: 'details',
        label: t('common.detail', 2),
        icon: 'i-mdi-details',
    },
    {
        slot: 'content',
        label: t('common.content'),
        icon: 'i-mdi-file-edit',
    },
];

const route = useRoute();
const router = useRouter();

const selectedTab = computed({
    get() {
        const index = items.findIndex((item) => item.slot === route.query.tab);
        if (index === -1) {
            return 0;
        }

        return index;
    },
    set(value) {
        // Hash is specified here to prevent the page from scrolling to the top
        router.replace({ query: { tab: items[value]?.slot }, hash: '#' });
    },
});

const categoriesLoading = ref(false);
</script>

<template>
    <UForm
        :schema="schema"
        :state="state"
        class="min-h-dscreen flex w-full max-w-full flex-1 flex-col overflow-y-auto"
        @submit="onSubmitThrottle"
    >
        <UDashboardNavbar :title="$t('pages.documents.templates.edit.title')">
            <template #right>
                <UButton
                    color="black"
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

        <UDashboardPanelContent class="p-0">
            <UTabs
                v-model="selectedTab"
                :items="items"
                class="flex flex-1 flex-col"
                :ui="{
                    wrapper: 'space-y-0 overflow-y-hidden',
                    container: 'flex flex-1 flex-col overflow-y-hidden',
                    base: 'flex flex-1 flex-col overflow-y-hidden',
                    list: { rounded: '' },
                }"
            >
                <template #details>
                    <UContainer class="mt-2 w-full overflow-y-scroll">
                        <div>
                            <UFormGroup name="weight" :label="`${$t('common.template', 1)} ${$t('common.weight')}`">
                                <UInput
                                    v-model="state.weight"
                                    type="number"
                                    name="weight"
                                    :min="0"
                                    :max="999999"
                                    :placeholder="$t('common.weight')"
                                />
                            </UFormGroup>

                            <UFormGroup name="title" :label="`${$t('common.template')} ${$t('common.title')}`" required>
                                <UTextarea v-model="state.title" name="title" :rows="1" :placeholder="$t('common.title')" />
                            </UFormGroup>

                            <UFormGroup
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
                            </UFormGroup>

                            <UFormGroup name="color" :label="$t('common.color')" class="flex-1 flex-row" required>
                                <div class="flex flex-1 gap-1">
                                    <ColorPickerTW v-model="state.color" class="flex-1" />
                                </div>
                            </UFormGroup>

                            <UFormGroup name="icon" :label="$t('common.icon')" class="flex-1">
                                <div class="flex flex-1 gap-1">
                                    <IconSelectMenu v-model="state.icon" class="flex-1" :fallback-icon="FileOutlineIcon" />

                                    <UButton icon="i-mdi-backspace" @click="state.icon = undefined" />
                                </div>
                            </UFormGroup>
                        </div>

                        <div class="my-2">
                            <h2 class="text-sm">{{ $t('common.template') }} {{ $t('common.access') }}</h2>

                            <AccessManager
                                v-model:jobs="state.jobAccess"
                                :target-id="templateId ?? 0"
                                :access-types="accessTypes"
                                :access-roles="
                                    enumToAccessLevelEnums(AccessLevel, 'enums.docstore.AccessLevel').filter(
                                        (e) => e.value === AccessLevel.VIEW || e.value === AccessLevel.EDIT,
                                    )
                                "
                            />
                        </div>

                        <div class="my-2">
                            <UAccordion
                                :items="[
                                    { slot: 'schema', label: $t('common.requirements', 2), icon: 'i-mdi-asterisk' },
                                    {
                                        slot: 'workflow',
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
                                :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.docstore.AccessLevel')"
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
                            :external="true"
                            link-target="_blank"
                        />

                        <UFormGroup name="contentTitle" :label="`${$t('common.content')} ${$t('common.title')}`" required>
                            <UTextarea v-model="state.contentTitle" name="contentTitle" :rows="2" />
                        </UFormGroup>

                        <UFormGroup name="category" :label="$t('common.category', 1)">
                            <ClientOnly>
                                <UInputMenu
                                    v-model="state.category"
                                    option-attribute="name"
                                    :search-attributes="['name']"
                                    block
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
                                    search-lazy
                                    :search-placeholder="$t('common.search_field')"
                                >
                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>

                                    <template #empty> {{ $t('common.not_found', [$t('common.category', 2)]) }} </template>
                                </UInputMenu>
                            </ClientOnly>
                        </UFormGroup>

                        <UFormGroup name="contentState" :label="`${$t('common.content')} ${$t('common.state')}`">
                            <UTextarea v-model="state.contentState" name="contentState" :rows="2" />
                        </UFormGroup>

                        <UFormGroup
                            name="content"
                            :label="`${$t('common.content')} ${$t('common.template')}`"
                            required
                            class="flex flex-1 flex-col overflow-y-hidden"
                            :ui="{ container: 'flex flex-1 overflow-y-hidden' }"
                        >
                            <ClientOnly>
                                <TiptapEditor
                                    v-model="state.content"
                                    class="mx-auto w-full max-w-screen-xl flex-1 overflow-y-hidden"
                                />
                            </ClientOnly>
                        </UFormGroup>
                    </UContainer>
                </template>
            </UTabs>
        </UDashboardPanelContent>
    </UForm>
</template>
