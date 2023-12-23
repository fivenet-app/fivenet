<script lang="ts" setup>
import { useDebounceFn } from '@vueuse/core';
import { useSound } from '@vueuse/sound';
import { AccountMultiplePlusIcon, CloseOctagonIcon, DotsVerticalIcon, MapMarkerIcon } from 'mdi-vue3';
import AssignDispatchModal from '~/components/centrum/dispatches/AssignDispatchModal.vue';
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

const openMessage = ref(false);
</script>

<template>
    <tr class="transition-colors hover:bg-neutral/5 even:bg-base-800">
        <DispatchDetails :dispatch="dispatch" :open="openDetails" @close="openDetails = false" @goto="$emit('goto', $event)" />
        <AssignDispatchModal v-if="openAssign" :open="openAssign" :dispatch="dispatch" @close="openAssign = false" />
        <DispatchStatusUpdateModal
            v-if="openStatus"
            :open="openStatus"
            :dispatch-id="dispatch.id"
            @close="openStatus = false"
        />

        <td
            class="relative items-center whitespace-nowrap pl-0 py-1 pr-0 text-left text-sm font-medium sm:pr-0.5 justify-start"
        >
            <template v-if="!hideActions">
                <button
                    type="button"
                    class="text-primary-400 hover:text-primary-600"
                    :title="$t('common.assign')"
                    @click="openAssign = true"
                >
                    <AccountMultiplePlusIcon class="w-6 h-auto ml-auto mr-1.5" aria-hidden="true" />
                </button>
                <button
                    type="button"
                    class="text-primary-400 hover:text-primary-600"
                    :title="$t('common.go_to_location')"
                    @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
                >
                    <MapMarkerIcon class="w-6 h-auto ml-auto mr-1.5" aria-hidden="true" />
                </button>
                <button
                    type="button"
                    class="text-primary-400 hover:text-primary-600"
                    :title="$t('common.status')"
                    @click="openStatus = true"
                >
                    <CloseOctagonIcon class="w-6 h-auto ml-auto mr-1.5" aria-hidden="true" />
                </button>
            </template>
            <button
                type="button"
                class="text-primary-400 hover:text-primary-600"
                :title="$t('common.detail', 2)"
                @click="openDetails = true"
            >
                <DotsVerticalIcon class="w-6 h-auto ml-auto mr-1.5" />
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
        <td class="px-1 py-1 text-sm text-gray-300">
            <p class="break-all max-h-22" :class="openMessage ? '' : 'line-clamp-1 max-w-sm'">
                {{ dispatch.message }}
            </p>
            <button
                v-if="dispatch.message.length > 40"
                type="button"
                class="flex justify-center px-1 py-1 text-sm font-semibold transition-colors rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 bg-accent-500 hover:bg-accent-400 focus-visible:outline-accent-500"
                @click="openMessage = !openMessage"
            >
                {{ openMessage ? $t('common.read_less') : $t('common.read_more') }}
            </button>
        </td>
    </tr>
</template>
