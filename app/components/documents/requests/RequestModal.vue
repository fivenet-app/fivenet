<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { checkDocAccess } from '~/components/documents/helpers';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { getDocumentsDocumentsClient } from '~~/gen/ts/clients';
import { AccessLevel, type DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import type { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ListDocumentReqsResponse } from '~~/gen/ts/services/documents/documents';
import RequestListEntry from './RequestListEntry.vue';

const props = defineProps<{
    access: DocumentAccess;
    doc: DocumentShort;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'refresh'): void;
}>();

const { attr, can, activeChar } = useAuth();

const notifications = useNotificationsStore();

const documentsDocumentsClient = await getDocumentsDocumentsClient();

type RequestType = { key: DocActivityType; attrKey: string };
const requestTypes = [
    { key: props.doc.closed ? DocActivityType.REQUESTED_OPENING : DocActivityType.REQUESTED_CLOSURE, attrKey: 'Closure' },
    { key: DocActivityType.REQUESTED_UPDATE, attrKey: 'Update' },
    { key: DocActivityType.REQUESTED_OWNER_CHANGE, attrKey: 'OwnerChange' },
    { key: DocActivityType.REQUESTED_DELETION, attrKey: 'Deletion' },
] as RequestType[];

const availableRequestTypes = computed<RequestType[]>(() =>
    requestTypes.filter((rt) => attr('documents.DocumentsService/CreateDocumentReq', 'Types', rt.attrKey).value),
);

const schema = z.object({
    reason: z.coerce.string().min(3).max(255),
    requestType: z.enum(DocActivityType).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    requestType: availableRequestTypes.value[0]?.key ?? undefined,
});

const offset = ref(0);

const {
    data: requests,
    status,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.doc.id}-requests-${offset.value}`, () => listDocumnetReqs(props.doc.id));

async function listDocumnetReqs(documentId: number): Promise<ListDocumentReqsResponse> {
    try {
        const call = documentsDocumentsClient.listDocumentReqs({
            pagination: {
                offset: offset.value,
            },
            documentId,
        });
        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function createDocumentRequest(values: Schema): Promise<void> {
    if (values.requestType === undefined) {
        return;
    }

    try {
        const call = documentsDocumentsClient.createDocumentReq({
            documentId: props.doc.id,
            reason: values.reason,
            requestType: values.requestType,
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.documents.requests.created.title' },
            description: { key: 'notifications.documents.requests.created.content' },
            type: NotificationType.SUCCESS,
        });

        // The request should have triggered the "owner override hatch" (creator not part of job anymore)
        if (
            values.requestType === DocActivityType.REQUESTED_OWNER_CHANGE &&
            response.request?.requestType === DocActivityType.OWNER_CHANGED &&
            !!response.request.accepted
        ) {
            emit('refresh');
        }

        refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canDo = computed(() => ({
    create:
        props.doc.creatorId !== activeChar.value?.userId &&
        availableRequestTypes.value.length > 0 &&
        can('documents.DocumentsService/CreateDocumentReq').value &&
        checkDocAccess(props.access, props.doc.creator, AccessLevel.VIEW, undefined, props.doc?.creatorJob),
    update:
        can('documents.DocumentsService/CreateDocumentReq').value &&
        checkDocAccess(
            props.access,
            props.doc.creator,
            AccessLevel.EDIT,
            'documents.DocumentsService/CreateDocumentReq',
            props.doc?.creatorJob,
        ),
    delete:
        can('documents.DocumentsService/DeleteDocumentReq').value &&
        checkDocAccess(
            props.access,
            props.doc.creator,
            AccessLevel.EDIT,
            'documents.DocumentsService/DeleteDocumentReq',
            props.doc?.creatorJob,
        ),
}));

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createDocumentRequest(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UDrawer :title="$t('common.request', 2)" :overlay="false" :close="{ onClick: () => $emit('close', false) }">
        <template #title>
            <span>{{ $t('common.request', 2) }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <div v-if="canDo.create" class="flex flex-row gap-2 md:flex-col">
                    <UFormField class="flex-1" name="requestsType" :label="$t('common.type', 2)" required>
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.requestType"
                                class="w-full"
                                :items="availableRequestTypes"
                                value-key="key"
                                :placeholder="$t('common.type')"
                                :search-input="{ placeholder: $t('common.search_field') }"
                            >
                                <template v-if="state.requestType" #default>
                                    {{ $t(`enums.documents.DocActivityType.${DocActivityType[state.requestType]}`, 2) }}
                                </template>

                                <template #item-label="{ item }">
                                    {{ $t(`enums.documents.DocActivityType.${DocActivityType[item.key]}`, 2) }}
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.type', 2)]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormField>

                    <UFormField name="reason" :label="$t('common.reason')" required>
                        <UTextarea v-model="state.reason" :placeholder="$t('common.reason')" :rows="4" class="w-full" />
                    </UFormField>
                </div>

                <USeparator class="my-2" />

                <ul v-if="isRequestPending(status)" class="mb-6 divide-y divide-default rounded-md" role="list">
                    <li v-for="idx in 2" :key="idx" class="flex justify-between gap-x-4 py-4">
                        <div class="flex min-w-0 gap-x-2 px-2">
                            <div class="min-w-0 flex-auto">
                                <p class="text-base leading-6 font-semibold text-toned">
                                    <USkeleton class="h-8 w-[325px]" />
                                </p>
                                <p class="mt-1 flex gap-1 text-sm leading-5">
                                    <USkeleton class="h-6 w-[350px]" />
                                </p>
                                <p class="mt-1 flex gap-1 text-sm leading-5">
                                    <USkeleton class="h-6 w-[175px]" />
                                </p>
                            </div>
                        </div>
                        <div class="flex shrink-0 items-center gap-x-6 px-2">
                            <div class="hidden gap-1 text-sm sm:flex sm:flex-col sm:items-end">
                                <div class="inline-flex gap-1">
                                    <USkeleton class="h-8 w-[250px]" />
                                </div>
                                <div>
                                    <USkeleton class="h-8 w-[200px]" />
                                </div>
                            </div>
                            <div class="flex items-center gap-2">
                                <USkeleton class="h-8 w-[63px]" />

                                <USkeleton class="h-8 w-[63px]" />
                            </div>
                        </div>
                    </li>
                </ul>

                <DataErrorBlock
                    v-else-if="error"
                    :title="$t('common.unable_to_load', [$t('common.request', 2)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataNoDataBlock
                    v-else-if="!requests || requests.requests.length === 0"
                    icon="i-mdi-frequently-asked-questions"
                    :message="$t('common.not_found', [$t('common.request', 2)])"
                />

                <ul v-else class="mb-6 divide-y divide-default rounded-md" role="list">
                    <RequestListEntry
                        v-for="request in requests.requests"
                        :key="request.id"
                        :request="request"
                        :can-update="canDo.update"
                        :can-delete="canDo.delete"
                        @refresh-requests="refresh()"
                    />
                </ul>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    v-if="canDo.create"
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.add')"
                    @click="() => formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UDrawer>
</template>
