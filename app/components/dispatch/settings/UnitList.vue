<script lang="ts" setup>
import { UButton, UFieldGroup, UIcon, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { h } from 'vue';
import UnitAttributes from '~/components/dispatch/partials/UnitAttributes.vue';
import UnitCreateOrUpdateModal from '~/components/dispatch/settings/UnitCreateOrUpdateModal.vue';
import ColorPicker from '~/components/partials/ColorPicker.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import { fallbackIconName, resolveIconComponent } from '~/components/partials/icons';
import Pagination from '~/components/partials/Pagination.vue';
import { getCentrumUnitsClient } from '~~/gen/ts/clients';
import type { Unit } from '~~/gen/ts/resources/centrum/units/units';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { useDraggable } from 'vue-draggable-plus';

const { t } = useI18n();

const { can } = useAuth();

const overlay = useOverlay();

const notifications = useNotificationsStore();

const appConfig = useAppConfig();

const centrumUnitsClient = await getCentrumUnitsClient();

const {
    data: units,
    status,
    refresh,
    error,
} = useLazyAsyncData('centrum-units', () => listUnits(), {
    default: () => [] as Unit[],
});

async function listUnits(): Promise<Unit[]> {
    try {
        const call = centrumUnitsClient.listUnits({
            status: [],
        });
        const { response } = await call;

        return response.units;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function deleteUnit(id: number): Promise<void> {
    try {
        const call = centrumUnitsClient.deleteUnit({
            unitId: id,
        });
        await call;

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

async function reorderUnits(units: Unit[]): Promise<void> {
    if (!units.length) return;

    try {
        const call = centrumUnitsClient.reorderUnits({
            unitIds: units.map((item) => item.id),
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

const { moveUp, moveDown } = useListReorder(units, {
    onMove: () => (orderChanged.value = true),
});

useDraggable(tableBodyRef, units, {
    animation: 150,
    handle: '.handle-choice',
    draggable: 'tr',
    onUpdate: () => (orderChanged.value = true),
});

const columns = computed<TableColumn<Unit>[]>(() => [
    {
        id: 'actions',
        cell: ({ row }) =>
            h('div', [
                h(
                    'div',
                    {
                        class: 'inline-flex items-center gap-1',
                    },
                    [
                        h(UTooltip, { text: t('common.draggable') }, [
                            h(UIcon, {
                                class: 'handle-choice size-6 cursor-move',
                                name: 'i-mdi-drag-horizontal',
                            }),
                        ]),
                        h(UFieldGroup, { orientation: 'vertical' }, [
                            h(UButton, {
                                size: 'xs',
                                variant: 'link',
                                icon: 'i-mdi-arrow-up',
                                onClick: () => moveUp(row.index),
                            }),
                            h(UButton, {
                                size: 'xs',
                                variant: 'link',
                                icon: 'i-mdi-arrow-down',
                                onClick: () => moveDown(row.index),
                            }),
                        ]),
                    ],
                ),
                h(
                    UTooltip,
                    {
                        text: t('common.update'),
                        vIf: can('centrum.UnitsService/CreateOrUpdateUnit').value,
                    },
                    [
                        h(UButton, {
                            variant: 'link',
                            icon: 'i-mdi-pencil',
                            onClick: () => {
                                unitCreateOrUpdate.open({
                                    unit: row.original,
                                    onUpdated: async () => refresh(),
                                });
                            },
                        }),
                    ],
                ),
                h(
                    UTooltip,
                    {
                        text: t('common.delete'),
                        vIf: can('centrum.UnitsService/DeleteUnit').value,
                    },
                    [
                        h(UButton, {
                            variant: 'link',
                            icon: 'i-mdi-delete',
                            color: 'error',
                            onClick: () => {
                                confirmModal.open({
                                    confirm: async () => deleteUnit(row.original.id),
                                });
                            },
                        }),
                    ],
                ),
            ]),
    },
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
        cell: ({ row }) => h('span', { class: 'text-highlighted' }, row.original.name),
    },
    {
        accessorKey: 'initials',
        header: ({ column }) => {
            const isSorted = column.getIsSorted();

            return h(UButton, {
                color: 'neutral',
                variant: 'ghost',
                label: t('common.initials'),
                icon: isSorted
                    ? isSorted === 'asc'
                        ? appConfig.custom.icons.sortAsc
                        : appConfig.custom.icons.sortDesc
                    : appConfig.custom.icons.sort,
                class: '-mx-2.5',
                onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
            });
        },
        cell: ({ row }) => h('span', {}, row.original.initials),
    },
    {
        accessorKey: 'description',
        header: t('common.description'),
        cell: ({ row }) => h('span', {}, row.original.description),
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
        accessorKey: 'attributes',
        header: t('common.attributes', 2),
        cell: ({ row }) => h(UnitAttributes, { attributes: row.original.attributes }),
    },
    {
        accessorKey: 'homePostal',
        header: t('common.department_postal'),
        cell: ({ row }) => h('span', {}, row.original.homePostal ?? t('common.na')),
    },
]);

const unitCreateOrUpdate = overlay.create(UnitCreateOrUpdateModal);
const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('common.unit', 2)">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/dispatch" />

                    <UTooltip v-if="orderChanged" :text="$t('common.save', 1)">
                        <UButton
                            color="primary"
                            variant="outline"
                            icon="i-mdi-content-save"
                            @click="() => reorderUnits(units)"
                        />
                    </UTooltip>

                    <UTooltip v-if="can('centrum.CentrumService/Stream').value" :text="$t('common.setting', 2)">
                        <UButton icon="i-mdi-settings" to="/dispatch/settings">
                            <span class="hidden truncate sm:block">
                                {{ $t('common.setting', 2) }}
                            </span>
                        </UButton>
                    </UTooltip>

                    <UButton
                        v-if="can('centrum.UnitsService/CreateOrUpdateUnit').value"
                        color="neutral"
                        variant="outline"
                        trailing-icon="i-mdi-plus"
                        @click="
                            unitCreateOrUpdate.open({
                                onCreated: async () => refresh(),
                                onUpdated: async () => refresh(),
                            })
                        "
                    >
                        <span class="hidden truncate sm:block">
                            {{ $t('components.centrum.units.create_unit') }}
                        </span>
                    </UButton>
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.unable_to_load', [$t('common.unit', 2)])"
                :error="error"
                :retry="refresh"
            />

            <UTable
                v-else
                ref="tableRef"
                class="flex-1"
                :loading="isRequestPending(status)"
                :columns="columns"
                :data="units"
                :empty="$t('common.not_found', [$t('common.unit', 2)])"
                :pagination-options="{ manualPagination: true }"
                sticky
                :ui="{ tbody: 'unit-list-table' }"
            />
        </template>

        <template #footer>
            <Pagination hide-buttons hide-text :status="status" :refresh="refresh" />
        </template>
    </UDashboardPanel>
</template>
