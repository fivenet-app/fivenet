<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { z } from 'zod';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
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
            name: z.string().min(1).max(64),
            color: z.string().length(7),
        })
        .array()
        .max(15)
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
</script>

<template>
    <UModal>
        <template #title>
            <h3 class="text-2xl leading-6 font-semibold">
                {{ $t('components.citizens.citizen_labels.title') }}
            </h3>
        </template>

        <template #body>
            <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
                <UFormField
                    v-if="state && can('citizens.CitizensService/ManageLabels').value"
                    class="grid items-center gap-2"
                    name="citizenAttributes.list"
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
                                <ColorPickerClient
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
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block @click="$emit('close', false)">
                    {{ $t('common.close', 1) }}
                </UButton>

                <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                    {{ $t('common.save') }}
                </UButton>
            </UButtonGroup>
        </template>
    </UModal>
</template>
