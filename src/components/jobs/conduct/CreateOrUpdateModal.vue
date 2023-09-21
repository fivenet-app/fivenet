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
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { digits, max, min, required } from '@vee-validate/rules';
import { useThrottleFn, watchDebounced } from '@vueuse/core';
import { CheckIcon, CloseIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { ConductEntry, ConductType } from '~~/gen/ts/resources/jobs/conduct';
import { User, UserShort } from '~~/gen/ts/resources/users/users';

const props = defineProps<{
    open: boolean;
    entry?: ConductEntry;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'created', entry: ConductEntry): void;
}>();

const { $grpc } = useNuxtApp();

async function conductCreateEntry(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const expiresAt = values.expiresAt ? toTimestamp(new Date(values.expiresAt)) : undefined;

            const call = $grpc.getJobsClient().conductCreateEntry({
                entry: {
                    id: props.entry?.id ?? 0n,
                    job: props.entry?.job ?? '',
                    type: values.type,
                    message: values.message,
                    creatorId: 1,
                    targetUserId: values.targetUser!,
                    expiresAt: expiresAt,
                },
            });
            const { response } = await call;

            resetForm();
            emit('created', response.entry!);
            emit('close');

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const entriesChars = ref<User[]>([]);
const queryTargets = ref<string>('');

async function listColleagues(): Promise<User[]> {
    return new Promise(async (res, rej) => {
        try {
            const req = {
                pagination: {
                    offset: 0n,
                },
                searchName: queryTargets.value,
            };

            const call = $grpc.getJobsClient().colleaguesList(req);
            const { response } = await call;

            return res(response.users);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watchDebounced(queryTargets, async () => (entriesChars.value = await listColleagues()), {
    debounce: 600,
    maxWait: 1400,
});

const cTypes = ref<{ status: ConductType; selected?: boolean }[]>([
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

interface FormData {
    targetUser?: number;
    type: ConductType;
    message: string;
    expiresAt?: string;
}

const { handleSubmit, meta, setValues, setFieldValue, resetForm } = useForm<FormData>({
    validationSchema: {
        targetUser: { required: true },
        type: { required: true },
        message: { required: true, min: 3, max: 2000 },
        expiresAt: { required: false },
    },
    initialValues: {
        type: ConductType.NEUTRAL,
    },
});

watch(props, () => {
    if (props.entry) targetUser.value = props.entry.targetUser!;

    setValues({
        targetUser: props.entry?.targetUserId,
        type: props.entry?.type ?? ConductType.NEUTRAL,
        message: props.entry?.message,
        expiresAt: props.entry?.expiresAt ? toDatetimeLocal(toDate(props.entry?.expiresAt)) : undefined,
    });
});

onMounted(async () => {
    entriesChars.value = await listColleagues();
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await conductCreateEntry(values).finally(() => setTimeout(() => (canSubmit.value = true), 350)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
            <div class="fixed inset-0" />

            <div class="fixed inset-0 overflow-hidden">
                <div class="absolute inset-0 overflow-hidden">
                    <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-6xl pl-10 sm:pl-16">
                        <TransitionChild
                            as="template"
                            enter="transform transition ease-in-out duration-150 sm:duration-300"
                            enter-from="translate-x-full"
                            enter-to="translate-x-0"
                            leave="transform transition ease-in-out duration-150 sm:duration-300"
                            leave-from="translate-x-0"
                            leave-to="translate-x-full"
                        >
                            <DialogPanel class="pointer-events-auto w-screen max-w-6xl">
                                <form
                                    @submit.prevent="onSubmitThrottle"
                                    class="flex h-full flex-col divide-y divide-gray-200 bg-gray-900 shadow-xl"
                                >
                                    <div class="h-0 flex-1 overflow-y-auto">
                                        <div class="bg-primary-700 px-4 py-6 sm:px-6">
                                            <div class="flex items-center justify-between">
                                                <DialogTitle class="text-base font-semibold leading-6 text-white">
                                                    <template v-if="entry === undefined">
                                                        {{ $t('components.jobs.conduct.CreateOrUpdateModal.create_title') }}
                                                    </template>
                                                    <template v-else>
                                                        {{ $t('components.jobs.conduct.CreateOrUpdateModal.update_title') }}
                                                    </template>
                                                </DialogTitle>
                                                <div class="ml-3 flex h-7 items-center">
                                                    <button
                                                        type="button"
                                                        class="rounded-md bg-gray-100 text-gray-500 hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-white"
                                                        @click="$emit('close')"
                                                    >
                                                        <span class="sr-only">{{ $t('common.close') }}</span>
                                                        <CloseIcon class="h-6 w-6" aria-hidden="true" />
                                                    </button>
                                                </div>
                                            </div>
                                            <div class="mt-1">
                                                <p class="text-sm text-primary-300">
                                                    {{ $t('components.centrum.create_dispatch.sub_title') }}
                                                </p>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4 sm:px-6">
                                                <div class="mt-1">
                                                    <dl class="border-b border-white/10 divide-y divide-white/10">
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                <label
                                                                    for="type"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.type') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    as="div"
                                                                    name="type"
                                                                    :placeholder="$t('common.type')"
                                                                    :label="$t('common.type')"
                                                                    v-slot="{ field }"
                                                                >
                                                                    <select
                                                                        v-bind="field"
                                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    >
                                                                        <option
                                                                            v-for="mtype in cTypes"
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
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                <label
                                                                    for="targetUser"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.target') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    as="div"
                                                                    name="targetUser"
                                                                    :placeholder="$t('common.target')"
                                                                    :label="$t('common.target')"
                                                                >
                                                                    <Combobox as="div" v-model="targetUser" class="w-full mt-2">
                                                                        <div class="relative">
                                                                            <ComboboxButton as="div">
                                                                                <ComboboxInput
                                                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                                    @change="queryTargets = $event.target.value"
                                                                                    :display-value="
                                                                                        (char: any) =>
                                                                                            char
                                                                                                ? `${char.firstname} ${char.lastname}`
                                                                                                : 'N/A'
                                                                                    "
                                                                                    :placeholder="$t('common.target')"
                                                                                    :label="$t('common.target')"
                                                                                />
                                                                            </ComboboxButton>

                                                                            <ComboboxOptions
                                                                                class="absolute z-10 w-full py-1 mt-1 overflow-auto text-base rounded-md bg-base-700 max-h-44 sm:text-sm"
                                                                            >
                                                                                <ComboboxOption
                                                                                    v-for="char in entriesChars"
                                                                                    :key="char.identifier"
                                                                                    :value="char"
                                                                                    as="char"
                                                                                    v-slot="{ active, selected }"
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
                                                                                            {{ char.firstname }}
                                                                                            {{ char.lastname }}
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
                                                                                                class="w-5 h-5"
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
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                <label
                                                                    for="message"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.message') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    as="textarea"
                                                                    name="message"
                                                                    class="h-36 block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.message')"
                                                                    :label="$t('common.message')"
                                                                />
                                                                <VeeErrorMessage
                                                                    name="message"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
                                                                <label
                                                                    for="expiresAt"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.expires_at') }}?
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    type="datetime-local"
                                                                    name="expiresAt"
                                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.expires_at')"
                                                                    :label="$t('common.expires_at')"
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
                                        <span class="isolate inline-flex rounded-md shadow-sm pr-4 w-full">
                                            <button
                                                type="submit"
                                                class="w-full relative inline-flex items-center rounded-l-md py-2.5 px-3.5 text-sm font-semibold text-neutral"
                                                :disabled="!meta.valid || !canSubmit"
                                                :class="[
                                                    !meta.valid || !canSubmit
                                                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                                ]"
                                            >
                                                {{ $t('common.create') }}
                                            </button>
                                            <button
                                                type="button"
                                                class="w-full relative -ml-px inline-flex items-center rounded-r-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 hover:text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-10"
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
