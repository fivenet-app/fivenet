<script lang="ts" setup>
import { CalculatorIcon } from 'mdi-vue3';
import { type SelectedPenalty } from '~/components/quickbuttons/penaltycalculator/PenaltyCalculator.vue';
import { LawBook } from '~~/gen/ts/resources/laws/laws';

const props = defineProps<{
    lawBooks: LawBook[];
    selectedLaws: SelectedPenalty[];
}>();

function getNameForLawBookId(id: string): string | undefined {
    return props.lawBooks.filter((b) => b.id === id)[0].name;
}
</script>

<template>
    <button
        v-if="selectedLaws.length === 0"
        type="button"
        disabled
        class="relative block w-full rounded-lg border-2 border-dashed border-gray-300 p-12 text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
    >
        <CalculatorIcon class="mx-auto size-12 text-neutral" aria-hidden="true" />
        <span class="mt-2 block text-sm font-semibold text-gray-300">
            {{ $t('common.none_selected', [`${$t('common.crime')}`]) }}
        </span>
    </button>
    <table v-else class="min-w-full divide-y divide-base-600">
        <thead>
            <tr>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-1">
                    {{ $t('common.crime') }}
                </th>
                <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                    {{ $t('common.fine') }}
                </th>
                <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                    {{ $t('common.detention_time') }}
                </th>
                <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                    {{ $t('common.traffic_infraction_points', 2) }}
                </th>
                <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-neutral">
                    {{ $t('common.other') }}
                </th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 text-right text-sm font-semibold text-neutral sm:pr-0">
                    {{ $t('common.count') }}
                </th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="p in selectedLaws" :key="`${p.law.lawbookId}-${p.law.id}`">
                <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
                    {{ getNameForLawBookId(p.law.lawbookId) }} - {{ p.law.name }}
                </td>
                <td class="whitespace-nowrap p-1 text-left text-accent-200">${{ p.law.fine ? p.law.fine * p.count : 0 }}</td>
                <td class="whitespace-nowrap p-1 text-left text-accent-200">
                    {{ p.law.detentionTime ? p.law.detentionTime * p.count : 0 }}
                </td>
                <td class="whitespace-nowrap p-1 text-left text-accent-200">
                    {{ p.law.stvoPoints ? p.law.stvoPoints * p.count : 0 }}
                </td>
                <td class="break-all p-1 text-left text-sm text-accent-200">
                    {{ p.law.description }}
                </td>
                <td class="w-20 min-w-20 p-1 text-left text-accent-200">
                    {{ p.count }}
                </td>
            </tr>
        </tbody>
    </table>
</template>
