<script lang="ts" setup>
import { SelectedPenalty } from '~/utils/penalty';

defineProps<{
    selectedPenalties: Array<SelectedPenalty>,
}>();
</script>

<template>
    <button v-if="selectedPenalties.length == 0" type="button" disabled
            class="relative block w-full p-12 text-center border-2 border-gray-300 border-dashed rounded-lg hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
            <UserIcon class="w-12 h-12 mx-auto text-neutral" />
            <span class="block mt-2 text-sm font-semibold text-gray-300">
                Keine Verbrechen ausgewählt
            </span>
        </button>
    <table v-else class="min-w-full divide-y divide-base-600">
        <thead>
            <tr>
                <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-neutral sm:pl-0">
                    Straftat
                </th>
                <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                    Geldstrafe
                </th>
                <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                    Haftzeit
                </th>
                <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                    StVo-Punkte
                </th>
                <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-neutral">
                    Sonstige
                </th>
                <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0 text-right text-sm font-semibold text-neutral">
                    Anzahl
                </th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="p in selectedPenalties" :key="p.penalty.name">
                <td class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
                    {{ p.penalty.category }} - {{ p.penalty.name }}
                </td>
                <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
                    ${{ p.penalty.fine ? p.penalty.fine * p.count : 0 }}
                </td>
                <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
                    {{ p.penalty.detentionTime ? p.penalty.detentionTime * p.count : 0 }}
                </td>
                <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
                    {{ p.penalty.stvoPoints ? p.penalty.stvoPoints * p.count : 0 }}
                </td>
                <td class="px-2 py-2 text-sm text-base-200">
                    {{ p.penalty.description }}
                </td>
                <td class="w-20 min-w-20 px-2 py-2 text-sm text-base-200">
                    {{ p.count }}
                </td>
            </tr>
        </tbody>
    </table>
</template>
