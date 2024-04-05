<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { unitStatusToBGColor, unitStatuses } from '~/components/centrum/helpers';
import { useCentrumStore } from '~/store/centrum';
import { useNotificatorStore } from '~/store/notificator';
import { StatusUnit, Unit } from '~~/gen/ts/resources/centrum/units';

const props = defineProps<{
    unit: Unit;
    status?: StatusUnit;
    location?: Coordinate;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useSlideover();

const centrumStore = useCentrumStore();
const { settings } = storeToRefs(centrumStore);

const notifications = useNotificatorStore();

const status = computed<number>(() => props.status ?? props.unit?.status?.status ?? StatusUnit.UNKNOWN);

interface FormData {
    status: number;
    code?: string;
    reason?: string;
}

async function updateUnitStatus(id: string, values: FormData): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().updateUnitStatus({
            unitId: id,
            status: values.status,
            code: values.code,
            reason: values.reason,
        });
        await call;

        setFieldValue('status', values.status.valueOf());

        notifications.add({
            title: { key: 'notifications.centrum.sidebar.unit_status_updated.title', parameters: {} },
            description: { key: 'notifications.centrum.sidebar.unit_status_updated.content', parameters: {} },
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

watch(props, () => {
    if (props.status) {
        setFieldValue('status', props.status.valueOf());
    }
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await updateUnitStatus(props.unit.id, values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

function updateReasonField(value: string): void {
    if (value.length === 0) {
        return;
    }

    setFieldValue('reason', value);
}
</script>

<template>
    <USlideover>
        <UCard
            class="flex flex-col flex-1"
            :ui="{ body: { base: 'flex-1' }, ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.centrum.update_unit_status.title') }}: {{ unit.name }} ({{ unit.initials }})
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <UForm :state="{}" class="flex h-full flex-col divide-y divide-gray-200" @submit.prevent="onSubmitThrottle">
                    <div class="flex flex-1 flex-col justify-between">
                        <div class="divide-y divide-gray-200 px-2 sm:px-6">
                            <div>
                                <dl class="divide-y divide-neutral/10 border-b border-neutral/10">
                                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                        <dt class="text-sm font-medium leading-6">
                                            <label for="status" class="block text-sm font-medium leading-6">
                                                {{ $t('common.status') }}
                                            </label>
                                        </dt>
                                        <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                            <VeeField
                                                v-slot="{ field }"
                                                name="status"
                                                as="div"
                                                class="grid w-full grid-cols-2 gap-0.5"
                                                :placeholder="$t('common.status')"
                                                :label="$t('common.status')"
                                            >
                                                <UButton
                                                    v-for="item in unitStatuses"
                                                    :key="item.name"
                                                    class="group my-0.5 flex w-full flex-col items-center rounded-md p-1.5 text-xs font-medium hover:bg-primary-100/10 hover:transition-all"
                                                    :class="[
                                                        field.value == item.status
                                                            ? 'disabled bg-base-500 hover:bg-base-400'
                                                            : item.status
                                                              ? unitStatusToBGColor(item.status)
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
                                                                ? $t(`enums.centrum.StatusUnit.${StatusUnit[item.status ?? 0]}`)
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
                                        <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                            <VeeField
                                                type="text"
                                                name="code"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                        <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                            <VeeField
                                                type="text"
                                                name="reason"
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                :placeholder="$t('common.reason')"
                                                :label="$t('common.reason')"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                            />
                                            <VeeErrorMessage name="reason" as="p" class="mt-2 text-sm text-error-400" />
                                        </dd>
                                    </div>
                                    <div
                                        v-if="settings?.predefinedStatus && settings?.predefinedStatus.unitStatus.length > 0"
                                        class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                                    >
                                        <dt class="text-sm font-medium leading-6">
                                            <label for="unitStatus" class="block text-sm font-medium leading-6">
                                                {{ $t('common.predefined', 2) }}
                                                {{ $t('common.reason', 2) }}
                                            </label>
                                        </dt>
                                        <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                            <select
                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                @focusin="focusTablet(true)"
                                                @focusout="focusTablet(false)"
                                                @change="updateReasonField(($event.target as HTMLSelectElement).value)"
                                            >
                                                <option value=""></option>
                                                <option
                                                    v-for="(preStatus, idx) in settings?.predefinedStatus.unitStatus"
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
                        </div>
                    </div>
                </UForm>
            </div>

            <template #footer>
                <UButton
                    class="inline-flex w-full items-center rounded-l-md px-3 py-2 text-sm font-semibold shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 sm:col-start-2"
                    :disabled="!meta.valid || !canSubmit"
                    :loading="!canSubmit"
                    @click="onSubmitThrottle"
                >
                    {{ $t('common.update') }}
                </UButton>
                <UButton @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </USlideover>
</template>
