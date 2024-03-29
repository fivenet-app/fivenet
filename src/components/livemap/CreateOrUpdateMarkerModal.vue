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
// eslint-disable-next-line camelcase
import { digits, max, max_value, min, min_value, required } from '@vee-validate/rules';
import { useThrottleFn, useTimeoutFn } from '@vueuse/core';
import { CheckIcon, ChevronDownIcon, CloseIcon, HelpIcon, LoadingIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import type { DefineComponent } from 'vue';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import { useLivemapStore } from '~/store/livemap';
import { type MarkerMarker, MarkerType } from '~~/gen/ts/resources/livemap/livemap';
import { markerIcons } from '~/components/livemap/helpers';

const props = defineProps<{
    open: boolean;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);
const { addOrpdateMarkerMarker } = livemapStore;

interface FormData {
    name: string;
    description: string;
    expiresAt?: string;
    markerType: MarkerType.CIRCLE;
    circleRadius: number;
    circleOpacity?: number;
}

const color = ref('#ee4b2b');
const queryIcon = ref<string>('');
const selectedIcon = shallowRef<DefineComponent>(HelpIcon);

async function createMarker(values: FormData): Promise<void> {
    const expiresAt = values.expiresAt && values.expiresAt !== '' ? toTimestamp(fromString(values.expiresAt)) : undefined;

    try {
        const marker: MarkerMarker = {
            info: {
                id: '0',
                job: '',
                jobLabel: '',
                name: values.name,
                description: values.description,
                x: location.value?.x ?? 0,
                y: location.value?.y ?? 0,
                color: color.value,
            },
            expiresAt,
            type: values.markerType,
        };

        if (values.markerType === MarkerType.CIRCLE) {
            marker.data = {
                data: {
                    oneofKind: 'circle',
                    circle: {
                        radius: values.circleRadius,
                        oapcity: values.circleOpacity ?? 3,
                    },
                },
            };
        } else if (values.markerType === MarkerType.ICON) {
            marker.data = {
                data: {
                    oneofKind: 'icon',
                    icon: {
                        icon: selectedIcon.value.name ?? '',
                    },
                },
            };
        }

        const call = $grpc.getLivemapperClient().createOrUpdateMarker({
            marker,
        });
        const { response } = await call;

        if (response.marker !== undefined) {
            addOrpdateMarkerMarker(response.marker);
        }

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

const defaultExpiresAt = ref<Date>(new Date());
defaultExpiresAt.value.setTime(defaultExpiresAt.value.getTime() + 1 * 60 * 60 * 1000);

const { handleSubmit, meta, values, setValues } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 255 },
        description: { required: false, min: 6, max: 512 },
        circleRadius: { required: false, min_value: 5, max_value: 250 },
        circleOpacity: { required: false, min_value: 1, max_value: 75 },
        expiresAt: { required: false },
    },
    initialValues: {
        expiresAt: toDatetimeLocal(defaultExpiresAt.value),
        circleRadius: 50,
        circleOpacity: 15,
    },
    validateOnMount: true,
});

async function setMarker(): Promise<void> {
    defaultExpiresAt.value.setTime(defaultExpiresAt.value.getTime() + 1 * 60 * 60 * 1000);

    setValues({
        expiresAt: toDatetimeLocal(defaultExpiresAt.value),
        markerType: MarkerType.CIRCLE,
        circleRadius: 50,
        circleOpacity: 3,
    });
}

onBeforeMount(async () => setMarker());
watch(props, () => setMarker());

