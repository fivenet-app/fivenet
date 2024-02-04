<script lang="ts" setup>
import QualificationsList from '~/components/jobs/qualifications/QualificationsList.vue';

useHead({
    title: 'pages.jobs.qualifications.title',
});
definePageMeta({
    title: 'pages.jobs.qualifications.title',
    requiresAuth: true,
    permission: 'JobsQualificationsService.ListTrainings',
});
</script>

<template>
    <div>
        <div class="py-2 pb-14">
            <div class="px-1 sm:px-2 lg:px-4">
                <div class="grid-col-2 grid gap-2">
                    <div class="sm:flex-auto">
                        <div class="overflow-hidden rounded-md bg-neutral shadow">
                            <div class="border-b border-gray-200 bg-neutral px-4 py-5 sm:px-6">
                                <div class="-ml-4 -mt-4 flex flex-wrap items-center justify-between sm:flex-nowrap">
                                    <div class="ml-4 mt-4">
                                        <h3 class="text-base font-semibold leading-6 text-gray-900">
                                            Available Qualifications
                                        </h3>
                                        <p class="mt-1 text-sm text-gray-500">
                                            Lorem ipsum dolor sit amet consectetur adipisicing elit quam corrupti consectetur.
                                        </p>
                                    </div>
                                    <div class="ml-4 mt-4 flex-shrink-0">
                                        <button
                                            type="button"
                                            class="relative inline-flex items-center rounded-md bg-primary-500 px-3 py-2 text-sm font-semibold text-neutral shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600"
                                        >
                                            Create new Qualification
                                        </button>
                                    </div>
                                </div>
                            </div>

                            <QualificationsList />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
