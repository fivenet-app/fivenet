<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import { useCompletorStore } from '~/stores/completor';
import type { ManageLabelsResponse } from '~~/gen/ts/services/citizens/citizens';

const { $grpc } = useNuxtApp();

const { can } = useAuth();

const { isOpen } = useModal();

const completorStore = useCompletorStore();

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
        const { response } = await $grpc.citizens.citizens.manageLabels({
            labels: values.labels ?? [],
        });

        state.labels = response.labels;

        isOpen.value = false;

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
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('components.citizens.citizen_labels.title') }}
                        </h3>

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <UFormGroup
                    v-if="state && can('citizens.CitizensService/ManageLabels').value"
                    class="grid items-center gap-2"
                    name="citizenAttributes.list"
                    :ui="{ container: '' }"
                >
                    <div class="flex flex-col gap-1">
                        <div v-for="(_, idx) in state.labels" :key="idx" class="flex items-center gap-1">
                            <UFormGroup class="flex-1" :name="`labels.${idx}.name`">
                                <UInput
                                    v-model="state.labels[idx]!.name"
                                    class="w-full flex-1"
                                    :name="`labels.${idx}.name`"
                                    type="text"
                                    :placeholder="$t('common.label', 1)"
                                />
                            </UFormGroup>

                            <UFormGroup :name="`labels.${idx}.color`">
                                <ColorPickerClient
                                    v-model="state.labels[idx]!.color"
                                    class="min-w-16"
                                    :name="`labels.${idx}.color`"
                                />
                            </UFormGroup>

                            <UButton
                                :ui="{ rounded: 'rounded-full' }"
                                :disabled="!canSubmit"
                                icon="i-mdi-close"
                                @click="state.labels.splice(idx, 1)"
                            />
                        </div>
                    </div>

                    <UButton
                        :class="state.labels.length ? 'mt-2' : ''"
                        :ui="{ rounded: 'rounded-full' }"
                        :disabled="!canSubmit"
                        icon="i-mdi-plus"
                        @click="state.labels.push({ id: 0, name: '', color: '#ffffff' })"
                    />
                </UFormGroup>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="black" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton class="flex-1" type="submit" block :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
