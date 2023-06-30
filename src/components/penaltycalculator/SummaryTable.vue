<script lang="ts" setup>
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiCalculator } from '@mdi/js';
import { SelectedPenalty } from '~/components/penaltycalculator/Calculator.vue';
import { LawBook } from '../../../gen/ts/resources/laws/laws';

const props = defineProps<{
    lawBooks: Array<LawBook>;
    selectedLaws: Array<SelectedPenalty>;
}>();

function getNameForLawBookId(id: bigint): string | undefined {
    return props.lawBooks.filter((b) => b.id === id)[0].name;
}
</script>

<template>
    <button
        v-if="selectedLaws.length === 0"
        type="button"
        disabled
        class="relative block w-full p-12 text-center border-2 border-gray-300 border-dashed rounded-lg hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
    >
        <SvgIcon class="w-12 h-12 mx-auto text-neutral" type="mdi" :path="mdiCalculator" />
        <span class="block mt-2 text-sm font-semibold text-gray-300">
            {{ $t('common.none_selected', [`${$t('common.crime')}`]) }}
        </span>
    </button>
    <table v-else class="min-w-full divide-y divide-base-600">
        <thead>
            <tr>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                    {{ $t('common.crime') }}
                </th>
                <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                    {{ $t('common.fine') }}
                </th>
                <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                    {{ $t('common.detention_time') }}
                </th>
                <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                    {{ $t('common.traffic_infraction_points', 2) }}
                </th>
                <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                    {{ $t('common.other') }}
                </th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                    {{ $t('common.count') }}
                </th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="p in selectedLaws" :key="`${p.law.lawbookId}-${p.law.id.toString()}`">
                <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
                    {{ getNameForLawBookId(p.law.lawbookId) }} - {{ p.law.name }}
                </td>
                <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">${{ p.law.fine ? p.law.fine * p.count : 0 }}</td>
                <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
                    {{ p.law.detentionTime ? p.law.detentionTime * p.count : 0 }}
                </td>
                <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
                    {{ p.law.stvoPoints ? p.law.stvoPoints * p.count : 0 }}
                </td>
                <td class="px-2 py-2 text-sm text-base-200">
                    {{ p.law.description }}
                </td>
                <td class="w-20 min-w-20 px-2 py-2 text-sm text-base-200">
                    {{ p.count }}
                </td>
            </tr>
        </tbody>
    </table>
</template>
