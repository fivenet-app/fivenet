<script lang="ts" setup>
import { AccountIcon, MapMarkerIcon } from 'mdi-vue3';
import Details from '~/components/centrum/dispatches/Details.vue';
import { dispatchStatusToBGColor } from '~/components/centrum/helpers';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useCentrumStore } from '~/store/centrum';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';

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
const { ownUnitId } = storeToRefs(centrumStore);

const expiresAt = props.dispatch.units.find((u) => u.unitId === ownUnitId.value)?.expiresAt;
const dispatchBackground = computed(() => dispatchStatusToBGColor(props.dispatch.status?.status ?? 0));

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
        <Details :open="open" @close="open = false" @goto="$emit('goto', $event)" :dispatch="dispatch" />

        <dt class="text-sm font-medium leading-6 text-neutral">
            <div class="flex h-6 items-center">
                <input
                    type="checkbox"
                    name="selected"
                    :checked="checked"
                    class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-600 h-6 w-6"
                    @change="
                        checked = !checked;
                        $emit('selected', checked);
                    "
                />
                <IDCopyBadge class="ml-2" prefix="DSP" :id="dispatch.id" :action="() => (open = true)" />
            </div>
            <div v-if="expiresAt" class="mt-1 text-neutral text-sm">
                {{ $t('common.expires_in') }}:
                {{ useLocaleTimeAgo(toDate(expiresAt), { showSecond: true, updateInterval: 1000 }).value }}
            </div>
        </dt>
        <dd class="mt-1 text-sm leading-6 text-gray-400 sm:col-span-2 sm:mt-0">
            <ul role="list" class="border divide-y rounded-md divide-base-200 border-base-200">
                <li class="flex items-center py-3 pl-3 pr-4 text-sm">
                    <span class="font-medium">{{ $t('common.sent_by') }}</span
                    >:
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
                    <span class="font-medium">{{ $t('common.message') }}</span
                    >: <span class="ml-1 truncate">{{ dispatch.message }}</span>
                </li>
                <li class="py-3 pl-3 pr-4 text-sm">
                    <span class="block">
                        {{ $t('common.postal') }}:
                        {{ dispatch.postal ?? $t('common.na') }}
                    </span>
                    <button
                        type="button"
                        class="inline-flex items-center text-primary-400 hover:text-primary-600"
                        @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
                    >
                        <MapMarkerIcon class="w-5 h-5" aria-hidden="true" />
                        <span class="ml-1">
                            {{ $t('common.go_to_location') }}
                        </span>
                    </button>
                </li>
                <li class="flex items-center py-3 pl-3 pr-4 text-sm">
                    <span class="font-medium">{{ $t('common.status') }}</span
                    >:
                    <span class="ml-1 text-neutral" :class="dispatchBackground">{{
                        $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`)
                    }}</span>
                </li>
                <li class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                    <div class="flex items-center flex-1">
                        <AccountIcon class="flex-shrink-0 w-5 h-5 text-base-400 mr-1" aria-hidden="true" />
                        <span class="font-medium mr-1">{{ $t('common.units', 2) }}:</span>
                        <span v-if="dispatch.units.length === 0">{{ $t('common.member', 0) }}</span>
                        <span v-else class="flex-1 ml-2 truncate grid grid-cols-2 gap-1">
                            <template v-for="unit in dispatch.units">
                                <UnitInfoPopover :unit="unit.unit" :initials-only="true" :badge="true" :assignment="unit" />
                            </template>
                        </span>
                    </div>
                </li>
            </ul>
        </dd>
    </div>
</template>
