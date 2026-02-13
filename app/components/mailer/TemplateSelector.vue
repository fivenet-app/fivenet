<script lang="ts" setup>
import type { Editor, JSONContent } from '@tiptap/core';
import { useMailerStore } from '~/stores/mailer';
import { getMailerMailerClient } from '~~/gen/ts/clients';
import { Struct } from '~~/gen/ts/google/protobuf/struct';
import type { Template } from '~~/gen/ts/resources/mailer/templates/template';
import type { ListTemplatesResponse } from '~~/gen/ts/services/mailer/mailer';

const props = defineProps<{
    editor: Editor | undefined;
}>();

const mailerMailerClient = await getMailerMailerClient();

defineOptions({
    inheritAttrs: false,
});

const mailerStore = useMailerStore();
const { selectedEmail } = storeToRefs(mailerStore);

const { data: templates } = useLazyAsyncData(`mailer-templates:${selectedEmail.value!.id}`, () => listTemplates());

async function listTemplates(): Promise<ListTemplatesResponse> {
    try {
        const call = mailerMailerClient.listTemplates({
            emailId: selectedEmail.value!.id,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const selectedTemplate = ref<Template | undefined>(undefined);

function selectTemplate(template: Template | undefined): void {
    if (!props.editor || !template) return;

    if (template.content?.tiptapJson) {
        const content = Struct.toJson(template.content.tiptapJson) as JSONContent;
        props.editor.commands.insertContent(content);
    } else if (template.content?.rawHtml) {
        props.editor.commands.insertContent(template.content.rawHtml);
    }

    selectedTemplate.value = undefined;
}
</script>

<template>
    <ClientOnly v-if="editor && templates?.templates && templates?.templates.length > 0">
        <USelectMenu
            v-if="!selectedTemplate"
            v-model="selectedTemplate"
            class="mb-1 min-w-48"
            :items="templates?.templates"
            label-key="title"
            :placeholder="$t('common.template')"
            size="sm"
            v-bind="$attrs"
        >
            <template #empty>
                {{ $t('common.not_found', [$t('common.template', 2)]) }}
            </template>
        </USelectMenu>

        <UFieldGroup v-else v-bind="$attrs">
            <UButton
                :label="$t('components.partials.tiptap_editor.insert')"
                icon="i-mdi-plus"
                @click="() => selectTemplate(selectedTemplate)"
            />

            <UButton :label="$t('common.cancel')" color="error" icon="i-mdi-cancel" @click="selectedTemplate = undefined" />
        </UFieldGroup>
    </ClientOnly>
</template>
