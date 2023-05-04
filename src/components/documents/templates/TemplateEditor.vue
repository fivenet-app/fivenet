<script lang="ts" setup>
import { ErrorMessage, Field, useForm } from 'vee-validate';
import { object, string } from 'yup';
import { toTypedSchema } from '@vee-validate/yup';
import { CreateTemplateRequest, GetTemplateRequest, UpdateTemplateRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { RpcError } from 'grpc-web';
import { DocumentTemplate, ObjectSpecs, TemplateRequirements, TemplateSchema } from '@fivenet/gen/resources/documents/templates_pb';
import TemplateSchemaEditor from './TemplateSchemaEditor.vue';
import { TemplateSchemaEditorValue, ObjectSpecsValue } from './TemplateSchemaEditor.vue';
import { useNotificationsStore } from '~/store/notifications';
import { JobGrade } from '@fivenet/gen/resources/jobs/jobs_pb';
import { CompleteJobsRequest } from '@fivenet/gen/services/completor/completor_pb';
import { watchDebounced } from '@vueuse/core';
import { Combobox, ComboboxInput, ComboboxButton, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { CheckIcon } from '@heroicons/vue/24/solid';

const { $grpc } = useNuxtApp();

const notifications = useNotificationsStore();

const props = defineProps({
    templateId: {
        type: Number,
        required: false,
    }
});

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

const entriesRank = ref<JobGrade[]>([]);
const filteredRank = ref<JobGrade[]>([]);
const queryRank = ref('');
const selectedRank = ref<JobGrade>();

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
        const tpl = new DocumentTemplate();
        tpl.setTitle(title.value);
        tpl.setDescription(description.value);
        tpl.setContentTitle(contentTitle.value);
        tpl.setContent(content.value);
        if (selectedRank.value) tpl.setJobGrade(selectedRank.value.getGrade());

        const tRequirements = new TemplateRequirements();
        tRequirements.setUsers(createObjectSpec(schema.value.users));
        tRequirements.setDocuments(createObjectSpec(schema.value.documents));
        tRequirements.setVehicles(createObjectSpec(schema.value.vehicles));

        const tSchema = new TemplateSchema();
        tSchema.setRequirements(tRequirements);

        tpl.setSchema(tSchema);

        req.setTemplate(tpl);

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
        const tpl = new DocumentTemplate();
        tpl.setTitle(title.value);
        tpl.setDescription(description.value);
        tpl.setContentTitle(contentTitle.value);
        tpl.setContent(content.value);
        if (selectedRank.value) tpl.setJobGrade(selectedRank.value.getGrade());

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
            title: string().required().min(3).max(24),
            description: string().required(),
            contentTitle: string().required().min(3).max(24),
            content: string().required().min(6).max(70),
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
            <label for="rank" class="block font-medium text-sm mt-2">
                {{ $t('common.min') }} {{ $t('common.rank') }}
            </label>
            <Combobox as="div" v-model="selectedRank">
                <div class="relative">
                    <ComboboxButton as="div">
                        <ComboboxInput
                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            @change="queryRank = $event.target.value" :display-value="(rank: any) => rank?.getLabel()" />
                    </ComboboxButton>

                    <ComboboxOptions v-if="entriesRank.length > 0"
                        class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-60 sm:text-sm">
                        <ComboboxOption v-for="rank in entriesRank" :key="rank.getGrade()" :value="rank" as="minimumrank"
                            v-slot="{ active, selected }">
                            <li
                                :class="['relative cursor-default select-none py-2 pl-8 pr-4 text-neutral', active ? 'bg-primary-500' : '']">
                                <span :class="['block truncate', selected && 'font-semibold']">
                                    {{ rank.getLabel() }}
                                </span>

                                <span v-if="selected"
                                    :class="[active ? 'text-neutral' : 'text-primary-500', 'absolute inset-y-0 left-0 flex items-center pl-1.5']">
                                    <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                </span>
                            </li>
                        </ComboboxOption>
                    </ComboboxOptions>
                </div>
            </Combobox>
            <div>
                <TemplateSchemaEditor v-model="schema" class="mt-2" />
            </div>
            <button type="submit"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300 mt-4">
                {{ $t('common.create') }}
            </button>
        </form>
    </div>
</template>
