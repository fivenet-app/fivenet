<script lang="ts" setup>
import { type NodeViewProps, NodeViewWrapper } from '@tiptap/vue-3';
import type { SelectedPenalty } from './helpers';
import PenaltyStats from './PenaltyStats.vue';
import PenaltySummaryTable from './PenaltySummaryTable.vue';

defineProps<NodeViewProps>();

const completorStore = useCompletorStore();
const { data: lawBooks, status, refresh, error } = useLazyAsyncData(`lawbooks`, () => completorStore.listLawBooks());

const selectedPenalties = useState<SelectedPenalty[]>('quickButton:penaltyCalculator:selected', () => [] as SelectedPenalty[]);
const reduction = useState<number>('quickButton:penaltyCalculator:reduction', () => 0);

const summary = computed(() => ({
    fine: selectedPenalties.value.reduce((acc, curr) => acc + (curr.law.fine ? curr.law.fine * curr.count : 0), 0),
    detentionTime: selectedPenalties.value.reduce(
        (acc, curr) => acc + (curr.law.detentionTime ? curr.law.detentionTime * curr.count : 0),
        0,
    ),
    stvoPoints: selectedPenalties.value.reduce(
        (acc, curr) => acc + (curr.law.stvoPoints ? curr.law.stvoPoints * curr.count : 0),
        0,
    ),
    count: selectedPenalties.value.reduce((acc, curr) => acc + curr.count, 0),
}));

// TODO
</script>

<template>
    <NodeViewWrapper>
        <div class="flex flex-col gap-2">
            <PenaltySummaryTable :law-books="lawBooks ?? []" :selected-laws="selectedPenalties" :reduction="reduction" />

            <PenaltyStats :summary="summary" :reduction="reduction" compact />
        </div>
    </NodeViewWrapper>
</template>
