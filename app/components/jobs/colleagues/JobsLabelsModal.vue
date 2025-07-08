<script lang="ts" setup>
import type { FormSubmitEvent } from '#ui/types';
import { VueDraggable } from 'vue-draggable-plus';
import { z } from 'zod';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetColleagueLabelsResponse, ManageLabelsResponse } from '~~/gen/ts/services/jobs/jobs';

const { $grpc } = useNuxtApp();

const { isOpen } = useModal();

const notifications = useNotificationsStore();

const schema = z.object({
    labels: z
        .object({
            id: z.coerce.number(),
            name: z.string().min(1).max(64),
            color: z.string().length(7),
            order: z.coerce.number().nonnegative().default(0),
        })
        .array()
        .max(15)
        .default([]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    labels: [],
});

async function getColleagueLabels(): Promise<GetColleagueLabelsResponse> {
    try {
        const { response } = await $grpc.jobs.jobs.getColleagueLabels({});

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const {
    data: labels,
    pending: loading,
    error,
    refresh,
} = useLazyAsyncData('jobs-colleagues-labels', () => getColleagueLabels());

async function manageLabels(values: Schema): Promise<ManageLabelsResponse> {
    try {
        const { response } = await $grpc.jobs.jobs.manageLabels({
            labels: values.labels,
        });

        state.labels = response.labels;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

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

watch(labels, () => (state.labels = labels.value?.labels ?? []));
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

                        <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                    </div>
                </template>

                <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.label', 2)])" />
                <DataErrorBlock v-else-if="error" :error="error" :retry="refresh" />

                <UFormGroup v-else class="grid items-center gap-2" name="list" :ui="{ container: '' }">
                    <div class="flex flex-col gap-1">
                        <VueDraggable v-model="state.labels" class="flex flex-col gap-2">
                            <div v-for="(_, idx) in state.labels" :key="idx" class="flex items-center gap-1">
                                <UTooltip :text="$t('common.draggable')">
                                    <UIcon class="size-6" name="i-mdi-drag-horizontal" />
                                </UTooltip>

                                <UFormGroup class="flex-1" :name="`labels.${idx}.name`">
                                    <UInput
                                        v-model="state.labels[idx]!.name"
                                        class="w-full flex-1"
                                        :name="`labels.${idx}.name`"
                                        type="text"
                                        :placeholder="$t('common.label', 1)"
                                    />
                                </UFormGroup>

                                <UFormGroup :name="`${idx}.color`">
                                    <ColorPickerClient
                                        v-model="state.labels[idx]!.color"
                                        class="min-w-16"
                                        :name="`${idx}.color`"
                                    />
                                </UFormGroup>

                                <UButton
                                    :ui="{ rounded: 'rounded-full' }"
                                    :disabled="!canSubmit"
                                    icon="i-mdi-close"
                                    @click="state.labels.splice(idx, 1)"
                                />
                            </div>
                        </VueDraggable>
                    </div>

                    <UButton
                        :class="state.labels.length ? 'mt-2' : ''"
                        :ui="{ rounded: 'rounded-full' }"
                        :disabled="!canSubmit"
                        icon="i-mdi-plus"
                        @click="state.labels.push({ id: 0, name: '', color: '#ffffff', order: 0 })"
                    />
                </UFormGroup>

                <template #footer>
                    <UButtonGroup class="inline-flex w-full">
                        <UButton class="flex-1" color="black" block @click="isOpen = false">
                            {{ $t('common.close', 1) }}
                        </UButton>

                        <UButton
                            class="flex-1"
                            type="submit"
                            block
                            :loading="loading || !canSubmit"
                            :disabled="!canSubmit || !!error"
                        >
                            {{ $t('common.save') }}
                        </UButton>
                    </UButtonGroup>
                </template>
            </UCard>
        </UForm>
    </UModal>
</template>
