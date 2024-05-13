<script lang="ts" setup>
import { z } from 'zod';
import type { FormSubmitEvent } from '#ui/types';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    documentId: string;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const accessLevels = [
    { level: AccessLevel.VIEW },
    { level: AccessLevel.COMMENT },
    { level: AccessLevel.STATUS },
    { level: AccessLevel.ACCESS },
    { level: AccessLevel.EDIT },
];

const schema = z.object({
    reason: z.string().min(3).max(255),
    accessLevel: z.nativeEnum(AccessLevel),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    reason: '',
    accessLevel: AccessLevel.VIEW,
});

async function createDocumentRequest(values: Schema): Promise<void> {
    try {
        const call = $grpc.getDocStoreClient().createDocumentReq({
            documentId: props.documentId,
            requestType: DocActivityType.REQUESTED_ACCESS,
            reason: values.reason,
            data: {
                data: {
                    oneofKind: 'accessRequested',
                    accessRequested: {
                        level: values.accessLevel,
                    },
                },
            },
        });
        await call;

        notifications.add({
            title: { key: 'notifications.docstore.requests.created.title' },
            description: { key: 'notifications.docstore.requests.created.content' },
            type: NotificationType.SUCCESS,
        });

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

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
                            {{ $t('common.request') }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup name="reason" :label="$t('common.reason')">
                        <UInput
                            v-model="state.reason"
                            type="text"
                            :placeholder="$t('common.reason')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                    </UFormGroup>

                    <UFormGroup name="requestType" :label="$t('common.access')" class="flex-1">
                        <USelectMenu
                            v-model="state.accessLevel"
                            :options="accessLevels"
                            :placeholder="$t('common.access')"
                            :searchable-placeholder="$t('common.search_field')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        >
                            <template #label>
                                <span v-if="state.accessLevel" class="truncate">{{
                                    $t(`enums.docstore.AccessLevel.${AccessLevel[state.accessLevel]}`)
                                }}</span>
                            </template>
                            <template #option="{ option }">
                                <span class="truncate">{{ $t(`enums.docstore.AccessLevel.${AccessLevel[option]}`) }}</span>
                            </template>
                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                            </template>
                        </USelectMenu>
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.request', 2) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
