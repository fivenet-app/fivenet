<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import DocumentRequestsList from '~/components/documents/requests/DocumentRequestsList.vue';
import type { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel, type DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { checkDocAccess } from '~/components/documents/helpers';

const props = defineProps<{
    access: DocumentAccess;
    doc: DocumentShort;
}>();

const emits = defineEmits<{
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const notifications = useNotificatorStore();

type RequestType = { key: DocActivityType; attrKey: string };
const requestTypes = [
    { key: props.doc.closed ? DocActivityType.REQUESTED_OPENING : DocActivityType.REQUESTED_CLOSURE, attrKey: 'Closure' },
    { key: DocActivityType.REQUESTED_UPDATE, attrKey: 'Update' },
    { key: DocActivityType.REQUESTED_OWNER_CHANGE, attrKey: 'OwnerChange' },
    { key: DocActivityType.REQUESTED_DELETION, attrKey: 'Deletion' },
] as RequestType[];

const availableRequestTypes = computed<RequestType[]>(() =>
    requestTypes.filter((rt) => attr('DocStoreService.CreateDocumentReq', 'Types', rt.attrKey)),
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

async function createDocumentRequest(values: Schema): Promise<void> {
    if (values.requestType === undefined) {
        return;
    }

    try {
        const call = $grpc.getDocStoreClient().createDocumentReq({
            documentId: props.doc.id,
            reason: values.reason,
            requestType: values.requestType,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.docstore.requests.created.title' },
            description: { key: 'notifications.docstore.requests.created.content' },
            type: 'success',
        });

        emits('refresh');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const canCreate =
    props.doc.creatorId !== activeChar.value?.userId &&
    availableRequestTypes.value.length > 0 &&
    can('DocStoreService.CreateDocumentReq') &&
    checkDocAccess(props.access, props.doc.creator, AccessLevel.VIEW);

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
                    <template v-if="canCreate">
                        <UFormGroup name="reason" :label="$t('common.reason')">
                            <UInput
                                v-model="state.reason"
                                type="text"
                                :placeholder="$t('common.reason')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                        </UFormGroup>

                        <div class="my-2">
                            <UFormGroup name="requestsType" :label="$t('common.type', 2)" class="flex-1">
                                <USelectMenu
                                    v-model="state.requestType"
                                    :options="availableRequestTypes"
                                    value-attribute="key"
                                    :placeholder="
                                        state.requestType
                                            ? $t(`enums.docstore.DocActivityType.${DocActivityType[state.requestType ?? 0]}`, 2)
                                            : $t('common.na')
                                    "
                                >
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
                            </UFormGroup>
                        </div>
                    </template>

                    <DocumentRequestsList :doc="doc" :access="access" @refresh="$emit('refresh')" />
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton
                            v-if="canCreate"
                            type="submit"
                            block
                            class="flex-1"
                            :disabled="!canSubmit"
                            :loading="!canSubmit"
                        >
                            {{ $t('common.add') }}
                        </UButton>

                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
