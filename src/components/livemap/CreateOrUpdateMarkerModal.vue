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
// eslint-disable-next-line camelcase
import { digits, max, max_value, min, min_value, required } from '@vee-validate/rules';
import { useThrottleFn } from '@vueuse/core';
import { CheckIcon, ChevronDownIcon, CloseIcon, HelpIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import type { DefineComponent } from 'vue';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import { useLivemapStore } from '~/store/livemap';
import { Marker, MarkerType } from '~~/gen/ts/resources/livemap/livemap';
import { markerIcons } from '~/components/livemap/helpers';

defineProps<{
    open: boolean;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

interface FormData {
    name: string;
    description: string;
    markerType: MarkerType.CIRCLE;
    circleRadius: number;
    circleOpacity: number;
}

const color = ref('#EE4B2B');
const selectedIcon = shallowRef<DefineComponent>(HelpIcon);

async function createMarker(values: FormData): Promise<void> {
    try {
        const marker: Marker = {
            info: {
                id: '0',
                job: '',
                jobLabel: '',
                name: values.name,
                description: values.description,
                x: location.value?.x ?? 0,
                y: location.value?.y ?? 0,
                color: color.value.replaceAll('#', ''),
            },
            type: values.markerType,
        };

        if (values.markerType === MarkerType.CIRCLE) {
            marker.data = {
                data: {
                    oneofKind: 'circle',
                    circle: {
                        radius: values.circleRadius,
                        oapcity: values.circleOpacity,
                    },
                },
            };
        } else if (values.markerType === MarkerType.ICON) {
            marker.data = {
                data: {
                    oneofKind: 'icon',
                    icon: {
                        icon: selectedIcon.value.name,
                    },
                },
            };
        }

        const call = $grpc.getLivemapperClient().createOrUpdateMarker({
            marker,
        });
        await call;

        emit('close');
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

const markerTypes = ref<{ status: MarkerType; selected?: boolean }[]>([
    { status: MarkerType.CIRCLE },
    { status: MarkerType.DOT },
    { status: MarkerType.ICON },
]);

defineRule('required', required);
defineRule('digits', digits);
defineRule('min', min);
defineRule('max', max);
defineRule('min_value', min_value);
defineRule('max_value', max_value);

const { handleSubmit, meta, values, setValues } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 255 },
        description: { required: false, min: 6, max: 512 },
        circleRadius: { required: false, min_value: 5, max_value: 250 },
        circleOpacity: { required: false, min_value: 1, max_value: 75 },
    },
    initialValues: {
        circleRadius: 50,
        circleOpacity: 3,
    },
    validateOnMount: true,
});

async function setMarker(): Promise<void> {
    setValues({
        markerType: MarkerType.CIRCLE,
        circleRadius: 50,
        circleOpacity: 3,
    });
}

