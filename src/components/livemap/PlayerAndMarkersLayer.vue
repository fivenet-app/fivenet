<script lang="ts" setup>
import { LControl, LLayerGroup } from '@vue-leaflet/vue-leaflet';
import { LeafletMouseEvent } from 'leaflet';
import 'leaflet-contextmenu';
import 'leaflet-contextmenu/dist/leaflet.contextmenu.min.css';
import 'leaflet/dist/leaflet.css';
import CreateOrUpdateModal from '~/components/centrum/dispatches/CreateOrUpdateModal.vue';
import { useAuthStore } from '~/store/auth';
import { useLivemapStore } from '~/store/livemap';
import { UserMarker } from '~~/gen/ts/resources/livemap/livemap';
import BaseMap from './BaseMap.vue';
import PlayerMarker from './PlayerMarker.vue';

withDefaults(
    defineProps<{
        centerSelectedMarker?: boolean;
        filterPlayers?: boolean;
    }>(),
    {
        centerSelectedMarker: true,
        filterPlayers: true,
    },
);

defineEmits<{
    (e: 'markerSelected', marker: UserMarker): void;
}>();

const { t } = useI18n();

const livemapStore = useLivemapStore();
const { markers, jobs } = storeToRefs(livemapStore);
const { startStream, stopStream } = livemapStore;

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const location = ref<{ x: number; y: number }>({ x: 0, y: 0 });
defineExpose({
    location,
});

interface ContextmenuItem {
    text: string;
    callback: (e: LeafletMouseEvent) => void;
}

const mapOptions = {
    contextmenu: true,
    contextmenuWidth: 150,
    contextmenuItems: [] as ContextmenuItem[],
};

if (can('CentrumService.CreateDispatch')) {
    mapOptions.contextmenuItems.push({
        text: t('components.centrum.create_dispatch.title'),
        callback: (e: LeafletMouseEvent) => {
            goto({ x: e.latlng.lat, y: e.latlng.lng });
            createDispatchOpen.value = true;
        },
    });
}

const playerQuery = ref<string>('');
const playerMarkersFiltered = computed(() =>
    markers.value.users.filter((m) => (m.user?.firstname + ' ' + m.user?.lastname).includes(playerQuery.value)),
);

watch(location, () => goto({ x: location.value?.x, y: location.value?.y }));

onBeforeMount(async () => startStream());

onBeforeUnmount(() => stopStream());

const livemapComponent = ref<InstanceType<typeof BaseMap>>();

const createOrUpdateModal = ref<InstanceType<typeof CreateOrUpdateModal>>();
const createDispatchOpen = ref(false);

function goto(e: { x: number; y: number }) {
    if (createOrUpdateModal.value) {
        createOrUpdateModal.value.location = { x: e.x, y: e.y };
    }

    if (livemapComponent.value) {
        location.value = { x: e.x, y: e.y };
    }
}
</script>

<template>
    <LLayerGroup
        v-for="job in jobs.users"
        :key="job.name"
        :name="`${$t('common.employee', 2)} ${job.label}`"
        layer-type="overlay"
        :visible="true"
    >
        <PlayerMarker
            v-for="marker in playerMarkersFiltered.filter((p) => p.user?.job === job.name)"
            :key="marker.marker!.id?.toString()"
            :marker="marker"
            :active-char="activeChar"
            @selected="$emit('markerSelected', marker)"
        />
    </LLayerGroup>

    <LControl position="topleft" v-if="filterPlayers">
        <div class="form-control flex flex-col gap-2">
            <div v-if="filterPlayers">
                <input
                    v-model="playerQuery"
                    class="w-full"
                    type="text"
                    name="searchPlayer"
                    :placeholder="`${$t('common.employee', 1)} ${$t('common.filter')}`"
                />
            </div>
        </div>
    </LControl>
</template>
