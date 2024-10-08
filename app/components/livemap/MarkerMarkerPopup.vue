<script lang="ts" setup>
import { LPopup } from '@vue-leaflet/vue-leaflet';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/store/livemap';
import type { MarkerMarker } from '~~/gen/ts/resources/livemap/livemap';
import MarkerCreateOrUpdateSlideover from './MarkerCreateOrUpdateSlideover.vue';

defineProps<{
    marker: MarkerMarker;
}>();

const modal = useModal();
const slideover = useSlideover();

const livemapStore = useLivemapStore();
const { deleteMarkerMarker, goto } = livemapStore;

async function deleteMarker(id: string): Promise<void> {
    try {
        const call = getGRPCLivemapperClient().deleteMarker({
            id,
        });
        await call;

        deleteMarkerMarker(id);
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <LPopup :options="{ closeButton: true }">
        <div class="flex flex-col gap-2">
            <div class="grid grid-cols-2 gap-1">
                <UButton
                    v-if="marker.info?.x !== undefined && marker.info?.y !== undefined"
                    variant="link"
                    icon="i-mdi-map-marker"
                    @click="goto({ x: marker.info?.x, y: marker.info?.y })"
                >
                    <span class="truncate">
                        {{ $t('common.mark') }}
                    </span>
                </UButton>

                <UButton
                    v-if="can('LivemapperService.CreateOrUpdateMarker').value"
                    :title="$t('common.edit')"
                    variant="link"
                    icon="i-mdi-pencil"
                    @click="
                        slideover.open(MarkerCreateOrUpdateSlideover, {
                            marker: marker,
                        })
                    "
                >
                    <span class="truncate">
                        {{ $t('common.edit') }}
                    </span>
                </UButton>

                <UButton
                    v-if="can('LivemapperService.DeleteMarker').value"
                    :title="$t('common.delete')"
                    variant="link"
                    icon="i-mdi-trash-can"
                    color="red"
                    @click="
                        modal.open(ConfirmModal, {
                            confirm: async () => deleteMarker(marker.info!.id),
                        })
                    "
                >
                    <span class="truncate">
                        {{ $t('common.delete') }}
                    </span>
                </UButton>
            </div>

            <p class="inline-flex items-center gap-1">
                <span class="font-semibold"> {{ $t('common.marker') }}:</span>
                <span>{{ marker.info?.name }}</span>
            </p>

            <ul role="list">
                <li>
                    <span class="font-semibold">{{ $t('common.job') }}:</span>
                    {{ marker.info?.jobLabel ?? $t('common.na') }}
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.description') }}:</span>
                    {{ marker.info?.description ?? $t('common.na') }}
                </li>
                <li class="inline-flex gap-1">
                    <span class="font-semibold">{{ $t('common.expires_at') }}:</span>
                    <GenericTime v-if="marker.expiresAt" :value="marker.expiresAt" />
                    <span v-else>{{ $t('common.na') }}</span>
                </li>
                <li class="inline-flex gap-1">
                    <span class="flex-initial">
                        <span class="font-semibold">{{ $t('common.sent_by') }}:</span>
                    </span>
                    <span class="flex-1">
                        <CitizenInfoPopover v-if="marker.creator" :user="marker.creator" />
                        <template v-else>
                            {{ $t('common.unknown') }}
                        </template>
                    </span>
                </li>
            </ul>
        </div>
    </LPopup>
</template>
