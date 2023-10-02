<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { GroupIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { CentrumMode, Settings } from '~~/gen/ts/resources/dispatch/settings';

defineProps<{
    open: boolean;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const { data: settings, pending, refresh, error } = useLazyAsyncData('rector-centrum-settings', () => getCentrumSettings());

async function getCentrumSettings(): Promise<Settings> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().getSettings({});
            const { response } = await call;

            return res(response);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const modes = ref<{ mode: CentrumMode; selected?: boolean }[]>([
    { mode: CentrumMode.MANUAL },
    { mode: CentrumMode.SIMPLIFIED },
    { mode: CentrumMode.CENTRAL_COMMAND },
    { mode: CentrumMode.AUTO_ROUND_ROBIN },
]);

function setSettingsValues(): void {
    if (!settings.value) return;

    setValues({
        enabled: settings.value.enabled,
        mode: settings.value.mode,
        fallbackMode: settings.value.fallbackMode,
    });
}

watch(settings, () => {
    setSettingsValues();
});

async function createOrUpdateUnit(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCentrumClient().updateSettings({
                job: '',
                enabled: values.enabled,
                mode: values.mode,
                fallbackMode: values.fallbackMode,
            });
            await call;

            refresh();

            emit('close');

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

defineRule('required', required);

interface FormData {
    enabled: boolean;
    mode: CentrumMode;
    fallbackMode: CentrumMode;
}

const { handleSubmit, meta, setValues } = useForm<FormData>({
    validationSchema: {
        enabled: { required: false },
        mode: { required: true },
        fallbackMode: { required: true },
    },
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await createOrUpdateUnit(values).finally(() => setTimeout(() => (canSubmit.value = true), 350)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
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

            <div class="fixed inset-0 z-10 overflow-y-auto">
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
                            class="relative transform overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6"
                        >
                            <form @submit.prevent="onSubmitThrottle">
                                <div>
                                    <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-success-100">
                                        <GroupIcon class="h-6 w-6 text-success-600" aria-hidden="true" />
                                    </div>
                                    <div class="mt-3 text-center sm:mt-5">
                                        <DialogTitle as="h3" class="text-base font-semibold leading-6 text-white">
                                            {{ $t('components.centrum.units.update_settings') }}
                                        </DialogTitle>
                                        <div class="mt-2">
                                            <div class="text-sm text-gray-100">
                                                <div class="flex-1 form-control">
                                                    <label
                                                        for="enabled"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.enabled') }}
                                                    </label>
                                                    <VeeField
                                                        name="enabled"
                                                        type="checkbox"
                                                        class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-600 h-6 w-6"
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
                                                <div class="flex-1 form-control">
                                                    <label for="mode" class="block text-sm font-medium leading-6 text-neutral">
                                                        {{ $t('common.mode') }}
                                                    </label>
                                                    <VeeField
                                                        name="mode"
                                                        as="div"
                                                        :placeholder="$t('common.mode')"
                                                        :label="$t('common.mode')"
                                                        v-slot="{ field }"
                                                    >
                                                        <select
                                                            v-bind="field"
                                                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        >
                                                            <option
                                                                v-for="mode in modes"
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
                                                <div class="flex-1 form-control">
                                                    <label
                                                        for="fallbackMode"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        Fallback {{ $t('common.mode') }}
                                                    </label>
                                                    <VeeField
                                                        name="fallbackMode"
                                                        as="div"
                                                        :placeholder="$t('common.mode')"
                                                        :label="$t('common.mode')"
                                                        v-slot="{ field }"
                                                    >
                                                        <select
                                                            v-bind="field"
                                                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        >
                                                            <option
                                                                v-for="mode in modes"
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
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="mt-5 sm:mt-6 sm:grid sm:grid-flow-row-dense sm:grid-cols-2 sm:gap-3">
                                    <button
                                        v-if="can('CentrumService.UpdateSettings')"
                                        type="submit"
                                        class="flex justify-center w-full rounded-md px-3 py-2 text-sm font-semibold text-white shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 sm:col-start-2"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                        ]"
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="animate-spin h-5 w-5 mr-2" />
                                        </template>
                                        {{ $t('common.update') }}
                                    </button>
                                    <button
                                        type="button"
                                        class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:col-start-1 sm:mt-0"
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
