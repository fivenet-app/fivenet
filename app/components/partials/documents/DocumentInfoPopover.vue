<script lang="ts" setup>
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { ClassProp } from '~/utils/types';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import type { Document, DocumentShort } from '~~/gen/ts/resources/documents/documents';
import CategoryBadge from './CategoryBadge.vue';

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        documentId?: number;
        document?: Document | DocumentShort;
        hideTrailing?: boolean;
        hideCategory?: boolean;
        showId?: boolean;
        loadOnOpen?: boolean;
        buttonClass?: ClassProp;
        disableTooltip?: boolean;
    }>(),
    {
        documentId: undefined,
        document: undefined,
        hideTrailing: false,
        hideCategory: false,
        showId: false,
        loadOnOpen: false,
        buttonClass: 'items-center',
        disableTooltip: false,
    },
);

const { can } = useAuth();

const { popover } = useAppConfig();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

const documentId = computed(() => props.documentId ?? props.document?.id ?? 0);

const { data, refresh, status, error } = useLazyAsyncData(
    `document-info-${documentId.value}`,
    () => getDocument(documentId.value),
    {
        immediate: !props.loadOnOpen,
    },
);

async function getDocument(id: number): Promise<Document> {
    const call = documentsDocumentsClient.getDocument({
        documentId: id,
        infoOnly: true,
    });
    const { response } = await call;

    return response.document!;
}

const document = computed(() => data.value || props.document);

// Invalidate the data when the documentId changes
watch(documentId, (val) => {
    if (val === documentId.value) {
        return;
    }

    data.value = undefined;
});

const opened = ref(false);
watchOnce(opened, async () => {
    if (props.document) {
        useTimeoutFn(async () => refresh(), popover.waitTime);
    } else {
        refresh();
    }
});
</script>

<template>
    <template v-if="!document && !documentId">
        <slot name="title" :document="document" :loading="isRequestPending(status)">
            <span class="inline-flex items-center">
                {{ $t('common.na') }}
            </span>
        </slot>
    </template>

    <UPopover v-else :ui="{ content: 'max-w-xl' }">
        <UButton
            class="line-clamp-2 inline-flex gap-1 p-1 break-words whitespace-normal"
            :class="buttonClass"
            variant="link"
            :trailing-icon="!hideTrailing ? 'i-mdi-chevron-down' : undefined"
            v-bind="$attrs"
            @click="opened = true"
        >
            <slot name="title" :document="document" :loading="isRequestPending(status)">
                <template v-if="!document && isRequestPending(status)">
                    <IDCopyBadge v-if="showId" :id="documentId" prefix="DOC" hide-icon :disable-tooltip="disableTooltip" />
                    <USkeleton v-else class="h-8 w-full min-w-[125px]" />
                </template>

                <template v-else>
                    <IDCopyBadge v-if="showId" :id="documentId" prefix="DOC" hide-icon :disable-tooltip="disableTooltip" />
                    <CategoryBadge v-if="document?.category && !hideCategory" :category="document?.category" />

                    <span class="text-left" v-bind="$attrs">{{ document?.title }} </span>
                </template>
            </slot>
        </UButton>

        <template #content>
            <div class="flex flex-col gap-2 p-4">
                <div class="inline-flex w-full gap-1">
                    <IDCopyBadge
                        :id="documentId ?? document?.id ?? 0"
                        prefix="DOC"
                        :title="{ key: 'notifications.document_view.copy_document_id.title', parameters: {} }"
                        :content="{ key: 'notifications.document_view.copy_document_id.content', parameters: {} }"
                        size="xs"
                        variant="link"
                    />

                    <UTooltip v-if="can('documents.DocumentsService/ListDocuments').value" :text="$t('common.open')">
                        <UButton
                            variant="link"
                            icon="i-mdi-eye"
                            :to="{ name: 'documents-id', params: { id: documentId ?? document?.id ?? 0 } }"
                        >
                            {{ $t('common.open') }}
                        </UButton>
                    </UTooltip>
                </div>

                <div v-if="error">
                    <DataErrorBlock
                        :title="$t('common.unable_to_load', [$t('common.document', 2)])"
                        :error="error"
                        :retry="refresh"
                    />
                </div>

                <div v-else-if="isRequestPending(status) && !document" class="flex flex-col gap-2 text-highlighted">
                    <USkeleton class="h-8 w-[250px]" />

                    <div class="flex flex-row items-center gap-2">
                        <USkeleton class="h-7 w-[75px]" />
                        <USkeleton class="h-6 w-[250px]" />
                    </div>
                </div>

                <div v-else-if="document" class="flex flex-col gap-2 text-highlighted">
                    <UButton
                        variant="link"
                        :to="{ name: 'documents-id', params: { id: document.id ?? 0 } }"
                        :ui="{ base: 'p-0' }"
                    >
                        <CategoryBadge v-if="document?.category" :category="document?.category" size="xs" />

                        <span class="line-clamp-1 text-lg hover:line-clamp-3">{{ document.title }}</span>
                    </UButton>

                    <div>
                        <UBadge v-if="document.state" class="inline-flex gap-1" size="xs">
                            <UIcon class="size-5" name="i-mdi-note-check" />
                            <span>
                                {{ document.state }}
                            </span>
                        </UBadge>
                    </div>

                    <div class="flex flex-row flex-wrap gap-2">
                        <div v-if="document.createdAt" class="flex flex-row items-center gap-1 text-sm">
                            <span>{{ $t('common.created_at') }}:</span>
                            <GenericTime :value="document.createdAt" />
                        </div>

                        <div v-if="document.updatedAt" class="flex flex-row items-center gap-1 text-sm">
                            <span>{{ $t('common.updated') }}:</span>
                            <GenericTime :value="document.updatedAt" ago />
                        </div>

                        <div v-if="document.deletedAt" class="flex flex-row items-center gap-1 font-bold">
                            <UIcon class="mr-1.5 size-5 shrink-0" name="i-mdi-delete" />
                            <span>{{ $t('common.deleted') }}</span>
                        </div>
                    </div>

                    <div v-if="document.creator" class="flex flex-row items-center justify-start gap-1 text-sm">
                        <span>{{ $t('common.creator') }}:</span>
                        <CitizenInfoPopover :user="document.creator" />
                    </div>
                </div>
            </div>
        </template>
    </UPopover>
</template>
