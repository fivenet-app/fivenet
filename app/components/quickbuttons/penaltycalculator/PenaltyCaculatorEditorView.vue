<script lang="ts" setup>
import { NodeViewWrapper, type NodeViewProps } from '@tiptap/vue-3';
import PenaltyCalculatorDrawer from '~/components/quickbuttons/penaltycalculator/PenaltyCalculatorDrawer.vue';
import type { DocumentData } from '~~/gen/ts/resources/documents/data/data';
import PenaltyStats from './PenaltyStats.vue';
import PenaltySummaryTable from './PenaltySummaryTable.vue';
import {
    calculatePenaltySummary,
    resolvePenaltyCalculatorSelection,
    toPenaltyCalculatorData,
    type SelectedPenalty,
} from './helpers';

const props = defineProps<NodeViewProps>();

const completorStore = useCompletorStore();
const { data: lawBooks, refresh } = useLazyAsyncData(`lawbooks`, () => completorStore.listLawBooks());

const selectedPenalties = useState<SelectedPenalty[]>('quickButton:penaltyCalculator:selected', () => [] as SelectedPenalty[]);
const reduction = useState<number>('quickButton:penaltyCalculator:reduction', () => 0);

const localDocumentData = ref<DocumentData | undefined>();
const documentData = inject<Ref<DocumentData | undefined>>('documents:editor:data', localDocumentData);
const disablePenaltyCalculatorBlockEditing = inject<Ref<boolean>>('tiptap:disablePenaltyCalculatorBlockEditing', ref(false));

const overlay = useOverlay();
const penaltyCalculatorDrawer = overlay.create(PenaltyCalculatorDrawer);

const isActionsOpen = ref(false);
const confirmDelete = ref(false);

let confirmDeleteTimer: ReturnType<typeof setTimeout> | undefined = undefined;

const persistedPenaltyData = computed(() => documentData.value?.penaltyCalculator);
const displayedSelectedPenalties = computed(() =>
    resolvePenaltyCalculatorSelection(persistedPenaltyData.value, lawBooks.value ?? []),
);
const displayedReduction = computed(() => persistedPenaltyData.value?.reduction ?? 0);
const summary = computed(() => calculatePenaltySummary(displayedSelectedPenalties.value));

function resetDeleteConfirm(): void {
    if (confirmDeleteTimer) clearTimeout(confirmDeleteTimer);
    confirmDelete.value = false;
}

function ensureDocumentDataContainer(): DocumentData {
    if (!documentData.value) documentData.value = {};
    return documentData.value;
}

async function openEditDrawer(): Promise<void> {
    if (!props.editor.isEditable) return;
    if (!lawBooks.value || lawBooks.value.length === 0) await refresh();

    const snapshotSelected = [...selectedPenalties.value];
    const snapshotReduction = reduction.value;

    selectedPenalties.value = resolvePenaltyCalculatorSelection(documentData.value?.penaltyCalculator, lawBooks.value ?? []);
    reduction.value = documentData.value?.penaltyCalculator?.reduction ?? 0;

    penaltyCalculatorDrawer.open({
        requireExplicitSave: true,
        onSave: () => {
            ensureDocumentDataContainer().penaltyCalculator = toPenaltyCalculatorData(selectedPenalties.value, reduction.value);
        },
        onCancel: () => {
            selectedPenalties.value = snapshotSelected;
            reduction.value = snapshotReduction;
        },
    });
}

function deleteNodeWithConfirm(): void {
    if (!props.editor.isEditable) return;

    if (!confirmDelete.value) {
        confirmDelete.value = true;
        if (confirmDeleteTimer) clearTimeout(confirmDeleteTimer);
        confirmDeleteTimer = setTimeout(() => {
            confirmDelete.value = false;
            confirmDeleteTimer = undefined;
        }, 3000);
        return;
    }

    resetDeleteConfirm();
    if (documentData.value?.penaltyCalculator) documentData.value.penaltyCalculator = undefined;
    props.deleteNode();
}

onBeforeUnmount(() => resetDeleteConfirm());
</script>

<template>
    <NodeViewWrapper>
        <div
            class="relative my-2 flex flex-col gap-2"
            @mouseenter="isActionsOpen = true"
            @mouseleave="
                isActionsOpen = false;
                resetDeleteConfirm();
            "
            @click="isActionsOpen = true"
        >
            <div
                v-if="editor.isEditable && !disablePenaltyCalculatorBlockEditing && (isActionsOpen || selected)"
                class="absolute top-2 right-2 z-20 flex items-center gap-1 rounded-md border border-neutral-300 bg-white/90 p-1 shadow-sm dark:border-neutral-700 dark:bg-neutral-900/90"
            >
                <UButton size="xs" color="neutral" variant="soft" icon="i-mdi-pencil" @click="openEditDrawer" />
                <UButton
                    size="xs"
                    :color="confirmDelete ? 'error' : 'neutral'"
                    :variant="confirmDelete ? 'solid' : 'soft'"
                    :icon="confirmDelete ? 'i-mdi-check' : 'i-mdi-delete'"
                    @click="deleteNodeWithConfirm"
                />
            </div>

            <PenaltySummaryTable
                :law-books="lawBooks ?? []"
                :selected-laws="displayedSelectedPenalties"
                :reduction="displayedReduction"
                disable-line-clamp
            />

            <PenaltyStats :summary="summary" :reduction="displayedReduction" compact />
        </div>
    </NodeViewWrapper>
</template>
