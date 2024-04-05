<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { CentrumMode, Settings } from '~~/gen/ts/resources/centrum/settings';

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const { data: settings, refresh } = useLazyAsyncData('rector-centrum-settings', () => getCentrumSettings());

async function getCentrumSettings(): Promise<Settings> {
    try {
        const call = $grpc.getCentrumClient().getSettings({});
        const { response } = await call;

        return response.settings!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const modes = ref<{ mode: CentrumMode; selected?: boolean }[]>([
    { mode: CentrumMode.MANUAL },
    { mode: CentrumMode.SIMPLIFIED },
    { mode: CentrumMode.CENTRAL_COMMAND },
    { mode: CentrumMode.AUTO_ROUND_ROBIN },
]);

interface FormData {
    enabled: boolean;
    mode: CentrumMode;
    fallbackMode: CentrumMode;
    unitStatus: string[];
    dispatchStatus: string[];
}

async function updateSettings(values: FormData): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().updateSettings({
            settings: {
                job: '',
                enabled: values.enabled,
                mode: values.mode,
                fallbackMode: values.fallbackMode,
                predefinedStatus: {
                    dispatchStatus: values.dispatchStatus.filter((s) => s.trim().length > 0),
                    unitStatus: values.unitStatus.filter((s) => s.trim().length > 0),
                },
            },
        });
        await call;

        refresh();

        isOpen.value = false;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta, setValues } = useForm<FormData>({
    validationSchema: {
        enabled: { required: false },
        mode: { required: true },
        fallbackMode: { required: true },
        unitStatus: { max: 255 },
        dispatchStatus: { max: 255 },
    },
    validateOnMount: true,
});

function setSettingsValues(): void {
    if (!settings.value) {
        return;
    }

    setValues({
        enabled: settings.value.enabled,
        mode: settings.value.mode,
        fallbackMode: settings.value.fallbackMode,
        unitStatus: settings.value.predefinedStatus?.unitStatus ?? [],
        dispatchStatus: settings.value.predefinedStatus?.dispatchStatus ?? [],
    });
}

watch(settings, () => setSettingsValues());

setSettingsValues();

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await updateSettings(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const { remove: usRemove, push: usPush, fields: usFields } = useFieldArray<string>('unitStatus');

const { remove: dspRemove, push: dspPush, fields: dspFields } = useFieldArray<string>('dispatchStatus');
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.centrum.units.update_settings') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <UForm :state="{}">
                    <div class="text-sm text-gray-100">
                        <div class="flex-1">
                            <label for="enabled" class="block text-sm font-medium leading-6">
                                {{ $t('common.enabled') }}
                            </label>
                            <VeeField
                                name="enabled"
                                type="checkbox"
                                class="size-5 rounded border-gray-300 text-primary-600 focus:ring-primary-600"
                                :placeholder="$t('common.enabled')"
                                :label="$t('common.enabled')"
                                :value="true"
                                :disabled="!can('SuperUser')"
                            />
                            <VeeErrorMessage name="enabled" as="p" class="mt-2 text-sm text-error-400" />
                        </div>

                        <div class="flex-1">
                            <label for="mode" class="block text-sm font-medium leading-6">
                                {{ $t('common.mode') }}
                            </label>
                            <VeeField
                                v-slot="{ field }"
                                name="mode"
                                as="div"
                                :placeholder="$t('common.mode')"
                                :label="$t('common.mode')"
                            >
                                <select
                                    v-bind="field"
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                >
                                    <option
                                        v-for="mode in modes"
                                        :key="mode.mode"
                                        :selected="settings !== null && mode.mode === settings.mode"
                                        :value="mode.mode"
                                    >
                                        {{ $t(`enums.centrum.CentrumMode.${CentrumMode[mode.mode ?? 0]}`) }}
                                    </option>
                                </select>
                            </VeeField>
                            <VeeErrorMessage name="mode" as="p" class="mt-2 text-sm text-error-400" />
                        </div>

                        <div class="flex-1">
                            <label for="fallbackMode" class="block text-sm font-medium leading-6">
                                {{ $t('common.fallback_mode') }}
                            </label>
                            <VeeField
                                v-slot="{ field }"
                                name="fallbackMode"
                                as="div"
                                :placeholder="$t('common.mode')"
                                :label="$t('common.mode')"
                            >
                                <select
                                    v-bind="field"
                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    @focusin="focusTablet(true)"
                                    @focusout="focusTablet(false)"
                                >
                                    <option
                                        v-for="mode in modes"
                                        :key="mode.mode"
                                        :selected="settings !== null && mode.mode === settings.fallbackMode"
                                        :value="mode.mode"
                                    >
                                        {{ $t(`enums.centrum.CentrumMode.${CentrumMode[mode.mode ?? 0]}`) }}
                                    </option>
                                </select>
                            </VeeField>
                            <VeeErrorMessage name="fallbackMode" as="p" class="mt-2 text-sm text-error-400" />
                        </div>

                        <!-- Predefined Unit Status Reason -->
                        <div class="flex-1">
                            <label for="unitStatus" class="block text-sm font-medium leading-6">
                                {{ `${$t('common.units')} ${$t('common.status')}` }}
                            </label>
                            <div class="flex flex-col gap-1">
                                <div v-for="(field, idx) in usFields" :key="field.key" class="flex items-center gap-1">
                                    <VeeField
                                        :name="`unitStatus[${idx}]`"
                                        type="text"
                                        class="block w-full flex-1 rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.reason')"
                                        :label="$t('common.reason')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />

                                    <UButton :ui="{ rounded: 'rounded-full' }" icon="i-mdi-close" @click="usRemove(idx)" />
                                </div>
                            </div>
                            <UButton
                                :ui="{ rounded: 'rounded-full' }"
                                icon="i-mdi-plus"
                                :disabled="!canSubmit || usFields.length >= 4"
                                @click="usPush('')"
                            />
                        </div>

                        <!-- Predefined Dispatch Status Reason -->
                        <div class="flex-1">
                            <label for="dispatchStatus" class="block text-sm font-medium leading-6">
                                {{ `${$t('common.dispatches')} ${$t('common.status')}` }}
                            </label>
                            <div class="flex flex-col gap-1">
                                <div v-for="(field, idx) in dspFields" :key="field.key" class="flex items-center gap-1">
                                    <VeeField
                                        :name="`dispatchStatus[${idx}]`"
                                        type="text"
                                        class="block w-full flex-1 rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.reason')"
                                        :label="$t('common.reason')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />

                                    <UButton :ui="{ rounded: 'rounded-full' }" icon="i-mdi-close" @click="dspRemove(idx)" />
                                </div>
                            </div>
                            <UButton
                                :ui="{ rounded: 'rounded-full' }"
                                icon="i-mdi-plus"
                                :disabled="!canSubmit || dspFields.length >= 4"
                                @click="dspPush('')"
                            />
                        </div>
                    </div>
                </UForm>
            </div>

            <template #footer>
                <div class="gap-2 sm:flex">
                    <UButton class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton
                        v-if="can('CentrumService.UpdateSettings')"
                        class="flex-1"
                        :disabled="!meta.valid || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle"
                    >
                        {{ $t('common.update') }}
                    </UButton>
                </div>
            </template>
        </UCard>
    </UModal>
</template>
