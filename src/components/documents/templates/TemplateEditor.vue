<script lang="ts" setup>
import { ErrorMessage, Field, useForm } from 'vee-validate';
import { object, string } from 'yup';
import { toTypedSchema } from '@vee-validate/yup';
import { CreateTemplateRequest, UpdateTemplateRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { RpcError } from 'grpc-web';
import { DocumentTemplate, TemplateSchema } from '@fivenet/gen/resources/documents/templates_pb';

const { $grpc } = useNuxtApp();

const router = useRouter();

const props = defineProps({
    templateId: {
        type: Number,
        required: false,
    }
});

async function createTemplate(title: string, description: string, contentTitle: string, content: string, schema: string): Promise<void> {
    if (props.templateId) {
        return updateTemplate(title, description, contentTitle, content, schema);
    }

    return new Promise(async (res, rej) => {
        const req = new CreateTemplateRequest();
        const tpl = new DocumentTemplate();
        tpl.setTitle(title);
        tpl.setDescription(description);
        tpl.setContentTitle(contentTitle);
        tpl.setContent(content);
        const tplSchema = new TemplateSchema();
        tpl.setSchema(tplSchema);

        req.setTemplate(tpl);

        try {
            const resp = await $grpc.getDocStoreClient().
                createTemplate(req, null);

            await router.push({ name: 'documents-templates-id', params: { id: resp.getId() } });
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function updateTemplate(title: string, description: string, contentTitle: string, content: string, schema: string): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new UpdateTemplateRequest();
        const tpl = new DocumentTemplate();
        tpl.setTitle(title);
        tpl.setDescription(description);
        tpl.setContentTitle(contentTitle);
        tpl.setContent(content);
        const tplSchema = new TemplateSchema();
        tpl.setSchema(tplSchema);

        try {
            const resp = await $grpc.getDocStoreClient().
                updateTemplate(req, null);

            await router.push({ name: 'documents-templates-id', params: { id: resp.getId() } });
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

const onSubmit = handleSubmit(async (values): Promise<void> => await createTemplate(values.title, values.description, values.contentTitle, values.content, values.schema!));
</script>

<template>
    <div>
        <form @submit="onSubmit">
            <label for="title" class="block text-sm font-medium leading-6 text-gray-100">Title</label>
            <div class="mt-2">
                <Field as="textarea" rows="4" name="title" id="title"
                    class="block w-full rounded-md border-0 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:py-1.5 sm:text-sm sm:leading-6" />
                <ErrorMessage name="title" as="p" class="mt-2 text-sm text-red-500" />
            </div>
            <label for="title" class="block text-sm font-medium leading-6 text-gray-100">Description</label>
            <div class="mt-2">
                <Field as="textarea" rows="4" name="description" id="description"
                    class="block w-full rounded-md border-0 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:py-1.5 sm:text-sm sm:leading-6" />
                <ErrorMessage name="description" as="p" class="mt-2 text-sm text-red-500" />
            </div>
            <label for="contentTitle" class="block text-sm font-medium leading-6 text-gray-100">Content Title</label>
            <div class="mt-2">
                <Field as="textarea" rows="4" name="contentTitle" id="contentTitle"
                    class="block w-full rounded-md border-0 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:py-1.5 sm:text-sm sm:leading-6" />
                <ErrorMessage name="contentTitle" as="p" class="mt-2 text-sm text-red-500" />
            </div>
            <label for="content" class="block text-sm font-medium leading-6 text-gray-100">Content</label>
            <div class="mt-2">
                <Field as="textarea" rows="4" name="content" id="content"
                    class="block w-full rounded-md border-0 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:py-1.5 sm:text-sm sm:leading-6" />
                <ErrorMessage name="content" as="p" class="mt-2 text-sm text-red-500" />
            </div>
            <label for="schema" class="block text-sm font-medium leading-6 text-gray-100">Schema</label>
            <div class="mt-2">
                <Field as="textarea" rows="4" name="schema" id="schema"
                    class="block w-full rounded-md border-0 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:py-1.5 sm:text-sm sm:leading-6" />
                <ErrorMessage name="schema" as="p" class="mt-2 text-sm text-red-500" />
            </div>
            <button type="submit"
                class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-primary-600 text-neutral hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
                Create
            </button>
        </form>
    </div>
</template>
