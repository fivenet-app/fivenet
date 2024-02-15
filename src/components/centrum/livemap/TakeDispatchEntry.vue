<script lang="ts" setup>
import { AccountIcon, MapMarkerIcon } from 'mdi-vue3';
import DispatchDetails from '~/components/centrum/dispatches/DispatchDetails.vue';
import { dispatchStatusToBGColor } from '~/components/centrum/helpers';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useCentrumStore } from '~/store/centrum';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';
import DispatchAttributes from '~/components/centrum/partials/DispatchAttributes.vue';

const props = withDefaults(
    defineProps<{
        dispatch: Dispatch;
        preselected?: boolean;
    }>(),
    {
        preselected: true,
    },
);

const emits = defineEmits<{
    (e: 'selected', state: boolean): void;
    (e: 'goto', loc: Coordinate): void;
}>();

const centrumStore = useCentrumStore();
const { ownUnitId, timeCorrection } = storeToRefs(centrumStore);

const expiresAt = props.dispatch.units.find((u) => u.unitId === ownUnitId.value)?.expiresAt;
const dispatchBackground = computed(() => dispatchStatusToBGColor(props.dispatch.status?.status));

const checked = ref(false);

onBeforeMount(() => {
    if (props.preselected === true) {
        checked.value = props.preselected;
        emits('selected', true);
    }
});

const open = ref(false);
</script>

<template>
    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
        <DispatchDetails :open="open" :dispatch="dispatch" @close="open = false" @goto="$emit('goto', $event)" />

        <dt class="text-sm font-medium leading-6 text-neutral">
            <div class="flex h-6 items-center">
                <input
                    type="checkbox"
                    name="selected"
                    :checked="checked"
                    class="h-4 h-5 w-4 w-5 rounded border-gray-300 text-primary-600 focus:ring-primary-600"
                    @change="
                        checked = !checked;
                        $emit('selected', checked);
                    "
                />
                <IDCopyBadge :id="dispatch.id" class="ml-2" prefix="DSP" :action="() => (open = true)" />
            </div>
            <div v-if="expiresAt" class="mt-1 flex flex-col text-sm text-neutral">
                <span class="font-semibold">{{ $t('common.expires_in') }}:</span>
                <span>{{
                    useLocaleTimeAgo(toDate(expiresAt, timeCorrection), { showSecond: true, updateInterval: 1_000 }).value
                }}</span>
            </div>
        </dt>
        <dd class="mt-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
            <ul role="list" class="divide-y divide-base-200 rounded-md border border-base-200">
                <li class="flex items-center py-3 pl-3 pr-4 text-sm">
                    <span class="font-medium">{{ $t('common.sent_by') }}:</span>
                    <span class="ml-1">
                        <template v-if="dispatch.anon">
                            {{ $t('common.anon') }}
                        </template>
                        <CitizenInfoPopover v-else-if="dispatch.creator" :user="dispatch.creator" />
                        <template v-else>
                            {{ $t('common.unknown') }}
                        </template>
                    </span>
                </li>
                <li class="flex items-center py-3 pl-3 pr-4 text-sm">
                    <span class="font-medium">{{ $t('common.message') }}:</span>
                    <span class="ml-1 truncate">{{ dispatch.message }}</span>
                </li>
                <li class="py-3 pl-3 pr-4 text-sm">
                    <span class="block">
                        <span class="font-medium">{{ $t('common.postal') }}:</span>
                        {{ dispatch.postal ?? $t('common.na') }}
                    </span>
                    <button
                        type="button"
                        class="inline-flex items-center text-primary-400 hover:text-primary-600"
                        @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
                    >
                        <MapMarkerIcon class="h-5 w-5" aria-hidden="true" />
                        <span class="ml-1">
                            {{ $t('common.go_to_location') }}
                        </span>
                    </button>
                </li>
                <li class="flex items-center py-3 pl-3 pr-4 text-sm">
                    <span class="font-medium">{{ $t('common.status') }}:</span>
                    <span class="ml-1 text-neutral" :class="dispatchBackground">{{
                        $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`)
                    }}</span>
                </li>
                <li v-if="dispatch.attributes" class="flex items-center py-3 pl-3 pr-4 text-sm">
                    <span class="font-medium">{{ $t('common.attributes') }}:</span>
                    <span class="ml-1">
                        <DispatchAttributes :attributes="dispatch.attributes" />
                    </span>
                </li>
                <li class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                    <div class="flex flex-1 items-center">
                        <AccountIcon class="mr-1 h-5 w-5 flex-shrink-0" aria-hidden="true" />
                        <span class="mr-1 font-medium">{{ $t('common.units', 2) }}:</span>
                        <span v-if="dispatch.units.length === 0">{{ $t('common.member', 0) }}</span>
                        <span v-else class="ml-2 grid flex-1 grid-cols-2 gap-1 truncate">
                            <UnitInfoPopover
                                v-for="unit in dispatch.units"
                                :key="unit.unitId"
                                :unit="unit.unit"
                                :initials-only="true"
                                :badge="true"
                                :assignment="unit"
                            />
                        </span>
                    </div>
                </li>
            </ul>
        </dd>
    </div>
</template>
