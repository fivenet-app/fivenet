<script lang="ts" setup>
import { ChevronRightIcon } from 'mdi-vue3';
import Time from '~/components/partials/elements/Time.vue';
import { Qualification } from '~~/gen/ts/resources/jobs/qualifications';

defineProps<{
    qualification: Qualification;
}>();
</script>

<template>
    <li class="relative flex justify-between px-4 py-5">
        <div class="flex min-w-0 gap-x-4">
            <div class="min-w-0 flex-auto">
                <p class="text-sm font-semibold leading-6 text-gray-900">
                    <a href="training.id">
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        {{ qualification.title }}
                    </a>
                </p>
                <p class="mt-1 flex text-xs leading-5 text-gray-500">
                    {{ qualification.description }}
                </p>
            </div>
        </div>
        <div class="flex shrink-0 items-center gap-x-4">
            <div class="hidden sm:flex sm:flex-col sm:items-end">
                <p class="text-sm leading-6 text-gray-900">{{ $t('common.rank') }}: {{ qualification.jobAccess }}</p>
                <p v-if="qualification.createdAt" class="mt-1 text-xs leading-5 text-gray-500">
                    {{ $t('common.created_at') }} <Time :value="qualification.createdAt" />
                </p>
            </div>
            <ChevronRightIcon class="h-5 w-5 flex-none text-gray-400" aria-hidden="true" />
        </div>
    </li>
</template>
