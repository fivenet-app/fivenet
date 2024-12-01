<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { z } from 'zod';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import type { GetColleagueLabelsResponse, ManageColleagueLabelsResponse } from '~~/gen/ts/services/jobs/jobs';

const { isOpen } = useModal();

const schema = z
    .object({
        id: z.string(),
        name: z.string().min(1).max(64),
        color: z.string().length(7),
    })
    .array()
    .max(15);

type Schema = z.output<typeof schema>;

const state = ref<Schema>([]);

async function getColleagueLabels(): Promise<GetColleagueLabelsResponse> {
    try {
        const { response } = await getGRPCJobsClient().getColleagueLabels({});

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { data: labels } = useLazyAsyncData('jobs-colleagues-labels', () => getColleagueLabels());

async function manageColleagueLabels(values: Schema): Promise<ManageColleagueLabelsResponse> {
    try {
        const { response } = await getGRPCJobsClient().manageColleagueLabels({
            labels: values,
        });

        state.value = response.labels;

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
    await manageColleagueLabels(event.data).finally(() => useTimeoutFn(() => (canSubmit.value = true), 400));
}, 1000);

watch(labels, () => (state.value = labels.value?.labels ?? []));
</script>

<template>
    <UModal :ui="{ width: 'w-full sm:max-w-5xl' }">
        <UForm :schema="schema" :state="state" @submit="onSubmitThrottle">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('common.label', 2) }}
                        </h3>

                        <UButton color="gray" variant="ghost" icon="i-mdi-window-close" class="-my-1" @click="isOpen = false" />
                    </div>
                </template>

                <UFormGroup name="list" class="grid items-center gap-2" :ui="{ container: '' }">
                    <div class="flex flex-col gap-1">
                        <div v-for="(_, idx) in state" :key="idx" class="flex items-center gap-1">
                            <UFormGroup :name="`${idx}.name`" class="flex-1">
                                <UInput
                                    v-model="state[idx]!.name"
                                    :name="`${idx}.name`"
                                    type="text"
                                    class="w-full flex-1"
                                    :placeholder="$t('common.label', 1)"
                                />
                            </UFormGroup>

                            <UFormGroup :name="`${idx}.color`">
                                <ColorPickerClient v-model="state[idx]!.color" :name="`${idx}.color`" class="min-w-16" />
                            </UFormGroup>

                            <UButton
                                :ui="{ rounded: 'rounded-full' }"
                                :disabled="!canSubmit"
                                icon="i-mdi-close"
                                @click="state.splice(idx, 1)"
                            />
                        </div>
                    </div>

                    <UButton
                        :ui="{ rounded: 'rounded-full' }"
                        :disabled="!canSubmit"
                        icon="i-mdi-plus"
                        :class="state.length ? 'mt-2' : ''"
                        @click="state.push({ id: '0', name: '', color: '#ffffff' })"
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
