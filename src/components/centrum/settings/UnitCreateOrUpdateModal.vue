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
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { max, min, required } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { CheckIcon, ChevronDownIcon, CloseIcon, GroupIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import { Unit } from '~~/gen/ts/resources/centrum/units';

const props = defineProps<{
    open: boolean;
    unit?: Unit;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'created', unit: Unit): void;
    (e: 'updated', unit: Unit): void;
}>();

const { $grpc } = useNuxtApp();

interface FormData {
    name: string;
    initials: string;
    description: string;
    color: string;
    attributes: string[];
    homePostal?: string;
}

const availableAttributes: string[] = ['static', 'no_dispatch_auto_assign'];
const selectedAttributes = ref<string[]>([]);

const color = ref('#000000');

async function createOrUpdateUnit(values: FormData): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().createOrUpdateUnit({
            unit: {
                id: props.unit?.id ?? '0',
                job: '',
                name: values.name,
                initials: values.initials,
                color: color.value,
                description: values.description,
                attributes: {
                    list: selectedAttributes.value,
                },
                users: [],
                homePostal: values.homePostal,
            },
        });
        const { response } = await call;

        if (props.unit?.id === undefined) {
            emit('created', response.unit!);
        } else {
            emit('updated', response.unit!);
        }

        emit('close');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

defineRule('required', required);
defineRule('min', min);
defineRule('max', max);

const { handleSubmit, meta, setFieldValue, setValues } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 24 },
        initials: { required: true, min: 2, max: 4 },
        description: { required: false, max: 255 },
        homePostal: { required: false, max: 48 },
    },
    validateOnMount: true,
});

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await createOrUpdateUnit(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
);
const onSubmitThrottle = useThrottleFn(async (e) => {
    canSubmit.value = false;
    await onSubmit(e);
}, 1000);

async function updateUnitInForm(): Promise<void> {
    if (props.unit !== undefined) {
        setValues({
            name: props.unit.name,
            initials: props.unit.initials,
            description: props.unit.description,
            homePostal: props.unit.homePostal,
        });

        color.value = props.unit.color;

        selectedAttributes.value = props.unit.attributes?.list ?? [];
    }
}

watch(props, async () => updateUnitInForm());

