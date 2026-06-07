<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import AccessManager from '~/components/partials/access/AccessManager.vue';
import { enumToAccessLevelEnums } from '~/components/partials/access/helpers';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import { getCentrumUnitsClient } from '~~/gen/ts/clients';
import { UnitAccessLevel } from '~~/gen/ts/resources/centrum/units/access/access';
import { type Unit, UnitAttribute } from '~~/gen/ts/resources/centrum/units/units';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { defaultUnitIcon } from '../helpers';

const props = defineProps<{
    unit?: Unit;
}>();

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'created', unit: Unit): void;
    (e: 'updated', unit: Unit): void;
}>();

const { t } = useI18n();

const notifications = useNotificationsStore();

const centrumUnitsClient = await getCentrumUnitsClient();

const availableAttributes = ref<{ label: string; value: UnitAttribute }[]>([
    { label: t(`enums.centrum.UnitAttribute.${UnitAttribute[UnitAttribute.STATIC]}`, 2), value: UnitAttribute.STATIC },
    {
        label: t(`enums.centrum.UnitAttribute.${UnitAttribute[UnitAttribute.NO_DISPATCH_AUTO_ASSIGN]}`, 2),
        value: UnitAttribute.NO_DISPATCH_AUTO_ASSIGN,
    },
]);

const { maxAccessEntries } = useAppConfig();

const schema = z.object({
    name: z.coerce.string().min(3).max(24),
    initials: z.coerce.string().min(2).max(4),
    description: z.union([z.string().min(1).max(255), z.string().length(0).optional()]),
    color: z.coerce.string().length(7),
    icon: z.string().max(128).optional(),
    homePostal: z.union([z.string().trim().min(1).max(48), z.string().trim().length(0).optional()]),
    attributes: z.enum(UnitAttribute).array().max(5).default([]),
    access: z.object({
        jobs: jobsAccessEntries(t).max(maxAccessEntries).default([]),
        qualifications: qualificationAccessEntries(t).max(maxAccessEntries).default([]),
    }),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    name: '',
    initials: '',
    description: '',
    color: '#000000',
    icon: defaultUnitIcon,
    attributes: [],
    access: {
        jobs: [],
        qualifications: [],
    },
});

const { hasUnsavedChanges, confirmLeave, syncSnapshot } = useSnapshotChanges(state);

