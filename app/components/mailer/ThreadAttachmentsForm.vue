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
        const call = getGRPCDocStoreClient().listDocuments(req);
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
                <template v-if="attachments[idx]?.data.oneofKind === 'documentId'">
                    <UFormGroup :name="`attachments.${idx}.data.documentId`" class="flex-1">
                        <USelectMenu
                            value-attribute="id"
                            option-attribute="title"
                            :disabled="!canSubmit"
                            class="w-full flex-1"
                            :searchable="listDocuments"
                            searchable-lazy
                            :placeholder="$t('common.document')"
                            :model-value="attachments[idx].data.documentId"
                            @update:model-value="attachments[idx] = { data: { oneofKind: 'documentId', documentId: $event } }"
                        >
                            <template #label="{ selected }">
                                <template v-if="selected">
                                    {{ `DOC-${selected.id}: ${selected?.title}` }}
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
            :ui="{ rounded: 'rounded-full' }"
            icon="i-mdi-plus"
            :disabled="!canSubmit || attachments.length >= 3"
            :class="attachments.length ? 'mt-2' : ''"
            @click="
                attachments.push({
                    data: { oneofKind: 'documentId', documentId: 0 },
                })
            "
        />
    </UFormGroup>
</template>
