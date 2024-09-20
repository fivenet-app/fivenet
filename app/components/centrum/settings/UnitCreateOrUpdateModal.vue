<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import { Unit } from '~~/gen/ts/resources/centrum/units';

const props = defineProps<{
    unit?: Unit;
}>();

const emit = defineEmits<{
    (e: 'created', unit: Unit): void;
    (e: 'updated', unit: Unit): void;
}>();

const { isOpen } = useModal();

const availableAttributes: string[] = ['static', 'no_dispatch_auto_assign'];

const schema = z.object({
    name: z.string().min(3).max(24),
    initials: z.string().min(2).max(4),
    description: z.union([z.string().min(1).max(255), z.string().length(0).optional()]),
    color: z.string().length(7),
    homePostal: z.union([z.string().min(1).max(48), z.string().length(0).optional()]),
    attributes: z.string().array().max(5),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    initials: '',
    description: '',
    color: '#000000',
    attributes: [],
});

const selectedAttributes = ref<string[]>([]);

async function createOrUpdateUnit(values: Schema): Promise<void> {
    try {
        const call = getGRPCCentrumClient().createOrUpdateUnit({
            unit: {
                id: props.unit?.id ?? '0',
                job: '',
                name: values.name,
                initials: values.initials,
                description: values.description,
                color: values.color,
                attributes: {
                    list: values.attributes,
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
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateUnit(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

async function updateUnitInForm(): Promise<void> {
    if (props.unit === undefined) {
        return;
    }

    state.name = props.unit.name;
    state.initials = props.unit.initials;
    state.description = props.unit.description;
    state.color = props.unit.color;
    state.attributes = props.unit.attributes?.list ?? [];
    state.homePostal = props.unit.homePostal;
}

watch(props, async () => updateUnitInForm());

onMounted(async () => updateUnitInForm());
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
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
                    <UFormGroup name="name" :label="$t('common.name')" class="flex-1">
                        <UInput v-model="state.name" name="name" type="text" :placeholder="$t('common.name')" />
                    </UFormGroup>

                    <UFormGroup name="initials" :label="$t('common.initials')" class="flex-1">
                        <UInput v-model="state.initials" name="initials" type="text" :placeholder="$t('common.initials')" />
                    </UFormGroup>

                    <UFormGroup name="description" :label="$t('common.description')" class="flex-1">
                        <UInput
                            v-model="state.description"
                            name="description"
                            type="text"
                            :placeholder="$t('common.description')"
                        />
                    </UFormGroup>

                    <UFormGroup name="attributes" :label="$t('common.attributes', 2)" class="flex-1">
                        <USelectMenu
                            v-model="state.attributes"
                            multiple
                            nullable
                            :options="availableAttributes"
                            :placeholder="selectedAttributes ? selectedAttributes.join(', ') : $t('common.na')"
                            :searchable-placeholder="$t('common.search_field')"
                        >
                            <template #option-empty="{ query: search }">
                                <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                            </template>
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                            </template>
                        </USelectMenu>
                    </UFormGroup>

                    <UFormGroup name="color" :label="$t('common.color')" class="flex-1">
                        <ColorPicker v-model="state.color" />
                    </UFormGroup>

                    <UFormGroup
                        name="homePostal"
                        :label="`${$t('common.department')} ${$t('common.postal_code')}`"
                        class="flex-1"
                    >
                        <UInput
                            v-model="state.homePostal"
                            name="homePostal"
                            type="text"
                            :placeholder="`${$t('common.department')} ${$t('common.postal_code')}`"
                        />
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :loading="!canSubmit" :disabled="!canSubmit">
                            <template v-if="unit && unit?.id">
                                {{ $t('components.centrum.units.update_unit') }}
                            </template>
                            <template v-else>
                                {{ $t('components.centrum.units.create_unit') }}
                            </template>
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
