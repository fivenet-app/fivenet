<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import ConfirmModal from '~/components/partials/ConfirmModal.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useLivemapStore } from '~/stores/livemap';
import { getLivemapLivemapClient } from '~~/gen/ts/clients';
import type { MarkerMarker } from '~~/gen/ts/resources/livemap/marker_marker';
import { availableIcons, fallbackIcon } from '../partials/icons';
import MarkerCreateOrUpdateSlideover from './MarkerCreateOrUpdateSlideover.vue';

defineProps<{
    marker: MarkerMarker;
}>();

const { can } = useAuth();

const overlay = useOverlay();

const livemapStore = useLivemapStore();
const { deleteMarkerMarker, goto } = livemapStore;

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

const confirmModal = overlay.create(ConfirmModal);
const markerCreateOrUpdateSlideover = overlay.create(MarkerCreateOrUpdateSlideover);
</script>

<template>
    <LPopup class="min-w-[175px]" :options="{ closeButton: false }">
        <UCard class="-my-[13px] -mr-[24px] -ml-[20px] flex min-w-[200px] flex-col">
            <template #header>
                <div class="grid grid-cols-2 gap-2">
                    <UTooltip v-if="marker.x !== undefined && marker.y !== undefined" :text="$t('common.mark')">
                        <UButton variant="link" icon="i-mdi-map-marker" @click="goto({ x: marker.x, y: marker.y })">
                            <span class="truncate">
                                {{ $t('common.mark') }}
                            </span>
                        </UButton>
                    </UTooltip>

                    <UTooltip v-if="can('livemap.LivemapService/CreateOrUpdateMarker').value" :text="$t('common.edit')">
                        <UButton
                            variant="link"
                            icon="i-mdi-pencil"
                            @click="
                                markerCreateOrUpdateSlideover.open({
                                    marker: marker,
                                })
                            "
                        >
                            <span class="truncate">
                                {{ $t('common.edit') }}
                            </span>
                        </UButton>
                    </UTooltip>

                    <UTooltip v-if="can('livemap.LivemapService/DeleteMarker').value" :text="$t('common.delete')">
                        <UButton
                            variant="link"
                            icon="i-mdi-delete"
                            color="error"
                            @click="
                                confirmModal.open({
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
            </template>

            <p class="inline-flex items-center gap-1">
                <span class="font-semibold"> {{ $t('common.marker') }} {{ marker.name }}</span>

                <template v-if="marker.data?.data.oneofKind === 'icon'">
                    <component
                        :is="
                            availableIcons.find(
                                (icon) =>
                                    marker.data?.data.oneofKind === 'icon' &&
                                    icon.name === convertDynamicIconNameToComponent(marker.data?.data.icon.icon),
                            )?.component ?? fallbackIcon.component
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
        </UCard>
    </LPopup>
</template>
