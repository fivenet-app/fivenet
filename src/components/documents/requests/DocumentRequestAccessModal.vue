<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import { useNotificatorStore } from '~/store/notificator';

const props = defineProps<{
    documentId: string;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

interface FormData {
    reason: string;
}

async function createDocumentRequest(values: FormData): Promise<void> {
    try {
        const call = $grpc.getDocStoreClient().createDocumentReq({
            documentId: props.documentId,
            requestType: DocActivityType.REQUESTED_ACCESS,
            reason: values.reason,
            data: {
                data: {
                    oneofKind: 'accessRequested',
                    accessRequested: {
                        level: selectedAccessLevel.value,
                    },
                },
            },
        });
        await call;

        notifications.add({
            title: { key: 'notifications.docstore.requests.created.title' },
            description: { key: 'notifications.docstore.requests.created.content' },
            type: 'success',
        });

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        reason: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await createDocumentRequest(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const accessLevels = [AccessLevel.VIEW, AccessLevel.COMMENT, AccessLevel.STATUS, AccessLevel.ACCESS, AccessLevel.EDIT];
const selectedAccessLevel = ref<AccessLevel>(AccessLevel.VIEW);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.request') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <UForm :state="{}" @submit="onSubmitThrottle">
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
                        <div class="flex-1">
                            <label for="requestsType" class="block text-sm font-medium leading-6">
                                {{ $t('common.access') }}
                            </label>
                            <VeeField
                                type="text"
                                name="requestsType"
                                :placeholder="$t('common.type', 2)"
                                :label="$t('common.type', 2)"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            >
                                <USelectMenu
                                    v-model="selectedAccessLevel"
                                    :options="accessLevels"
                                    :placeholder="
                                        selectedAccessLevel
                                            ? $t(`enums.docstore.AccessLevel.${AccessLevel[selectedAccessLevel]}`)
                                            : $t('common.na')
                                    "
                                >
                                    <template #label>
                                        <span v-if="selectedAccessLevel" class="truncate">{{
                                            $t(`enums.docstore.AccessLevel.${AccessLevel[selectedAccessLevel]}`)
                                        }}</span>
                                    </template>
                                    <template #option="{ option }">
                                        <span class="truncate">{{
                                            $t(`enums.docstore.AccessLevel.${AccessLevel[option]}`)
                                        }}</span>
                                    </template>
                                    <template #option-empty="{ query: search }">
                                        <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                    </template>
                                    <template #empty>
                                        {{ $t('common.not_found', [$t('common.attributes', 1)]) }}
                                    </template>
                                </USelectMenu>
                            </VeeField>
                            <VeeErrorMessage name="requestsType" as="p" class="mt-2 text-sm text-error-400" />
                        </div>
                    </div>
                </UForm>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton
                        block
                        class="flex-1"
                        :disabled="!meta.valid || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle"
                    >
                        {{ $t('common.request', 2) }}
                    </UButton>

                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
