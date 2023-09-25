<script lang="ts" setup>
import { TrainingModule } from '~~/gen/ts/resources/jobs/training';
import ListEntry from './ListEntry.vue';

const trainings = ref<TrainingModule[]>([
    {
        id: 0n,
        job: 'ambulance',
        title: 'Test Training Title',
        description: 'Test training module',
        minimumGrade: 0,
        open: true,
    },
]);
</script>

<template>
    <ul role="list" class="divide-y divide-gray-100">
        <ListEntry v-for="training in trainings" :training="training" />
    </ul>
</template>
