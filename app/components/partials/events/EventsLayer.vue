<script lang="ts" setup>
const settingsStore = useSettingsStore();
const { eventsDisabled, eventsShowSnowflakes } = storeToRefs(settingsStore);
</script>

<template>
    <LazyPartialsEventsSnowflakesContainer v-if="!eventsDisabled && eventsShowSnowflakes" />
</template>
