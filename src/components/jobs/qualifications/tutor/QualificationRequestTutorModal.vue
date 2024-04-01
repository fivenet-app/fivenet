<script lang="ts" setup>
import {
    Dialog,
    DialogPanel,
    DialogTitle,
    Listbox,
    ListboxButton,
    ListboxOption,
    ListboxOptions,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { CheckIcon, ChevronDownIcon, CloseIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import { RequestStatus, type QualificationRequest } from '~~/gen/ts/resources/qualifications/qualifications';
import type { CreateOrUpdateQualificationRequestResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = withDefaults(
    defineProps<{
        request: QualificationRequest;
        status?: RequestStatus;
    }>(),
    {
        status: RequestStatus.PENDING,
    },
);

const emits = defineEmits<{
    (e: 'close'): void;
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

interface FormData {
    status: RequestStatus;
    approverComment: string;
}

async function createOrUpdateQualificationRequest(
    qualificationId: string,
    userId: number,
    values: FormData,
): Promise<CreateOrUpdateQualificationRequestResponse> {
    try {
        const call = $grpc.getQualificationsClient().createOrUpdateQualificationRequest({
            request: {
                qualificationId,
                userId,
                status: values.status,
                approverComment: values.approverComment,
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: 'success',
        });

        emits('refresh');
        emits('close');

        return response;
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
        approverComment: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
    initialValues: {
        status: props.status,
    },
});

watch(props, () => setFieldValue('status', props.status ?? RequestStatus.PENDING));

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<any> =>
        await createOrUpdateQualificationRequest(props.request.qualificationId, props.request.userId, values).finally(() =>
            useTimeoutFn(() => (canSubmit.value = true), 400),
        ),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const availableStatus = [RequestStatus.ACCEPTED, RequestStatus.DENIED, RequestStatus.PENDING];
</script>

<template>
    <TransitionRoot as="template" :show="!!request">
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
                <div class="fixed inset-0 bg-base-900/75 transition-opacity" />
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
                            class="relative h-112 w-full overflow-hidden rounded-lg bg-base-800 px-4 pb-4 pt-5 text-left text-neutral transition-all sm:my-8 sm:max-w-2xl sm:p-6"
                        >
                            <div class="absolute right-0 top-0 block pr-4 pt-4">
                                <button
                                    type="button"
                                    class="rounded-md bg-neutral-50 text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="$emit('close')"
                                >
                                    <span class="sr-only">{{ $t('common.close') }}</span>
                                    <CloseIcon class="size-5" aria-hidden="true" />
                                </button>
                            </div>
                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                {{ $t('components.qualifications.request_modal.title') }}
                            </DialogTitle>
                            <form @submit.prevent="onSubmitThrottle">
                                <div class="my-2">
                                    <div class="flex-1">
                                        <label for="status" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.status') }}
                                        </label>
                                        <VeeField
                                            v-slot="{ field }"
                                            as="div"
                                            name="status"
                                            :placeholder="$t('common.status')"
                                            :label="$t('common.status')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        >
                                            <Listbox v-bind="field" as="div">
                                                <div class="relative">
                                                    <ListboxButton
                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    >
                                                        <span class="block truncate">
                                                            {{
                                                                $t(
                                                                    `enums.qualifications.RequestStatus.${RequestStatus[availableStatus.find((t) => t === field.value) ?? 0]}`,
                                                                )
                                                            }}
                                                        </span>
                                                        <span
                                                            class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2"
                                                        >
                                                            <ChevronDownIcon class="size-5 text-gray-400" aria-hidden="true" />
                                                        </span>
                                                    </ListboxButton>

                                                    <transition
                                                        leave-active-class="transition duration-100 ease-in"
                                                        leave-from-class="opacity-100"
                                                        leave-to-class="opacity-0"
                                                    >
                                                        <ListboxOptions
                                                            class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                                        >
                                                            <ListboxOption
                                                                v-for="stat in availableStatus"
                                                                :key="stat"
                                                                v-slot="{ active, selected }"
                                                                as="template"
                                                                :value="stat"
                                                            >
                                                                <li
                                                                    :class="[
                                                                        active ? 'bg-primary-500' : '',
                                                                        'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                                    ]"
                                                                >
                                                                    <span
                                                                        :class="[
                                                                            selected ? 'font-semibold' : 'font-normal',
                                                                            'block truncate',
                                                                        ]"
                                                                    >
                                                                        {{
                                                                            $t(
                                                                                `enums.qualifications.RequestStatus.${RequestStatus[stat]}`,
                                                                            )
                                                                        }}
                                                                    </span>

                                                                    <span
                                                                        v-if="selected"
                                                                        :class="[
                                                                            active ? 'text-neutral' : 'text-primary-500',
                                                                            'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                                        ]"
                                                                    >
                                                                        <CheckIcon class="size-5" aria-hidden="true" />
                                                                    </span>
                                                                </li>
                                                            </ListboxOption>
                                                        </ListboxOptions>
                                                    </transition>
                                                </div>
                                            </Listbox>
                                        </VeeField>
                                        <VeeErrorMessage name="status" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>

                                    <div class="flex-1">
                                        <label for="approverComment" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.message') }}
                                        </label>
                                        <VeeField
                                            as="textarea"
                                            name="approverComment"
                                            class="block h-36 w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.message')"
                                            :label="$t('common.message')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="approverComment" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="absolute bottom-0 left-0 flex w-full">
                                    <button
                                        type="button"
                                        class="flex-1 rounded-md bg-neutral-50 px-3.5 py-2.5 text-sm font-semibold text-gray-900 hover:bg-gray-200"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </button>
                                    <button
                                        type="submit"
                                        class="flex flex-1 justify-center rounded-md px-3.5 py-2.5 text-sm font-semibold text-neutral"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400',
                                        ]"
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="mr-2 size-5 animate-spin" aria-hidden="true" />
                                        </template>
                                        {{ $t('common.submit') }}
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
