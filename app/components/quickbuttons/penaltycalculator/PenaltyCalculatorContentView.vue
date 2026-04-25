<script lang="ts" setup>
import type { ComputedRef, Ref } from 'vue';
import type { DocumentData, PenaltyCalculatorData } from '~~/gen/ts/resources/documents/data/data';
import PenaltyStats from './PenaltyStats.vue';
import PenaltySummaryTable from './PenaltySummaryTable.vue';
import { calculatePenaltySummary, resolvePenaltyCalculatorSelection } from './helpers';

const props = defineProps<{
    data?: PenaltyCalculatorData;
}>();

const completorStore = useCompletorStore();
const { data: lawBooks } = useLazyAsyncData(`lawbooks`, () => completorStore.listLawBooks());

const documentData = inject<ComputedRef<DocumentData | undefined> | Ref<DocumentData | undefined>>('documents:content:data');

const data = computed(() => props.data ?? documentData?.value?.penaltyCalculator);
const selectedPenalties = computed(() => resolvePenaltyCalculatorSelection(data.value, lawBooks.value ?? []));
const reduction = computed(() => data.value?.reduction ?? 0);
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
