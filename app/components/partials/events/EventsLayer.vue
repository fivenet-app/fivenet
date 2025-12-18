<script lang="ts" setup>
const settingsStore = useSettingsStore();
const { eventsDisabled } = storeToRefs(settingsStore);

// Use client date to show any event overlays
const now = new Date();
const showSnowflakes = now.getMonth() + 1 === 12 && now.getDate() >= 21 && now.getDate() <= 26;
</script>

<template>
    <LazyPartialsEventsSnowflakesContainer v-if="!eventsDisabled && showSnowflakes" />
</template>
