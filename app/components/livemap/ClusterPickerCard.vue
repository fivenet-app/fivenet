<script setup lang="ts">
import type { Dispatch } from '~~/gen/ts/resources/centrum/dispatches';
import type { MarkerMarker } from '~~/gen/ts/resources/livemap/marker_marker';
import type { UserMarker } from '~~/gen/ts/resources/livemap/user_marker';

const props = withDefaults(
    defineProps<{
        hits: Hit[];
        hiddenCount?: number;
    }>(),
    {
        hiddenCount: 0,
    },
);

type Hit = {
    userMarker?: UserMarker;
    dispatchMarker?: Dispatch;
    markerMarker?: MarkerMarker;

    openPopup: () => void;
};

const { t } = useI18n();

const livemapStore = useLivemapStore();
const { selectedMarker } = storeToRefs(livemapStore);

const groupedMarkers = computed(() =>
    [
        { key: 'userMarkers', label: t('common.user', 2), items: props.hits.filter((h) => h.userMarker) },
        { key: 'dispatchMarkers', label: t('common.dispatch', 2), items: props.hits.filter((h) => h.dispatchMarker) },
        { key: 'markerMarkers', label: t('common.marker', 2), items: props.hits.filter((h) => h.markerMarker) },
    ].filter((group) => group.items.length > 0),
);

function clicked(h: Hit): void {
    h.openPopup();

    if (h.userMarker) selectedMarker.value = h.userMarker;
}
</script>

<template>
    <UCard
        class="-my-[13px] -mr-[24px] -ml-[20px] flex min-w-[200px] flex-col"
        :ui="{ header: 'mx-auto p-1 sm:px-2', body: 'p-1 sm:p-2 xl:mx-auto max-h-[90%]', footer: 'mx-auto p-1 sm:px-2' }"
    >
        <template #header>
            <div class="font-semibold">{{ $t('common.choose_one') }} ({{ hits.length }})</div>
        </template>

        <div class="space-y-1 divide-y divide-y-0 divide-default overflow-y-auto">
            <ul v-for="(group, idx) in groupedMarkers" :key="idx" class="space-y-1">
                <li class="text-center font-semibold">{{ group.label }}</li>
                <li v-for="(h, i) in group.items" :key="i">
                    <UButton block color="neutral" variant="soft" @click="() => clicked(h)">
                        <template v-if="h.markerMarker">
                            <UIcon
                                v-if="h.markerMarker?.data?.data.oneofKind === 'icon'"
                                :name="convertComponentIconNameToDynamic(h.markerMarker.data?.data.icon.icon)"
                                class="size-5"
                                :style="{ color: h.markerMarker.color ?? 'currentColor' }"
                            />

                            {{ h.markerMarker.name }}
                        </template>
                        <template v-else-if="h.userMarker">
                            {{ h.userMarker.user?.firstname }}
                            {{ h.userMarker.user?.lastname }} ({{ h.userMarker.jobLabel }})
                        </template>
                        <template v-else-if="h.dispatchMarker"> DSP-{{ h.dispatchMarker.id }} </template>
                        <template v-else>
                            {{ $t('common.unknown') }}
                        </template>
                    </UButton>
                </li>
            </ul>
        </div>

        <template v-if="hiddenCount > 0" #footer>
            <div class="font-semibold">{{ $t('common.hidden_markers') }}: {{ hiddenCount }}</div>
        </template>
    </UCard>
</template>
