<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/stores/livemap';
import { getLivemapLivemapClient } from '~~/gen/ts/clients';
import { MarkerType } from '~~/gen/ts/resources/livemap/marker_marker';

const { t } = useI18n();

const { can } = useAuth();

const livemapStore = useLivemapStore();
const { deleteMarkerMarker, goto } = livemapStore;
const { markersMarkers } = storeToRefs(livemapStore);

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

const modal = useOverlay();

const columns = [
    {
        key: 'actions',
        label: t('common.action', 2),
        sortable: false,
    },
    {
        key: 'createdAt',
        label: t('common.created'),
    },
    {
        key: 'expiresAt',
        label: t('common.expires_at'),
    },
    {
        key: 'name',
        label: t('common.name'),
    },
    {
        key: 'type',
        label: t('common.type'),
    },
    {
        key: 'description',
        label: t('common.description'),
    },
    {
        key: 'creator',
        label: t('common.creator'),
    },
    {
        key: 'job',
        label: t('common.job'),
    },
];
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
                :rows="Array.from(markersMarkers.values())"
                :empty-state="{
                    icon: 'i-mdi-map-marker',
                    label: $t('common.not_found', [$t('common.marker', 2)]),
                }"
            >
                <template #actions-data="{ row: marker }">
                    <div :key="marker.id">
                        <UTooltip :text="$t('common.mark')">
                            <UButton variant="link" icon="i-mdi-map-marker" @click="goto({ x: marker.x, y: marker.y })" />
                        </UTooltip>

                        <UTooltip :text="$t('common.delete')">
                            <UButton
                                v-if="can('livemap.LivemapService/DeleteMarker').value"
                                variant="link"
                                icon="i-mdi-delete"
                                color="error"
                                @click="
                                    modal.open(ConfirmModal, {
                                        confirm: async () => deleteMarker(marker.id),
                                    })
                                "
                            />
                        </UTooltip>
                    </div>
                </template>

                <template #createdAt-data="{ row: marker }">
                    <GenericTime :value="marker.createdAt" type="compact" />
                </template>

                <template #expiresAt-data="{ row: marker }">
                    <GenericTime v-if="marker.expiresAt" :value="marker.expiresAt" type="compact" />
                    <span v-else>
                        {{ $t('common.na') }}
                    </span>
                </template>

                <template #name-data="{ row: marker }">
                    {{ marker.name }}
                </template>

                <template #type-data="{ row: marker }">
                    {{ $t(`enums.livemap.MarkerType.${MarkerType[marker.type]}`) }}
                </template>

                <template #description-data="{ row: marker }">
                    <p class="max-h-14 truncate overflow-y-scroll break-words">
                        {{ marker.description ?? $t('common.na') }}
                    </p>
                </template>

                <template #creator-data="{ row: marker }">
                    <span v-if="marker.creator">
                        <CitizenInfoPopover :user="marker.creator" :trailing="false" />
                    </span>
                    <span v-else>
                        {{ $t('common.unknown') }}
                    </span>
                </template>

                <template #job-data="{ row: marker }">
                    {{ marker.creator?.jobLabel ?? $t('common.na') }}
                </template>
            </UTable>

            <div class="flex-1" />
        </div>
    </div>
</template>
