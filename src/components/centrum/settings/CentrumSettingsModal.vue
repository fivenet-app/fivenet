<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { CloseIcon, GroupIcon, LoadingIcon, PlusIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { CentrumMode, Settings } from '~~/gen/ts/resources/centrum/settings';

const props = defineProps<{
    open: boolean;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

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

        emit('close');
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

watch(props, () => setSettingsValues());
watch(settings, () => setSettingsValues());

setSettingsValues();

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await updateSettings(values).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const { remove: usRemove, push: usPush, fields: usFields } = useFieldArray<string>('unitStatus');

const { remove: dspRemove, push: dspPush, fields: dspFields } = useFieldArray<string>('dispatchStatus');
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="$emit('close')">
            <TransitionChild
                as="template"
                enter="ease-out duration-300"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="ease-in duration-200"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
            </TransitionChild>

            <div class="fixed inset-0 z-30 overflow-y-auto">
                <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
                    <TransitionChild
                        as="template"
                        enter="ease-out duration-300"
                        enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                        enter-to="opacity-100 translate-y-0 sm:scale-100"
                        leave="ease-in duration-200"
                        leave-from="opacity-100 translate-y-0 sm:scale-100"
                        leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                    >
                        <DialogPanel
                            class="relative w-full transform overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:max-w-lg sm:p-6"
                        >
                            <div class="absolute right-0 top-0 block pr-4 pt-4">
                                <button
                                    type="button"
                                    class="rounded-md bg-neutral text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="$emit('close')"
                                >
                                    <span class="sr-only">{{ $t('common.close') }}</span>
                                    <CloseIcon class="h-5 w-5" aria-hidden="true" />
                                </button>
                            </div>
                            <form @submit.prevent="onSubmitThrottle">
                                <div>
                                    <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-success-100">
                                        <GroupIcon class="h-5 w-5 text-success-600" aria-hidden="true" />
                                    </div>
                                    <div class="mt-3 text-center sm:mt-5">
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6 text-neutral">
                                            {{ $t('components.centrum.units.update_settings') }}
                                        </DialogTitle>
                                        <div class="mt-2">
                                            <div class="text-sm text-gray-100">
                                                <div class="form-control flex-1">
                                                    <label
                                                        for="enabled"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.enabled') }}
                                                    </label>
                                                    <VeeField
                                                        name="enabled"
                                                        type="checkbox"
                                                        class="h-4 h-5 w-4 w-5 rounded border-gray-300 text-primary-600 focus:ring-primary-600"
                                                        :placeholder="$t('common.enabled')"
                                                        :label="$t('common.enabled')"
                                                        :value="true"
                                                        :disabled="!can('SuperUser')"
                                                    />
                                                    <VeeErrorMessage
                                                        name="enabled"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                                <div class="form-control flex-1">
                                                    <label for="mode" class="block text-sm font-medium leading-6 text-neutral">
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
                                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                            @focusin="focusTablet(true)"
                                                            @focusout="focusTablet(false)"
                                                        >
                                                            <option
                                                                v-for="mode in modes"
                                                                :key="mode.mode"
                                                                :selected="settings !== null && mode.mode === settings.mode"
                                                                :value="mode.mode"
                                                            >
                                                                {{
                                                                    $t(
                                                                        `enums.centrum.CentrumMode.${
                                                                            CentrumMode[mode.mode ?? 0]
                                                                        }`,
                                                                    )
                                                                }}
                                                            </option>
                                                        </select>
                                                    </VeeField>
                                                    <VeeErrorMessage name="mode" as="p" class="mt-2 text-sm text-error-400" />
                                                </div>
                                                <div class="form-control flex-1">
                                                    <label
                                                        for="fallbackMode"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
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
                                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                            @focusin="focusTablet(true)"
                                                            @focusout="focusTablet(false)"
                                                        >
                                                            <option
                                                                v-for="mode in modes"
                                                                :key="mode.mode"
                                                                :selected="
                                                                    settings !== null && mode.mode === settings.fallbackMode
                                                                "
                                                                :value="mode.mode"
                                                            >
                                                                {{
                                                                    $t(
                                                                        `enums.centrum.CentrumMode.${
                                                                            CentrumMode[mode.mode ?? 0]
                                                                        }`,
                                                                    )
                                                                }}
                                                            </option>
                                                        </select>
                                                    </VeeField>
                                                    <VeeErrorMessage
                                                        name="fallbackMode"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                                <!-- Predefined Unit Status Reason -->
                                                <div class="form-control flex-1">
                                                    <label
                                                        for="unitStatus"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ `${$t('common.units')} ${$t('common.status')}` }}
                                                    </label>
                                                    <div class="flex flex-col gap-1">
                                                        <div
                                                            v-for="(field, idx) in usFields"
                                                            :key="field.key"
                                                            class="flex gap-1 items-center"
                                                        >
                                                            <VeeField
                                                                :name="`unitStatus[${idx}]`"
                                                                type="text"
                                                                class="flex-1 block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                :placeholder="$t('common.reason')"
                                                                :label="$t('common.reason')"
                                                                @focusin="focusTablet(true)"
                                                                @focusout="focusTablet(false)"
                                                            />

                                                            <button
                                                                type="button"
                                                                class="rounded-full bg-primary-500 p-1.5 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                                                @click="usRemove(idx)"
                                                            >
                                                                <CloseIcon class="h-5 w-5" aria-hidden="true" />
                                                            </button>
                                                        </div>
                                                    </div>
                                                    <button
                                                        type="button"
                                                        class="mt-2 rounded-full p-1.5 text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                                        :disabled="!canSubmit || usFields.length >= 4"
                                                        :class="
                                                            !canSubmit || usFields.length >= 4
                                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500'
                                                        "
                                                        @click="usPush('')"
                                                    >
                                                        <PlusIcon class="h-5 w-5" aria-hidden="true" />
                                                    </button>
                                                </div>
                                                <!-- Predefined Dispatch Status Reason -->
                                                <div class="form-control flex-1">
                                                    <label
                                                        for="dispatchStatus"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ `${$t('common.dispatches')} ${$t('common.status')}` }}
                                                    </label>
                                                    <div class="flex flex-col gap-1">
                                                        <div
                                                            v-for="(field, idx) in dspFields"
                                                            :key="field.key"
                                                            class="flex gap-1 items-center"
                                                        >
                                                            <VeeField
                                                                :name="`dispatchStatus[${idx}]`"
                                                                type="text"
                                                                class="flex-1 block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                :placeholder="$t('common.reason')"
                                                                :label="$t('common.reason')"
                                                                @focusin="focusTablet(true)"
                                                                @focusout="focusTablet(false)"
                                                            />

                                                            <button
                                                                type="button"
                                                                class="rounded-full bg-primary-500 p-1.5 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                                                @click="dspRemove(idx)"
                                                            >
                                                                <CloseIcon class="h-5 w-5" aria-hidden="true" />
                                                            </button>
                                                        </div>
                                                    </div>
                                                    <button
                                                        type="button"
                                                        class="mt-2 rounded-full p-1.5 text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                                        :disabled="!canSubmit || dspFields.length >= 4"
                                                        :class="
                                                            !canSubmit || dspFields.length >= 4
                                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500'
                                                        "
                                                        @click="dspPush('')"
                                                    >
                                                        <PlusIcon class="h-5 w-5" aria-hidden="true" />
                                                    </button>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                                    <button
                                        v-if="can('CentrumService.UpdateSettings')"
                                        type="submit"
                                        class="flex w-full items-center rounded-md px-3 py-2 text-sm font-semibold text-neutral shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 sm:col-start-2"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                        ]"
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="mr-2 h-5 w-5 animate-spin" aria-hidden="true" />
                                        </template>
                                        {{ $t('common.update') }}
                                    </button>
                                    <button
                                        type="button"
                                        class="mt-3 inline-flex w-full items-center rounded-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-200 sm:col-start-1 sm:mt-0"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close') }}
                                    </button>
                                </div>
                            </form>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
