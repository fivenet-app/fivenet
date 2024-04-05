<script lang="ts" setup>
import { LPopup } from '@vue-leaflet/vue-leaflet';
import { type MarkerMarker } from '~~/gen/ts/resources/livemap/livemap';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/store/livemap';
import ConfirmModal from '../partials/ConfirmModal.vue';

defineProps<{
    marker: MarkerMarker;
}>();

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const { $grpc } = useNuxtApp();

const modal = useModal();

const livemapStore = useLivemapStore();
const { deleteMarkerMarker } = livemapStore;

async function deleteMarker(id: string): Promise<void> {
    try {
        const call = $grpc.getLivemapperClient().deleteMarker({
            id,
        });
        await call;

        deleteMarkerMarker(id);
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <LPopup :options="{ closeButton: true }">
        <div class="mb-1 flex items-center gap-2">
            <UButton
                v-if="marker.info?.x && marker.info?.y"
                variant="link"
                icon="i-mdi-map-marker"
                @click="$emit('goto', { x: marker.info?.x, y: marker.info?.y })"
            >
                {{ $t('common.mark') }}
            </UButton>
            <UButton
                v-if="can('LivemapperService.DeleteMarker')"
                :title="$t('common.delete')"
                variant="link"
                icon="i-mdi-trash-can"
                @click="
                    modal.open(ConfirmModal, {
                        confirm: async () => deleteMarker(marker.info!.id),
                    })
                "
            >
                {{ $t('common.delete') }}
            </UButton>
        </div>
        <span class="font-semibold"> {{ $t('common.marker') }}: {{ marker.info?.name }} </span>
        <ul role="list" class="flex flex-col">
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
    </LPopup>
</template>
