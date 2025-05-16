<script lang="ts" setup>
import type { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import type { MessageAttachment } from '~~/gen/ts/resources/mailer/message';
import type { ListDocumentsRequest } from '~~/gen/ts/services/docstore/docstore';

const props = defineProps<{
    modelValue: MessageAttachment[];
    canSubmit: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', attachments: MessageAttachment[]): void;
}>();

const { $grpc } = useNuxtApp();

const attachments = useVModel(props, 'modelValue', emit);

async function listDocuments(search: string): Promise<DocumentShort[]> {
    const req: ListDocumentsRequest = {
        pagination: {
            offset: 0,
            pageSize: 10,
        },
        search: search ?? '',
        documentIds: [],
        categoryIds: [],
        creatorIds: [],
    };

    try {
        const call = $grpc.docstore.docStore.listDocuments(req);
        const { response } = await call;

        return response.documents;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <UFormGroup name="attachments" :label="$t('common.attachment', 2)">
        <div class="flex flex-col gap-1">
            <div v-for="(_, idx) in attachments" :key="idx" class="flex items-center gap-1">
                <template v-if="attachments[idx]?.data.oneofKind === 'document'">
                    <UFormGroup class="flex-1" :name="`attachments.${idx}.data.documentId`">
                        <USelectMenu
                            class="w-full flex-1"
                            option-attribute="title"
                            :disabled="!canSubmit"
                            :searchable="listDocuments"
                            searchable-lazy
                            :placeholder="$t('common.document')"
                            :model-value="attachments[idx].data.document.id > 0 ? attachments[idx].data.document : undefined"
                            @update:model-value="attachments[idx] = { data: { oneofKind: 'document', document: $event } }"
                        >
                            <template #label="{ selected }">
                                <template v-if="selected">
                                    {{ `DOC-${attachments[idx].data.document.id}: ${attachments[idx].data.document?.title}` }}
                                </template>
                            </template>

                            <template #option="{ option: document }">
                                {{ `DOC-${document.id}: ${document?.title}` }}
                            </template>

                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>

                            <template #empty> {{ $t('common.not_found', [$t('common.document', 2)]) }} </template>
                        </USelectMenu>
                    </UFormGroup>
                </template>

                <UButton
                    :ui="{ rounded: 'rounded-full' }"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    @click="attachments.splice(idx, 1)"
                />
            </div>
        </div>

        <UButton
            :class="attachments.length ? 'mt-2' : ''"
            :ui="{ rounded: 'rounded-full' }"
            icon="i-mdi-plus"
            :disabled="!canSubmit || attachments.length >= 3"
            @click="
                attachments.push({
                    data: { oneofKind: 'document', document: { id: 0 } },
                })
            "
        />

        <UAlert
            class="mt-2"
            icon="i-mdi-information-outline"
            :description="$t('components.mailer.ThreadAttachmentsForm.document_title_warning')"
        />
    </UFormGroup>
</template>
