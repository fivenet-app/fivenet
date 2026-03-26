<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import { getCitizensLabelsClient } from '~~/gen/ts/clients';
import type { Label, Settings } from '~~/gen/ts/resources/citizens/labels/labels';
import type { CreateOrUpdateLabelResponse, GetLabelResponse } from '~~/gen/ts/services/citizens/labels';

const props = defineProps<{
    labelId?: number;
}>();

const emits = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'refresh'): void;
}>();

const citizensLabelsClient = await getCitizensLabelsClient();

const schema = z.object({
    id: z.coerce.number(),
    name: z.coerce.string().min(1).max(64),
    color: z.coerce.string().length(7),
    icon: z.coerce.string().max(255).optional(),
    settings: z.custom<Settings>(),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    id: 0,
    name: '',
    color: '#ffffff',
    icon: undefined,
    settings: {
        requiresExpiration: false,
    },
});

const { data, status, error, refresh } = useLazyAsyncData(
    `citizens-label-${props.labelId}`,
    () => getCitizenLabel(props.labelId!),
    {
        immediate: !!props.labelId,
        watch: [() => props.labelId],
    },
);

async function getCitizenLabel(labelId: number): Promise<GetLabelResponse> {
    try {
        const { response } = await citizensLabelsClient.getLabel({ id: labelId });

        if (!response?.label) return response;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

function setFromData(label: Label | undefined): void {
    if (!label) return;

    state.id = label.id;
    state.name = label.name;
    state.color = label.color;
    state.icon = label.icon;
    state.settings.requiresExpiration = label.settings?.requiresExpiration ?? false;
}

watch(data, () => setFromData(data.value?.label));
setFromData(data.value?.label);

async function createOrUpdateLabel(values: Schema): Promise<CreateOrUpdateLabelResponse> {
    try {
        const { response } = await citizensLabelsClient.createOrUpdateLabel({
            label: {
                id: values.id ?? 0,
                name: values.name ?? '',
                color: values.color ?? '#ffffff',
                icon: values.icon,
                settings: {
                    requiresExpiration: values.settings.requiresExpiration,
                },
                access: {
                    jobs: [],
                },
            },
        });

        if (!response?.label) return response;

        const label = response.label;

        state.id = label.id;
        state.name = label.name;
        state.color = label.color;

        emits('refresh');
        emits('close', false);

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await createOrUpdateLabel(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.citizens.citizen_labels.title')">
        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.label', 2)])" />
            <DataErrorBlock v-else-if="error" :error="error" :retry="refresh" />

            <UForm v-else ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField class="flex-1" name="label.name" :label="$t('common.name')">
                    <UInput
                        v-model="state.name"
                        class="w-full"
                        name="label.name"
                        type="text"
                        :placeholder="$t('common.name')"
                    />
                </UFormField>

                <UFormField name="label.color" :label="$t('common.color')">
                    <ColorPicker v-model="state.color" class="w-full" name="label.color" />
                </UFormField>

                <UFormField name="label.icon" :label="$t('common.icon')">
                    <IconSelectMenu v-model="state.icon" name="label.icon" :hex-color="state.color" class="w-full" />
                </UFormField>

                <!--
                <UFormField name="label.settings" :label="$t('common.settings')">
                    <UFormField
                        name="label.settings.requiresExpiration"
                        :label="$t('components.citizens.citizen_labels.settings.requires_expiration')"
                    >
                        <USwitch v-model="state.settings.requiresExpiration" name="label.settings.requiresExpiration" />
                    </UFormField>
                    <div>TODO</div>
                </UFormField>

                <UFormField name="access" :label="$t('common.access')">
                    <div>TODO</div>
                </UFormField>
                -->
            </UForm>
        </template>

        <template #footer>
            <UFieldGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit"
                    :loading="!canSubmit"
                    :label="$t('common.save')"
                    @click="() => formRef?.submit()"
                />
            </UFieldGroup>
        </template>
    </UModal>
</template>
