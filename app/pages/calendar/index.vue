<script setup lang="ts">
import type { TypedRouteLocationRaw } from '@typed-router';
import { routeForCalendarView } from '~/utils/calendar-view';

useHead({
    title: 'common.calendar',
});

definePageMeta({
    title: 'common.calendar',
    requiresAuth: true,
    redirect: routeForCalendarView('month', new Date()) as unknown as TypedRouteLocationRaw<string>,
});
</script>

<template>
    <div />
</template>
