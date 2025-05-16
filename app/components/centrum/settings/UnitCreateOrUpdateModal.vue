<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import { useNotificatorStore } from '~/stores/notificator';
import { UnitAccessLevel, type UnitJobAccess, type UnitQualificationAccess } from '~~/gen/ts/resources/centrum/access';
import { UnitAttribute } from '~~/gen/ts/resources/centrum/attributes';
import type { Unit } from '~~/gen/ts/resources/centrum/units';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const props = defineProps<{
    unit?: Unit;
}>();

const emit = defineEmits<{
    (e: 'created', unit: Unit): void;
    (e: 'updated', unit: Unit): void;
}>();

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificatorStore();

const availableAttributes = ref<{ type: UnitAttribute }[]>([
    { type: UnitAttribute.STATIC },
    { type: UnitAttribute.NO_DISPATCH_AUTO_ASSIGN },
]);

const { maxAccessEntries } = useAppConfig();

const schema = z.object({
    name: z.string().min(3).max(24),
    initials: z.string().min(2).max(4),
    description: z.union([z.string().min(1).max(255), z.string().length(0).optional()]),
    color: z.string().length(7),
    homePostal: z.union([z.string().min(1).max(48), z.string().length(0).optional()]),
    attributes: z.nativeEnum(UnitAttribute).array().max(5),
    access: z.object({
        jobs: z.custom<UnitJobAccess>().array().max(maxAccessEntries),
        qualifications: z.custom<UnitQualificationAccess>().array().max(maxAccessEntries),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    initials: '',
    description: '',
    color: '#000000',
    attributes: [],
    access: {
        jobs: [],
        qualifications: [],
    },
});

const selectedAttributes = ref<string[]>([]);

async function createOrUpdateUnit(values: Schema): Promise<void> {
    try {
        const call = $grpc.centrum.centrum.createOrUpdateUnit({
            unit: {
                id: props.unit?.id ?? 0,
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
                access: values.access,
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successfull.title', parameters: {} },
            description: { key: 'notifications.action_successfull.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

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
    state.access = {
        jobs: props.unit.access?.jobs ?? [],
        qualifications: props.unit.access?.qualifications ?? [],
    };
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

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <div>
                    <UFormGroup class="flex-1" name="name" :label="$t('common.name')">
                        <UInput v-model="state.name" name="name" type="text" :placeholder="$t('common.name')" />
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="initials" :label="$t('common.initials')">
                        <UInput v-model="state.initials" name="initials" type="text" :placeholder="$t('common.initials')" />
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="description" :label="$t('common.description')">
                        <UInput
                            v-model="state.description"
                            name="description"
                            type="text"
                            :placeholder="$t('common.description')"
                        />
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="attributes" :label="$t('common.attributes', 2)">
                        <ClientOnly>
                            <USelectMenu
                                v-model="state.attributes"
                                multiple
                                nullable
                                value-attribute="type"
                                :options="availableAttributes"
                                :placeholder="selectedAttributes ? selectedAttributes.join(', ') : $t('common.na')"
                                :searchable-placeholder="$t('common.search_field')"
                            >
                                <template #option="{ option }">
                                    <span class="truncate">{{
                                        $t(`enums.centrum.UnitAttribute.${UnitAttribute[option.type]}`, 2)
                                    }}</span>
                                </template>

                                <template #option-empty="{ query: search }">
                                    <q>{{ search }}</q> {{ $t('common.query_not_found') }}
                                </template>

                                <template #empty>
                                    {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormGroup>

                    <UFormGroup class="flex-1" name="color" :label="$t('common.color')">
                        <ColorPickerClient v-model="state.color" />
                    </UFormGroup>

                    <UFormGroup
                        class="flex-1"
                        name="homePostal"
                        :label="`${$t('common.department')} ${$t('common.postal_code')}`"
                    >
                        <UInput
                            v-model="state.homePostal"
                            name="homePostal"
                            type="text"
                            :placeholder="`${$t('common.department')} ${$t('common.postal_code')}`"
                        />
                    </UFormGroup>

                    <UFormGroup name="access" :label="$t('common.access')">
                        <AccessManager
                            v-model:jobs="state.access.jobs"
                            v-model:qualifications="state.access.qualifications"
                            :target-id="unit?.id ?? 0"
                            :access-roles="
                                enumToAccessLevelEnums(UnitAccessLevel, 'enums.centrum.UnitAccessLevel').filter(
                                    (a) => a.value > 1,
                                )
                            "
                            :access-types="[
                                { type: 'job', name: $t('common.job', 2) },
                                { type: 'qualification', name: $t('common.qualification', 2) },
                            ]"
                        />
                    </UFormGroup>
                </div>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="black" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton class="flex-1" type="submit" block :loading="!canSubmit" :disabled="!canSubmit">
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
