<script lang="ts" setup>
import type { PenaltyCalculatorData } from '~~/gen/ts/resources/documents/data/data';
import PenaltyStats from './PenaltyStats.vue';
import PenaltySummaryTable from './PenaltySummaryTable.vue';
import { calculatePenaltySummary, resolvePenaltyCalculatorSelection } from './helpers';

const props = defineProps<{
    data?: PenaltyCalculatorData;
}>();

const completorStore = useCompletorStore();
const { data: lawBooks } = useLazyAsyncData(`lawbooks`, () => completorStore.listLawBooks());

const selectedPenalties = computed(() => resolvePenaltyCalculatorSelection(props.data, lawBooks.value ?? []));
const reduction = computed(() => props.data?.reduction ?? 0);
const summary = computed(() => calculatePenaltySummary(selectedPenalties.value));
</script>

<template>
    <div class="my-2 flex flex-col gap-2">
        <PenaltySummaryTable
            :law-books="lawBooks ?? []"
            :selected-laws="selectedPenalties"
            :reduction="reduction"
            disable-line-clamp
            :ui="{
                tr: 'border-neutral-700',
            }"
        />
        <PenaltyStats :summary="summary" :reduction="reduction" compact hide-warning />
    </div>
</template>
