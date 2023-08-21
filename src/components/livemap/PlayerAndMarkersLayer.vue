<script lang="ts" setup>
import { LControl, LLayerGroup } from '@vue-leaflet/vue-leaflet';
import 'leaflet-contextmenu';
import 'leaflet-contextmenu/dist/leaflet.contextmenu.min.css';
import 'leaflet/dist/leaflet.css';
import { useAuthStore } from '~/store/auth';
import { useLivemapStore } from '~/store/livemap';
import { UserMarker } from '~~/gen/ts/resources/livemap/livemap';
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

const livemapStore = useLivemapStore();
const { markers, jobs } = storeToRefs(livemapStore);
const { startStream, stopStream } = livemapStore;

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);

const playerQuery = ref<string>('');
const playerMarkersFiltered = computed(() =>
    markers.value.users.filter((m) => (m.user?.firstname + ' ' + m.user?.lastname).includes(playerQuery.value)),
);

onBeforeMount(async () => startStream());

onBeforeUnmount(() => stopStream());
</script>

<template>
    <LLayerGroup
        v-for="job in jobs.users"
        :key="`users_${job.name}`"
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
            <input
                v-model="playerQuery"
                class="w-full p-0.5 px-1 bg-clip-padding rounded-md border-2 border-black/20"
                type="text"
                name="searchPlayer"
                :placeholder="`${$t('common.employee', 1)} ${$t('common.filter')}`"
            />
        </div>
    </LControl>
</template>
