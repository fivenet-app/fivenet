<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import { useNotificatorStore } from '~/stores/notificator';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { CheckDomainAvailabilityResponse, RegisterDomainResponse } from '~~/gen/ts/services/internet/domain';

const props = defineProps<{
    domain: { tldId: number; search: string };
    status?: CheckDomainAvailabilityResponse | undefined;
}>();

defineEmits<{
    (e: 'cancel'): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

const schema = z.object({
    transferCode: z.string().length(6).optional(),
});

type Schema = z.output<typeof schema>;

const state = reactive({
    transferCode: undefined,
});

async function registerDomain(values: Schema): Promise<RegisterDomainResponse> {
    try {
        const call = $grpc.internet.domain.registerDomain({
            tldId: props.domain.tldId,
            name: props.domain.search,
            transferCode: values.transferCode,
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;

    await registerDomain(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);
</script>

<template>
    <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
        <UFormGroup v-if="status?.transferable" :label="$t('components.internet.transfer_code')">
            <UInput v-model="state.transferCode" type="text" />
        </UFormGroup>

        <!-- Even without a transfer code "are you sure" confirmation -->
        <UFormGroup :label="$t('components.internet.pages.nic_registrar.register_form.submit')">
            <div class="flex gap-2">
                <UButton :label="$t('common.yes')" type="submit" />

                <UButton :label="$t('common.cancel')" color="red" class="flex-1" @click="$emit('cancel')" />
            </div>
        </UFormGroup>
    </UForm>
</template>