onBeforeMount(async () => setMarker());

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> => await createMarker(values).finally(() => setTimeout(() => (canSubmit.value = true), 400)),
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
                    <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-2xl pl-10 sm:pl-16">
                        <TransitionChild
                            as="template"
                            enter="transform transition ease-in-out duration-150 sm:duration-300"
                            enter-from="translate-x-full"
                            enter-to="translate-x-0"
                            leave="transform transition ease-in-out duration-150 sm:duration-300"
                            leave-from="translate-x-0"
                            leave-to="translate-x-full"
                        >
                            <DialogPanel class="pointer-events-auto w-screen max-w-3xl">
                                <form
                                    class="flex h-full flex-col divide-y divide-gray-200 bg-gray-900 shadow-xl"
                                    @submit.prevent="onSubmitThrottle"
                                >
                                    <div class="h-0 flex-1 overflow-y-auto">
                                        <div class="bg-primary-700 px-4 py-6 sm:px-6">
                                            <div class="flex items-center justify-between">
                                                <DialogTitle class="text-base font-semibold leading-6 text-neutral">
                                                    {{ $t('components.livemap.create_marker.title') }}
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
                                            <div class="mt-1">
                                                <p class="text-sm text-primary-300">
                                                    {{ $t('components.livemap.create_marker.sub_title') }}
                                                </p>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-4 sm:px-6">
                                                <div class="mt-1">
                                                    <dl class="divide-y divide-neutral/10 border-b border-neutral/10">
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                <label
                                                                    for="name"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.name') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    type="text"
                                                                    name="name"
                                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.name')"
                                                                    :label="$t('common.name')"
                                                                    @focusin="focusTablet(true)"
                                                                    @focusout="focusTablet(false)"
                                                                />
                                                                <VeeErrorMessage
                                                                    name="name"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                <label
                                                                    for="description"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.description') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    type="text"
                                                                    name="description"
                                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                <label
                                                                    for="color"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.color') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <ColorInput
                                                                    v-model="color"
                                                                    disable-alpha
                                                                    format="hex"
                                                                    position="top"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                <label
                                                                    for="markerType"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.marker') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    v-slot="{ field }"
                                                                    as="div"
                                                                    name="markerType"
                                                                    :placeholder="$t('common.marker')"
                                                                    :label="$t('common.marker')"
                                                                >
                                                                    <select
                                                                        v-bind="field"
                                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    >
                                                                        <option
                                                                            v-for="mtype in markerTypes"
                                                                            :key="mtype.status"
                                                                            :selected="mtype.selected"
                                                                            :value="mtype.status"
                                                                        >
                                                                            {{
                                                                                $t(
                                                                                    `enums.livemap.MarkerType.${
                                                                                        MarkerType[
                                                                                            mtype.status ?? (0 as number)
                                                                                        ]
                                                                                    }`,
                                                                                )
                                                                            }}
                                                                        </option>
                                                                    </select>
                                                                </VeeField>
                                                                <VeeErrorMessage
                                                                    name="markerType"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div
                                                            v-if="values.markerType === MarkerType.CIRCLE"
                                                            class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                                                        >
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                <label
                                                                    for="circleRadius"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.radius') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    type="number"
                                                                    name="circleRadius"
                                                                    min="5"
                                                                    max="250"
                                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.radius')"
                                                                    :label="$t('common.radius')"
                                                                />
                                                                <VeeErrorMessage
                                                                    name="circleRadius"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div
                                                            v-else-if="values.markerType === MarkerType.ICON"
                                                            class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                                                        >
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                <label
                                                                    for="icon"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.radius') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    as="div"
                                                                    name="icon"
                                                                    :placeholder="$t('common.icon')"
                                                                    :label="$t('common.icon')"
                                                                >
                                                                    <Listbox v-model="selectedIcon" as="div" class="mt-2">
                                                                        <div class="relative">
                                                                            <ListboxButton
                                                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                            >
                                                                                <span
                                                                                    class="block inline-flex items-center truncate"
                                                                                >
                                                                                    <component
                                                                                        :is="selectedIcon"
                                                                                        v-if="selectedIcon"
                                                                                        class="mr-1 h-5 w-5"
                                                                                        :style="{ color: color }"
                                                                                    />
                                                                                    {{
                                                                                        (selectedIcon.name ?? 'N/A').replace(
                                                                                            'Icon',
                                                                                            '',
                                                                                        )
                                                                                    }}
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
                                                                                        v-for="icon in markerIcons"
                                                                                        v-slot="{ active, selected }"
                                                                                        :key="icon.name"
                                                                                        as="template"
                                                                                        :value="icon"
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
                                                                                                    'block inline-flex items-center truncate',
                                                                                                ]"
                                                                                            >
                                                                                                <component
                                                                                                    :is="icon"
                                                                                                    class="mr-1 h-5 w-5"
                                                                                                />
                                                                                                {{
                                                                                                    icon.name.replace(
                                                                                                        'Icon',
                                                                                                        '',
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
                                                                    name="icon"
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
                                                class="relative flex w-full justify-center rounded-l-md px-3.5 py-2.5 text-sm font-semibold text-neutral"
                                                :disabled="!meta.valid || !canSubmit"
                                                :class="[
                                                    !meta.valid || !canSubmit
                                                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                                ]"
                                            >
                                                <template v-if="!canSubmit">
                                                    <LoadingIcon class="mr-2 h-5 w-5 animate-spin" />
                                                </template>
                                                {{ $t('common.create') }}
                                            </button>
                                            <button
                                                type="button"
                                                class="relative -ml-px inline-flex w-full items-center rounded-r-md bg-neutral px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 hover:text-gray-900"
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
