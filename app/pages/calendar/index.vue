<script setup lang="ts">
import { routeForCalendarView } from '~/utils/calendar-view';

useHead({
    title: 'common.calendar',
});

definePageMeta({
    title: 'common.calendar',
    requiresAuth: true,
    redirect: () => routeForCalendarView('month', new Date()),
});
</script>

<template>
    <div />
</template>
