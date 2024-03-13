<script lang="ts" setup>
import { ChevronRightIcon } from 'mdi-vue3';
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
                <p class="text-sm font-semibold leading-6 text-gray-100">
                    <NuxtLink :to="{ name: 'jobs-qualifications-id', params: { id: qualification.id } }">
                        <span class="absolute inset-x-0 -top-px bottom-0" />
                        {{ qualification.abbreviation }}: {{ qualification.title }}
                    </NuxtLink>
                </p>
                <p class="mt-1 flex text-xs leading-5 text-gray-300">
                    <span class="inline-flex gap-1">
                        <span class="font-semibold">{{ $t('common.score') }}: {{ qualification.result?.score }}</span>
                        <span v-if="qualification.result?.summary">
                            ({{ $t('common.summary') }}: {{ qualification.result?.summary }})</span
                        >
                    </span>
                </p>
            </div>
        </div>
        <div class="flex shrink-0 items-center gap-x-4">
            <div class="hidden sm:flex sm:flex-col sm:items-end">
                <p v-if="qualification.result?.createdAt" class="mt-1 text-xs leading-5 text-gray-300">
                    {{ $t('common.created_at') }} <GenericTime :value="qualification.result.createdAt" />
                </p>
            </div>
            <ChevronRightIcon class="h-5 w-5 flex-none text-gray-300" aria-hidden="true" />
        </div>
    </li>
</template>
