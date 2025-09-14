<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import { dispatchStatusToBadgeColor, dispatchStatuses } from '~/components/centrum/helpers';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { useCentrumStore } from '~/stores/centrum';
import { getCentrumCentrumClient } from '~~/gen/ts/clients';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    dispatchId: number;
    status?: StatusDispatch;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const centrumStore = useCentrumStore();
const { settings } = storeToRefs(centrumStore);

const notifications = useNotificationsStore();

const centrumCentrumClient = await getCentrumCentrumClient();

const schema = z.object({
    status: z.nativeEnum(StatusDispatch),
    code: z.union([z.coerce.string().min(1).max(20), z.coerce.string().length(0).optional()]),
    reason: z.union([z.coerce.string().min(3).max(255), z.coerce.string().length(0).optional()]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    status: props.status ?? StatusDispatch.NEW,
});

async function updateDispatchStatus(dispatchId: number, values: Schema): Promise<void> {
    try {
        const call = centrumCentrumClient.updateDispatchStatus({
            dispatchId: dispatchId,
            status: values.status,
            code: values.code,
            reason: values.reason,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.centrum.sidebar.dispatch_status_updated.title', parameters: {} },
            description: { key: 'notifications.centrum.sidebar.dispatch_status_updated.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await updateDispatchStatus(props.dispatchId, event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

watch(props, () => {
    state.status = props.status ?? StatusDispatch.NEW;
});

function updateReasonField(value: string): void {
    if (value.length === 0) return;

    state.reason = value;
}

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.centrum.update_dispatch_status.title')" :overlay="false">
        <template #actions>
            <IDCopyBadge :id="dispatchId" class="ml-2" prefix="DSP" />
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <dl class="divide-y divide-default">
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="status">
                                {{ $t('common.status') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField name="status">
                                <div class="grid w-full grid-cols-2 gap-0.5">
                                    <UButton
                                        v-for="item in dispatchStatuses"
                                        :key="item.name"
                                        class="group my-0.5 flex w-full flex-col items-center rounded-md p-1.5 text-xs font-medium hover:transition-all"
                                        :class="state.status == item.status ? 'bg-neutral-500 hover:bg-neutral-400' : ''"
                                        :color="dispatchStatusToBadgeColor(item.status)"
                                        :disabled="state.status == item.status"
                                        :icon="item.icon"
                                        :label="
                                            item.status
                                                ? $t(`enums.centrum.StatusDispatch.${StatusDispatch[item.status ?? 0]}`)
                                                : $t(item.name)
                                        "
                                        @click="state.status = item.status ?? StatusDispatch.NEW"
                                    />
                                </div>
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="code">
                                {{ $t('common.code') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField class="flex-1" name="code">
                                <UInput
                                    v-model="state.code"
                                    type="text"
                                    name="code"
                                    :placeholder="$t('common.code')"
                                    class="w-full"
                                />
                            </UFormField>
                        </dd>
                    </div>

                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="reason">
                                {{ $t('common.reason') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <UFormField class="flex-1" name="reason" required>
                                <UInput v-model="state.reason" type="text" :placeholder="$t('common.reason')" class="w-full" />
                            </UFormField>
                        </dd>
                    </div>

                    <div
                        v-if="settings?.predefinedStatus && settings?.predefinedStatus.dispatchStatus.length > 0"
                        class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                    >
                        <dt class="text-sm leading-6 font-medium">
                            <label class="block text-sm leading-6 font-medium" for="dispatchStatus">
                                {{ $t('common.predefined', 2) }}
                                {{ $t('common.reason', 2) }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <ClientOnly>
                                <USelectMenu
                                    name="dispatchStatus"
                                    :items="['&nbsp;', ...settings?.predefinedStatus.dispatchStatus]"
                                    :search-input="{ placeholder: $t('common.search_field') }"
                                    @update:model-value="($event) => updateReasonField($event)"
                                >
                                    <template #item-label="{ item }">
                                        {{ item !== '' ? item : '&nbsp;' }}
                                    </template>
                                </USelectMenu>
                            </ClientOnly>
                        </dd>
                    </div>
                </dl>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.update')"
                    @click="() => formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>
