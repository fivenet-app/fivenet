<script lang="ts" setup>
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue';
// eslint-disable-next-line camelcase
import { max, max_value, min, min_value, numeric, required } from '@vee-validate/rules';
import { CheckIcon, ChevronDownIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useNotificatorStore } from '~/store/notificator';
import { ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import type { CreateOrUpdateQualificationResultResponse } from '~~/gen/ts/services/qualifications/qualifications';

const props = defineProps<{
    qualificationId: string;
    userId: number;
}>();

const emits = defineEmits<{
    (e: 'refresh'): void;
}>();

const { isOpen } = useModal();

const { $grpc } = useNuxtApp();

const notifications = useNotificatorStore();

interface FormData {
    status: ResultStatus;
    score: number;
    summary: string;
}

async function createOrUpdateQualificationRequest(
    qualificationId: string,
    userId: number,
    values: FormData,
): Promise<CreateOrUpdateQualificationResultResponse> {
    try {
        const call = $grpc.getQualificationsClient().createOrUpdateQualificationResult({
            result: {
                id: '0',
                qualificationId,
                userId,
                status: values.status,
                score: values.score,
                summary: values.summary,
                creatorId: 0,
                creatorJob: '',
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: 'success',
        });

        emits('refresh');
        isOpen.value = false;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);
defineRule('min_value', min_value);
defineRule('max_value', max_value);
defineRule('numeric', numeric);

const { handleSubmit, meta } = useForm<FormData>({
    validationSchema: {
        status: { required: true },
        score: { required: true, min_value: 0, max_value: 100, numeric: true },
        summary: { required: true, min: 3, max: 255 },
    },
    validateOnMount: true,
    initialValues: {
        status: ResultStatus.PENDING,
        score: 0,
    },
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<any> =>
        await createOrUpdateQualificationRequest(props.qualificationId, props.userId, values).finally(() =>
            useTimeoutFn(() => (canSubmit.value = true), 400),
        ),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

const availableStatus = [ResultStatus.SUCCESSFUL, ResultStatus.FAILED, ResultStatus.PENDING];
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
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                    >
                                        <span class="block truncate">
                                            {{
                                                $t(
                                                    `enums.qualifications.ResultStatus.${ResultStatus[availableStatus.find((t) => t === field.value) ?? 0]}`,
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
                                                        {{ $t(`enums.qualifications.ResultStatus.${ResultStatus[stat]}`) }}
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
                        <label for="score" class="block text-sm font-medium leading-6">
                            {{ $t('common.score') }}
                        </label>
                        <VeeField
                            type="number"
                            name="score"
                            min="0"
                            max="100"
                            class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            :placeholder="$t('common.score')"
                            :label="$t('common.score')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="score" as="p" class="mt-2 text-sm text-error-400" />
                    </div>

                    <div class="flex-1">
                        <label for="summary" class="block text-sm font-medium leading-6">
                            {{ $t('common.summary') }}
                        </label>
                        <VeeField
                            as="textarea"
                            name="summary"
                            class="block h-24 w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                            :placeholder="$t('common.summary')"
                            :label="$t('common.summary')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
                        />
                        <VeeErrorMessage name="summary" as="p" class="mt-2 text-sm text-error-400" />
                    </div>
                </UForm>
            </div>

            <template #footer>
                <UButton block @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>

                <UButton :disabled="!meta.valid || !canSubmit" :loading="!canSubmit" @click="onSubmitThrottle">
                    {{ $t('common.submit') }}
                </UButton>
            </template>
        </UCard>
    </UModal>
</template>
