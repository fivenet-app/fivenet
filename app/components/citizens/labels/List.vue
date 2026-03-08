<script lang="ts" setup>
import { UButton } from '#components';
import type { TableColumn } from '@nuxt/ui';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useCompletorStore } from '~/stores/completor';
import type { Label } from '~~/gen/ts/resources/users/labels/labels';
import DataErrorBlock from '../../partials/data/DataErrorBlock.vue';
import DataPendingBlock from '../../partials/data/DataPendingBlock.vue';
import CreateOrUpdateModal from './CreateOrUpdateModal.vue';

const { can } = useAuth();

const { t } = useI18n();

const overlay = useOverlay();

const completorStore = useCompletorStore();

const {
    data: labels,
    status,
    error,
    refresh,
} = useLazyAsyncData('citizens-labels', () => completorStore.completeCitizenLabels(''));

const createOrUpdateModal = overlay.create(CreateOrUpdateModal);

const appConfig = useAppConfig();

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

                    <UTooltip v-if="can('citizens.CitizensService/ManageLabels').value" :text="$t('common.edit')">
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
