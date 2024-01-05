<script lang="ts" setup>
import { Qualification } from '~~/gen/ts/resources/jobs/qualifications';
import QualificationsListEntry from '~/components/jobs/qualifications/QualificationsListEntry.vue';

const trainings = ref<Qualification[]>([
    {
        id: '0',
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
        <QualificationsListEntry v-for="training in trainings" :key="training.id" :qualification="training" />
    </ul>
</template>
