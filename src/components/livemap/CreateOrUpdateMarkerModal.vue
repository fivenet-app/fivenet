<script lang="ts" setup>
import { Combobox, ComboboxButton, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/vue';
import { digits, max, max_value, min, min_value, required } from '@vee-validate/rules';
import { defineRule } from 'vee-validate';
import { CheckIcon, ChevronDownIcon } from 'mdi-vue3';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import { useLivemapStore } from '~/store/livemap';
import { type MarkerMarker, MarkerType } from '~~/gen/ts/resources/livemap/livemap';
import { markerIcons } from '~/components/livemap/helpers';

const props = defineProps<{
    location?: Coordinate;
}>();

const emit = defineEmits<{
    (e: 'close'): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useSlideover();

const livemapStore = useLivemapStore();
const { location: storeLocation } = storeToRefs(livemapStore);
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
const selectedIcon = ref<string>('i-mdi-help');

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
                x: props.location ? props.location.x : storeLocation.value?.x ?? 0,
                y: props.location ? props.location.y : storeLocation.value?.y ?? 0,
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
                        icon: selectedIcon.value,
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
        isOpen.value = false;
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

onBeforeMount(() => setMarker());

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
    <USlideover>
        <UCard
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex-1 max-h-[calc(100vh-(2*var(--header-height)))] overflow-y-auto',
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        {{ $t('components.livemap.create_marker.title') }}
                    </h3>

                    <UButton
                        color="gray"
                        variant="ghost"
                        icon="i-mdi-window-close"
                        class="-my-1"
                        @click="
                            $emit('close');
                            isOpen = false;
                        "
                    />
                </div>
            </template>

            <div>
                <div class="flex flex-1 flex-col justify-between">
                    <div class="divide-y divide-gray-200 px-2 sm:px-6">
                        <div class="mt-1">
                            <dl class="divide-neutral/10 border-neutral/10 divide-y border-b">
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        <label for="name" class="block text-sm font-medium leading-6">
                                            {{ $t('common.name') }}
                                        </label>
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                        <VeeField
                                            type="text"
                                            name="name"
                                            :placeholder="$t('common.name')"
                                            :label="$t('common.name')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        <label for="description" class="block text-sm font-medium leading-6">
                                            {{ $t('common.description') }}
                                        </label>
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                        <VeeField
                                            type="text"
                                            name="description"
                                            :placeholder="$t('common.description')"
                                            :label="$t('common.description')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        <label for="expiresAt" class="block text-sm font-medium leading-6">
                                            {{ $t('common.expires_at') }}
                                        </label>
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                        <VeeField
                                            type="datetime-local"
                                            name="expiresAt"
                                            :placeholder="$t('common.expires_at')"
                                            :label="$t('common.expires_at')"
                                            @focusin="focusTablet(true)"
                                            @focusout="focusTablet(false)"
                                        />
                                        <VeeErrorMessage name="expiresAt" as="p" class="mt-2 text-sm text-error-400" />
                                    </dd>
                                </div>

                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        <label for="color" class="block text-sm font-medium leading-6">
                                            {{ $t('common.color') }}
                                        </label>
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                        <ColorInput v-model="color" disable-alpha format="hex" position="top" />
                                    </dd>
                                </div>
                                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                                    <dt class="text-sm font-medium leading-6">
                                        <label for="markerType" class="block text-sm font-medium leading-6">
                                            {{ $t('common.marker') }}
                                        </label>
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                        <VeeField
                                            v-slot="{ field }"
                                            as="div"
                                            name="markerType"
                                            :placeholder="$t('common.marker')"
                                            :label="$t('common.marker')"
                                        >
                                            <select v-bind="field" @focusin="focusTablet(true)" @focusout="focusTablet(false)">
                                                <option
                                                    v-for="mtype in markerTypes"
                                                    :key="mtype.status"
                                                    :selected="mtype.selected"
                                                    :value="mtype.status"
                                                >
                                                    {{
                                                        $t(
                                                            `enums.livemap.MarkerType.${
                                                                MarkerType[mtype.status ?? (0 as number)]
                                                            }`,
                                                        )
                                                    }}
                                                </option>
                                            </select>
                                        </VeeField>
                                        <VeeErrorMessage name="markerType" as="p" class="mt-2 text-sm text-error-400" />
                                    </dd>
                                </div>
                                <div
                                    v-if="values.markerType === MarkerType.CIRCLE"
                                    class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                                >
                                    <dt class="text-sm font-medium leading-6">
                                        <label for="circleRadius" class="block text-sm font-medium leading-6">
                                            {{ $t('common.radius') }}
                                        </label>
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
                                        <VeeField
                                            type="number"
                                            name="circleRadius"
                                            min="5"
                                            max="250"
                                            :placeholder="$t('common.radius')"
                                            :label="$t('common.radius')"
                                        />
                                        <VeeErrorMessage name="circleRadius" as="p" class="mt-2 text-sm text-error-400" />
                                    </dd>
                                </div>
                                <div
                                    v-else-if="values.markerType === MarkerType.ICON"
                                    class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0"
                                >
                                    <dt class="text-sm font-medium leading-6">
                                        <label for="icon" class="block text-sm font-medium leading-6">
                                            {{ $t('common.icon') }}
                                        </label>
                                    </dt>
                                    <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
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
                                                        class="placeholder:text-accent-200 block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                    >
                                                        <ComboboxInput
                                                            autocomplete="off"
                                                            :placeholder="$t('common.icon')"
                                                            @change="queryIcon = $event.target.value"
                                                            @focusin="focusTablet(true)"
                                                            @focusout="focusTablet(false)"
                                                        />

                                                        <span
                                                            class="pointer-events-none absolute inset-y-0 right-0 flex pr-2 pt-4"
                                                        >
                                                            <ChevronDownIcon class="size-5 text-gray-400" />
                                                        </span>

                                                        <span class="mt-1 inline-flex items-center truncate">
                                                            <UIcon
                                                                v-if="selectedIcon"
                                                                :name="selectedIcon"
                                                                class="mr-1 size-5"
                                                                :style="{ color: color }"
                                                            />
                                                            {{ (selectedIcon ?? 'N/A').replace('Icon', '') }}
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
                                                                v-for="icon in markerIcons.filter((icon) =>
                                                                    icon.includes(queryIcon),
                                                                )"
                                                                v-slot="{ active, selected }"
                                                                :key="icon"
                                                                as="template"
                                                                :value="icon"
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
                                                                            'inline-flex items-center truncate',
                                                                        ]"
                                                                    >
                                                                        <UIcon :name="icon" class="mr-1 size-5" />
                                                                        {{ icon.replace('i-mdi-', '') }}
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
                                                            </ComboboxOption>
                                                        </ComboboxOptions>
                                                    </transition>
                                                </div>
                                            </Combobox>
                                        </VeeField>
                                        <VeeErrorMessage name="icon" as="p" class="mt-2 text-sm text-error-400" />
                                    </dd>
                                </div>
                            </dl>
                        </div>
                    </div>
                </div>
            </div>

            <template #footer>
                <UButtonGroup class="inline-flex w-full">
                    <UButton
                        color="black"
                        block
                        class="flex-1"
                        @click="
                            $emit('close');
                            isOpen = false;
                        "
                    >
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton
                        block
                        class="flex-1"
                        :disabled="!meta.valid || !canSubmit"
                        :loading="!canSubmit"
                        @click="onSubmitThrottle"
                    >
                        {{ $t('common.create') }}
                    </UButton>
                </UButtonGroup>
            </template>
        </UCard>
    </USlideover>
</template>
