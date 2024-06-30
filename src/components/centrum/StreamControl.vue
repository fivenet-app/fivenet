<script lang="ts" setup>
import { useCentrumStore } from '~/store/centrum';
import { useLivemapStore } from '~/store/livemap';

const centrumStore = useCentrumStore();
const { abort: centrumAbort } = storeToRefs(centrumStore);

const livemapStore = useLivemapStore();
const { abort: livemapAbort } = storeToRefs(livemapStore);

async function start(): Promise<void> {
    if (can('CentrumService.Stream').value) {
        centrumStore.startStream();
    }

    if (can('LivemapperService.Stream').value) {
        livemapStore.startStream();
    }
}

async function stop(): Promise<void> {
    Promise.all([centrumStore.stopStream(), livemapStore.stopStream()]);
}

const running = computed(() => livemapAbort.value || centrumAbort.value);
</script>

<template>
    <div>
        <UButton v-if="running" variant="soft" icon="i-mdi-pause" class="flex-initial" @click="stop()"> Pause </UButton>
        <UButton v-else variant="soft" icon="i-mdi-play" class="flex-initial" @click="start()"> Start </UButton>
    </div>
</template>