onBeforeMount(async () => updateUnitInForm());
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
                                            <span v-if="unit && unit?.id">
                                                {{ $t('components.centrum.units.update_unit') }}
                                            </span>
                                            <span v-else>
                                                {{ $t('components.centrum.units.create_unit') }}
                                            </span>
                                        </DialogTitle>
                                        <div class="mt-2">
                                            <div class="text-sm text-gray-100">
                                                <div class="form-control flex-1">
                                                    <label for="name" class="block text-sm font-medium leading-6 text-neutral">
                                                        {{ $t('common.name') }}
                                                    </label>
                                                    <VeeField
                                                        name="name"
                                                        type="text"
                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :placeholder="$t('common.name')"
                                                        :label="$t('common.name')"
                                                        @focusin="focusTablet(true)"
                                                        @focusout="focusTablet(false)"
                                                    />
                                                    <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
                                                </div>
                                                <div class="form-control flex-1">
                                                    <label
                                                        for="initials"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.initials') }}
                                                    </label>
                                                    <VeeField
                                                        name="initials"
                                                        type="text"
                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :placeholder="$t('common.initials')"
                                                        :label="$t('common.initials')"
                                                        @focusin="focusTablet(true)"
                                                        @focusout="focusTablet(false)"
                                                    />
                                                    <VeeErrorMessage
                                                        name="initials"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                                <div class="form-control flex-1">
                                                    <label
                                                        for="description"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.description') }}
                                                    </label>
                                                    <VeeField
                                                        name="description"
                                                        type="text"
                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :placeholder="$t('common.description')"
                                                        :label="$t('common.description')"
                                                        @focusin="focusTablet(true)"
                                                        @focusout="focusTablet(false)"
                                                    />
                                                    <VeeErrorMessage
                                                        name="description"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                                <div class="form-control flex-1">
                                                    <label
                                                        for="attributes"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ $t('common.attributes', 2) }}
                                                    </label>
                                                    <VeeField
                                                        name="attributes"
                                                        type="text"
                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :placeholder="$t('common.attributes', 2)"
                                                        :label="$t('common.attributes', 2)"
                                                        @focusin="focusTablet(true)"
                                                        @focusout="focusTablet(false)"
                                                    >
                                                        <Listbox v-model="selectedAttributes" as="div" nullable multiple>
                                                            <div class="relative">
                                                                <ListboxButton
                                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                >
                                                                    <span class="block truncate">
                                                                        <template v-if="selectedAttributes.length > 0">
                                                                            <span
                                                                                v-for="attr in selectedAttributes"
                                                                                :key="attr"
                                                                                class="mr-1"
                                                                            >
                                                                                {{
                                                                                    $t(
                                                                                        `components.centrum.units.attributes.${attr}`,
                                                                                    )
                                                                                }}
                                                                            </span>
                                                                        </template>
                                                                        <template v-else>
                                                                            {{ $t('common.none_selected') }}
                                                                        </template>
                                                                    </span>
                                                                    <span
                                                                        class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2"
                                                                    >
                                                                        <ChevronDownIcon
                                                                            class="h-5 w-5 text-gray-400"
                                                                            aria-hidden="true"
                                                                        />
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
                                                                            v-for="attr in availableAttributes"
                                                                            :key="attr"
                                                                            v-slot="{ active, selected }"
                                                                            as="template"
                                                                            :value="attr"
                                                                        >
                                                                            <li
                                                                                :class="[
                                                                                    active ? 'bg-primary-500' : '',
                                                                                    'relative cursor-default select-none py-2 pl-8 pr-4 text-neutral',
                                                                                ]"
                                                                            >
                                                                                <span
                                                                                    :class="[
                                                                                        selected
                                                                                            ? 'font-semibold'
                                                                                            : 'font-normal',
                                                                                        'block truncate',
                                                                                    ]"
                                                                                >
                                                                                    {{
                                                                                        $t(
                                                                                            `components.centrum.units.attributes.${attr}`,
                                                                                        )
                                                                                    }}
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
                                                                        </ListboxOption>
                                                                    </ListboxOptions>
                                                                </transition>
                                                            </div>
                                                        </Listbox>
                                                    </VeeField>
                                                    <VeeErrorMessage
                                                        name="attributes"
                                                        as="p"
                                                        class="mt-2 text-sm text-error-400"
                                                    />
                                                </div>
                                                <div class="form-control flex-1">
                                                    <label for="color" class="block text-sm font-medium leading-6 text-neutral">
                                                        {{ $t('common.color') }}
                                                    </label>
                                                    <ColorInput
                                                        v-model="color"
                                                        disable-alpha
                                                        format="hex"
                                                        position="top"
                                                        @change="setFieldValue('color', $event)"
                                                    />
                                                </div>
                                                <div class="form-control flex-1">
                                                    <label
                                                        for="homePostal"
                                                        class="block text-sm font-medium leading-6 text-neutral"
                                                    >
                                                        {{ `${$t('common.department')} ${$t('common.postal_code')}` }}
                                                    </label>
                                                    <VeeField
                                                        name="homePostal"
                                                        type="text"
                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                        :placeholder="`${$t('common.department')} ${$t('common.postal_code')}`"
                                                        :label="`${$t('common.department')} ${$t('common.postal_code')}`"
                                                        @focusin="focusTablet(true)"
                                                        @focusout="focusTablet(false)"
                                                    />
                                                    <VeeErrorMessage
                                                        name="homePostal"
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
                                        type="submit"
                                        class="inline-flex w-full justify-center rounded-md px-3 py-2 text-sm font-semibold text-neutral shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 sm:col-start-2"
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
                                        <span v-if="unit && unit?.id">
                                            {{ $t('components.centrum.units.update_unit') }}
                                        </span>
                                        <span v-else>
                                            {{ $t('components.centrum.units.create_unit') }}
                                        </span>
                                    </button>
                                    <button
                                        type="button"
                                        class="mt-3 inline-flex w-full justify-center rounded-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-200 sm:col-start-1 sm:mt-0"
                                        @click="$emit('close')"
                                    >
                                        {{ $t('common.cancel') }}
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
