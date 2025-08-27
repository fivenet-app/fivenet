<script lang="ts" setup>
import type { TableColumn } from '@nuxt/ui';
import { h } from 'vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/stores/livemap';
import { getLivemapLivemapClient } from '~~/gen/ts/clients';
import { type MarkerMarker, MarkerType } from '~~/gen/ts/resources/livemap/marker_marker';

const { t } = useI18n();

const { can } = useAuth();

const livemapStore = useLivemapStore();
const { deleteMarkerMarker, goto } = livemapStore;
const { markersMarkers } = storeToRefs(livemapStore);

const overlay = useOverlay();

const livemapLivemapClient = await getLivemapLivemapClient();

async function deleteMarker(id: number): Promise<void> {
    try {
        const call = livemapLivemapClient.deleteMarker({
            id,
        });
        await call;

        deleteMarkerMarker(id);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const columns = computed(
    () =>
        [
            {
                id: 'actions',
                cell: ({ row }) =>
                    h('div', [
                        h('UTooltip', { text: t('common.mark') }, () =>
                            h('UButton', {
                                variant: 'link',
                                icon: 'i-mdi-map-marker',
                                onClick: () => goto({ x: row.original.x, y: row.original.y }),
                            }),
                        ),
                        h('UTooltip', { text: t('common.delete') }, () =>
                            can('livemap.LivemapService/DeleteMarker').value
                                ? h('UButton', {
                                      variant: 'link',
                                      icon: 'i-mdi-delete',
                                      color: 'error',
                                      onClick: () =>
                                          confirmModal.open({
                                              confirm: async () => deleteMarker(row.original.id),
                                          }),
                                  })
                                : null,
                        ),
                    ]),
            },
            {
                accessorKey: 'createdAt',
                header: t('common.created'),
                cell: ({ row }) => h(GenericTime, { value: row.original.createdAt, type: 'compact' }),
            },
            {
                accessorKey: 'expiresAt',
                header: t('common.expires_at'),
                cell: ({ row }) =>
                    row.original.expiresAt
                        ? h(GenericTime, { value: row.original.expiresAt, type: 'compact' })
                        : t('common.na'),
            },
            {
                accessorKey: 'name',
                header: t('common.name'),
                cell: ({ row }) => row.original.name,
            },
            {
                accessorKey: 'type',
                header: t('common.type'),
                cell: ({ row }) => t(`enums.livemap.MarkerType.${MarkerType[row.original.type]}`),
            },
            {
                accessorKey: 'description',
                header: t('common.description'),
                cell: ({ row }) =>
                    h(
                        'p',
                        { class: 'max-h-14 truncate overflow-y-scroll break-words' },
                        row.original.description ?? t('common.na'),
                    ),
            },
            {
                accessorKey: 'creator',
                header: t('common.creator'),
                cell: ({ row }) =>
                    row.original.creator
                        ? h(CitizenInfoPopover, { user: row.original.creator, trailing: false })
                        : t('common.unknown'),
            },
            {
                accessorKey: 'job',
                header: t('common.job'),
                cell: ({ row }) => row.original.creator?.jobLabel ?? t('common.na'),
            },
        ] as TableColumn<MarkerMarker>[],
);

const confirmModal = overlay.create(ConfirmModal);
</script>

<template>
    <div class="flex h-full grow flex-col px-1">
        <div class="flex justify-between">
            <h2 class="inline-flex flex-1 items-center text-base leading-6 font-semibold text-gray-100">
                {{ $t('common.marker', 2) }}
            </h2>

            <h2 class="text-base font-semibold text-gray-100">
                {{ $t('common.count') }}:
                {{ [...markersMarkers.values()].length }}
            </h2>
        </div>

        <div class="flex flex-1 flex-col overflow-x-auto overflow-y-auto">
            <UTable
                class="overflow-x-visible"
                :columns="columns"
                :data="Array.from(markersMarkers.values())"
                :empty="$t('common.not_found', [$t('common.marker', 2)])"
                :pagination-options="{ manualPagination: true }"
                :sorting-options="{ manualSorting: true }"
            />

            <div class="flex-1" />
        </div>
    </div>
</template>
