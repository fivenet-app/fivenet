<script lang="ts" setup>
import { ChevronRightIcon, LockIcon, LockOpenVariantIcon } from 'mdi-vue3';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
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
                    <NuxtLink :to="{ name: 'jobs-qualifications-id', params: { id: qualification.id } }">
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        {{ qualification.abbreviation }}: {{ qualification.title }}
                    </NuxtLink>
                </p>
                <p class="mt-1 flex text-xs leading-5 text-gray-500"></p>
            </div>
        </div>
        <div class="flex shrink-0 items-center gap-x-4">
            <div class="hidden sm:flex sm:flex-col sm:items-end">
                <div v-if="qualification.closed" class="flex flex-initial flex-row gap-1 rounded-full bg-error-100 px-2 py-1">
                    <LockIcon class="h-5 w-5 text-error-400" aria-hidden="true" />
                    <span class="text-sm font-medium text-error-700">
                        {{ $t('common.close', 2) }}
                    </span>
                </div>
                <div v-else class="flex flex-initial flex-row gap-1 rounded-full bg-success-100 px-2 py-1">
                    <LockOpenVariantIcon class="h-5 w-5 text-success-500" aria-hidden="true" />
                    <span class="text-sm font-medium text-success-700">
                        {{ $t('common.open', 2) }}
                    </span>
                </div>
                <p v-if="qualification.createdAt" class="mt-1 text-xs leading-5 text-gray-500">
                    {{ $t('common.created_at') }} <GenericTime :value="qualification.createdAt" />
                </p>
            </div>
            <ChevronRightIcon class="h-5 w-5 flex-none text-gray-400" aria-hidden="true" />
        </div>
    </li>
</template>
