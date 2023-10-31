<script lang="ts" setup>
import { Qualification } from '~~/gen/ts/resources/jobs/qualifications';
import ListEntry from '~/components/jobs/qualifications/ListEntry.vue';

const trainings = ref<Qualification[]>([
    {
        id: 0n,
        job: 'ambulance',
        title: 'Test Training Title',
        description: 'Test training module',
        open: true,
        abbreviation: 'TR1',
        jobAccess: [],
        weight: 0,
    },
]);
</script>

<template>
    <ul role="list" class="divide-y divide-gray-100">
        <ListEntry v-for="training in trainings" :key="training.id.toString()" :qualification="training" />
    </ul>
</template>
