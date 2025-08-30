<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { VueDraggable } from 'vue-draggable-plus';
import { z } from 'zod';
import ColorPickerClient from '~/components/partials/ColorPicker.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { getJobsJobsClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetColleagueLabelsResponse, ManageLabelsResponse } from '~~/gen/ts/services/jobs/jobs';

const emit = defineEmits<{
    (e: 'close', v: boolean): void;
}>();

const notifications = useNotificationsStore();

const jobsJobsClient = await getJobsJobsClient();

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
        const { response } = await jobsJobsClient.getColleagueLabels({});

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { data: labels, status, error, refresh } = useLazyAsyncData('jobs-colleagues-labels', () => getColleagueLabels());

async function manageLabels(values: Schema): Promise<ManageLabelsResponse> {
    try {
        const { response } = await jobsJobsClient.manageLabels({
            labels: values.labels,
        });

        state.labels = response.labels;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

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

watch(labels, () => (state.labels = labels.value?.labels ?? []));

const { moveUp, moveDown } = useListReorder(toRef(state, 'labels'));

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UModal :title="$t('common.label', 2)">
        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.label', 2)])" />
                <DataErrorBlock v-else-if="error" :error="error" :retry="refresh" />

                <UFormField v-else class="grid items-center gap-2" name="list" :ui="{ container: '' }">
                    <div class="flex flex-col gap-1">
                        <VueDraggable
                            v-model="state.labels"
                            class="flex flex-col gap-2"
                            :disabled="!canSubmit"
                            handle=".handle"
                        >
                            <div v-for="(_, idx) in state.labels" :key="idx" class="flex items-center gap-1">
                                <div class="inline-flex items-center gap-1">
                                    <UTooltip :text="$t('common.draggable')">
                                        <UIcon class="handle size-6 cursor-move" name="i-mdi-drag-horizontal" />
                                    </UTooltip>

                                    <UButtonGroup>
                                        <UButton size="xs" variant="link" icon="i-mdi-arrow-up" @click="moveUp(idx)" />
                                        <UButton size="xs" variant="link" icon="i-mdi-arrow-down" @click="moveDown(idx)" />
                                    </UButtonGroup>
                                </div>

                                <UFormField class="flex-1" :name="`labels.${idx}.name`">
                                    <UInput
                                        v-model="state.labels[idx]!.name"
                                        class="w-full flex-1"
                                        :name="`labels.${idx}.name`"
                                        type="text"
                                        :placeholder="$t('common.label', 1)"
                                    />
                                </UFormField>

                                <UFormField :name="`${idx}.color`">
                                    <ColorPickerClient
                                        v-model="state.labels[idx]!.color"
                                        class="min-w-16"
                                        :name="`${idx}.color`"
                                    />
                                </UFormField>

                                <UButton :disabled="!canSubmit" icon="i-mdi-close" @click="state.labels.splice(idx, 1)" />
                            </div>
                        </VueDraggable>
                    </div>

                    <UButton
                        :class="state.labels.length ? 'mt-2' : ''"
                        :disabled="!canSubmit"
                        icon="i-mdi-plus"
                        @click="state.labels.push({ id: 0, name: '', color: '#ffffff', order: 0 })"
                    />
                </UFormField>
            </UForm>
        </template>

        <template #footer>
            <UButtonGroup class="inline-flex w-full">
                <UButton class="flex-1" color="neutral" block :label="$t('common.close', 1)" @click="$emit('close', false)" />

                <UButton
                    class="flex-1"
                    block
                    :disabled="!canSubmit || !!error"
                    :loading="isRequestPending(status) || !canSubmit"
                    :label="$t('common.save')"
                    @click="formRef?.submit()"
                />
            </UButtonGroup>
        </template>
    </UModal>
</template>