const canSubmit = ref(true);
const onSubmit = handleSubmit(
    async (values): Promise<void> =>
        await createMarker(values).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400)),
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
                            enter="transform transition ease-in-out duration-100 sm:duration-200"
                            enter-from="translate-x-full"
                            enter-to="translate-x-0"
                            leave="transform transition ease-in-out duration-100 sm:duration-200"
                            leave-from="translate-x-0"
                            leave-to="translate-x-full"
                        >
                            <DialogPanel class="pointer-events-auto w-screen max-w-3xl">
                                <form
                                    class="flex h-full flex-col divide-y divide-gray-200 bg-primary-900 shadow-xl"
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
                                                        <CloseIcon class="size-5" aria-hidden="true" />
                                                    </button>
                                                </div>
                                            </div>
                                            <div class="mt-1">
                                                <p class="text-sm text-primary-300">
                                                    {{ $t('components.livemap.create_marker.subtitle') }}
                                                </p>
                                            </div>
                                        </div>
                                        <div class="flex flex-1 flex-col justify-between">
                                            <div class="divide-y divide-gray-200 px-2 sm:px-6">
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
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    type="text"
                                                                    name="name"
                                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    type="text"
                                                                    name="description"
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
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-neutral">
                                                                <label
                                                                    for="expiresAt"
                                                                    class="block text-sm font-medium leading-6 text-neutral"
                                                                >
                                                                    {{ $t('common.expires_at') }}
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
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
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
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
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
                                                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                        @focusin="focusTablet(true)"
                                                                        @focusout="focusTablet(false)"
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
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    type="number"
                                                                    name="circleRadius"
                                                                    min="5"
                                                                    max="250"
                                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                                                    {{ $t('common.icon') }}
                                                                </label>
                                                            </dt>
                                                            <dd
                                                                class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0"
                                                            >
                                                                <VeeField
                                                                    as="div"
                                                                    name="icon"
                                                                    :placeholder="$t('common.icon')"
                                                                    :label="$t('common.icon')"
                                                                >
                                                                    <Combobox v-model="selectedIcon" as="div" class="mt-2">
                                                                        <div class="relative">
                                                                            <ComboboxButton
                                                                                as="div"
                                                                                class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                            >
                                                                                <ComboboxInput
                                                                                    autocomplete="off"
                                                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 text-neutral placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                                    :placeholder="$t('common.icon')"
                                                                                    @change="queryIcon = $event.target.value"
                                                                                    @focusin="focusTablet(true)"
                                                                                    @focusout="focusTablet(false)"
                                                                                />

                                                                                <span
                                                                                    class="pointer-events-none absolute inset-y-0 right-0 flex pr-2 pt-4"
                                                                                >
                                                                                    <ChevronDownIcon
                                                                                        class="size-5 text-gray-400"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                </span>

                                                                                <span
                                                                                    class="mt-1 inline-flex items-center truncate"
                                                                                >
                                                                                    <component
                                                                                        :is="selectedIcon"
                                                                                        v-if="selectedIcon"
                                                                                        class="mr-1 size-5"
                                                                                        :style="{ color: color }"
                                                                                        aria-hidden="true"
                                                                                    />
                                                                                    {{
                                                                                        (selectedIcon.name ?? 'N/A').replace(
                                                                                            'Icon',
                                                                                            '',
                                                                                        )
                                                                                    }}
                                                                                </span>
                                                                            </ComboboxButton>

                                                                            <transition
                                                                                leave-active-class="transition duration-100 ease-in"
                                                                                leave-from-class="opacity-100"
                                                                                leave-to-class="opacity-0"
                                                                            >
                                                                                <ComboboxOptions
                                                                                    class="absolute z-10 mt-1 max-h-44 w-full overflow-auto rounded-md bg-base-700 py-1 text-base sm:text-sm"
                                                                                >
                                                                                    <ComboboxOption
                                                                                        v-for="icon in markerIcons.filter(
                                                                                            (icon) =>
                                                                                                icon?.name?.includes(queryIcon),
                                                                                        )"
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
                                                                                                    'inline-flex items-center truncate',
                                                                                                ]"
                                                                                            >
                                                                                                <component
                                                                                                    :is="icon"
                                                                                                    class="mr-1 size-5"
                                                                                                    aria-hidden="true"
                                                                                                />
                                                                                                {{
                                                                                                    icon?.name?.replace(
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
                                                                                                    class="size-5"
                                                                                                    aria-hidden="true"
                                                                                                />
                                                                                            </span>
                                                                                        </li>
                                                                                    </ComboboxOption>
                                                                                </ComboboxOptions>
                                                                            </transition>
                                                                        </div>
                                                                    </Combobox>
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
                                    <div class="flex shrink-0 justify-end p-4">
                                        <span class="isolate inline-flex w-full rounded-md pr-4 shadow-sm">
                                            <button
                                                type="submit"
                                                class="relative flex w-full rounded-l-md px-3.5 py-2.5 text-sm font-semibold text-neutral"
                                                :disabled="!meta.valid || !canSubmit"
                                                :class="[
                                                    !meta.valid || !canSubmit
                                                        ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                        : 'bg-primary-500 hover:bg-primary-400 focus-visible:outline-primary-500',
                                                ]"
                                            >
                                                <template v-if="!canSubmit">
                                                    <LoadingIcon class="mr-2 size-5 animate-spin" aria-hidden="true" />
                                                </template>
                                                {{ $t('common.create') }}
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
