<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { dispatchStatusToBGColor, dispatchStatuses } from '~/components/centrum/helpers';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { useCentrumStore } from '~/store/centrum';
import { useNotificatorStore } from '~/store/notificator';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

const props = defineProps<{
    dispatchId: string;
    status?: StatusDispatch;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const centrumStore = useCentrumStore();
const { settings } = storeToRefs(centrumStore);

const notifications = useNotificatorStore();

const status = computed<number>(() => props.status ?? StatusDispatch.NEW);

interface FormData {
    status: number;
    code?: string;
    reason?: string;
}

async function updateDispatchStatus(dispatchId: string, values: FormData): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().updateDispatchStatus({
            dispatchId,
            status: values.status,
            code: values.code,
            reason: values.reason,
        });
        await call;

        setFieldValue('status', values.status.valueOf());

        notifications.add({
            title: { key: 'notifications.centrum.sidebar.dispatch_status_updated.title', parameters: {} },
            description: { key: 'notifications.centrum.sidebar.dispatch_status_updated.content', parameters: {} },
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

const { handleSubmit, meta, setFieldValue } = useForm<FormData>({
    validationSchema: {
        status: { required: true },
        code: { required: false },
        reason: { required: false, min: 3, max: 255 },
    },
    initialValues: {
        status: status.value,
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await updateDispatchStatus(props.dispatchId, values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

watch(props, () => {
    if (props.status) {
        setFieldValue('status', props.status.valueOf());
    }
});

function updateReasonField(value: string): void {
    setFieldValue('reason', value);
}
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="inline-flex items-center text-2xl font-semibold leading-6">
                        {{ $t('components.centrum.update_dispatch_status.title') }}:
                        <IDCopyBadge :id="dispatchId" class="ml-2" prefix="DSP" />
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <dl class="divide-neutral/10 divide-y">
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm font-medium leading-6">
                            <label for="status" class="block text-sm font-medium leading-6">
                                {{ $t('common.status') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <VeeField
                                v-slot="{ field }"
                                name="status"
                                as="div"
                                class="grid w-full grid-cols-2 gap-0.5"
                                :placeholder="$t('common.status')"
                                :label="$t('common.status')"
                            >
                                <UButton
                                    v-for="(item, idx) in dispatchStatuses"
                                    :key="item.name"
                                    class="hover:bg-primary-100/10 group my-0.5 flex w-full flex-col items-center rounded-md p-1.5 text-xs font-medium hover:transition-all"
                                    :class="[
                                        idx >= dispatchStatuses.length - 1 ? 'col-span-2' : '',
                                        field.value == item.status
                                            ? 'disabled bg-base-500 hover:bg-base-400'
                                            : item.status
                                              ? dispatchStatusToBGColor(item.status)
                                              : '',
                                        ,
                                    ]"
                                    :disabled="field.value == item.status"
                                    @click="setFieldValue('status', item.status?.valueOf() ?? 0)"
                                >
                                    <UIcon :name="item.icon" class="size-5 shrink-0" />
                                    <span class="mt-1">
                                        {{
                                            item.status
                                                ? $t(`enums.centrum.StatusDispatch.${StatusDispatch[item.status ?? 0]}`)
                                                : $t(item.name)
                                        }}
                                    </span>
                                </UButton>
                            </VeeField>
                            <VeeErrorMessage name="status" as="p" class="mt-2 text-sm text-error-400" />
                        </dd>
                    </div>
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm font-medium leading-6">
                            <label for="code" class="block text-sm font-medium leading-6">
                                {{ $t('common.code') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <VeeField
                                type="text"
                                name="code"
                                :placeholder="$t('common.code')"
                                :label="$t('common.code')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                            <VeeErrorMessage name="code" as="p" class="mt-2 text-sm text-error-400" />
                        </dd>
                    </div>
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm font-medium leading-6">
                            <label for="reason" class="block text-sm font-medium leading-6">
                                {{ $t('common.reason') }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <VeeField
                                type="text"
                                name="reason"
                                :placeholder="$t('common.reason')"
                                :label="$t('common.reason')"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            />
                            <VeeErrorMessage name="reason" as="p" class="mt-2 text-sm text-error-400" />
                        </dd>
                    </div>

                    <div
                        v-if="settings?.predefinedStatus && settings?.predefinedStatus.dispatchStatus.length > 0"
                        class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                    >
                        <dt class="text-sm font-medium leading-6">
                            <label for="dispatchStatus" class="block text-sm font-medium leading-6">
                                {{ $t('common.predefined', 2) }}
                                {{ $t('common.reason', 2) }}
                            </label>
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <select
                                name="dispatchStatus"
                                class="mt-1 block w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-1 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                                @change="updateReasonField(($event.target as HTMLSelectElement).value)"
                            >
                                <option value=""></option>
                                <option
                                    v-for="(preStatus, idx) in settings?.predefinedStatus.dispatchStatus"
                                    :key="idx"
                                    :value="preStatus"
                                >
                                    {{ preStatus }}
                                </option>
                            </select>
                        </dd>
                    </div>
                </dl>
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
                        {{ $t('common.update') }}
                    </UButton>

                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
