<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { getDocumentsApprovalClient } from '~~/gen/ts/clients';
import { ApprovalRuleKind, OnEditBehavior, type ApprovalPolicy } from '~~/gen/ts/resources/documents/approval';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import PolicyEditor from './PolicyEditor.vue';

const props = defineProps<{
    documentId: number;
}>();

const emits = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const policy = defineModel<ApprovalPolicy | undefined>();

const notifications = useNotificationsStore();

const approvalClient = await getDocumentsApprovalClient();

const schema = z.object({
    ruleKind: z.enum(ApprovalRuleKind).default(ApprovalRuleKind.REQUIRE_ALL),
    onEditBehavior: z.enum(OnEditBehavior).default(OnEditBehavior.KEEP_PROGRESS),
    requiredCount: z.number().min(1).max(10).default(2),
    signatureRequired: z.boolean().default(false),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    ruleKind: ApprovalRuleKind.REQUIRE_ALL,
    onEditBehavior: OnEditBehavior.KEEP_PROGRESS,
    requiredCount: 2,
    signatureRequired: false,
});

function setFromProps(): void {
    if (!policy.value) {
        state.ruleKind = ApprovalRuleKind.REQUIRE_ALL;
        state.onEditBehavior = OnEditBehavior.KEEP_PROGRESS;
        state.requiredCount = 2;
        state.signatureRequired = false;
        return;
    }

    state.ruleKind = policy.value.ruleKind;
    state.onEditBehavior = policy.value.onEditBehavior;
    state.requiredCount =
        policy.value.requiredCount === undefined || policy.value.requiredCount < 0 ? 1 : policy.value.requiredCount;
    state.signatureRequired = policy.value.signatureRequired;
}

setFromProps();
watch(policy, () => setFromProps());

async function upsertPolicy(values: Schema): Promise<void> {
    const call = approvalClient.upsertApprovalPolicy({
        policy: {
            documentId: props.documentId,
            ruleKind: values.ruleKind,
            onEditBehavior: values.onEditBehavior,
            requiredCount: values.ruleKind === ApprovalRuleKind.QUORUM_ANY ? values.requiredCount : undefined,
            signatureRequired: values.signatureRequired,

            assignedCount: 0,
            approvedCount: 0,
            declinedCount: 0,
            pendingCount: 0,
            anyDeclined: false,
        },
    });
    const { response } = await call;

    policy.value = response.policy;

    emits('close', true);

    notifications.add({
        title: { key: 'notifications.action_successful.title', parameters: {} },
        description: { key: 'notifications.action_successful.content', parameters: {} },
        type: NotificationType.SUCCESS,
    });
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await upsertPolicy(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UDrawer
        :title="$t('common.approve')"
        :overlay="false"
        handle-only
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{ container: 'flex-1', content: 'min-h-[50%]', title: 'flex flex-row gap-2', body: 'h-full' }"
    >
        <template #title>
            <div class="inline-flex flex-1 items-center gap-1">
                <span>{{ $t('common.approve') }}</span>
                <UIcon name="i-mdi-slash-forward" />
                <span>{{ $t('common.policy') }}</span>
            </div>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="mx-auto w-full max-w-[80%] min-w-3/4">
                <UCard :ui="{ body: 'p-4 sm:p-4', footer: 'p-4 sm:px-4' }">
                    <template #header>
                        <div class="flex items-center justify-between gap-1">
                            <h3 class="flex-1 text-base leading-6 font-semibold">
                                {{ $t('components.documents.approval.policy_form.title') }}
                            </h3>
                        </div>
                    </template>

                    <UForm ref="formRef" :schema="schema" :state="state" class="flex flex-col gap-2" @submit="onSubmitThrottle">
                        <PolicyEditor v-model="state" />
                    </UForm>

                    <template #footer>
                        <UButtonGroup class="w-full flex-1">
                            <UButton
                                :disabled="!canSubmit"
                                block
                                class="w-full"
                                :label="$t('common.save')"
                                @click="formRef?.submit()"
                            />
                        </UButtonGroup>
                    </template>
                </UCard>
            </div>
        </template>

        <template #footer>
            <div class="mx-auto flex w-full max-w-[80%] min-w-3/4 flex-1 flex-col">
                <UButtonGroup class="w-full flex-1">
                    <UButton
                        color="neutral"
                        variant="subtle"
                        icon="i-mdi-arrow-back"
                        block
                        :label="$t('common.back')"
                        @click="() => $emit('close', false)"
                    />
                </UButtonGroup>
            </div>
        </template>
    </UDrawer>
</template>
