<script lang="ts" setup>
import { ApprovalRuleKind, OnEditBehavior } from '~~/gen/ts/resources/documents/approval';

defineProps<{
    disabled?: boolean;
}>();

const policy = defineModel<{
    ruleKind: ApprovalRuleKind;
    onEditBehavior: OnEditBehavior;
    requiredCount?: number;
    signatureRequired: boolean;
    selfApproveAllowed: boolean;
}>({
    default: () => ({
        ruleKind: ApprovalRuleKind.REQUIRE_ALL,
        onEditBehavior: OnEditBehavior.KEEP_PROGRESS,
        requiredCount: 2,
        signatureRequired: false,
        selfApproveAllowed: false,
    }),
});

const { t } = useI18n();

const ruleKinds = computed(() => [
    {
        label: t(`enums.documents.ApprovalRuleKind.${ApprovalRuleKind[ApprovalRuleKind.REQUIRE_ALL]}`),
        value: ApprovalRuleKind.REQUIRE_ALL,
    },
    {
        label: t(`enums.documents.ApprovalRuleKind.${ApprovalRuleKind[ApprovalRuleKind.QUORUM_ANY]}`),
        value: ApprovalRuleKind.QUORUM_ANY,
    },
]);

const onEditBehaviors = computed(() => [
    {
        label: t(`enums.documents.OnEditBehavior.${OnEditBehavior[OnEditBehavior.KEEP_PROGRESS]}`),
        value: OnEditBehavior.KEEP_PROGRESS,
    },
    { label: t(`enums.documents.OnEditBehavior.${OnEditBehavior[OnEditBehavior.RESET]}`), value: OnEditBehavior.RESET },
]);
</script>

<template>
    <UFormField name="ruleKind" :label="$t('components.documents.approval.policy_form.rule_kind')">
        <USelectMenu
            v-model="policy.ruleKind"
            :items="ruleKinds"
            value-key="value"
            label-key="label"
            class="w-full"
            :disabled="disabled"
        >
            <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
        </USelectMenu>
    </UFormField>

    <UFormField
        v-if="policy.ruleKind === ApprovalRuleKind.QUORUM_ANY"
        name="requiredCount"
        :label="$t('components.documents.approval.policy_form.required_approvals')"
    >
        <UInputNumber v-model="policy.requiredCount" :min="0" :max="10" class="w-full" :disabled="disabled" />
    </UFormField>

    <UFormField name="signatureRequired" :label="$t('components.documents.approval.signature_required')">
        <USwitch v-model="policy.signatureRequired" :disabled="disabled" />
    </UFormField>

    <UFormField
        name="selfApproveAllowed"
        :label="$t('components.documents.approval.self_approve_allowed.title')"
        :description="$t('components.documents.approval.self_approve_allowed.description')"
    >
        <USwitch v-model="policy.selfApproveAllowed" :disabled="disabled" />
    </UFormField>

    <UFormField name="onEditBehavior" :label="$t('components.documents.approval.policy_form.on_edit_behavior')">
        <USelectMenu
            v-model="policy.onEditBehavior"
            :items="onEditBehaviors"
            value-key="value"
            label-key="label"
            class="w-full"
            :disabled="disabled"
        >
            <template #empty> {{ $t('common.not_found', [$t('common.type', 2)]) }} </template>
        </USelectMenu>
    </UFormField>
</template>
