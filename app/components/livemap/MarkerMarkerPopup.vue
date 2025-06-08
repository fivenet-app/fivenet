<script lang="ts" setup>
import { LPopup } from '@vue-leaflet/vue-leaflet';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/stores/livemap';
import type { MarkerMarker } from '~~/gen/ts/resources/livemap/livemap';
import { availableIcons, fallbackIcon } from '../partials/icons';
import MarkerCreateOrUpdateSlideover from './MarkerCreateOrUpdateSlideover.vue';

defineProps<{
    marker: MarkerMarker;
}>();

const { $grpc } = useNuxtApp();

const { can } = useAuth();

const modal = useModal();
const slideover = useSlideover();

const livemapStore = useLivemapStore();
const { deleteMarkerMarker, goto } = livemapStore;

async function deleteMarker(id: number): Promise<void> {
    try {
        const call = $grpc.livemap.livemap.deleteMarker({
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
        <div class="flex flex-col gap-2 overflow-x-hidden">
            <div class="grid grid-cols-2 gap-1">
                <UTooltip v-if="marker.x !== undefined && marker.y !== undefined" :text="$t('common.mark')">
                    <UButton variant="link" icon="i-mdi-map-marker" @click="goto({ x: marker.x, y: marker.y })">
                        <span class="truncate">
                            {{ $t('common.mark') }}
                        </span>
                    </UButton>
                </UTooltip>

                <UTooltip v-if="can('livemap.LivemapService.CreateOrUpdateMarker').value" :text="$t('common.edit')">
                    <UButton
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
                </UTooltip>

                <UTooltip v-if="can('livemap.LivemapService.DeleteMarker').value" :text="$t('common.delete')">
                    <UButton
                        variant="link"
                        icon="i-mdi-delete"
                        color="error"
                        @click="
                            modal.open(ConfirmModal, {
                                confirm: async () => deleteMarker(marker.id),
                            })
                        "
                    >
                        <span class="truncate">
                            {{ $t('common.delete') }}
                        </span>
                    </UButton>
                </UTooltip>
            </div>

            <p class="inline-flex items-center gap-1">
                <span class="font-semibold"> {{ $t('common.marker') }}:</span>
                <span class="flex-1">{{ marker.name }}</span>

                <template v-if="marker.data?.data.oneofKind === 'icon'">
                    <component
                        :is="
                            availableIcons.find(
                                (icon) =>
                                    marker.data?.data.oneofKind === 'icon' &&
                                    icon.name === convertDynamicIconNameToComponent(marker.data?.data.icon.icon),
                            )?.component ?? fallbackIcon.name
                        "
                        class="size-6"
                        :style="{ color: marker.color ?? 'currentColor' }"
                    />
                </template>
            </p>

            <ul role="list">
                <li>
                    <span class="font-semibold">{{ $t('common.job') }}:</span>
                    {{ marker.jobLabel ?? $t('common.na') }}
                </li>
                <li>
                    <span class="font-semibold">{{ $t('common.description') }}:</span>
                    {{ marker.description ?? $t('common.na') }}
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
