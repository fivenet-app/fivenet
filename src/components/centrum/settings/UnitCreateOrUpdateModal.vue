<script lang="ts" setup>
import { Listbox, ListboxButton, ListboxOption, ListboxOptions } from '@headlessui/vue';
import { max, min, required } from '@vee-validate/rules';
import { CheckIcon, ChevronDownIcon } from 'mdi-vue3';
import { defineRule } from 'vee-validate';
import ColorInput from 'vue-color-input/dist/color-input.esm';
import { Unit } from '~~/gen/ts/resources/centrum/units';

const props = defineProps<{
    unit?: Unit;
}>();

const emit = defineEmits<{
    (e: 'created', unit: Unit): void;
    (e: 'updated', unit: Unit): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

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

async function createOrUpdateUnit(values: FormData): Promise<void> {
    try {
        const call = $grpc.getCentrumClient().createOrUpdateUnit({
            unit: {
                id: props.unit?.id ?? '0',
                job: '',
                name: values.name,
                initials: values.initials,
                color: values.color,
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

        isOpen.value = false;
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
            color: props.unit.color,
            homePostal: props.unit.homePostal,
        });

        selectedAttributes.value = props.unit.attributes?.list ?? [];
    }
}

watch(props, async () => updateUnitInForm());

onBeforeMount(async () => updateUnitInForm());
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="text-2xl font-semibold leading-6">
                        <template v-if="unit && unit?.id">
                            {{ $t('components.centrum.units.update_unit') }}
                        </template>
                        <template v-else>
                            {{ $t('components.centrum.units.create_unit') }}
                        </template>
                    </h3>

                    <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <UForm :state="{}" @submit.prevent="onSubmitThrottle">
                    <div class="text-center">
                        <div>
                            <div class="text-sm text-gray-100">
                                <div class="flex-1">
                                    <label for="name" class="block text-sm font-medium leading-6">
                                        {{ $t('common.name') }}
                                    </label>
                                    <VeeField
                                        name="name"
                                        type="text"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.name')"
                                        :label="$t('common.name')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                    <VeeErrorMessage name="name" as="p" class="mt-2 text-sm text-error-400" />
                                </div>
                                <div class="flex-1">
                                    <label for="initials" class="block text-sm font-medium leading-6">
                                        {{ $t('common.initials') }}
                                    </label>
                                    <VeeField
                                        name="initials"
                                        type="text"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.initials')"
                                        :label="$t('common.initials')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                    <VeeErrorMessage name="initials" as="p" class="mt-2 text-sm text-error-400" />
                                </div>
                                <div class="flex-1">
                                    <label for="description" class="block text-sm font-medium leading-6">
                                        {{ $t('common.description') }}
                                    </label>
                                    <VeeField
                                        name="description"
                                        type="text"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.description')"
                                        :label="$t('common.description')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                    <VeeErrorMessage name="description" as="p" class="mt-2 text-sm text-error-400" />
                                </div>
                                <div class="flex-1">
                                    <label for="attributes" class="block text-sm font-medium leading-6">
                                        {{ $t('common.attributes', 2) }}
                                    </label>
                                    <VeeField
                                        name="attributes"
                                        type="text"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="$t('common.attributes', 2)"
                                        :label="$t('common.attributes', 2)"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    >
                                        <Listbox v-model="selectedAttributes" as="div" nullable multiple>
                                            <div class="relative">
                                                <ListboxButton
                                                    class="block w-full rounded-md border-0 bg-base-700 py-1.5 pl-3 text-left placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                                >
                                                    <span class="block truncate">
                                                        <template v-if="selectedAttributes.length > 0">
                                                            <span v-for="attr in selectedAttributes" :key="attr" class="mr-1">
                                                                {{ $t(`components.centrum.units.attributes.${attr}`) }}
                                                            </span>
                                                        </template>
                                                        <template v-else>
                                                            {{ $t('common.none_selected') }}
                                                        </template>
                                                    </span>
                                                    <span
                                                        class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2"
                                                    >
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
                                                            v-for="attr in availableAttributes"
                                                            :key="attr"
                                                            v-slot="{ active, selected }"
                                                            as="template"
                                                            :value="attr"
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
                                                                    {{ $t(`components.centrum.units.attributes.${attr}`) }}
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
                                    <VeeErrorMessage name="attributes" as="p" class="mt-2 text-sm text-error-400" />
                                </div>
                                <div class="flex-1">
                                    <label for="color" class="block text-sm font-medium leading-6">
                                        {{ $t('common.color') }}
                                    </label>
                                    <ColorInput
                                        :model-value="unit?.color ?? '#000000'"
                                        disable-alpha
                                        format="hex"
                                        position="top"
                                        @change="setFieldValue('color', $event)"
                                    />
                                </div>
                                <div class="flex-1">
                                    <label for="homePostal" class="block text-sm font-medium leading-6">
                                        {{ `${$t('common.department')} ${$t('common.postal_code')}` }}
                                    </label>
                                    <VeeField
                                        name="homePostal"
                                        type="text"
                                        class="block w-full rounded-md border-0 bg-base-700 py-1.5 placeholder:text-accent-200 focus:ring-2 focus:ring-inset focus:ring-base-300 sm:text-sm sm:leading-6"
                                        :placeholder="`${$t('common.department')} ${$t('common.postal_code')}`"
                                        :label="`${$t('common.department')} ${$t('common.postal_code')}`"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                    <VeeErrorMessage name="homePostal" as="p" class="mt-2 text-sm text-error-400" />
                                </div>
                            </div>
                        </div>
                    </div>
                </UForm>
            </div>

            <template #footer>
                <div class="gap-2 sm:flex">
                    <UButton class="flex-1" @click="isOpen = false">
                        {{ $t('common.close', 1) }}
                    </UButton>
                    <UButton
                        class="flex-1"
                        :loading="!canSubmit"
                        :disabled="!meta.valid || !canSubmit"
                        @click="onSubmitThrottle"
                    >
                        <template v-if="unit && unit?.id">
                            {{ $t('components.centrum.units.update_unit') }}
                        </template>
                        <template v-else>
                            {{ $t('components.centrum.units.create_unit') }}
                        </template>
                    </UButton>
                </div>
            </template>
        </UCard>
    </UModal>
</template>
