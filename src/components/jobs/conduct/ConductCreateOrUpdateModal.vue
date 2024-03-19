<script lang="ts" setup>
import {
    Combobox,
    ComboboxButton,
    ComboboxInput,
    ComboboxOption,
    ComboboxOptions,
    Dialog,
    DialogPanel,
    DialogTitle,
    TransitionChild,
    TransitionRoot,
} from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { digits, max, min, required } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn, watchDebounced } from '@vueuse/core';
import { CheckIcon, CloseIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useJobsStore } from '~/store/jobs';
import { ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import { UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    open: boolean;
    entry?: ConductEntry;
    userId?: number;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'created', entry: ConductEntry): void;
    (e: 'update', entry: ConductEntry): void;
}>();

const { $grpc } = useNuxtApp();

interface FormData {
    targetUser?: number;
    type: ConductType;
    message: string;
    expiresAt?: string;
}

async function conductCreateOrUpdateEntry(values: FormData, id?: string): Promise<void> {
    try {
        const expiresAt = values.expiresAt ? toTimestamp(fromString(values.expiresAt)) : undefined;

        const req = {
            entry: {
                id: id ?? '0',
                job: '',
                type: values.type,
                message: values.message,
                creatorId: 1,
                targetUserId: values.targetUser!,
                expiresAt,
            },
        };

        if (id === undefined) {
            const call = $grpc.getJobsConductClient().createConductEntry(req);
            const { response } = await call;

            emit('created', response.entry!);
        } else {
            const call = $grpc.getJobsConductClient().updateConductEntry(req);
            const { response } = await call;

            emit('update', response.entry!);
        }

        emit('close');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const queryTargets = ref<string>('');

const jobsStore = useJobsStore();
const { data, refresh } = useLazyAsyncData(`jobs-colleagues-0-${queryTargets.value}`, () =>
    jobsStore.listColleagues({
        pagination: { offset: 0n },
        searchName: queryTargets.value,
        userId: props.userId,
    }),
);

watchDebounced(queryTargets, async () => refresh(), {
    debounce: 600,
    maxWait: 1400,
});

const cTypes = ref<{ status: ConductType; selected?: boolean }[]>([
    { status: ConductType.NOTE },
    { status: ConductType.NEUTRAL },
    { status: ConductType.POSITIVE },
    { status: ConductType.NEGATIVE },
    { status: ConductType.WARNING },
    { status: ConductType.SUSPENSION },
]);

const targetUser = ref<UserShort | undefined>();
watch(targetUser, () => {
    if (targetUser.value) {
        setFieldValue('targetUser', targetUser.value.userId);
    } else {
        setFieldValue('targetUser', undefined);
    }
});

defineRule('required', required);
defineRule('digits', digits);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta, setValues, setFieldValue, resetForm } = useForm<FormData>({
    validationSchema: {
        targetUser: { required: true },
        type: { required: true },
        message: { required: true, min: 3, max: 2000 },
        expiresAt: { required: false },
    },
    initialValues: {
        type: ConductType.NOTE,
    },
    validateOnMount: true,
});

watch(props, () => {
    resetForm();

    if (props.entry) {
        targetUser.value = props.entry.targetUser;
    }

    setValues({
        targetUser: props.entry?.targetUserId,
        type: props.entry?.type ?? ConductType.NOTE,
        message: props.entry?.message,
        expiresAt: props.entry?.expiresAt ? toDatetimeLocal(toDate(props.entry?.expiresAt)) : undefined,
    });
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await conductCreateOrUpdateEntry(values, props.entry?.id).finally(() =>
            useTimeoutFn(() => (canSubmit.value = true), 400),
        ),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-30" @close="$emit('close')">
            <div class="fixed inset-0" />

            <div class="fixed inset-0 overflow-hidden">
                <div class="absolute inset-0 overflow-hidden">
                    <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-6xl pl-10 sm:pl-16">
                        <TransitionChild
                            as="template"
                            enter="transform transition ease-in-out duration-100 sm:duration-200"
                            enter-from="translate-x-full"
                            enter-to="translate-x-0"
                            leave="transform transition ease-in-out duration-100 sm:duration-200"
                            leave-from="translate-x-0"
                            leave-to="translate-x-full"
                        >
                            <DialogPanel class="pointer-events-auto w-screen max-w-6xl">
                                <form
                                    class="flex h-full flex-col divide-y divide-gray-200 bg-primary-900 shadow-xl"
                                    @submit.prevent="onSubmitThrottle"
                                >
                                    <div class="h-0 flex-1 overflow-y-auto">
                                        <div class="bg-primary-700 px-4 py-6 sm:px-6">
                                            <div class="flex items-center justify-between">
                                                <DialogTitle class="text-base font-semibold leading-6 text-neutral">
                                                    {{
                                                        entry === undefined
                                                            ? $t('components.jobs.conduct.CreateOrUpdateModal.create.title')
                                                            : $t('components.jobs.conduct.CreateOrUpdateModal.update.title')
                                                    }}
                                                </DialogTitle>
                                                <div class="ml-3 flex h-7 items-center">
                                                    <button
                                                        type="button"
                                                        class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-neutral"
                                                        @click="$emit('close')"
                                                    >
                                                        <span class="sr-only">{{ $t('common.close') }}</span>
                                                        <CloseIcon class="h-5 w-5" aria-hidden="true" />
                                                    </button>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-2 sm:px-6">
                                                <div class="mt-1">
                                                    <dl class="divide-y divide-neutral/10 border-b border-neutral/10">
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                <label
                                                                    for="type"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.type') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    v-slot="{ field }"
                                                                    as="div"
                                                                    name="type"
                                                                    :placeholder="$t('common.type')"
                                                                    :label="$t('common.type')"
                                                                >
                                                                    <select
                                                                        v-bind="field"
                                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                        @focusin="focusTablet(true)"
                                                                        @focusout="focusTablet(false)"
                                                                    >
                                                                        <option
                                                                            v-for="mtype in cTypes"
                                                                            :key="mtype.status"
                                                                            :selected="mtype.selected"
                                                                            :value="mtype.status"
                                                                        >
                                                                            {{
                                                                                $t(
                                                                                    `enums.jobs.ConductType.${
                                                                                        ConductType[
                                                                                            mtype.status ?? (0 as number)
                                                                                        ]
                                                                                    }`,
                                                                                )
                                                                            }}
                                                                        </option>
                                                                    </select>
                                                                </VeeField>
                                                                <VeeErrorMessage
                                                                    name="type"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                <label
                                                                    for="targetUser"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.target') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    as="div"
                                                                    name="targetUser"
                                                                    :placeholder="$t('common.target')"
                                                                    :label="$t('common.target')"
                                                                >
                                                                    <Combobox v-model="targetUser" as="div" class="mt-2 w-full">
                                                                        <div class="relative">
                                                                            <ComboboxButton as="div">
                                                                                <ComboboxInput
                                                                                    autocomplete="off"
                                                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                                    :display-value="
                                                                                        (char: any) =>
                                                                                            char
                                                                                                ? `${char.firstname} ${char.lastname}`
                                                                                                : $t('common.na')
                                                                                    "
                                                                                    :placeholder="$t('common.target')"
                                                                                    :label="$t('common.target')"
                                                                                    @change="queryTargets = $event.target.value"
                                                                                    @focusin="focusTablet(true)"
                                                                                    @focusout="focusTablet(false)"
                                                                                />
                                                                            </ComboboxButton>

                                                                            <ComboboxOptions
                                                                                class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                                                            >
                                                                                <ComboboxOption
                                                                                    v-for="colleague in data?.colleagues"
                                                                                    :key="colleague.identifier"
                                                                                    v-slot="{ active, selected }"
                                                                                    :value="colleague"
                                                                                    as="char"
                                                                                >
                                                                                    <li
                                                                                        :class="[
                                                                                            'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                                                            active ? 'bg-primary-500' : '',
                                                                                        ]"
                                                                                    >
                                                                                        <span
                                                                                            :class="[
                                                                                                'block truncate',
                                                                                                selected && 'font-semibold',
                                                                                            ]"
                                                                                        >
                                                                                            {{ colleague.firstname }}
                                                                                            {{ colleague.lastname }}
                                                                                        </span>

                                                                                        <span
                                                                                            v-if="selected"
                                                                                            :class="[
                                                                                                active
                                                                                                    ? 'text-neutral'
                                                                                                    : 'text-primary-500',
                                                                                                'absolute inset-y-0 left-0 flex items-center pl-1.5',
                                                                                            ]"
                                                                                        >
                                                                                            <CheckIcon
                                                                                                class="h-5 w-5"
                                                                                                aria-hidden="true"
                                                                                            />
                                                                                        </span>
                                                                                    </li>
                                                                                </ComboboxOption>
                                                                            </ComboboxOptions>
                                                                        </div>
                                                                    </Combobox>
                                                                </VeeField>
                                                                <VeeErrorMessage
                                                                    name="targetUser"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                <label
                                                                    for="message"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.message') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    as="textarea"
                                                                    name="message"
                                                                    class="block h-36 w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.message')"
                                                                    :label="$t('common.message')"
                                                                    @focusin="focusTablet(true)"
                                                                    @focusout="focusTablet(false)"
                                                                />
                                                                <VeeErrorMessage
                                                                    name="message"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                <label
                                                                    for="expiresAt"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.expires_at') }}?
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    type="datetime-local"
                                                                    name="expiresAt"
                                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.expires_at')"
                                                                    :label="$t('common.expires_at')"
                                                                    @focusin="focusTablet(true)"
                                                                    @focusout="focusTablet(false)"
                                                                />
                                                                <VeeErrorMessage
                                                                    name="expiresAt"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                    </dl>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="flex flex-shrink-0 justify-end px-4 py-4">
                                        <span class="isolate inline-flex w-full rounded-md pr-4 shadow-sm">
                                            <button
                                                type="submit"
                                                class="relative flex w-full items-center rounded-l-md px-3.5 py-2.5 text-sm font-semibold text-neutral"
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
                                                {{ entry?.id === undefined ? $t('common.create') : $t('common.update') }}
                                            </button>
                                            <button
                                                type="button"
                                                class="relative -ml-px inline-flex w-full items-center rounded-r-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-200 hover:text-gray-900"
                                                @click="$emit('close')"
                                            >
                                                {{ $t('common.close', 1) }}
                                            </button>
                                        </span>
                                    </div>
                                </form>
                            </DialogPanel>
                        </TransitionChild>
                    </div>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
