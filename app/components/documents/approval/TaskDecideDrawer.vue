<script lang="ts" setup>
import { z } from 'zod';
import { getDocumentsApprovalClient } from '~~/gen/ts/clients';
import { ApprovalTaskStatus } from '~~/gen/ts/resources/documents/approval';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    documentId: number;
    policyId: number;
    approve: boolean;
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

async function decideApproval(approve: boolean) {
    try {
        const call = approvalClient.decideApproval({
            documentId: props.documentId,
            policyId: props.policyId,
            newStatus: approve ? ApprovalTaskStatus.APPROVED : ApprovalTaskStatus.DECLINED,
            comment: '',
        });
        await call;

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
                <UForm :schema="schema" :state="state" class="flex flex-1 flex-col gap-4">
                    <UFormField name="reason" :label="$t('common.reason')">
                        <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" class="w-full" />
                    </UFormField>

                    <UFormField>
                        <UButton
                            v-if="approve"
                            type="submit"
                            color="success"
                            block
                            size="lg"
                            :label="$t('common.approve')"
                            @click="decideApproval(true)"
                        />
                        <UButton
                            v-else
                            color="red"
                            icon="i-mdi-close-bold"
                            block
                            size="lg"
                            :label="$t('common.decline')"
                            @click="decideApproval(false)"
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
