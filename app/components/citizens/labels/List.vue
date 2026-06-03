<script lang="ts" setup>
import { UButton, UIcon, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import Pagination from '~/components/partials/Pagination.vue';
import TableSortButton from '~/components/partials/TableSortButton.vue';
import { getCitizensLabelsClient } from '~~/gen/ts/clients';
import type { Label } from '~~/gen/ts/resources/citizens/labels/labels';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import DataErrorBlock from '../../partials/data/DataErrorBlock.vue';
import CreateOrUpdateModal from './CreateOrUpdateModal.vue';
import { useDraggable } from 'vue-draggable-plus';
import ReorderButtons from '~/components/partials/ReorderButtons.vue';
import DraggableHandle from '~/components/partials/DraggableHandle.vue';

const { can } = useAuth();

const { t } = useI18n();

const overlay = useOverlay();

const notifications = useNotificationsStore();

const formatDuration = useDurationFormatter();

const citizensLabelsClient = await getCitizensLabelsClient();

const {
    data: labels,
    status,
    error,
    refresh,
} = useLazyAsyncData('citizens-labels', () => listLabels(), {
    default: () => [] as Label[],
});

async function listLabels(): Promise<Label[]> {
    try {
        const { response } = await citizensLabelsClient.listLabels({
            ownJobOnly: true,
        });

        return response?.labels ?? [];
    } catch (e) {
        handleGRPCError(e as RpcError);

        return [];
    }
}

const createOrUpdateModal = overlay.create(CreateOrUpdateModal);
const deleteConfirmModal = overlay.create(ConfirmModal);

async function deleteLabel(labelId: number): Promise<void> {
    try {
        await citizensLabelsClient.deleteLabel({
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

async function reorderLabels(labels: Label[]) {
    if (!labels.length) return;

    try {
        const call = citizensLabelsClient.reorderLabels({
            labelIds: labels.map((item) => item.id),
        });
        await call;

        orderChanged.value = false;

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

const orderChanged = ref(false);
const tableRef = useTemplateRef('tableRef');
const tableBodyRef = computed<HTMLElement | null>(() => {
    const rootEl = tableRef.value?.$el as HTMLElement | undefined;
    return rootEl?.querySelector('tbody.unit-list-table') ?? null;
});

const { moveUp, moveDown } = useListReorder(labels, {
    onMove: () => (orderChanged.value = true),
});

useDraggable(tableBodyRef, labels, {
    animation: 150,
    handle: '.handle-choice',
    draggable: 'tr',
    onUpdate: () => (orderChanged.value = true),
});

const columns = computed<TableColumn<Label>[]>(() => [
    {
        id: 'actions',
        cell: ({ row }) =>
            h('div', [
                can('citizens.LabelsService/CreateOrUpdateLabel').value
                    ? h(
                          'div',
                          {
                              class: 'inline-flex items-center gap-1',
                          },
                          [
                              h(DraggableHandle, {
                                  handleClass: 'handle-choice',
                              }),
                              h(ReorderButtons, {
                                  idx: row.index,
                                  moveUp: moveUp,
                                  moveDown: moveDown,
                              }),
                          ],
                      )
                    : undefined,
                can('citizens.LabelsService/CreateOrUpdateLabel').value
                    ? h(
                          UTooltip,
                          { text: t('common.edit') },
                          h(UButton, {
                              color: 'primary',
                              variant: 'link',
                              icon: 'i-mdi-pencil',
                              onClick: () => {
                                  createOrUpdateModal.open({
                                      labelId: row.original.id,
                                      onRefresh: () => refresh(),
                                  });
                              },
                          }),
                      )
                    : undefined,
                can('citizens.LabelsService/DeleteLabel').value
                    ? h(
                          UTooltip,
                          { text: row.original.deletedAt ? t('common.restore') : t('common.delete') },
                          h(UButton, {
                              color: !row.original.deletedAt ? 'error' : 'success',
                              variant: 'link',
                              icon: !row.original.deletedAt ? 'i-mdi-delete' : 'i-mdi-restore',
                              onClick: () => {
                                  deleteConfirmModal.open({
                                      confirm: () => row.original.id && deleteLabel(row.original.id),
                                  });
                              },
                          }),
                      )
                    : undefined,
            ]),
    },
    {
        accessorKey: 'name',
        header: ({ column }) => {
            return h(TableSortButton, {
                column,
                label: t('common.name'),
            });
        },
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
    {
        accessorKey: 'expiration',
        header: t('common.expiration'),
        cell: ({ row }) =>
            h(
                'span',
                !row.original.settings?.requiresExpiration
                    ? t('common.no')
                    : h(
                          'span',
                          `${t('common.yes')} (${t('common.min')}: ${row.original.settings.minDuration ? formatDuration(row.original.settings.minDuration) : t('common.na')}, ${t('common.max')}: ${row.original.settings.maxDuration ? formatDuration(row.original.settings.maxDuration) : t('common.na')})`,
                      ),
            ),
    },
]);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('components.citizens.citizen_labels.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton to="/citizens" />

                    <UTooltip v-if="orderChanged" :text="$t('common.save', 1)">
                        <UButton
                            color="primary"
                            variant="outline"
                            icon="i-mdi-content-save"
                            @click="() => reorderLabels(labels)"
                        />
                    </UTooltip>

                    <UTooltip v-if="can('citizens.LabelsService/CreateOrUpdateLabel').value" :text="$t('common.edit')">
                        <UButton
                            color="neutral"
                            variant="outline"
                            trailing-icon="i-mdi-plus"
                            @click="
                                createOrUpdateModal.open({
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
            </UDashboardNavbar>
        </template>

        <template #body>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('components.citizens.citizen_labels.title')])"
                :error="error"
                :retry="refresh"
            />

            <UTable
                v-else
                ref="tableRef"
                class="flex-1"
                :loading="isRequestPending(status)"
                :columns="columns"
                :data="labels"
                :empty="$t('common.not_found', [$t('components.citizens.citizen_labels.title')])"
                :pagination-options="{ manualPagination: true }"
                sticky
                :ui="{ tbody: 'unit-list-table' }"
            />
        </template>

        <template #footer>
            <Pagination :status="status" :refresh="refresh" hide-buttons hide-text />
        </template>
    </UDashboardPanel>
</template>
