<script lang="ts" setup>
import List from '~/components/jobs/trainings/List.vue';

useHead({
    title: 'pages.jobs.trainings.title',
});
definePageMeta({
    title: 'pages.jobs.trainings.title',
    requiresAuth: true,
    permission: 'Jobs.TrainingsList',
});
</script>

<template>
    <div class="flex">
        <div class="overflow-hidden rounded-md bg-white shadow">
            <div class="border-b border-gray-200 bg-white px-4 py-5 sm:px-6">
                <div class="-ml-4 -mt-4 flex flex-wrap items-center justify-between sm:flex-nowrap">
                    <div class="ml-4 mt-4">
                        <h3 class="text-base font-semibold leading-6 text-gray-900">Available Trainings</h3>
                        <p class="mt-1 text-sm text-gray-500">
                            Lorem ipsum dolor sit amet consectetur adipisicing elit quam corrupti consectetur.
                        </p>
                    </div>
                    <div class="ml-4 mt-4 flex-shrink-0">
                        <button
                            type="button"
                            class="relative inline-flex items-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                        >
                            Create new job
                        </button>
                    </div>
                </div>
            </div>
            <List />
        </div>
    </div>
</template>
