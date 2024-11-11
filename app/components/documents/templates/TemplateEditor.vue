<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import SingleHint from '~/components/SingleHint.vue';
import TemplateSchemaEditor, { type SchemaEditorValue } from '~/components/documents/templates/TemplateSchemaEditor.vue';
import type { ObjectSpecsValue } from '~/components/documents/templates/types';
import DocEditor from '~/components/partials/DocEditor.vue';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums, type AccessType } from '~/components/partials/access/helpers';
import { useAuthStore } from '~/store/auth';
import { useCompletorStore } from '~/store/completor';
import { useNotificatorStore } from '~/store/notificator';
import type { DocumentJobAccess, DocumentUserAccess } from '~~/gen/ts/resources/documents/access';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { ObjectSpecs, Template, TemplateJobAccess, TemplateRequirements } from '~~/gen/ts/resources/documents/templates';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { CreateTemplateRequest, UpdateTemplateRequest } from '~~/gen/ts/services/docstore/docstore';

const props = defineProps<{
    templateId?: string;
}>();

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
    contentTitle: z.string().min(3).max(2048),
    content: z.string().min(3).max(1500000),
    contentState: z.union([z.string().min(1).max(2048), z.string().length(0)]),
    category: z.custom<Category>().optional(),
    jobAccess: z.custom<TemplateJobAccess>().array().max(maxAccessEntries),
    contentAccess: z.object({
        jobs: z.custom<DocumentJobAccess>().array().max(maxAccessEntries),
        users: z.custom<DocumentUserAccess>().array().max(maxAccessEntries),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    weight: 0,
    title: '',
    description: '',
    contentTitle: '',
    content: '',
    contentState: '',
    category: undefined,
    jobAccess: [],
    contentAccess: {
        jobs: [],
        users: [],
    },
});

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
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

async function createOrUpdateTemplate(values: Schema, templateId?: string): Promise<void> {
    const tRequirements: TemplateRequirements = {
        users: createObjectSpec(schemaEditor.value.users),
        documents: createObjectSpec(schemaEditor.value.documents),
        vehicles: createObjectSpec(schemaEditor.value.vehicles),
    };

    const req: CreateTemplateRequest | UpdateTemplateRequest = {
        template: {
            id: templateId ?? '0',
            weight: values.weight,
            title: values.title,
            description: values.description,
            contentTitle: values.contentTitle,
            content: values.content,
            state: values.contentState,
            schema: {
                requirements: tRequirements,
            },
            contentAccess: values.contentAccess,
            jobAccess: values.jobAccess,
            category: state.category,
            creatorJob: '',
        },
    };

    try {
        if (templateId === undefined) {
            const call = getGRPCDocStoreClient().createTemplate(req);
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
            const call = getGRPCDocStoreClient().updateTemplate(req);
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
    state.contentTitle = tpl.contentTitle;
    state.content = tpl.content;
    state.contentState = tpl.state;
    state.category = tpl.category;
    state.contentAccess = tpl.contentAccess ?? {
        jobs: [],
        users: [],
    };
    state.jobAccess = tpl.jobAccess;

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

onMounted(async () => {
    if (props.templateId) {
        try {
            const call = getGRPCDocStoreClient().getTemplate({
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
            id: '0',
            targetId: props.templateId ?? '0',
            job: activeChar.value!.job,
            minimumGrade: game.startJobGrade,
            access: AccessLevel.VIEW,
        });
    }

    findCategories();
});

const categoriesLoading = ref(false);
</script>

<template>
    <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UDashboardNavbar :title="$t('pages.documents.templates.edit.title')">
            <template #right>
                <UButton
                    color="black"
                    icon="i-mdi-arrow-left"
                    :to="templateId ? { name: 'documents-templates-id', params: { id: templateId } } : `/documents/templates`"
                >
                    {{ $t('common.back') }}
                </UButton>

                <UButtonGroup class="inline-flex">
                    <UButton type="submit" trailing-icon="i-mdi-content-save" :disabled="!canSubmit" :loading="!canSubmit">
                        {{ templateId ? $t('common.save') : $t('common.create') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UDashboardNavbar>

        <UContainer class="w-full">
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

                <UFormGroup name="description" :label="`${$t('common.template')} ${$t('common.description')}`" required>
                    <UTextarea v-model="state.description" name="description" :rows="4" :label="$t('common.description')" />
                </UFormGroup>
            </div>

            <div class="my-2">
                <h2 class="text-sm">{{ $t('common.template') }} {{ $t('common.access') }}</h2>

                <AccessManager
                    v-model:jobs="state.jobAccess"
                    :target-id="templateId ?? '0'"
                    :access-types="accessTypes"
                    :access-roles="
                        enumToAccessLevelEnums(AccessLevel, 'enums.docstore.AccessLevel').filter(
                            (e) => e.value === AccessLevel.VIEW || e.value === AccessLevel.EDIT,
                        )
                    "
                />
            </div>

            <SingleHint
                class="my-2"
                hint-id="template_editor_templating"
                to="https://fivenet.app/user-guides/documents/templates"
                :external="true"
                link-target="_blank"
            />

            <div>
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

                <UFormGroup name="content" :label="`${$t('common.content')} ${$t('common.template')}`" required>
                    <ClientOnly>
                        <DocEditor v-model="state.content" split-screen />
                    </ClientOnly>
                </UFormGroup>
            </div>

            <div class="my-2">
                <h2 class="text-sm">{{ $t('common.requirements', 2) }}</h2>

                <TemplateSchemaEditor v-model="schemaEditor" class="mt-2" />
            </div>

            <div class="my-2">
                <h2 class="text-sm">{{ $t('common.content') }} {{ $t('common.access') }}</h2>

                <AccessManager
                    v-model:jobs="state.contentAccess.jobs"
                    v-model:users="state.contentAccess.users"
                    :target-id="templateId ?? '0'"
                    :access-types="contentAccessTypes"
                    :access-roles="enumToAccessLevelEnums(AccessLevel, 'enums.docstore.AccessLevel')"
                    :disabled="true"
                    :show-required="true"
                />
            </div>
        </UContainer>
    </UForm>
</template>
