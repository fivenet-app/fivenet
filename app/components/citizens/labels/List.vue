<script lang="ts" setup>
import { UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import { availableIcons } from '~/components/partials/icons';
import Pagination from '~/components/partials/Pagination.vue';
import { getCitizensLabelsClient } from '~~/gen/ts/clients';
import type { Label } from '~~/gen/ts/resources/citizens/labels/labels';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import DataErrorBlock from '../../partials/data/DataErrorBlock.vue';
import DataPendingBlock from '../../partials/data/DataPendingBlock.vue';
import CreateOrUpdateModal from './CreateOrUpdateModal.vue';

const { can } = useAuth();

const { t } = useI18n();

const overlay = useOverlay();

const notifications = useNotificationsStore();

const appConfig = useAppConfig();

const citizensLabelsClient = await getCitizensLabelsClient();

const { data: labels, status, error, refresh } = useLazyAsyncData('citizens-labels', () => listLabels());

async function listLabels(): Promise<Label[]> {
    try {
        const { response } = await citizensLabelsClient.listLabels({});

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

const columns = computed<TableColumn<Label>[]>(() => [
    {
        accessorKey: 'name',
        header: ({ column }) => {
            const isSorted = column.getIsSorted();

            return h(UButton, {
                color: 'neutral',
                variant: 'ghost',
                label: t('common.name'),
                icon: isSorted
                    ? isSorted === 'asc'
                        ? appConfig.custom.icons.sortAsc
                        : appConfig.custom.icons.sortDesc
                    : appConfig.custom.icons.sort,
                class: '-mx-2.5',
                onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
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
            h(availableIcons.find((item) => item.name === row.original.icon)?.component ?? 'span', {
                class: 'size-5',
                fill: row.original.color ?? 'currentColor',
            }),
    },
    {
        id: 'actions',
        header: '',
        cell: ({ row }) =>
            h('div', [
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
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.label', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('components.citizens.citizen_labels.title')])"
                :error="error"
                :retry="refresh"
            />

            <UTable
                v-else
                class="flex-1"
                :loading="isRequestPending(status)"
                :columns="columns"
                :data="labels"
                :empty="$t('common.not_found', [$t('common.label', 2)])"
                sticky
            />
        </template>

        <template #footer>
            <Pagination :status="status" :refresh="refresh" hide-buttons hide-text />
        </template>
    </UDashboardPanel>
</template>
