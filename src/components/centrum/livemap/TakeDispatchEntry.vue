<script lang="ts" setup>
import { AccountIcon, MapMarkerIcon } from 'mdi-vue3';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { useCentrumStore } from '~/store/centrum';
import { Dispatch, StatusDispatch } from '~~/gen/ts/resources/dispatch/dispatches';
import Details from '../dispatches/Details.vue';
import { dispatchStatusToBGColor } from '../helpers';

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
    (e: 'selected', id: bigint): void;
    (e: 'goto', loc: Coordinate): void;
}>();

const centrumStore = useCentrumStore();
const { ownUnitId } = storeToRefs(centrumStore);

const expiresAt = props.dispatch.units.find((u) => u.unitId === ownUnitId.value)?.expiresAt;
const dispatchBackground = computed(() => dispatchStatusToBGColor(props.dispatch.status?.status ?? 0));

onBeforeMount(() => {
    if (props.preselected === true) {
        emits('selected', props.dispatch.id);
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
                    :checked="preselected"
                    class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-600 h-6 w-6"
                    @change="$emit('selected', dispatch.id)"
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
                    <span class="font-medium">{{ $t('common.message') }}</span
                    >: {{ dispatch.message }}
                </li>
                <li class="flex items-center py-3 pl-3 pr-4 text-sm">
                    <span class="font-medium">{{ $t('common.status') }}</span
                    >:
                    <span class="ml-1 text-neutral" :class="dispatchBackground">{{
                        $t(`enums.centrum.StatusDispatch.${StatusDispatch[dispatch.status?.status ?? 0]}`)
                    }}</span>
                </li>
                <li class="py-3 pl-3 pr-4 text-sm">
                    <span class="block">
                        {{ $t('common.postal') }}:
                        {{ dispatch.postal ?? $t('common.na') }}
                    </span>
                    <button
                        v-if="dispatch.x && dispatch.y"
                        type="button"
                        class="inline-flex items-center text-primary-400 hover:text-primary-600"
                        @click="$emit('goto', { x: dispatch.x, y: dispatch.y })"
                    >
                        <MapMarkerIcon class="w-5 h-5 mr-1" aria-hidden="true" />
                        {{ $t('common.go_to_location') }}
                    </button>
                    <span v-else>{{ $t('common.no_location') }}</span>
                </li>
                <li class="flex items-center justify-between py-3 pl-3 pr-4 text-sm">
                    <div class="flex items-center flex-1">
                        <AccountIcon class="flex-shrink-0 w-5 h-5 text-base-400 mr-1" aria-hidden="true" />
                        <span class="font-medium mr-1">{{ $t('common.members') }}:</span>
                        <span v-if="dispatch.units.length === 0">{{ $t('common.member', 0) }}</span>
                        <span v-else class="flex-1 ml-2 truncate">
                            <span v-for="unit in dispatch.units">
                                {{ unit.unit?.name }}
                                ({{ unit.unit?.initials }})
                            </span>
                        </span>
                    </div>
                </li>
            </ul>
        </dd>
    </div>
</template>
