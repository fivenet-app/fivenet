<script lang="ts" setup>
const settingsStore = useSettingsStore();
const { eventsShowSnowflakes } = storeToRefs(settingsStore);
</script>

<template>
    <LazyPartialsEventsSnowflakesContainer v-if="eventsShowSnowflakes" />
</template>
