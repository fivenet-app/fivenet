<script setup lang="ts">
import { PX_PER_MINUTE } from '~/utils/calendar-view';

const now = useNow({ interval: 60_000 });

const top = computed(() => (now.value.getHours() * 60 + now.value.getMinutes()) * PX_PER_MINUTE);
</script>

<template>
    <div class="pointer-events-none absolute inset-x-0 z-10 border-t-2 border-error" :style="{ top: `${top}px` }">
        <div class="absolute -start-1 -top-[5px] size-2 rounded-full bg-error" />
    </div>
</template>
