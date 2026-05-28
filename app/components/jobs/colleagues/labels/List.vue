<script lang="ts" setup>
import type { FormSubmitEvent } from '@nuxt/ui';
import { VueDraggable } from 'vue-draggable-plus';
import { z } from 'zod';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import IconSelectMenu from '~/components/partials/IconSelectMenu.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import { getJobsColleaguesClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { GetColleagueLabelsResponse, ManageLabelsResponse } from '~~/gen/ts/services/jobs/colleagues';

const notifications = useNotificationsStore();

const { t } = useI18n();

const jobsColleaguesClient = await getJobsColleaguesClient();

const schema = z.object({
    labels: z
        .object({
            id: z.coerce.number(),
            name: z.coerce.string().min(1).max(64),
            color: z.coerce.string().length(7),
            icon: z.coerce.string().max(128).optional(),
            sortOrder: z.coerce.number().nonnegative().default(0),
        })
        .array()
        .max(50)
        .default([]),
});

type Schema = z.output<typeof schema>;

const state = reactive<Schema>({
    labels: [],
});

async function getColleagueLabels(): Promise<GetColleagueLabelsResponse> {
    try {
        const { response } = await jobsColleaguesClient.getColleagueLabels({});

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { data: labels, status, error, refresh } = useLazyAsyncData('jobs-colleagues-labels', () => getColleagueLabels());

async function manageLabels(values: Schema): Promise<ManageLabelsResponse> {
    try {
        const { response } = await jobsColleaguesClient.manageLabels({
            labels: values.labels,
        });

        state.labels = response.labels;

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

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

const breadcrumbs = computed(() => [
    {
        label: t('pages.jobs.colleagues.title'),
        icon: 'i-mdi-account-group',
        to: '/jobs/colleagues',
    },
    {
        label: t('pages.jobs.colleagues.labels.title'),
        icon: 'i-mdi-label',
    },
]);

const formRef = useTemplateRef('formRef');
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0' }">
        <template #header>
            <UDashboardToolbar>
                <template #left>
                    <UBreadcrumb :items="breadcrumbs" />
                </template>
            </UDashboardToolbar>

            <UDashboardToolbar>
                <template #left>
                    <UButton
                        block
                        :disabled="!canSubmit || !!error"
                        icon="i-mdi-content-save"
                        :loading="isRequestPending(status) || !canSubmit"
                        :label="$t('common.save')"
                        @click="formRef?.submit()"
                    />
                </template>

                <template #right>
                    <RefreshButton @click="() => refresh()" />
                    <UButton
                        to="/jobs/colleagues"
                        icon="i-mdi-arrow-left"
                        variant="subtle"
                        :label="$t('pages.jobs.colleagues.title')"
                    />
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <UForm ref="formRef" :schema="schema" :state="state" @submit="onSubmitThrottle">
                <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.label', 2)])" />
                <DataErrorBlock v-else-if="error" :error="error" :retry="refresh" />

                <UFormField v-else class="grid items-center gap-2" name="labels">
                    <div class="flex flex-col gap-1">
                        <VueDraggable
                            v-model="state.labels"
                            class="flex flex-col gap-2 divide-y divide-default"
                            :disabled="!canSubmit"
                            handle=".handle"
                        >
                            <div v-for="(_, idx) in state.labels" :key="idx" class="flex items-center gap-1 pb-2">
                                <div class="inline-flex items-center gap-1">
                                    <UTooltip :text="$t('common.draggable')">
                                        <UIcon class="handle size-6 cursor-move" name="i-mdi-drag-horizontal" />
                                    </UTooltip>

                                    <UFieldGroup orientation="vertical">
                                        <UButton size="xs" variant="link" icon="i-mdi-arrow-up" @click="moveUp(idx)" />
                                        <UButton size="xs" variant="link" icon="i-mdi-arrow-down" @click="moveDown(idx)" />
                                    </UFieldGroup>
                                </div>

                                <div class="flex flex-1 flex-col gap-1">
                                    <UFormField class="flex-1" :name="`labels.${idx}.name`" :label="$t('common.label', 1)">
                                        <UInput
                                            v-model="state.labels[idx]!.name"
                                            class="w-full flex-1"
                                            :name="`labels.${idx}.name`"
                                            type="text"
                                            :placeholder="$t('common.label', 1)"
                                        />
                                    </UFormField>

                                    <div class="flex flex-1 flex-row gap-2">
                                        <UFormField class="flex-1" :name="`labels.${idx}.color`" :label="$t('common.color')">
                                            <ColorPicker
                                                v-model="state.labels[idx]!.color"
                                                class="w-full"
                                                :name="`labels.${idx}.color`"
                                            />
                                        </UFormField>

                                        <UFormField class="flex-1" :name="`labels.${idx}.icon`" :label="$t('common.icon')">
                                            <IconSelectMenu
                                                v-model="state.labels[idx]!.icon"
                                                class="w-full"
                                                :name="`labels.${idx}.icon`"
                                                :hex-color="state.labels[idx]!.color"
                                                clear
                                            />
                                        </UFormField>
                                    </div>
                                </div>

                                <UButton
                                    color="red"
                                    :disabled="!canSubmit"
                                    icon="i-mdi-remove"
                                    @click="state.labels.splice(idx, 1)"
                                />
                            </div>
                        </VueDraggable>
                    </div>

                    <UButton
                        :class="state.labels.length ? 'mt-2' : ''"
                        :disabled="!canSubmit"
                        icon="i-mdi-plus"
                        @click="state.labels.push({ id: 0, name: '', color: '#5c7aff', sortOrder: 0 })"
                    />
                </UFormField>
            </UForm>
        </template>
    </UDashboardPanel>
</template>
