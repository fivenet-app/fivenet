<script lang="ts" setup>
import { ErrorMessage, Field, useForm } from 'vee-validate';
import { object, string } from 'yup';
import { toTypedSchema } from '@vee-validate/yup';
import { CreateTemplateRequest, GetTemplateRequest, UpdateTemplateRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { RpcError } from 'grpc-web';
import { Template, ObjectSpecs, TemplateRequirements, TemplateSchema, TemplateJobAccess } from '@fivenet/gen/resources/documents/templates_pb';
import TemplateSchemaEditor from './TemplateSchemaEditor.vue';
import { TemplateSchemaEditorValue, ObjectSpecsValue } from './TemplateSchemaEditor.vue';
import { useNotificationsStore } from '~/store/notifications';
import { Job, JobGrade } from '@fivenet/gen/resources/jobs/jobs_pb';
import { CompleteJobsRequest } from '@fivenet/gen/services/completor/completor_pb';
import { watchDebounced } from '@vueuse/core';
import DocumentAccessEntry from '~/components/documents/DocumentAccessEntry.vue';
import { ACCESS_LEVEL } from '@fivenet/gen/resources/documents/access_pb';
import {
    PlusIcon,
} from '@heroicons/vue/20/solid';
import { useAuthStore } from '~/store/auth';
import { DocumentAccess, DocumentJobAccess, DocumentUserAccess } from '@fivenet/gen/resources/documents/documents_pb';
import { ACCESS_LEVEL_Util } from '@fivenet/gen/resources/documents/access.pb_enums';

const { $grpc } = useNuxtApp();
const { t } = useI18n();
const authStore = useAuthStore();

const notifications = useNotificationsStore();

const props = defineProps({
    templateId: {
        type: Number,
        required: false,
    }
});

const activeChar = computed(() => authStore.getActiveChar);

const maxAccessEntries = 10;

const title = ref<string>('');
const description = ref<string>('');
const contentTitle = ref<string>('');
const content = ref<string>('');
const schema = ref<TemplateSchemaEditorValue>({
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
const access = ref<Map<number, { id: number, type: number, values: { job?: string, accessrole?: ACCESS_LEVEL, minimumrank?: number } }>>(new Map());

const accessTypes = [
    { id: 1, name: t('common.job', 2) },
];

function addAccessEntry(): void {
    if (access.value.size > maxAccessEntries - 1) {
        notifications.dispatchNotification({
            title: t('notifications.max_access_entry.title'),
            content: t('notifications.max_access_entry.content', [maxAccessEntries]),
            type: 'error'
        });
        return;
    }

    let id = access.value.size > 0 ? [...access.value.keys()].pop() as number + 1 : 0;
    access.value.set(id, {
        id,
        type: 1,
        values: {}
    })
}

function removeAccessEntry(event: {
    id: number
}): void {
    access.value.delete(event.id);
}

function updateAccessEntryType(event: {
    id: number,
    type: number
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.type = event.type;
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryName(event: {
    id: number,
    job?: Job,
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    if (event.job) {
        accessEntry.values.job = event.job.getName();

        access.value.set(event.id, accessEntry);
    }
}

function updateAccessEntryRank(event: {
    id: number,
    rank: JobGrade
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.values.minimumrank = event.rank.getGrade();
    access.value.set(event.id, accessEntry);
}

function updateAccessEntryAccess(event: {
    id: number,
    access: ACCESS_LEVEL
}): void {
    const accessEntry = access.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.values.accessrole = event.access;
    access.value.set(event.id, accessEntry);
}

const contentAccess = ref<Map<number, { id: number, type: number, values: { job?: string, char?: number, accessrole?: ACCESS_LEVEL, minimumrank?: number } }>>(new Map());

const contentAccessTypes = [
    { id: 0, name: t('common.citizen', 2) },
    { id: 1, name: t('common.job', 2) },
];

function addContentAccessEntry(): void {
    if (contentAccess.value.size > maxAccessEntries - 1) {
        notifications.dispatchNotification({
            title: t('notifications.max_access_entry.title'),
            content: t('notifications.max_access_entry.content', [maxAccessEntries]),
            type: 'error'
        });
        return;
    }

    let id = contentAccess.value.size > 0 ? [...contentAccess.value.keys()].pop() as number + 1 : 0;
    contentAccess.value.set(id, {
        id,
        type: 1,
        values: {}
    })
}

function removeContentAccessEntry(event: {
    id: number
}): void {
    contentAccess.value.delete(event.id);
}

function updateContentAccessEntryType(event: {
    id: number,
    type: number
}): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.type = event.type;
    contentAccess.value.set(event.id, accessEntry);
}

function updateContentAccessEntryName(event: {
    id: number,
    job?: Job,
}): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) return;

    if (event.job) {
        accessEntry.values.job = event.job.getName();

        contentAccess.value.set(event.id, accessEntry);
    }
}

function updateContentAccessEntryRank(event: {
    id: number,
    rank: JobGrade
}): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.values.minimumrank = event.rank.getGrade();
    contentAccess.value.set(event.id, accessEntry);
}

function updateContentAccessEntryAccess(event: {
    id: number,
    access: ACCESS_LEVEL
}): void {
    const accessEntry = contentAccess.value.get(event.id);
    if (!accessEntry) return;

    accessEntry.values.accessrole = event.access;
    contentAccess.value.set(event.id, accessEntry);
}

const entriesRank = ref<JobGrade[]>([]);
const filteredRank = ref<JobGrade[]>([]);
const queryRank = ref('');

function createObjectSpec(v: ObjectSpecsValue): ObjectSpecs {
    const o = new ObjectSpecs();
    o.setRequired(v.req ?? false);
    if (v.min > 0) {
        o.setMin(v.min);
    }
    if (v.max > 0) {
        o.setMax(v.max);
    }
    return o;
}

async function createTemplate(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new CreateTemplateRequest();
        const tpl = new Template();
        tpl.setTitle(title.value);
        tpl.setDescription(description.value);
        tpl.setContentTitle(contentTitle.value);
        tpl.setContent(content.value);

        const tRequirements = new TemplateRequirements();
        tRequirements.setUsers(createObjectSpec(schema.value.users));
        tRequirements.setDocuments(createObjectSpec(schema.value.documents));
        tRequirements.setVehicles(createObjectSpec(schema.value.vehicles));

        const tSchema = new TemplateSchema();
        tSchema.setRequirements(tRequirements);
        tpl.setSchema(tSchema);
        req.setTemplate(tpl);

        const jobAccesses = new Array<TemplateJobAccess>();
        access.value.forEach(entry => {
            if (entry.values.accessrole === undefined) return;

            if (entry.type === 1) {
                if (!entry.values.job) return;

                const job = new TemplateJobAccess();
                job.setJob(entry.values.job);
                job.setMinimumgrade(entry.values.minimumrank ? entry.values.minimumrank : 0);
                job.setAccess(ACCESS_LEVEL_Util.fromInt(entry.values.accessrole));

                jobAccesses.push(job);
            }
        });
        tpl.setJobAccessList(jobAccesses);

        const reqAccess = new DocumentAccess();
        contentAccess.value.forEach(entry => {
            if (entry.values.accessrole === undefined) return;

            if (entry.type === 0) {
                if (!entry.values.char) return;

                const user = new DocumentUserAccess();
                user.setUserId(entry.values.char);
                user.setAccess(ACCESS_LEVEL_Util.fromInt(entry.values.accessrole));

                reqAccess.addUsers(user);
            } else if (entry.type === 1) {
                if (!entry.values.job) return;

                const job = new DocumentJobAccess();
                job.setJob(entry.values.job);
                job.setMinimumgrade(entry.values.minimumrank ? entry.values.minimumrank : 0);
                job.setAccess(ACCESS_LEVEL_Util.fromInt(entry.values.accessrole));

                reqAccess.addJobs(job);
            }
        });
        tpl.setContentAccess(reqAccess);

        try {
            const resp = await $grpc.getDocStoreClient().
                createTemplate(req, null);

            notifications.dispatchNotification({
                title: 'Template: Created',
                content: 'Template created successfully.',
                type: 'success',
            });

            await navigateTo({ name: 'documents-templates-id', params: { id: resp.getId() } });

            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function updateTemplate(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new UpdateTemplateRequest();
        const tpl = new Template();
        tpl.setTitle(title.value);
        tpl.setDescription(description.value);
        tpl.setContentTitle(contentTitle.value);
        tpl.setContent(content.value);

        const tRequirements = new TemplateRequirements();
        tRequirements.setUsers(createObjectSpec(schema.value.users));
        tRequirements.setDocuments(createObjectSpec(schema.value.documents));
        tRequirements.setVehicles(createObjectSpec(schema.value.vehicles));

        const tSchema = new TemplateSchema();
        tSchema.setRequirements(tRequirements);

        tpl.setSchema(tSchema);

        try {
            const resp = await $grpc.getDocStoreClient().
                updateTemplate(req, null);

            notifications.dispatchNotification({
                title: 'Template: Updated',
                content: 'Template updated successfully.',
                type: 'success',
            });

            await navigateTo({ name: 'documents-templates-id', params: { id: resp.getId() } });

            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const { handleSubmit } = useForm({
    validationSchema: toTypedSchema(
        object({
            title: string().required().min(3).max(255),
            description: string().required().max(512),
            contentTitle: string().required().min(3).max(1024),
            content: string().required().min(6).max(15360),
            schema: object().optional(),
        }),
    ),
});

const onSubmit = handleSubmit(async (): Promise<void> => {
    if (props.templateId && props.templateId > 0) {
        return updateTemplate();
    } else {
        await createTemplate();
    }
});

onMounted(async () => {
    if (props.templateId) {
        const req = new GetTemplateRequest();
        req.setTemplateId(props.templateId);
        req.setRender(false);

        try {
            const resp = (await $grpc.getDocStoreClient().getTemplate(req, null)).getTemplate();
            if (!resp) return;

            title.value = resp.getTitle();
            description.value = resp.getDescription();
            contentTitle.value = resp.getContentTitle();
            content.value = resp.getContent();

            schema.value.users.req = resp.getSchema()?.getRequirements()?.getUsers()?.getRequired() ?? false;
            schema.value.users.min = resp.getSchema()?.getRequirements()?.getUsers()?.getMin() ?? 0;
            schema.value.users.max = resp.getSchema()?.getRequirements()?.getUsers()?.getMax() ?? 0;

            schema.value.documents.req = resp.getSchema()?.getRequirements()?.getDocuments()?.getRequired() ?? false;
            schema.value.documents.min = resp.getSchema()?.getRequirements()?.getDocuments()?.getMin() ?? 0;
            schema.value.documents.max = resp.getSchema()?.getRequirements()?.getDocuments()?.getMax() ?? 0;

            schema.value.vehicles.req = resp.getSchema()?.getRequirements()?.getVehicles()?.getRequired() ?? false;
            schema.value.vehicles.min = resp.getSchema()?.getRequirements()?.getVehicles()?.getMin() ?? 0;
            schema.value.vehicles.max = resp.getSchema()?.getRequirements()?.getVehicles()?.getMax() ?? 0;
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
        }
    } else {
        access.value.set(0, { id: 0, type: 1, values: { job: activeChar.value?.getJob(), minimumrank: 1, accessrole: ACCESS_LEVEL.VIEW } });
    }

    const req = new CompleteJobsRequest();
    req.setExactMatch(true);
    req.setCurrentJob(true);

    try {
        const resp = await $grpc.getCompletorClient().completeJobs(req, null);
        entriesRank.value = resp.getJobsList()[0].getGradesList();
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
    }
});

watchDebounced(queryRank, async () => filteredRank.value = entriesRank.value.filter(r => r.getLabel().startsWith(queryRank.value)), { debounce: 700, maxWait: 1850 });
</script>

<template>
    <div class="text-neutral">
        <form @submit="onSubmit">
            <label for="title" class="block font-medium text-sm mt-2">
                {{ $t('common.template') }} {{ $t('common.title') }}
            </label>
            <div>
                <Field as="textarea" rows="1" name="title" id="title"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    v-model="title" />
                <ErrorMessage name="title" as="p" class="mt-2 text-sm text-error-400" />
            </div>
            <label for="description" class="block font-medium text-sm mt-2">
                {{ $t('common.template') }} {{ $t('common.description') }}
            </label>
            <div>
                <Field as="textarea" rows="4" name="description" id="description"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    v-model="description" />
                <ErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
            </div>
            <div class="my-3">
                <h2 class="text-neutral">
                    {{ $t('common.template') }} {{ $t('common.access') }}
                </h2>
                <DocumentAccessEntry v-for="entry in access.values()" :key="entry.id" :init="entry"
                    :access-types="accessTypes" :access-roles="[ACCESS_LEVEL.VIEW, ACCESS_LEVEL.EDIT]"
                    @typeChange="updateAccessEntryType($event)" @nameChange="updateAccessEntryName($event)"
                    @rankChange="updateAccessEntryRank($event)" @accessChange="updateAccessEntryAccess($event)"
                    @deleteRequest="removeAccessEntry($event)" />
                <button type="button"
                    class="p-2 rounded-full bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    data-te-toggle="tooltip" :title="$t('components.documents.document_editor.add_permission')"
                    @click="addAccessEntry()">
                    <PlusIcon class="w-5 h-5" aria-hidden="true" />
                </button>
            </div>
            <label for="contentTitle" class="block font-medium text-sm mt-2">
                {{ $t('common.content') }} {{ $t('common.title') }}
            </label>
            <div>
                <Field as="textarea" rows="1" name="contentTitle" id="contentTitle"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    v-model="contentTitle" />
                <ErrorMessage name="contentTitle" as="p" class="mt-2 text-sm text-error-400" />
                <p class="text-neutral">
                    <NuxtLink :external="true" target="_blank" to="https://pkg.go.dev/html/template">
                        Golang {{ $t('common.template') }}
                    </NuxtLink>
                </p>
            </div>
            <label for="content" class="block font-medium text-sm mt-2">
                {{ $t('common.content') }} {{ $t('common.template') }}
            </label>
            <div>
                <Field as="textarea" rows="4" name="content" id="content"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    v-model="content" />
                <ErrorMessage name="content" as="p" class="mt-2 text-sm text-error-400" />
                <p class="text-neutral">
                    <NuxtLink :external="true" target="_blank" to="https://pkg.go.dev/html/template">
                        Golang {{ $t('common.template') }}
                    </NuxtLink>
                </p>
            </div>
            <div>
                <TemplateSchemaEditor v-model="schema" class="mt-2" />
            </div>
            <div class="my-3">
                <h2 class="text-neutral">
                    {{ $t('common.content') }} {{ $t('common.access') }}
                </h2>
                <DocumentAccessEntry v-for="entry in contentAccess.values()" :key="entry.id" :init="entry"
                    :access-types="contentAccessTypes" @typeChange="updateContentAccessEntryType($event)"
                    @nameChange="updateContentAccessEntryName($event)" @rankChange="updateContentAccessEntryRank($event)"
                    @accessChange="updateContentAccessEntryAccess($event)"
                    @deleteRequest="removeContentAccessEntry($event)" />
                <button type="button"
                    class="p-2 rounded-full bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                    data-te-toggle="tooltip" :title="$t('components.documents.document_editor.add_permission')"
                    @click="addContentAccessEntry()">
                    <PlusIcon class="w-5 h-5" aria-hidden="true" />
                </button>
            </div>
            <button type="submit"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300 mt-4">
                {{ $t('common.create') }}
            </button>
        </form>
    </div>
</template>
