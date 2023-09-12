<script lang="ts" setup>
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { digits, max, min, required } from '@vee-validate/rules';
import { CloseIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import { useLivemapStore } from '~/store/livemap';
import { MARKER_TYPE, Marker } from '~~/gen/ts/resources/livemap/livemap';

defineProps<{
    open: boolean;
}>();

const emits = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const livemapStore = useLivemapStore();
const { location } = storeToRefs(livemapStore);

async function createMarker(values: FormData): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            const marker: Marker = {
                info: {
                    id: BigInt(0),
                    job: '',
                    name: values.name,
                    description: values.description,
                    x: location.value?.x ?? 0,
                    y: location.value?.y ?? 0,
                    color: values.color.substring(1),
                },
                type: values.markerType,
            };

            if (values.markerType === MARKER_TYPE.CIRCLE) {
                marker.data = {
                    data: {
                        oneofKind: 'circle',
                        circle: {
                            radius: values.circleRadius,
                            oapcity: values.circleOpacity,
                        },
                    },
                };
            }

            const call = $grpc.getLivemapperClient().createOrUpdateMarker({
                marker: marker,
            });
            await call;

            emits('close');

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

const markerTypes = ref<{ status: MARKER_TYPE; selected?: boolean }[]>([{ status: MARKER_TYPE.CIRCLE }]);

defineRule('required', required);
defineRule('digits', digits);
defineRule('min', min);
defineRule('max', max);

interface FormData {
    name: string;
    description: string;
    markerType: MARKER_TYPE.CIRCLE;
    color: string;
    circleRadius: number;
    circleOpacity: number;
}

const { handleSubmit, values, setValues } = useForm<FormData>({
    validationSchema: {
        name: { required: true, min: 3, max: 255 },
        description: { required: false, min: 6, max: 512 },
    },
    initialValues: {
        color: '#EE4B2B',
        circleRadius: 50,
        circleOpacity: 5,
    },
});
setValues({
    markerType: MARKER_TYPE.CIRCLE,
    circleRadius: 50,
    circleOpacity: 5,
});

const onSubmit = handleSubmit(async (values): Promise<void> => await createMarker(values));
</script>

<template>
    <TransitionRoot as="template" :show="open">
        <Dialog as="div" class="relative z-10" @close="$emit('close')">
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
                                    @submit="onSubmit"
                                    class="flex h-full flex-col divide-y divide-gray-200 bg-gray-900 shadow-xl"
                                >
                                    <div class="h-0 flex-1 overflow-y-auto">
                                        <div class="bg-primary-700 px-4 py-6 sm:px-6">
                                            <div class="flex items-center justify-between">
                                                <DialogTitle class="text-base font-semibold leading-6 text-white">
                                                    {{ $t('components.livemap.create_marker.title') }}
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
                                                    {{ $t('components.livemap.create_marker.sub_title') }}
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
                                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.name')"
                                                                    :label="$t('common.name')"
                                                                />
                                                                <VeeErrorMessage
                                                                    name="name"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
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
                                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.description')"
                                                                    :label="$t('common.description')"
                                                                />
                                                                <VeeErrorMessage
                                                                    name="description"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
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
                                                                <VeeField
                                                                    type="color"
                                                                    name="color"
                                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.color')"
                                                                    :label="$t('common.color')"
                                                                />
                                                                <VeeErrorMessage
                                                                    name="color"
                                                                    as="p"
                                                                    class="mt-2 text-sm text-error-400"
                                                                />
                                                            </dd>
                                                        </div>
                                                        <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                            <dt class="text-sm font-medium leading-6 text-white">
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
                                                                    as="div"
                                                                    name="markerType"
                                                                    class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    :placeholder="$t('common.marker')"
                                                                    :label="$t('common.marker')"
                                                                    v-slot="{ field }"
                                                                >
                                                                    <select
                                                                        v-bind="field"
                                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                                    >
                                                                        <option
                                                                            v-for="mtype in markerTypes"
                                                                            :selected="mtype.selected"
                                                                            :value="mtype.status"
                                                                        >
                                                                            {{
                                                                                $t(
                                                                                    `enums.livemap.MARKER_TYPE.${
                                                                                        MARKER_TYPE[
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
                                                        <template v-if="values.markerType === MARKER_TYPE.CIRCLE">
                                                            <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                                                <dt class="text-sm font-medium leading-6 text-white">
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
                                                                        class="block w-full rounded-md border-0 py-1.5 bg-base-700 text-neutral placeholder:text-base-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
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
                                                        </template>
                                                    </dl>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="flex flex-shrink-0 justify-end px-4 py-4">
                                        <span class="isolate inline-flex rounded-md shadow-sm pr-4 w-full">
                                            <button
                                                type="submit"
                                                class="w-full relative inline-flex items-center rounded-l-md bg-primary-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-primary-400"
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
