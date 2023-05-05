<script lang="ts" setup>
import { PenaltiesSummary } from '~/utils/penalty';

defineProps<{
    summary: PenaltiesSummary,
}>();
</script>

<template>
    <div class="mx-auto max-w-7xl">
        <h3 class="text-base font-semibold leading-6 text-white-900">Ergebnis</h3>
        <div class="mt-5 grid grid-cols-1 gap-2 sm:grid-cols-4">
            <div class="bg-gray-900 px-4 py-6 sm:px-6 lg:px-8">
                <p class="text-sm font-medium leading-6 text-gray-400">
                    Geldsumme
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
                    Geldsumme
                </p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-4xl font-semibold tracking-tight text-white">
                        {{ summary.detentionTime }}
                    </span>
                    <span class="text-sm text-gray-400">Monate</span>
                </p>
            </div>
            <div class="bg-gray-900 px-4 py-6 sm:px-6 lg:px-8">
                <p class="text-sm font-medium leading-6 text-gray-400">
                    StVo-Punkte
                </p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-4xl font-semibold tracking-tight text-white">
                        {{ summary.stvoPoints }}
                    </span>
                </p>
            </div>
            <div class="bg-gray-900 px-4 py-6 sm:px-6 lg:px-8">
                <p class="text-sm font-medium leading-6 text-gray-400">
                    Gesamtanzahl
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
