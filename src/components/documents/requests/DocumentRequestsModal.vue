<script lang="ts" setup>
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue';
import { max, min, required } from '@vee-validate/rules';
import { CheckIcon, ChevronDownIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { DocActivityType } from '~~/gen/ts/resources/documents/activity';
import DocumentRequestsList from '~/components/documents/requests/DocumentRequestsList.vue';
import type { DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { AccessLevel, type DocumentAccess } from '~~/gen/ts/resources/documents/access';
import { checkDocAccess } from '~/components/documents/helpers';

const props = defineProps<{
    access: DocumentAccess;
    doc: DocumentShort;
}>();

defineEmits<{
    (e: 'refresh'): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const notifications = useNotificatorStore();

type RequestType = { key: DocActivityType; attrKey: string };
const requestTypes = [
    { key: props.doc.closed ? DocActivityType.REQUESTED_OPENING : DocActivityType.REQUESTED_CLOSURE, attrKey: 'Closure' },
    { key: DocActivityType.REQUESTED_UPDATE, attrKey: 'Update' },
    { key: DocActivityType.REQUESTED_OWNER_CHANGE, attrKey: 'OwnerChange' },
    { key: DocActivityType.REQUESTED_DELETION, attrKey: 'Deletion' },
] as RequestType[];

const availableRequestTypes = computed<RequestType[]>(() =>
    requestTypes.filter((rt) => attr('DocStoreService.CreateDocumentReq', 'Types', rt.attrKey)),
);

const selectedRequestType = ref<RequestType | undefined>(availableRequestTypes.value.at(0));

interface FormData {
    reason?: string;
}

async function createDocumentRequest(values: FormData): Promise<void> {
    if (selectedRequestType.value === undefined) {
        return;
    }

    try {
        const call = $grpc.getDocStoreClient().createDocumentReq({
            documentId: props.doc.id,
            requestType: selectedRequestType.value.key,
            reason: values.reason,
        });
        await call;

        notifications.add({
            title: { key: 'notifications.docstore.requests.created.title' },
            description: { key: 'notifications.docstore.requests.created.content' },
            type: 'success',
        });
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

const canCreate =
    props.doc.creatorId !== activeChar.value?.userId &&
    availableRequestTypes.value.length > 0 &&
    can('DocStoreService.CreateDocumentReq') &&
    checkDocAccess(props.access, props.doc.creator, AccessLevel.VIEW);

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await createDocumentRequest(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('common.request', 2) }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <template v-if="canCreate">
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
                    <div class="my-2 space-y-20">
                        <div class="flex-1">
                            <label for="requestsType" class="block text-sm font-medium leading-6">
                                {{ $t('common.type', 2) }}
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
                                <Listbox v-model="selectedRequestType" as="div">
                                    <div class="relative">
                                        <ListboxButton
                                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        >
                                            <span class="block truncate">
                                                {{
                                                    $t(
                                                        `enums.docstore.DocActivityType.${DocActivityType[selectedRequestType?.key ?? 0]}`,
                                                        2,
                                                    )
                                                }}
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
                                                    v-for="requestType in availableRequestTypes"
                                                    :key="requestType.key"
                                                    v-slot="{ active, selected }"
                                                    as="template"
                                                    :value="requestType"
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
                                                            {{
                                                                $t(
                                                                    `enums.docstore.DocActivityType.${DocActivityType[requestType.key]}`,
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
                </template>

                <DocumentRequestsList :doc="doc" :access="access" @refresh="$emit('refresh')" />
            </div>

            <template #footer>
                <UButtonGroup>
                    <UButton block @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton
                        v-if="canCreate"
                        block
                        :disabled="!meta.valid || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle($event)"
                    >
                        {{ $t('common.add') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