async function createOrUpdateUnit(values: Schema): Promise<void> {
    values.access.jobs.forEach((job) => job.id < 0 && (job.id = 0));
    values.access.qualifications.forEach((quali) => quali.id < 0 && (quali.id = 0));

    try {
        const call = centrumUnitsClient.createOrUpdateUnit({
            unit: {
                id: props.unit?.id ?? 0,
                job: '',
                name: values.name,
                initials: values.initials,
                description: values.description,
                color: values.color,
                icon: values.icon === '' ? undefined : values.icon,
                attributes: {
                    list: values.attributes,
                },
                users: [],
                homePostal: values.homePostal,
                access: values.access,
                sortOrder: props.unit?.sortOrder ?? 0,
            },
        });
        const { response } = await call;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        if (props.unit?.id === undefined) {
            emit('created', response.unit!);
        } else {
            emit('updated', response.unit!);
        }

        emit('close', false);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref<boolean>(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateUnit(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

async function updateUnitInForm(): Promise<void> {
    if (props.unit === undefined) return;

    state.name = props.unit.name;
    state.initials = props.unit.initials;
    state.description = props.unit.description;
    state.color = props.unit.color;
    state.icon = props.unit.icon;
    state.attributes = props.unit.attributes?.list ?? [];
    state.homePostal = props.unit.homePostal;
    state.access = {
        jobs: props.unit.access?.jobs ?? [],
        qualifications: props.unit.access?.qualifications ?? [],
    };

    syncSnapshot();
}

async function closeModal(): Promise<void> {
    if (!canSubmit.value) return;

    if (hasUnsavedChanges.value && !(await confirmLeave())) return;

    emit('close', false);
}

onBeforeMount(async () => updateUnitInForm());
watch(props, async () => updateUnitInForm());

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal
        :title="unit?.id ? $t('components.dispatch.units.update_unit') : $t('components.dispatch.units.create_unit')"
        :close="false"
        :dismissible="!hasUnsavedChanges && canSubmit"
    >
        <template #header>
            <div class="flex w-full items-center justify-between gap-1.5">
                <h3 class="font-semibold text-highlighted">
                    {{ unit?.id ? $t('components.dispatch.units.update_unit') : $t('components.dispatch.units.create_unit') }}
                </h3>

                <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-mdi-close"
                    :disabled="!canSubmit"
                    :aria-label="$t('common.close', 1)"
                    @click="closeModal"
                />
            </div>
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField class="flex-1" name="name" :label="$t('common.name')">
                    <UInput v-model="state.name" class="w-full" name="name" type="text" :placeholder="$t('common.name')" />
                </UFormField>

                <UFormField class="flex-1" name="initials" :label="$t('common.initials')">
                    <UInput
                        v-model="state.initials"
                        class="w-full"
                        name="initials"
                        type="text"
                        :placeholder="$t('common.initials')"
                    />
                </UFormField>

                <UFormField class="flex-1" name="description" :label="$t('common.description')">
                    <UInput
                        v-model="state.description"
                        class="w-full"
                        name="description"
                        type="text"
                        :placeholder="$t('common.description')"
                    />
                </UFormField>

                <UFormField class="flex-1" name="attributes" :label="$t('common.attributes', 2)">
                    <ClientOnly>
                        <USelectMenu
                            v-model="state.attributes"
                            class="w-full"
                            multiple
                            nullable
                            :items="availableAttributes"
                            :search-input="{ placeholder: $t('common.search_field') }"
                            value-key="value"
                        >
                            <template #empty>
                                {{ $t('common.not_found', [$t('common.attributes', 2)]) }}
                            </template>
                        </USelectMenu>
                    </ClientOnly>
                </UFormField>

                <UFormField class="flex-1" name="color" :label="$t('common.color')">
                    <ColorPicker v-model="state.color" class="w-full" />
                </UFormField>

                <UFormField class="flex-1" name="icon" :label="$t('common.icon')">
                    <IconSelectMenu v-model="state.icon" class="w-full" :hex-color="state.color" />
                </UFormField>

                <UFormField class="flex-1" name="homePostal" :label="`${$t('common.department')} ${$t('common.postal_code')}`">
                    <UInput
                        v-model="state.homePostal"
                        class="w-full"
                        name="homePostal"
                        type="text"
                        :placeholder="`${$t('common.department')} ${$t('common.postal_code')}`"
                    />
                </UFormField>

                <UFormField name="access" :label="$t('common.access')">
                    <AccessManager
                        v-model:jobs="state.access.jobs"
                        v-model:qualifications="state.access.qualifications"
                        :target-id="unit?.id ?? 0"
                        :access-roles="
                            enumToAccessLevelEnums(UnitAccessLevel, 'enums.centrum.UnitAccessLevel').filter((a) => a.value > 1)
                        "
                        :access-types="[
                            { label: $t('common.job', 2), value: 'job' },
                            { label: $t('common.qualification', 2), value: 'qualification' },
                        ]"
                        name="access"
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton
                    class="flex-1"
                    color="neutral"
                    block
                    :disabled="!canSubmit"
                    :label="$t('common.close', 1)"
                    @click="closeModal"
                />

                <UButton
                    class="flex-1"
                    block
                    :loading="!canSubmit"
                    :disabled="!canSubmit"
                    :label="
                        unit && unit?.id
                            ? $t('components.dispatch.units.update_unit')
                            : $t('components.dispatch.units.create_unit')
                    "
                    @click="formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
