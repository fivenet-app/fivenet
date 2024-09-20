<script lang="ts" setup>
import DispatchDetailsSlideover from '~/components/centrum/dispatches/DispatchDetailsSlideover.vue';
import { dispatchStatusToBGColor } from '~/components/centrum/helpers';
import DispatchAttributes from '~/components/centrum/partials/DispatchAttributes.vue';
import UnitInfoPopover from '~/components/centrum/units/UnitInfoPopover.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';
import type { Dispatch} from '~~/gen/ts/resources/centrum/dispatches';
import { StatusDispatch } from '~~/gen/ts/resources/centrum/dispatches';

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
}>();

const { goto } = useLivemapStore();

const centrumStore = useCentrumStore();
const { ownUnitId, timeCorrection } = storeToRefs(centrumStore);

const expiresAt = props.dispatch.units.find((u) => u.unitId === ownUnitId.value)?.expiresAt;
const dispatchBackground = computed(() => dispatchStatusToBGColor(props.dispatch.status?.status));

const slideover = useSlideover();

const checked = ref(false);

onBeforeMount(() => {
    if (props.preselected === true) {
        checked.value = props.preselected;
        emits('selected', true);
    }
});
</script>

<template>
    <div class="flex px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-2 sm:px-0">
        <dt class="flex flex-initial flex-col gap-1 text-sm font-medium leading-6">
            <div class="flex items-center">
                <UCheckbox v-model="checked" name="selected" @change="$emit('selected', checked)" />

                <IDCopyBadge
                    :id="dispatch.id"
                    class="ml-2"
                    prefix="DSP"
                    :action="
                        () =>
                            slideover.open(DispatchDetailsSlideover, {
                                dispatchId: dispatch.id,
                            })
                    "
                />
            </div>
            <div v-if="expiresAt" class="flex flex-col text-sm">
                <span class="font-semibold">{{ $t('common.expires_in') }}:</span>
                <span>{{
                    useLocaleTimeAgo(toDate(expiresAt, timeCorrection), { showSecond: true, updateInterval: 1_000 }).value
                }}</span>
            </div>
            <div v-if="expiresAt" class="flex flex-col text-sm">
                <span class="font-semibold">{{ $t('common.created') }}:</span>
                <span>{{
                    useLocaleTimeAgo(toDate(dispatch.createdAt, timeCorrection), { showSecond: true, updateInterval: 1_000 })
                        .value
                }}</span>
            </div>
        </dt>
        <dd class="mt-1 flex-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
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
                    <div class="sm:inline-flex sm:flex-row sm:gap-2">
                        <span>
                            <span class="font-medium">{{ $t('common.postal') }}:</span>
                            {{ dispatch.postal ?? $t('common.na') }}
                        </span>
                        <UButton
                            size="xs"
                            variant="link"
                            icon="i-mdi-map-marker"
                            :padded="false"
                            @click="goto({ x: dispatch.x, y: dispatch.y })"
                        >
                            {{ $t('common.go_to_location') }}
                        </UButton>
                    </div>
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
                    <div class="flex flex-1 items-center gap-1">
                        <span class="font-medium">{{ $t('common.unit', 2) }}:</span>
                        <span v-if="dispatch.units.length === 0">{{ $t('common.member', 0) }}</span>
                        <span v-else class="grid flex-1 grid-cols-2 gap-1 truncate">
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
