<script lang="ts" setup>
import { UButton, UTooltip } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { h } from 'vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/stores/livemap';
import { getLivemapLivemapClient } from '~~/gen/ts/clients';
import { type MarkerMarker, MarkerType } from '~~/gen/ts/resources/livemap/markers/marker_marker';
import MarkerCreateOrUpdateSlideover from '../MarkerCreateOrUpdateSlideover.vue';

const emits = defineEmits<{
    (e: 'editing', editing: boolean): void;
}>();

const { t } = useI18n();

const { can } = useAuth();

const livemapStore = useLivemapStore();
const { deleteMarkerMarker, gotoCoords } = livemapStore;
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

const editingMarker = ref(false);

const markerCreateOrUpdateSlideover = overlay.create(MarkerCreateOrUpdateSlideover);
const confirmModal = overlay.create(ConfirmModal);

const columns = computed(
    () =>
        [
            {
                id: 'actions',
                cell: ({ row }) =>
                    h('div', [
                        h(UTooltip, { text: t('common.mark') }, () =>
                            h(UButton, {
                                variant: 'link',
                                icon: 'i-mdi-map-marker',
                                onClick: () => gotoCoords({ x: row.original.x, y: row.original.y }),
                            }),
                        ),
                        can('livemap.LivemapService/CreateOrUpdateMarker').value
                            ? h(UTooltip, { text: t('common.edit') }, () =>
                                  h(UButton, {
                                      variant: 'link',
                                      icon: 'i-mdi-pencil',
                                      onClick: () => {
                                          emits('editing', true);
                                          markerCreateOrUpdateSlideover
                                              .open({
                                                  marker: row.original,
                                              })
                                              .finally(() => emits('editing', false));
                                      },
                                  }),
                              )
                            : null,
                        can('livemap.LivemapService/DeleteMarker').value
                            ? h(UTooltip, { text: t('common.delete') }, () =>
                                  h(UButton, {
                                      variant: 'link',
                                      icon: 'i-mdi-delete',
                                      color: 'error',
                                      onClick: () => {
                                          confirmModal.open({
                                              confirm: async () => deleteMarker(row.original.id),
                                          });
                                      },
                                  }),
                              )
                            : null,
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
        ] as TableColumn<MarkerMarker>[],
);

const querySearchRaw = ref('');
const querySearch = computed(() => querySearchRaw.value.trim().toLowerCase());

const filteredMarkers = computed(() =>
    Array.from(markersMarkers.value.values())
        .filter((marker) => querySearch.value === '' || marker.name.toLowerCase().includes(querySearch.value))
        .sort((a, b) => a.name.localeCompare(b.name)),
);
</script>

<template>
    <div class="flex h-full min-h-0 flex-col gap-2">
        <UFormField name="search">
            <UInput
                v-model="querySearchRaw"
                class="w-full"
                type="text"
                name="search"
                :placeholder="$t('common.search')"
                :ui="{ trailing: 'pe-1' }"
            >
                <template #trailing>
                    <UButton
                        v-if="querySearchRaw !== ''"
                        color="red"
                        variant="link"
                        icon="i-mdi-clear"
                        aria-controls="search"
                        @click="querySearchRaw = ''"
                    />
                </template>
            </UInput>
        </UFormField>

        <UTable
            class="min-h-0 flex-1"
            :columns="columns"
            :data="filteredMarkers"
            :empty="$t('common.not_found', [$t('common.marker', 2)])"
            :pagination-options="{ manualPagination: true }"
            :sorting-options="{ manualSorting: true }"
            sticky
            :virtualize="{ estimateSize: 40 }"
        />
    </div>
</template>
