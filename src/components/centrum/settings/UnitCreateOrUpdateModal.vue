<script lang="ts" setup>
import { max, min, required } from '@vee-validate/rules';
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
        <UForm :schema="{}" :state="{}" @submit="onSubmitThrottle">
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
                                        :placeholder="$t('common.name')"
                                        :label="$t('common.name')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                                <div class="flex-1">
                                    <label for="initials" class="block text-sm font-medium leading-6">
                                        {{ $t('common.initials') }}
                                    </label>
                                    <VeeField
                                        name="initials"
                                        type="text"
                                        :placeholder="$t('common.initials')"
                                        :label="$t('common.initials')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                                <div class="flex-1">
                                    <label for="description" class="block text-sm font-medium leading-6">
                                        {{ $t('common.description') }}
                                    </label>
                                    <VeeField
                                        name="description"
                                        type="text"
                                        :placeholder="$t('common.description')"
                                        :label="$t('common.description')"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                                <div class="flex-1">
                                    <label for="attributes" class="block text-sm font-medium leading-6">
                                        {{ $t('common.attributes', 2) }}
                                    </label>
                                    <VeeField
                                        name="attributes"
                                        :placeholder="$t('common.attributes', 2)"
                                        :label="$t('common.attributes', 2)"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    >
                                        <USelectMenu
                                            v-model="selectedAttributes"
                                            multiple
                                            nullable
                                            :options="availableAttributes"
                                            :placeholder="selectedAttributes ? selectedAttributes.join(', ') : $t('common.na')"
                                        >
                                            <template #option-empty="{ query: search }">
                                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                            </template>
                                            <template #empty>
                                                {{ $t('common.not_found', [$t('common.attributes', 1)]) }}
                                            </template>
                                        </USelectMenu>
                                    </VeeField>
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
                                        :placeholder="`${$t('common.department')} ${$t('common.postal_code')}`"
                                        :label="`${$t('common.department')} ${$t('common.postal_code')}`"
                                        @focusin="focusTablet(true)"
                                        @focusout="focusTablet(false)"
                                    />
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton type="submit" block class="flex-1" :loading="!canSubmit" :disabled="!canSubmit">
                            <template v-if="unit && unit?.id">
                                {{ $t('components.centrum.units.update_unit') }}
                            </template>
                            <template v-else>
                                {{ $t('components.centrum.units.create_unit') }}
                            </template>
                        </UButton>

                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
