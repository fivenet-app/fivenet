<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { LockQuestionIcon } from 'mdi-vue3';
import GenericAlert from '~/components/partials/elements/GenericAlert.vue';
import { useNotificatorStore } from '~/store/notificator';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';

const props = defineProps<{
    documentId: string;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

async function createDocumentRequest(reason: string): Promise<void> {
    try {
        const call = $grpc.getDocStoreClient().createDocumentReq({
            documentId: props.documentId,
            requestType: DocActivityType.REQUESTED_ACCESS,
            reason,
        });
        await call;

        notifications.dispatchNotification({
            title: { key: 'notifications.docstore.requests.created.title' },
            content: { key: 'notifications.docstore.requests.created.content' },
            type: 'success',
        });
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template v-if="can('DocStoreService.CreateDocumentReq') && attr('DocStoreService.CreateDocumentReq', 'Types', 'Access')">
    <div class="mx-auto max-w-md rounded-md">
        <GenericAlert
            title="You don't have access to this document"
            message="But you can request access to this document using the button below."
            type="info"
            :icon="markRaw(LockQuestionIcon)"
            callback-message="Request Access here"
            @clicked="createDocumentRequest('TEST123')"
        />
    </div>
</template>
