<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import { useCompletorStore } from '~/stores/completor';
import { getCitizensCitizensClient } from '~~/gen/ts/clients';
import type { ManageLabelsResponse } from '~~/gen/ts/services/citizens/citizens';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const { can } = useAuth();

const completorStore = useCompletorStore();

const citizensCitizensClient = await getCitizensCitizensClient();

const schema = z.object({
    labels: z
        .object({
            id: z.coerce.number(),
            name: z.coerce.string().min(1).max(64),
            color: z.coerce.string().length(7),
        })
        .array()
        .max(50)
        .default([]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    labels: [],
});

const { data: labels } = useLazyAsyncData('citizens-labels', () => completorStore.completeCitizenLabels(''));

async function manageLabels(values: Schema): Promise<ManageLabelsResponse> {
    try {
        const { response } = await citizensCitizensClient.manageLabels({
            labels: values.labels ?? [],
        });

        state.labels = response.labels;

        emit('close', false);

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const canSubmit = ref(true);
const onSubmitThrottle = useThrottleFn(async (event: FormSubmitEvent<Schema>) => {
    canSubmit.value = false;
    await manageLabels(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

watch(labels, () => (state.labels = labels.value ?? []));

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('components.citizens.citizen_labels.title')">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField
                    v-if="state && can('citizens.CitizensService/ManageLabels').value"
                    class="grid items-center gap-2"
                    name="labels"
                    :ui="{ container: '' }"
                >
                    <div class="flex flex-col gap-1">
                        <div v-for="(_, idx) in state.labels" :key="idx" class="flex items-center gap-1">
                            <UFormField class="flex-1" :name="`labels.${idx}.name`">
                                <UInput
                                    v-model="state.labels[idx]!.name"
                                    class="w-full flex-1"
                                    :name="`labels.${idx}.name`"
                                    type="text"
                                    :placeholder="$t('common.label', 1)"
                                />
                            </UFormField>

                            <UFormField :name="`labels.${idx}.color`">
                                <ColorPicker
                                    v-model="state.labels[idx]!.color"
                                    class="min-w-16"
                                    :name="`labels.${idx}.color`"
                                />
                            </UFormField>

                            <UButton :disabled="!canSubmit" icon="i-mdi-close" @click="state.labels.splice(idx, 1)" />
                        </div>
                    </div>

                    <UButton
                        :class="state.labels.length ? 'mt-2' : ''"
                        :disabled="!canSubmit"
                        icon="i-mdi-plus"
                        @click="state.labels.push({ id: 0, name: '', color: '#ffffff' })"
                    />
                </UFormField>
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
