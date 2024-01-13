<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
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
import { useThrottleFn } from '@vueuse/core';
import { CheckIcon, ChevronDownIcon, CloseIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import DocRequestsList from '~/components/documents/DocRequestsList.vue';
import type { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';

const props = defineProps<{
    open: boolean;
    doc: DocumentShort;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const notifications = useNotificatorStore();

const requestTypes = computed(() => [
    // DocActivityType.REQUESTED_ACCESS,
    props.doc.closed ? DocActivityType.REQUESTED_OPENING : DocActivityType.REQUESTED_CLOSURE,
    DocActivityType.REQUESTED_UPDATE,
    DocActivityType.REQUESTED_OWNER_CHANGE,
    DocActivityType.REQUESTED_DELETION,
]);

const selectedRequestType = ref(requestTypes.value[0]);

interface FormData {
    reason?: string;
}

async function createDocumentAction(values: FormData): Promise<void> {
    try {
        const call = $grpc.getDocStoreClient().createDocumentReq({
            documentId: props.doc.id,
            requestType: selectedRequestType.value,
            reason: values.reason,
        });
        await call;

        notifications.dispatchNotification({
            title: { key: 'notifications.docstore.requests.created.title' },
            content: { key: 'notifications.docstore.requests.created.content' },
            type: 'success',
        });

        emits('close');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('max', max);
defineRule('min', min);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        reason: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await createDocumentAction(values).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
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
                <div class="fixed inset-0 transition-opacity bg-opacity-75 bg-base-900" />
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
                            class="relative px-4 pt-5 pb-4 overflow-hidden text-left transition-all transform rounded-lg bg-base-800 text-neutral sm:my-8 w-full sm:max-w-5xl sm:p-6 sm:min-h-[28rem]"
                        >
                            <div class="absolute right-0 top-0 pr-4 pt-4 block">
                                <button
                                    type="button"
                                    class="rounded-md bg-neutral text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                                    @click="$emit('close')"
                                >
                                    <span class="sr-only">{{ $t('common.close') }}</span>
                                    <CloseIcon class="h-5 w-5" aria-hidden="true" />
                                </button>
                            </div>

                            <DialogTitle as="h3" class="text-base font-semibold leading-6">
                                {{ $t('common.request', 2) }}
                            </DialogTitle>

                            <form v-if="doc.creatorId !== activeChar?.userId" @submit.prevent="">
                                <div class="my-2 space-y-24">
                                    <div class="flex-1 form-control">
                                        <label for="reason" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.reason') }}
                                        </label>
                                        <VeeField
                                            type="text"
                                            name="reason"
                                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.reason')"
                                            :label="$t('common.reason')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="reason" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="my-2 space-y-20">
                                    <div class="flex-1 form-control">
                                        <label for="requestsType" class="block text-sm font-medium leading-6 text-neutral">
                                            {{ $t('common.type', 2) }}
                                        </label>
                                        <VeeField
                                            type="text"
                                            name="requestsType"
                                            class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                            :placeholder="$t('common.type', 2)"
                                            :label="$t('common.type', 2)"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        >
                                            <Listbox v-model="selectedRequestType" as="div">
                                                <div class="relative">
                                                    <ListboxButton
                                                        class="block pl-3 text-left w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    >
                                                        <span class="block truncate">
                                                            {{
                                                                $t(
                                                                    `enums.docstore.DocActivityType.${DocActivityType[selectedRequestType]}`,
                                                                    2,
                                                                )
                                                            }}
                                                        </span>
                                                        <span
                                                            class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none"
                                                        >
                                                            <ChevronDownIcon class="w-5 h-5 text-gray-400" aria-hidden="true" />
                                                        </span>
                                                    </ListboxButton>

                                                    <transition
                                                        leave-active-class="transition duration-100 ease-in"
                                                        leave-from-class="opacity-100"
                                                        leave-to-class="opacity-0"
                                                    >
                                                        <ListboxOptions
                                                            class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                                        >
                                                            <ListboxOption
                                                                v-for="requestType in requestTypes"
                                                                :key="requestType"
                                                                v-slot="{ active, selected }"
                                                                as="template"
                                                                :value="requestType"
                                                            >
                                                                <li
                                                                    :class="[
                                                                        active ? 'bg-primary-500' : '',
                                                                        'text-neutral relative cursor-default select-none py-2 pl-8 pr-4',
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
                                                                                `enums.docstore.DocActivityType.${DocActivityType[requestType]}`,
                                                                                2,
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
                                                                        <CheckIcon class="w-5 h-5" aria-hidden="true" />
                                                                    </span>
                                                                </li>
                                                            </ListboxOption>
                                                        </ListboxOptions>
                                                    </transition>
                                                </div>
                                            </Listbox>
                                        </VeeField>
                                        <VeeErrorMessage name="requestsType" as="p" class="mt-2 text-sm text-error-400" />
                                    </div>
                                </div>
                                <div class="absolute bottom-0 w-full left-0 flex">
                                    <button
                                        type="button"
                                        class="flex-1 rounded-bd bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.close', 1) }}
                                    </button>
                                    <button
                                        type="button"
                                        class="flex justify-center flex-1 rounded-bd py-2.5 px-3.5 text-sm font-semibold text-neutral"
                                        :disabled="!meta.valid || !canSubmit"
                                        :class="[
                                            !meta.valid || !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                        ]"
                                        @click="onSubmitThrottle($event)"
                                    >
                                        <template v-if="!canSubmit">
                                            <LoadingIcon class="animate-spin h-5 w-5 mr-2" />
                                        </template>
                                        {{ $t('common.add') }}
                                    </button>
                                </div>
                            </form>

                            <DocRequestsList :document-id="doc.id" @refresh="$emit('refresh')" />
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
