<script lang="ts" setup>
import { type PenaltiesSummary } from './Calculator.vue';

const props = defineProps<{
    summary: PenaltiesSummary;
}>();
</script>

<template>
    <div class="mx-auto max-w-7xl">
        <div class="mt-2 grid grid-cols-1 gap-2 sm:grid-cols-4">
            <div class="bg-gray-900 px-4 py-6 sm:px-6 lg:px-8">
                <p class="text-sm font-medium leading-6 text-gray-400">
                    {{ $t('common.fine') }}
                </p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-4xl font-semibold tracking-tight text-neutral">
                        {{ summary.fine.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ') }}
                    </span>
                    <span class="text-sm text-gray-400">$</span>
                </p>
            </div>
            <div class="bg-gray-900 px-4 py-6 sm:px-6 lg:px-8">
                <p class="text-sm font-medium leading-6 text-gray-400">
                    {{ $t('common.detention_time') }}
                </p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-4xl font-semibold tracking-tight text-neutral">
                        {{ summary.detentionTime }}
                    </span>
                    <span class="text-sm text-gray-400">
                        {{ $t('common.time_ago.month', summary.detentionTime.toString()) }}
                    </span>
                </p>
            </div>
            <div class="bg-gray-900 px-4 py-6 sm:px-6 lg:px-8">
                <p class="text-sm font-medium leading-6 text-gray-400">
                    {{ $t('common.traffic_infraction_points', 2) }}
                </p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-4xl font-semibold tracking-tight text-neutral">
                        {{ summary.stvoPoints }}
                    </span>
                </p>
            </div>
            <div class="bg-gray-900 px-4 py-6 sm:px-6 lg:px-8">
                <p class="text-sm font-medium leading-6 text-gray-400">
                    {{ $t('common.total_count') }}
                </p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-4xl font-semibold tracking-tight text-neutral">
                        {{ summary.count }}
                    </span>
                </p>
            </div>
        </div>
    </div>
</template>
