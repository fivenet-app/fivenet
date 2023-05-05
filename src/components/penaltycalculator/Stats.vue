<script lang="ts" setup>
import { PenaltiesSummary } from '~/utils/penalty';

defineProps<{
    summary: PenaltiesSummary,
}>();
</script>

<template>
    <div class="mx-auto max-w-7xl">
        <div class="mt-2 grid grid-cols-1 gap-2 sm:grid-cols-4">
            <div class="bg-gray-900 px-4 py-6 sm:px-6 lg:px-8">
                <p class="text-sm font-medium leading-6 text-gray-400">
                    {{ $t('components.penaltycalculator.fine') }}
                </p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-4xl font-semibold tracking-tight text-white">
                        {{ summary.fine.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ") }}
                    </span>
                    <span class="text-sm text-gray-400">$</span>
                </p>
            </div>
            <div class="bg-gray-900 px-4 py-6 sm:px-6 lg:px-8">
                <p class="text-sm font-medium leading-6 text-gray-400">
                    {{ $t('components.penaltycalculator.detention_time') }}
                </p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-4xl font-semibold tracking-tight text-white">
                        {{ summary.detentionTime }}
                    </span>
                    <span class="text-sm text-gray-400">
                        {{ $t('common.time_ago.month', summary.detentionTime) }}
                    </span>
                </p>
            </div>
            <div class="bg-gray-900 px-4 py-6 sm:px-6 lg:px-8">
                <p class="text-sm font-medium leading-6 text-gray-400">
                    {{ $t('components.penaltycalculator.stvo_points', 2) }}
                </p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-4xl font-semibold tracking-tight text-white">
                        {{ summary.stvoPoints }}
                    </span>
                </p>
            </div>
            <div class="bg-gray-900 px-4 py-6 sm:px-6 lg:px-8">
                <p class="text-sm font-medium leading-6 text-gray-400">
                    {{ $t('common.total_count') }}
                </p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-4xl font-semibold tracking-tight text-white">
                        {{ summary.count }}
                    </span>
                </p>
            </div>
        </div>
    </div>
</template>
