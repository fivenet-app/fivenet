<script lang="ts" setup>
import { UButton, UIcon, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DraggableHandle from '~/components/partials/DraggableHandle.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import ReorderButtons from '~/components/partials/ReorderButtons.vue';
import StatsModalClient from '~/components/jobs/colleagues/labels/StatsModal.client.vue';
import { getJobsColleaguesClient } from '~~/gen/ts/clients';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { Label } from '~~/gen/ts/resources/jobs/labels/labels';
import CreateOrUpdateDrawer from './CreateOrUpdateDrawer.vue';
import { useDraggable } from 'vue-draggable-plus';

const notifications = useNotificationsStore();

const { t } = useI18n();

const overlay = useOverlay();

const { attr, can } = useAuth();

const jobsColleaguesClient = await getJobsColleaguesClient();

const canCreateOrUpdateLabel = computed(() => can('jobs.ColleaguesService/CreateOrUpdateLabel').value);

const {
    data: labels,
    status,
    error,
    refresh,
} = useLazyAsyncData('jobs-colleagues-labels', () => getColleagueLabels(), {
    default: () => [] as Label[],
});

async function getColleagueLabels(): Promise<Label[]> {
    try {
        const { response } = await jobsColleaguesClient.getColleagueLabels({});

        return response?.labels ?? [];
    } catch (e) {
        handleGRPCError(e as RpcError);

        return [];
    }
}

const createOrUpdateDrawer = overlay.create(CreateOrUpdateDrawer);
const deleteConfirmModal = overlay.create(ConfirmModal);
const statsModal = overlay.create(StatsModalClient);

async function deleteLabel(labelId: number): Promise<void> {
    try {
        await jobsColleaguesClient.deleteLabel({
            id: labelId,
        });

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });

        refresh();
    } catch (e) {
        handleGRPCError(e as RpcError);
    }
}

async function reorderLabels(currentLabels: Label[]) {
    if (!currentLabels.length || !canCreateOrUpdateLabel.value) return;

    try {
        await jobsColleaguesClient.reorderLabels({
            labelIds: currentLabels.map((item) => item.id),
        });

        syncSnapshot();

        notifications.add({
            title: { key: 'notifications.action_successful.title', parameters: {} },
            description: { key: 'notifications.action_successful.content', parameters: {} },
            type: NotificationType.SUCCESS,
        });
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { snapshotDirty: orderChanged, syncSnapshot } = useSnapshotChanges(() => labels.value?.map((label) => label.id) ?? []);
const tableRef = useTemplateRef('tableRef');
const tableBodyRef = computed<HTMLElement | null>(() => {
    const rootEl = tableRef.value?.$el as HTMLElement | undefined;
    return rootEl?.querySelector('tbody.jobs-label-list-table') ?? null;
});

const { moveUp, moveDown } = useListReorder(labels, {
    onMove: () => undefined,
});

useDraggable(tableBodyRef, labels, {
    animation: 150,
    handle: '.handle',
    draggable: 'tr',
    disabled: !canCreateOrUpdateLabel.value,
    onUpdate: () => undefined,
});

watch(
    status,
    (newStatus) => {
        if (!isRequestPending(newStatus)) {
            syncSnapshot();
        }
    },
    { immediate: true },
);

const columns = computed(
    () =>
        [
            {
                id: 'actions',
                cell: ({ row }) =>
                    h(
                        'div',
                        canCreateOrUpdateLabel.value
                            ? [
                                  h(
                                      'div',
                                      {
                                          class: 'inline-flex items-center gap-1',
                                      },
                                      [
                                          h(DraggableHandle),
                                          h(ReorderButtons, {
                                              idx: row.index,
                                              moveUp: moveUp,
                                              moveDown: moveDown,
                                          }),
                                      ],
                                  ),
                                  h(
                                      UTooltip,
                                      { text: t('common.edit') },
                                      h(UButton, {
                                          color: 'primary',
                                          variant: 'link',
                                          icon: 'i-mdi-pencil',
                                          onClick: () => {
                                              createOrUpdateDrawer.open({
                                                  label: row.original,
                                                  onRefresh: () => refresh(),
                                              });
                                          },
                                      }),
                                  ),

                                  h(
                                      UTooltip,
                                      { text: t('common.delete') },
                                      h(UButton, {
                                          color: 'error',
                                          variant: 'link',
                                          icon: 'i-mdi-delete',
                                          onClick: () => {
                                              deleteConfirmModal.open({
                                                  confirm: () => row.original.id && deleteLabel(row.original.id),
                                              });
                                          },
                                      }),
                                  ),
                              ]
                            : [],
                    ),
            },
            {
                accessorKey: 'name',
                header: t('common.name'),
                meta: {
                    class: {
                        td: 'text-highlighted',
                    },
                },
            },
            {
                accessorKey: 'color',
                header: t('common.color'),
                cell: ({ row }) =>
                    h(ColorPicker, {
                        modelValue: row.original.color,
                        disabled: true,
                        hideLabel: true,
                    }),
            },
            {
                accessorKey: 'icon',
                header: t('common.icon'),
                cell: ({ row }) =>
                    row.original.icon
                        ? h(UIcon, {
                              class: 'size-5',
                              name: convertComponentIconNameToDynamic(row.original.icon),
                              style: {
                                  color: row.original.color ?? 'currentColor',
                              },
                          })
                        : undefined,
            },
        ] as TableColumn<Label>[],
);

const breadcrumbs = computed(() => [
    {
        label: t('pages.jobs.colleagues.title'),
        icon: 'i-mdi-account-group',
        to: '/jobs/colleagues',
    },
    {
        label: t('pages.jobs.colleagues.labels.title'),
        icon: 'i-mdi-tag',
    },
]);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardToolbar>
                <template #left>
                    <UBreadcrumb :items="breadcrumbs" />
                </template>

                <template #right>
                    <UTooltip v-if="orderChanged && canCreateOrUpdateLabel" :text="$t('common.save', 1)">
                        <UButton
                            color="primary"
                            variant="outline"
                            icon="i-mdi-content-save"
                            @click="() => reorderLabels(labels ?? [])"
                        />
                    </UTooltip>

                    <UTooltip v-if="canCreateOrUpdateLabel" :text="$t('common.create')">
                        <UButton
                            color="neutral"
                            variant="outline"
                            trailing-icon="i-mdi-plus"
                            @click="
                                createOrUpdateDrawer.open({
                                    onRefresh: () => refresh(),
                                })
                            "
                        >
                            <span class="hidden truncate sm:block">
                                {{ $t('common.label', 1) }}
                            </span>
                        </UButton>
                    </UTooltip>
                </template>
            </UDashboardToolbar>

            <UDashboardToolbar>
                <template #right>
                    <UTooltip
                        v-if="attr('jobs.ColleaguesService/GetColleague', 'Types', 'Labels').value"
                        :text="$t('common.stats')"
                    >
                        <UButton icon="i-mdi-chart-donut" color="neutral" @click="statsModal.open({})" />
                    </UTooltip>

                    <RefreshButton @click="() => refresh()" />
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataErrorBlock v-if="error" :error="error" :retry="refresh" />

            <UTable
                v-else
                ref="tableRef"
                class="flex-1"
                :columns="columns"
                :data="labels"
                :empty="$t('common.not_found', [$t('pages.jobs.colleagues.labels.title')])"
                :pagination-options="{ manualPagination: true }"
                sticky
                :ui="{ tbody: 'jobs-label-list-table' }"
            />
        </template>
    </UDashboardPanel>
</template>
