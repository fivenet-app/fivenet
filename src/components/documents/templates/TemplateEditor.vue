<script lang="ts" setup>
import { ErrorMessage, Field, useForm } from 'vee-validate';
import { object, string } from 'yup';
import { toTypedSchema } from '@vee-validate/yup';
import { CreateTemplateRequest, GetTemplateRequest, UpdateTemplateRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { RpcError } from 'grpc-web';
import { DocumentTemplate, ObjectSpecs, TemplateRequirements, TemplateSchema } from '@fivenet/gen/resources/documents/templates_pb';
import TemplateSchemaEditor from './TemplateSchemaEditor.vue';
import { TemplateSchemaEditorValue } from './TemplateSchemaEditor.vue';

const { $grpc } = useNuxtApp();

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

async function createTemplate(): Promise<void> {
    if (props.templateId) return updateTemplate();

    return new Promise(async (res, rej) => {
        const req = new CreateTemplateRequest();
        const tpl = new DocumentTemplate();
        tpl.setTitle(title.value);
        tpl.setDescription(description.value);
        tpl.setContentTitle(contentTitle.value);
        tpl.setContent(content.value);

        const tRequirements = new TemplateRequirements();
        tRequirements.setUsers((new ObjectSpecs).setRequired(schema.value.users.req).setMin(schema.value.users.min).setMax(schema.value.users.max));
        tRequirements.setDocuments((new ObjectSpecs).setRequired(schema.value.documents.req).setMin(schema.value.documents.min).setMax(schema.value.documents.max));
        tRequirements.setVehicles((new ObjectSpecs).setRequired(schema.value.vehicles.req).setMin(schema.value.vehicles.min).setMax(schema.value.vehicles.max));

        const tSchema = new TemplateSchema();
        tSchema.setRequirements(tRequirements);

        tpl.setSchema(tSchema);

        req.setTemplate(tpl);

        try {
            const resp = await $grpc.getDocStoreClient().
                createTemplate(req, null);

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

        const tRequirements = new TemplateRequirements();
        tRequirements.setUsers((new ObjectSpecs).setRequired(schema.value.users.req).setMin(schema.value.users.min).setMax(schema.value.users.max));
        tRequirements.setDocuments((new ObjectSpecs).setRequired(schema.value.documents.req).setMin(schema.value.documents.min).setMax(schema.value.documents.max));
        tRequirements.setVehicles((new ObjectSpecs).setRequired(schema.value.vehicles.req).setMin(schema.value.vehicles.min).setMax(schema.value.vehicles.max));

        const tSchema = new TemplateSchema();
        tSchema.setRequirements(tRequirements);

        tpl.setSchema(tSchema);

        try {
            const resp = await $grpc.getDocStoreClient().
                updateTemplate(req, null);

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
            schema: string(),
        }),
    ),
});

const onSubmit = handleSubmit(async (): Promise<void> => await createTemplate());

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
});
</script>

<template>
    <div class="text-neutral">
        <form @submit="onSubmit">
            <label for="title" class="block font-medium text-sm mt-2">Title</label>
            <div>
                <Field as="textarea" rows="1" name="title" id="title"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    v-model="title" />
                <ErrorMessage name="title" as="p" class="mt-2 text-sm text-error-400" />
            </div>
            <label for="description" class="block font-medium text-sm mt-2">Description</label>
            <div>
                <Field as="textarea" rows="4" name="description" id="description"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    v-model="description" />
                <ErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
            </div>
            <label for="contentTitle" class="block font-medium text-sm mt-2">Content Title</label>
            <div>
                <Field as="textarea" rows="1" name="contentTitle" id="contentTitle"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    v-model="contentTitle" />
                <ErrorMessage name="contentTitle" as="p" class="mt-2 text-sm text-error-400" />
            </div>
            <label for="content" class="block font-medium text-sm mt-2">Content</label>
            <div>
                <Field as="textarea" rows="4" name="content" id="content"
                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                    v-model="content" />
                <ErrorMessage name="content" as="p" class="mt-2 text-sm text-error-400" />
            </div>
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
