<script lang="ts" setup>
import type { PointExpression } from 'leaflet';
import { BellIcon } from 'mdi-vue3';
import {
    calculateDispatchZIndexOffset,
    checkDispatchAccess,
    dispatchStatusAnimate,
    dispatchStatusToFillColor,
} from '~/components/centrum/helpers';
import DispatchAttributes from '~/components/centrum/partials/DispatchAttributes.vue';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useLivemapStore } from '~/stores/livemap';
import { CentrumAccessLevel } from '~~/gen/ts/resources/centrum/access';
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';
import DispatchAssignModal from '../dispatches/DispatchAssignModal.vue';
import DispatchStatusBadge from '../partials/DispatchStatusBadge.vue';

const props = withDefaults(
    defineProps<{
        dispatch: Dispatch;
        size?: number;
    }>(),
    {
        size: 22,
    },
);

const emit = defineEmits<{
    (e: 'selected', dsp: Dispatch): void;
}>();

const overlay = useOverlay();

const dispatchAssignModal = overlay.create(DispatchAssignModal);

const { goto } = useLivemapStore();

const { selfAssign, canDo } = useCentrumStore();

const iconAnchor: PointExpression = [props.size / 2, props.size * 1.65];
const popupAnchor: PointExpression = [0, -(props.size * 1.7)];

function selected(_: string | number | string) {
    emit('selected', props.dispatch);
}

const dispatchClasses = computed(() => [
    dispatchStatusToFillColor(props.dispatch.status?.status),
    dispatchStatusAnimate(props.dispatch.status?.status) ? 'animate-wiggle' : '',
]);

const zIndexOffset = computed(() => calculateDispatchZIndexOffset(props.dispatch.status?.status));
</script>

<template>
    <LMarker
        :key="dispatch.id"
        :lat-lng="[dispatch.y, dispatch.x]"
        :name="dispatch.id.toString()"
        :z-index-offset="zIndexOffset"
    >
        <LIcon :icon-anchor="iconAnchor" :popup-anchor="popupAnchor" :icon-size="[size, size]">
            <div class="flex flex-col items-center uppercase">
                <span
                    class="inset-0 rounded-md border border-black/20 bg-neutral-50 bg-clip-padding whitespace-nowrap text-black hover:bg-[#f4f4f4]"
                >
                    DSP-{{ dispatch.id }}
                </span>
                <BellIcon class="size-full" :class="dispatchClasses" />
            </div>
        </LIcon>

        <LPopup class="min-w-[175px]" :options="{ closeButton: false }">
            <UCard
                class="-my-[13px] -mr-[24px] -ml-[20px] flex min-w-[200px] flex-col"
                :ui="{ header: 'p-1 sm:px-2', body: 'p-1 sm:p-2', footer: 'p-1 sm:px-2' }"
            >
                <template #header>
                    <div class="grid grid-cols-2 gap-2">
                        <UButton
                            v-if="dispatch?.x !== undefined && dispatch?.y !== undefined"
                            variant="link"
                            icon="i-mdi-map-marker"
                            :label="$t('common.mark')"
                            @click="goto({ x: dispatch?.x, y: dispatch?.y })"
                        />

                        <UButton
                            variant="link"
                            icon="i-mdi-car-emergency"
                            :label="$t('common.detail', 2)"
                            @click="selected(dispatch.id)"
                        />

                        <UButton
                            v-if="canDo('TakeControl') && checkDispatchAccess(dispatch.jobs, CentrumAccessLevel.DISPATCH)"
                            class="truncate"
                            icon="i-mdi-account-multiple-plus"
                            variant="link"
                            :label="$t('common.assign')"
                            @click="
                                dispatchAssignModal.open({
                                    dispatchId: dispatch.id,
                                })
                            "
                        />

                        <UButton
                            v-if="canDo('TakeDispatch') && checkDispatchAccess(dispatch.jobs, CentrumAccessLevel.PARTICIPATE)"
                            class="text-left"
                            icon="i-mdi-plus"
                            variant="link"
                            :label="$t('common.self_assign')"
                            @click="selfAssign(dispatch.id)"
                        />
                    </div>
                </template>

                <p class="inline-flex items-center gap-1">
                    <span class="font-semibold">{{ $t('common.dispatch', 1) }}</span>
                    <UButton
                        class="font-semibold"
                        :label="`DSP-${dispatch.id}`"
                        variant="link"
                        @click="selected(dispatch.id)"
                    />
                </p>

                <ul role="list">
                    <li class="inline-flex gap-1">
                        <span class="flex-initial font-semibold">{{ $t('common.job') }}:</span>
                        <span class="flex-1">
                            {{ dispatch.jobs?.jobs.map((j) => j.label ?? j.name).join(', ') }}
                        </span>
                    </li>
                    <li>
                        <span class="font-semibold">{{ $t('common.sent_at') }}:</span>
                        {{ $d(toDate(dispatch.createdAt), 'short') }}
                    </li>
                    <li class="inline-flex gap-1">
                        <span class="flex-initial">
                            <span class="font-semibold">{{ $t('common.sent_by') }}:</span>
                        </span>
                        <span class="flex-1">
                            <template v-if="dispatch.anon">
                                {{ $t('common.anon') }}
                            </template>
                            <CitizenInfoPopover v-else-if="dispatch.creator" :user="dispatch.creator" />
                            <template v-else>
                                {{ $t('common.unknown') }}
                            </template>
                        </span>
                    </li>
                    <li>
                        <span class="font-semibold">{{ $t('common.postal') }}:</span> {{ dispatch.postal ?? $t('common.na') }}
                    </li>
                    <li>
                        <span class="font-semibold">{{ $t('common.message') }}:</span> {{ dispatch.message }}
                    </li>
                    <li class="truncate">
                        <span class="font-semibold">{{ $t('common.description') }}:</span>
                        {{ dispatch.description ?? $t('common.na') }}
                    </li>
                    <li class="inline-flex gap-1">
                        <span class="font-semibold">{{ $t('common.status') }}:</span>
                        <DispatchStatusBadge :status="dispatch.status?.status" />
                    </li>
                    <li>
                        <span class="font-semibold">{{ $t('common.attributes', 2) }}:</span>
                        <DispatchAttributes class="ml-1" :attributes="dispatch.attributes" />
                    </li>
                    <li class="inline-flex gap-1">
                        <span class="font-semibold">{{ $t('common.unit') }}:</span>

                        <span v-if="dispatch.units.length === 0" class="italic">{{
                            $t('enums.centrum.StatusDispatch.UNASSIGNED')
                        }}</span>
                        <span v-else class="grid grid-cols-2 gap-1">
                            <UnitInfoPopover
                                v-for="unit in dispatch.units"
                                :key="unit.unitId"
                                :unit="unit.unit"
                                :initials-only="true"
                                :badge="true"
                                :assignment="unit"
                            />
                        </span>
                    </li>
                </ul>
            </UCard>
        </LPopup>
    </LMarker>
</template>
