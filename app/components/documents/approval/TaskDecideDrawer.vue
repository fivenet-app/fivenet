<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import SignaturePad from '~/components/partials/SignaturePad.vue';
import { getDocumentsApprovalClient } from '~~/gen/ts/clients';
import { type ApprovalPolicy, ApprovalTaskStatus } from '~~/gen/ts/resources/documents/approval';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    documentId: number;
    policy?: ApprovalPolicy;
    approve: boolean;
}>();

const emits = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const notifications = useNotificationsStore();

const approvalClient = await getDocumentsApprovalClient();

const schema = z.object({
    reason: z.string().max(255).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
});

const isOpen = ref(false);

watch(isOpen, (newVal) => {
    if (!newVal) emits('close', false);
});

async function onSubmit(values: FormSubmitEvent<Schema>) {
    const payloadSVG = saveSignature();

    try {
        console.log('Deciding approval with values:', values, 'and payloadSVG:', payloadSVG);
        const call = approvalClient.decideApproval({
            documentId: props.documentId,
            newStatus: props.approve ? ApprovalTaskStatus.APPROVED : ApprovalTaskStatus.DECLINED,
            comment: values.data.reason ?? '',
            payloadSvg: payloadSVG,
            stampId: undefined,
        });
        await call;

        emits('close', true);
        isOpen.value = false;

        state.reason = '';

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
    }
}

const signatureRef = useTemplateRef('signatureRef');

function saveSignature(): string | undefined {
    const sig = signatureRef.value?.signature?.saveSignature('image/svg+xml') ?? '';
    if (sig === '') return undefined;

    // atob? Yes, because supporting FiveM's NUI CEF version 103 is fun..
    return atob(sig.replace(/^data:image\/svg\+xml;base64,/, ''));
}
</script>

<template>
    <UDrawer
        v-model:open="isOpen"
        :title="$t('common.approve')"
        handle-only
        :ui="{ container: 'flex-1', title: 'flex flex-row gap-2', body: 'h-full' }"
    >
        <slot />

        <template #title>
            <span class="flex-1">{{ $t(approve ? 'common.approve' : 'common.decline') }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="isOpen = false" />
        </template>

        <template #body>
            <div class="mx-auto w-full max-w-[80%] min-w-3/4">
                <UForm :schema="schema" :state="state" class="flex flex-1 flex-col gap-4" @submit="onSubmit">
                    <UFormField
                        v-if="approve"
                        name="stampSVG"
                        :label="$t('common.signature')"
                        :required="policy?.signatureRequired"
                    >
                        <SignaturePad ref="signatureRef" />
                    </UFormField>

                    <UFormField name="reason" :label="$t('common.reason')" required>
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" class="w-full" />
                    </UFormField>

                    <UFormField>
                        <UButton
                            type="submit"
                            :color="approve ? 'success' : 'red'"
                            block
                            size="lg"
                            :label="approve ? $t('common.approve') : $t('common.decline')"
                            :icon="approve ? 'i-mdi-check-bold' : 'i-mdi-close-bold'"
                        />
                    </UFormField>
                </UForm>
            </div>
        </template>

        <template #footer>
            <div class="mx-auto flex w-full max-w-[80%] min-w-3/4 flex-1 flex-col gap-4">
                <UButtonGroup class="w-full flex-1">
                    <UButton class="flex-1" color="neutral" block :label="$t('common.cancel')" @click="isOpen = false" />
                </UButtonGroup>
            </div>
        </template>
    </UDrawer>
</template>
