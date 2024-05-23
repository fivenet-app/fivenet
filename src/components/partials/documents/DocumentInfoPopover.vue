<script lang="ts" setup>
import { type Document, type DocumentShort } from '~~/gen/ts/resources/documents/documents';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import DataErrorBlock from '../data/DataErrorBlock.vue';

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(
    defineProps<{
        documentId?: string;
        document?: Document | DocumentShort | undefined;
        trailing?: boolean;
    }>(),
    {
        trailing: true,
    },
);

const { popover } = useAppConfig();

const {
    data,
    refresh,
    pending: loading,
    error,
} = useLazyAsyncData(
    `document-info-${props.documentId ?? props.document?.id}`,
    () => getDocument(props.documentId ?? props.document?.id ?? '0'),
    { immediate: false },
);

async function getDocument(id: string): Promise<Document> {
    try {
        const call = getGRPCDocStoreClient().getDocument({
            documentId: id,
            infoOnly: true,
        });
        const { response } = await call;

        return response.document!;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const document = computed(() => props.document || data.value);

const opened = ref(false);
watchOnce(opened, async () => {
    if (!props.document) {
        refresh();
    } else {
        useTimeoutFn(async () => refresh(), popover.waitTime);
    }
});
</script>

<template>
    <template v-if="!document && !documentId">
        <span class="inline-flex items-center">
            {{ $t('common.na') }}
        </span>
    </template>
    <UPopover v-else :ui="{ container: 'max-w-[50%]' }">
        <UButton
            v-bind="$attrs"
            variant="link"
            :padded="false"
            class="line-clamp-2 inline-flex w-full items-center gap-1 whitespace-normal break-words p-px"
            :trailing-icon="trailing ? 'i-mdi-chevron-down' : undefined"
            @click="opened = true"
        >
            <slot name="title">
                <USkeleton v-if="!document && loading" class="h-8 w-[125px]" />

                <template v-else>
                    <UBadge v-if="document?.category">
                        {{ document.category.name }}
                    </UBadge>

                    <span> {{ document?.title }} </span>
                </template>
            </slot>
        </UButton>

        <template #panel>
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

                    <UButton
                        v-if="can('DocStoreService.GetDocument')"
                        variant="link"
                        icon="i-mdi-eye"
                        :to="{ name: 'documents-id', params: { id: documentId ?? document?.id ?? 0 } }"
                    >
                        {{ $t('common.open') }}
                    </UButton>
                </div>

                <div v-if="error">
                    <DataErrorBlock :title="$t('common.unable_to_load', [$t('common.document', 2)])" :retry="refresh" />
                </div>

                <div v-else-if="loading && !document" class="flex flex-col gap-2 text-gray-900 dark:text-white">
                    <USkeleton class="h-8 w-[250px]" />

                    <div class="flex flex-row items-center gap-2">
                        <USkeleton class="h-7 w-[75px]" />
                        <USkeleton class="h-6 w-[250px]" />
                    </div>
                </div>

                <div v-else-if="document" class="flex flex-col gap-2 text-gray-900 dark:text-white">
                    <UButton variant="link" :padded="false" :to="{ name: 'documents-id', params: { id: document.id ?? 0 } }">
                        <UBadge v-if="document?.category">
                            {{ document.category.name }}
                        </UBadge>

                        <span class="line-clamp-1 text-lg hover:line-clamp-2">{{ document.title }}</span>
                    </UButton>

                    <div class="flex flex-row flex-wrap gap-2">
                        <UBadge v-if="document.state" class="inline-flex gap-1" size="md">
                            <UIcon name="i-mdi-note-check" class="size-5" />
                            <span>
                                {{ document.state }}
                            </span>
                        </UBadge>

                        <div v-if="document.createdAt" class="flex flex-row items-center gap-1">
                            <span>{{ $t('common.created_at') }}:</span>
                            <GenericTime :value="document.createdAt" />
                        </div>

                        <div v-if="document.updatedAt" class="flex flex-row items-center gap-1">
                            <span>{{ $t('common.updated') }}:</span>
                            <GenericTime :value="document.updatedAt" :ago="true" />
                        </div>

                        <div v-if="document.deletedAt" class="flex flex-row items-center gap-1 font-bold">
                            <UIcon name="i-mdi-trash-can" class="mr-1.5 size-5 shrink-0" />
                            <span>{{ $t('common.deleted') }}</span>
                        </div>

                        <div v-if="document.creator" class="flex flex-row items-center justify-start gap-1">
                            <CitizenInfoPopover :user="document.creator" />
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </UPopover>
</template>
