<script lang="ts" setup>
import { GetTemplateRequest } from '@fivenet/gen/services/docstore/docstore_pb';
import { RpcError } from 'grpc-web';
import { useAuthStore } from '~/store/auth';
import { useClipboardStore } from '~/store/clipboard';

const { $grpc } = useNuxtApp();
const authStore = useAuthStore();
const clipboardStore = useClipboardStore();

const activeChar = computed(() => authStore.getActiveChar);

const props = defineProps({
    templateId: {
        type: Number,
        required: true,
    },
});

const { data: template, pending, refresh, error } = useLazyAsyncData(`documents-templates-${props.templateId}`, () => getTemplate());

async function getTemplate(): Promise<{ title: string; content: string; }> {
    return new Promise(async (res, rej) => {
        try {
            const req = new GetTemplateRequest();
            req.setTemplateId(props.templateId);
            req.setRender(true);

            const data = clipboardStore.getTemplateData();
            data.setActivechar(activeChar.value!);
            req.setData(JSON.stringify(data.toObject()));

            const resp = await $grpc.getDocStoreClient().
                getTemplate(req, null);

            return res({
                title: resp.getTemplate()?.getContentTitle()!,
                content: resp.getTemplate()?.getContent()!,
            });
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej();
        }
    });
}
</script>

<template>
    <!-- TODO -->
</template>
