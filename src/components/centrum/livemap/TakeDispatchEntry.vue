<script lang="ts" setup>
import { AccountIcon } from 'mdi-vue3';
import DispatchDetailsSlideover from '~/components/centrum/dispatches/DispatchDetailsSlideover.vue';
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
    <div class="flex px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-2 sm:px-0">
        <DispatchDetailsSlideover :open="open" :dispatch="dispatch" @close="open = false" @goto="$emit('goto', $event)" />

        <dt class="flex-initial text-sm font-medium leading-6">
            <div class="flex h-6 items-center">
                <UInput
                    type="checkbox"
                    name="selected"
                    :checked="checked"
                    class="text-primary-600 focus:ring-primary-600 size-5 rounded border-gray-300"
                    @change="
                        checked = !checked;
                        $emit('selected', checked);
                    "
                />
                <IDCopyBadge :id="dispatch.id" class="ml-2" prefix="DSP" :action="() => (open = true)" :hide-icon="true" />
            </div>
            <div v-if="expiresAt" class="mt-1 flex flex-col text-sm">
                <span class="font-semibold">{{ $t('common.expires_in') }}:</span>
                <span>{{
                    useLocaleTimeAgo(toDate(expiresAt, timeCorrection), { showSecond: true, updateInterval: 1_000 }).value
                }}</span>
            </div>
            <div v-if="expiresAt" class="mt-1 flex flex-col text-sm">
                <span class="font-semibold">{{ $t('common.created') }}:</span>
                <span>{{
                    useLocaleTimeAgo(toDate(dispatch.createdAt, timeCorrection), { showSecond: true, updateInterval: 1_000 })
                        .value
                }}</span>
            </div>
        </dt>
        <dd class="mt-1 flex-1 text-sm leading-6 text-gray-300 sm:col-span-2 sm:mt-0">
            <ul role="list" class="divide-y divide-base-200 rounded-md border border-base-200">
                <li class="flex items-center gap-1 py-3 pl-3 pr-4 text-sm">
                    <span class="font-medium">{{ $t('common.sent_by') }}:</span>
                    <span>
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
                    <UButton
                        size="xs"
                        variant="link"
                        icon="i-mdi-map-marker"
                        @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
                    >
                        {{ $t('common.go_to_location') }}
                    </UButton>
                </li>
                <li class="flex items-center gap-1 py-3 pl-3 pr-4 text-sm">
                    <span class="font-medium">{{ $t('common.status') }}:</span>
                    <span :class="dispatchBackground">{{
                        $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`)
                    }}</span>
                </li>
                <li v-if="dispatch.attributes" class="flex items-center gap-1 py-3 pl-3 pr-4 text-sm">
                    <span class="font-medium">{{ $t('common.attributes', 2) }}:</span>
                    <span>
                        <DispatchAttributes :attributes="dispatch.attributes" />
                    </span>
                </li>
                <li class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                    <div class="flex flex-1 items-center">
                        <AccountIcon class="mr-1 size-5 shrink-0" />
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
