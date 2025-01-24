<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { checkDocAccess } from '~/components/documents/helpers';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel, type DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import type { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ListDocumentReqsResponse } from '~~/gen/ts/services/docstore/docstore';
import DocumentRequestsListEntry from './DocumentRequestsListEntry.vue';

const props = defineProps<{
    access: DocumentAccess;
    doc: DocumentShort;
}>();

const emit = defineEmits<{
    (e: 'refresh'): void;
}>();

const { isOpen } = useModal();

const { attr, can, activeChar } = useAuth();

const notifications = useNotificatorStore();

type RequestType = { key: DocActivityType; attrKey: string };
const requestTypes = [
    { key: props.doc.closed ? DocActivityType.REQUESTED_OPENING : DocActivityType.REQUESTED_CLOSURE, attrKey: 'Closure' },
    { key: DocActivityType.REQUESTED_UPDATE, attrKey: 'Update' },
    { key: DocActivityType.REQUESTED_OWNER_CHANGE, attrKey: 'OwnerChange' },
    { key: DocActivityType.REQUESTED_DELETION, attrKey: 'Deletion' },
] as RequestType[];

const availableRequestTypes = computed<RequestType[]>(() =>
    requestTypes.filter((rt) => attr('DocStoreService.CreateDocumentReq', 'Types', rt.attrKey).value),
);

const schema = z.object({
    reason: z.string().min(3).max(255),
    requestType: z.nativeEnum(DocActivityType).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    requestType: availableRequestTypes.value[0]?.key ?? undefined,
});

const offset = ref(0);

const {
    data: requests,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`document-${props.doc.id}-requests-${offset.value}`, () => listDocumnetReqs(props.doc.id));

async function listDocumnetReqs(documentId: number): Promise<ListDocumentReqsResponse> {
    try {
        const call = getGRPCDocStoreClient().listDocumentReqs({
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
        const call = getGRPCDocStoreClient().createDocumentReq({
            documentId: props.doc.id,
            reason: values.reason,
            requestType: values.requestType,
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.docstore.requests.created.title' },
            description: { key: 'notifications.docstore.requests.created.content' },
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
        can('DocStoreService.CreateDocumentReq').value &&
        checkDocAccess(props.access, props.doc.creator, AccessLevel.VIEW),
    update:
        can('DocStoreService.CreateDocumentReq').value &&
        checkDocAccess(props.access, props.doc.creator, AccessLevel.EDIT, 'DocStoreService.CreateDocumentReq'),
    delete:
        can('DocStoreService.DeleteDocumentReq').value &&
        checkDocAccess(props.access, props.doc.creator, AccessLevel.EDIT, 'DocStoreService.DeleteDocumentReq'),
}));

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createDocumentRequest(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('common.request', 2) }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <template v-if="canDo.create">
                        <UFormGroup name="reason" :label="$t('common.reason')" required>
                            <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" />
                        </UFormGroup>

                        <div class="my-2">
                            <UFormGroup name="requestsType" :label="$t('common.type', 2)" class="flex-1">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="state.requestType"
                                        :options="availableRequestTypes"
                                        value-attribute="key"
                                        :placeholder="$t('common.type')"
                                        :searchable-placeholder="$t('common.search_field')"
                                    >
                                        <template #label>
                                            <span v-if="state.requestType" class="truncate">
                                                {{
                                                    $t(
                                                        `enums.docstore.DocActivityType.${DocActivityType[state.requestType]}`,
                                                        2,
                                                    )
                                                }}
                                            </span>
                                        </template>

                                        <template #option="{ option }">
                                            <span class="truncate">{{
                                                $t(`enums.docstore.DocActivityType.${DocActivityType[option.key]}`, 2)
                                            }}</span>
                                        </template>

                                        <template #option-empty="{ query: search }">
                                            <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                        </template>

                                        <template #empty>
                                            {{ $t('common.not_found', [$t('common.type', 2)]) }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormGroup>
                        </div>
                    </template>

                    <div>
                        <ul v-if="loading" role="list" class="mb-6 divide-y divide-gray-800 rounded-md dark:divide-gray-500">
                            <li v-for="idx in 2" :key="idx" class="flex justify-between gap-x-4 py-4">
                                <div class="flex min-w-0 gap-x-2 px-2">
                                    <div class="min-w-0 flex-auto">
                                        <p class="text-base font-semibold leading-6 text-gray-100">
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

                        <ul v-else role="list" class="mb-6 divide-y divide-gray-800 rounded-md dark:divide-gray-500">
                            <DocumentRequestsListEntry
                                v-for="request in requests.requests"
                                :key="request.id"
                                :request="request"
                                :can-update="canDo.update"
                                :can-delete="canDo.delete"
                                @refresh-requests="refresh()"
                            />
                        </ul>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton
                            v-if="canDo.create"
                            type="submit"
                            block
                            class="flex-1"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                        >
                            {{ $t('common.add') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
