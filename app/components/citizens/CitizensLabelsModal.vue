<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import { useCompletorStore } from '~/stores/completor';
import type { ManageCitizenLabelsResponse } from '~~/gen/ts/services/citizenstore/citizenstore';

const { $grpc } = useNuxtApp();

const { can } = useAuth();

const { isOpen } = useModal();

const completorStore = useCompletorStore();

const schema = z.object({
    labels: z
        .object({
            id: z.number(),
            name: z.string().min(1).max(64),
            color: z.string().length(7),
        })
        .array()
        .max(15),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    labels: [],
});

const { data: labels } = useLazyAsyncData('citizenstore-labels', () => completorStore.completeCitizenLabels(''));

async function manageCitizenLabels(values: Schema): Promise<ManageCitizenLabelsResponse> {
    try {
        const { response } = await $grpc.citizenstore.citizenStore.manageCitizenLabels({
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
    await manageCitizenLabels(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
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

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <UFormGroup
                    v-if="state && can('CitizenStoreService.ManageCitizenLabels').value"
                    name="citizenAttributes.list"
                    class="grid items-center gap-2"
                    :ui="{ container: '' }"
                >
                    <div class="flex flex-col gap-1">
                        <div v-for="(_, idx) in state.labels" :key="idx" class="flex items-center gap-1">
                            <UFormGroup :name="`labels.${idx}.name`" class="flex-1">
                                <UInput
                                    v-model="state.labels[idx]!.name"
                                    :name="`labels.${idx}.name`"
                                    type="text"
                                    class="w-full flex-1"
                                    :placeholder="$t('common.label', 1)"
                                />
                            </UFormGroup>

                            <UFormGroup :name="`labels.${idx}.color`">
                                <ColorPickerClient
                                    v-model="state.labels[idx]!.color"
                                    :name="`labels.${idx}.color`"
                                    class="min-w-16"
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
                        :ui="{ rounded: 'rounded-full' }"
                        :disabled="!canSubmit"
                        icon="i-mdi-plus"
                        :class="state.labels.length ? 'mt-2' : ''"
                        @click="state.labels.push({ id: 0, name: '', color: '#ffffff' })"
                    />
                </UFormGroup>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton color="black" block class="flex-1" @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton type="submit" block class="flex-1" :disabled="!canSubmit" :loading="!canSubmit">
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
