<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getDocumentsSigningClient } from '~~/gen/ts/clients';
import { SignatureBindingMode, type SignaturePolicy } from '~~/gen/ts/resources/documents/signing';
import SignaturePadDrawer from './SignaturePadDrawer.vue';

const props = defineProps<{
    documentId: number;
}>();

defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const overlay = useOverlay();

const { can } = useAuth();

const signingClient = await getDocumentsSigningClient();

const { data, status, error, refresh } = useLazyAsyncData(
    () => `approval-drawer-${props.documentId}`,
    () => getPolicy(),
);

async function getPolicy(): Promise<SignaturePolicy[]> {
    const call = signingClient.listSignaturePolicies({
        documentId: props.documentId,
    });
    const { response } = await call;

    return response.policies ?? [];
}

async function recomputeSignaturePolicyCounters() {
    try {
        const call = signingClient.recomputeSignatureStatus({
            documentId: props.documentId,
        });
        await call;

        refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
    }
}

const signaturePad = overlay.create(SignaturePadDrawer);
</script>

<template>
    <UDrawer
        :title="$t('common.sign')"
        :overlay="false"
        handle-only
        :close="{ onClick: () => $emit('close', false) }"
        :ui="{ container: 'flex-1', content: 'min-h-[50%]', title: 'flex flex-row gap-2', body: 'h-full' }"
    >
        <template #title>
            <span class="flex-1">{{ $t('common.sign') }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="$emit('close', false)" />
        </template>

        <template #body>
            <div class="mx-auto w-full max-w-[80%] min-w-3/4">
                <div class="flex flex-1 flex-col sm:flex-row sm:gap-4">
                    <DataPendingBlock
                        v-if="isRequestPending(status)"
                        :message="$t('common.loading', [$t('common.approvals')])"
                    />
                    <DataErrorBlock
                        v-else-if="error"
                        :title="$t('common.unable_to_load', [$t('common.approvals')])"
                        :error="error"
                        :retry="refresh"
                    />

                    <template v-else>
                        <div class="basis-1/4 gap-4 sm:gap-0">
                            <UCard :ui="{ header: 'p-2 sm:px-2', body: 'p-2 sm:p-2', footer: 'p-2 sm:px-2' }">
                                <template v-if="data" #header>
                                    <div class="flex flex-row flex-wrap gap-1">
                                        <UCollapsible class="flex flex-1 flex-col gap-1">
                                            <template #default>
                                                <div class="group flex flex-1 flex-row flex-wrap items-center gap-1">
                                                    <UTooltip :text="$t('common.expand_collapse')">
                                                        <UButton
                                                            class="order-last"
                                                            size="sm"
                                                            icon="i-mdi-chevron-double-down"
                                                            variant="link"
                                                            :ui="{
                                                                leadingIcon:
                                                                    'group-data-[state=open]:rotate-180 transition-transform duration-200',
                                                            }"
                                                        />
                                                    </UTooltip>

                                                    <p class="flex-1 shrink-0 text-lg font-medium">
                                                        <UTooltip :text="$t('common.approved')"
                                                            ><span class="text-success">{{ 0 }}</span></UTooltip
                                                        >/<UTooltip :text="$t('common.declined')"
                                                            ><span class="text-error">{{ 0 }}</span></UTooltip
                                                        >/<UTooltip :text="$t('enums.documents.SignatureTaskStatus.PENDING')"
                                                            ><span class="text-info">{{ 0 }}</span></UTooltip
                                                        >
                                                        {{ $t('common.signature') }}
                                                    </p>
                                                </div>
                                            </template>

                                            <template #content>
                                                <div class="flex flex-col gap-1">
                                                    <p class="text-muted-foreground text-sm">
                                                        {{ $t('common.approved') }}:
                                                        <span class="text-success">{{ 0 }}</span>
                                                    </p>
                                                    <p class="text-muted-foreground text-sm">
                                                        {{ $t('common.declined') }}:
                                                        <span class="text-error">{{ 0 }}</span>
                                                    </p>
                                                    <p class="text-muted-foreground text-sm">
                                                        {{ $t('enums.documents.SignatureTaskStatus.PENDING') }}:
                                                        <span class="text-info">{{ 0 }}</span>
                                                    </p>
                                                </div>
                                            </template>
                                        </UCollapsible>

                                        <div>
                                            <UTooltip :text="$t('common.refresh')">
                                                <UButton
                                                    icon="i-mdi-refresh"
                                                    size="sm"
                                                    variant="link"
                                                    :loading="!data && isRequestPending(status)"
                                                    @click="() => refresh()"
                                                />
                                            </UTooltip>
                                        </div>
                                    </div>
                                </template>

                                <div v-if="data" class="flex flex-col gap-2">
                                    <UBadge
                                        :label="$t(`enums.documents.SignatureBindingMode.${SignatureBindingMode[0]}`)"
                                        color="info"
                                        variant="outline"
                                    />
                                </div>

                                <div v-else>
                                    <DataNoDataBlock icon="i-mdi-approval" :type="$t('common.policy')" />
                                </div>

                                <template #footer>
                                    <UButtonGroup class="flex w-full flex-1">
                                        <UButton
                                            block
                                            :label="$t('common.policy')"
                                            :trailing-icon="data ? 'i-mdi-pencil' : 'i-mdi-plus'"
                                        />

                                        <UButton
                                            icon="i-mdi-calculator-variant"
                                            variant="outline"
                                            @click="() => recomputeSignaturePolicyCounters()"
                                        />
                                    </UButtonGroup>
                                </template>
                            </UCard>
                        </div>

                        <div class="flex flex-1 basis-3/4 flex-col gap-4">
                            {{ data }}
                        </div>
                    </template>
                </div>
            </div>
        </template>

        <template
            v-if="
                can('documents.SigningService/UpsertSignaturePolicy').value ||
                can('documents.SigningService/RevokeSignature').value
            "
            #footer
        >
            <div class="mx-auto flex w-full max-w-[80%] min-w-3/4 flex-1 flex-col gap-4">
                <UButtonGroup class="w-full flex-1">
                    <UButton
                        v-if="can('documents.SigningService/UpsertSignaturePolicy').value"
                        class="flex-1"
                        color="success"
                        icon="i-mdi-signature"
                        block
                        size="lg"
                        :label="$t('common.sign')"
                        @click="
                            signaturePad.open({
                                policyId: 0,
                                modelValue: '',
                                'onUpdate:modelValue': ($event) => console.log($event),
                            })
                        "
                    />
                </UButtonGroup>

                <UButtonGroup v-if="can('documents.SigningService/RevokeSignature').value" class="w-full flex-1">
                    <UButton
                        class="flex-1"
                        color="neutral"
                        block
                        :label="$t('common.close', 1)"
                        @click="$emit('close', false)"
                    />
                </UButtonGroup>
            </div>
        </template>
    </UDrawer>
</template>
