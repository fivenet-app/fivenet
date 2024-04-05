<script lang="ts" setup>
import { CheckIcon, ChevronDownIcon, CloseIcon } from 'mdi-vue3';
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
import { defineRule } from 'vee-validate';
import { AccessLevel } from '~~/gen/ts/resources/documents/access';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import { useNotificatorStore } from '~/store/notificator';

const props = defineProps<{
    documentId: string;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

interface FormData {
    reason: string;
}

async function createDocumentRequest(values: FormData): Promise<void> {
    try {
        const call = $grpc.getDocStoreClient().createDocumentReq({
            documentId: props.documentId,
            requestType: DocActivityType.REQUESTED_ACCESS,
            reason: values.reason,
            data: {
                data: {
                    oneofKind: 'accessRequested',
                    accessRequested: {
                        level: selectedAccessLevel.value,
                    },
                },
            },
        });
        await call;

        notifications.add({
            title: { key: 'notifications.docstore.requests.created.title' },
            description: { key: 'notifications.docstore.requests.created.content' },
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

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        reason: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await createDocumentRequest(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const accessLevels = [AccessLevel.VIEW, AccessLevel.COMMENT, AccessLevel.STATUS, AccessLevel.ACCESS, AccessLevel.EDIT];
const selectedAccessLevel = ref<AccessLevel>(AccessLevel.VIEW);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
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
                <UForm :state="{}" @submit.prevent="onSubmitThrottle">
                    <div class="my-2 space-y-24">
                        <div class="flex-1">
                            <label for="reason" class="block text-sm font-medium leading-6">
                                {{ $t('common.reason') }}
                            </label>
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
                        </div>
                    </div>
                    <div class="my-2">
                        <div class="flex-1">
                            <label for="requestsType" class="block text-sm font-medium leading-6">
                                {{ $t('common.access') }}
                            </label>
                            <VeeField
                                type="text"
                                name="requestsType"
                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                :placeholder="$t('common.type', 2)"
                                :label="$t('common.type', 2)"
                                @focusin="focusTablet(true)"
                                @focusout="focusTablet(false)"
                            >
                                <Listbox v-model="selectedAccessLevel" as="div">
                                    <div class="relative">
                                        <ListboxButton
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        >
                                            <span class="block truncate">
                                                {{ $t(`enums.docstore.AccessLevel.${AccessLevel[selectedAccessLevel]}`) }}
                                            </span>
                                            <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                                                <ChevronDownIcon class="size-5 text-gray-400" />
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
                                                    v-for="level in accessLevels"
                                                    :key="level"
                                                    v-slot="{ active, selected }"
                                                    as="template"
                                                    :value="level"
                                                >
                                                    <li
                                                        :class="[
                                                            active ? 'bg-primary-500' : '',
                                                            'relative cursor-default select-none py-2 pl-8 pr-4',
                                                        ]"
                                                    >
                                                        <span
                                                            :class="[
                                                                selected ? 'font-semibold' : 'font-normal',
                                                                'block truncate',
                                                            ]"
                                                        >
                                                            {{ $t(`enums.docstore.AccessLevel.${AccessLevel[level]}`, 2) }}
                                                        </span>

                                                        <span
                                                            v-if="selected"
                                                            :class="[
                                                                active ? 'text-neutral' : 'text-primary-500',
                                                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                            ]"
                                                        >
                                                            <CheckIcon class="size-5" />
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
                </UForm>
            </div>

            <template #footer>
                <div class="flex items-center">
                    <UButton @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton :disabled="!meta.valid || !canSubmit" :loading="!canSubmit" @click="onSubmitThrottle">
                        {{ $t('common.request', 2) }}
                    </UButton>
                </div>
            </template>
        </UCard>
    </UModal>
</template>
