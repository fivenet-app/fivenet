<script lang="ts" setup>
// Use client date to show any event overlays
const now = new Date();
const showSnowflakes = now.getMonth() + 1 === 12 && now.getDate() >= 21 && now.getDate() <= 26;
</script>

<template>
    <LazyPartialsEventsSnowflakesContainer v-if="showSnowflakes" />
</template>
