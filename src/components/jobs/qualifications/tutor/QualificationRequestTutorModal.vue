<script lang="ts" setup>
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue';
import { max, min, required } from '@vee-validate/rules';
import { CheckIcon, ChevronDownIcon } from 'mdi-vue3';
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

const { isOpen } = useModal();

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
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.qualifications.request_modal.title') }}
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="$emit('close')" />
                </div>
            </template>

            <div>
                <UForm :state="{}" @submit.prevent="onSubmitThrottle">
                    <div class="flex-1">
                        <label for="status" class="block text-sm font-medium leading-6">
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
                                        class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    >
                                        <span class="block truncate">
                                            {{
                                                $t(
                                                    `enums.qualifications.RequestStatus.${RequestStatus[availableStatus.find((t) => t === field.value) ?? 0]}`,
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
                                                v-for="stat in availableStatus"
                                                :key="stat"
                                                v-slot="{ active, selected }"
                                                as="template"
                                                :value="stat"
                                            >
                                                <li
                                                    :class="[
                                                        active ? 'bg-primary-500' : '',
                                                        'relative cursor-default select-none py-2 pl-8 pr-4',
                                                    ]"
                                                >
                                                    <span
                                                        :class="[selected ? 'font-semibold' : 'font-normal', 'block truncate']"
                                                    >
                                                        {{ $t(`enums.qualifications.RequestStatus.${RequestStatus[stat]}`) }}
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
                        <VeeErrorMessage name="status" as="p" class="mt-2 text-sm text-error-400" />
                    </div>

                    <div class="flex-1">
                        <label for="approverComment" class="block text-sm font-medium leading-6">
                            {{ $t('common.message') }}
                        </label>
                        <VeeField
                            as="textarea"
                            name="approverComment"
                            class="placeholder:text-accent-200 block h-36 w-full rounded-md border-0 bg-base-700 py-1.5 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            :placeholder="$t('common.message')"
                            :label="$t('common.message')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="approverComment" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                </UForm>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton color="black" block class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton
                        block
                        class="flex-1"
                        :disabled="!meta.valid || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle"
                    >
                        {{ $t('common.submit') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </UModal>
</template>
