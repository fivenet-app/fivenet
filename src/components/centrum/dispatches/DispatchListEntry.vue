<script lang="ts" setup>
import { useDebounceFn } from '@vueuse/core';
import { useSound } from '@vueuse/sound';
import { AccountMultiplePlusIcon, CloseOctagonIcon, DotsVerticalIcon, MapMarkerIcon } from 'mdi-vue3';
import DispatchAssignModal from '~/components/centrum/dispatches/DispatchAssignModal.vue';
import DispatchDetails from '~/components/centrum/dispatches/DispatchDetails.vue';
import DispatchStatusUpdateModal from '~/components/centrum/dispatches/DispatchStatusUpdateModal.vue';
import { dispatchStatusAnimate, dispatchStatusToBGColor } from '~/components/centrum/helpers';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

const props = withDefaults(
    defineProps<{
        dispatch: Dispatch;
        hideActions?: boolean;
    }>(),
    {
        hideActions: false,
    },
);

defineEmits<{
    (e: 'goto', loc: Coordinate): void;
}>();

const dispatchBackground = computed(() => dispatchStatusToBGColor(props.dispatch.status?.status));
const dispatchAnimated = computed(() => (dispatchStatusAnimate(props.dispatch.status?.status) ? 'animate-pulse' : ''));

const dispatchAssistanceSound = useSound('/sounds/centrum/morse-sos.mp3', {
    volume: 0.15,
});
const debouncedPlay = useDebounceFn(() => dispatchAssistanceSound.play(), 950);

watch(props, () => {
    if (props.dispatch.status?.status === StatusDispatch.NEED_ASSISTANCE) {
        debouncedPlay();
    }
});

const openDetails = ref(false);
const openAssign = ref(false);
const openStatus = ref(false);
</script>

<template>
    <tr class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <DispatchDetails :dispatch="dispatch" :open="openDetails" @close="openDetails = false" @goto="$emit('goto', $event)" />
        <DispatchAssignModal v-if="openAssign" :open="openAssign" :dispatch="dispatch" @close="openAssign = false" />
        <DispatchStatusUpdateModal
            v-if="openStatus"
            :open="openStatus"
            :dispatch-id="dispatch.id"
            @close="openStatus = false"
        />

        <td
            class="relative items-center justify-start whitespace-nowrap py-1 pl-0 pr-0 text-left text-sm font-medium sm:pr-0.5"
        >
            <button
                v-if="!hideActions"
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.assign')"
                @click="openAssign = true"
            >
                <AccountMultiplePlusIcon class="ml-auto mr-1.5 h-auto w-5" aria-hidden="true" />
            </button>
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.go_to_location')"
                @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
            >
                <MapMarkerIcon class="ml-auto mr-1.5 h-auto w-5" aria-hidden="true" />
            </button>
            <button
                v-if="!hideActions"
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.status')"
                @click="openStatus = true"
            >
                <CloseOctagonIcon class="ml-auto mr-1.5 h-auto w-5" aria-hidden="true" />
            </button>
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.detail', 2)"
                @click="openDetails = true"
            >
                <DotsVerticalIcon class="ml-auto mr-1.5 h-auto w-5" />
            </button>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            {{ dispatch.id }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <GenericTime :value="dispatch.createdAt" type="compact" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-100" :class="dispatchBackground">
            <span :class="dispatchAnimated">
                {{ $t(`enums.centrum.StatusDispatch.${StatusDispatch[props.dispatch.status?.status ?? 0]}`) }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            {{ dispatch.postal ?? $t('common.na') }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <span v-if="dispatch.units.length === 0" class="italic">{{ $t('enums.centrum.StatusDispatch.UNASSIGNED') }}</span>
            <span v-else class="mr-1 grid grid-cols-2 gap-1">
                <UnitInfoPopover
                    v-for="unit in dispatch.units"
                    :key="unit.unitId"
                    :unit="unit.unit"
                    :initials-only="true"
                    :badge="true"
                    :assignment="unit"
                />
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-sm text-gray-300">
            <span v-if="dispatch.anon">
                {{ $t('common.anon') }}
            </span>
            <span v-else-if="dispatch.creator">
                <CitizenInfoPopover :user="dispatch.creator" />
            </span>
            <span v-else>
                {{ $t('common.unknown') }}
            </span>
        </td>
        <td class="inline-flex min-w-36 items-center px-1 py-1 text-sm text-gray-300">
            <p class="line-clamp-2 hover:line-clamp-6">
                {{ dispatch.message }}
            </p>
        </td>
    </tr>
</template>
