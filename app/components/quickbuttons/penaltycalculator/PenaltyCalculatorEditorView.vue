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
const enablePenaltyCalculatorBlockRemoval = inject<Ref<boolean>>('tiptap:enablePenaltyCalculatorBlockRemoval', ref(false));

const overlay = useOverlay();
const penaltyCalculatorDrawer = overlay.create(PenaltyCalculatorDrawer);

const isActionsOpen = ref(false);
const confirmDelete = ref(false);
const confirmClear = ref(false);

const { start: startConfirmDeleteTimeout, stop: stopConfirmDeleteTimeout } = useTimeoutFn(
    () => {
        confirmDelete.value = false;
    },
    3000,
    { immediate: false },
);
const { start: startConfirmClearTimeout, stop: stopConfirmClearTimeout } = useTimeoutFn(
    () => {
        confirmClear.value = false;
    },
    3000,
    { immediate: false },
);

const persistedPenaltyData = computed(() => documentData.value?.penaltyCalculator);
const displayedSelectedPenalties = computed(() =>
    resolvePenaltyCalculatorSelection(persistedPenaltyData.value, lawBooks.value ?? []),
);
const displayedReduction = computed(() => persistedPenaltyData.value?.reduction ?? 0);
const summary = computed(() => calculatePenaltySummary(displayedSelectedPenalties.value));

function resetDeleteConfirm(): void {
    stopConfirmDeleteTimeout();
    confirmDelete.value = false;
}

function resetClearConfirm(): void {
    stopConfirmClearTimeout();
    confirmClear.value = false;
}

function resetConfirmActions(): void {
    resetDeleteConfirm();
    resetClearConfirm();
}

function clearWithConfirm(): void {
    if (!props.editor.isEditable) return;

    if (!confirmClear.value) {
        confirmClear.value = true;
        resetDeleteConfirm();
        stopConfirmClearTimeout();
        startConfirmClearTimeout();
        return;
    }

    resetClearConfirm();
    if (documentData.value?.penaltyCalculator) documentData.value.penaltyCalculator = undefined;
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
        resetClearConfirm();
        stopConfirmDeleteTimeout();
        startConfirmDeleteTimeout();
        return;
    }

    resetDeleteConfirm();
    if (documentData.value?.penaltyCalculator) documentData.value.penaltyCalculator = undefined;
    props.deleteNode();
}

onBeforeUnmount(() => resetConfirmActions());
</script>

<template>
    <NodeViewWrapper>
        <div
            class="relative my-2 flex flex-col gap-2"
            @mouseenter="isActionsOpen = true"
            @mouseleave="
                isActionsOpen = false;
                resetConfirmActions();
            "
            @click="isActionsOpen = true"
        >
            <UFieldGroup
                v-if="
                    editor.isEditable &&
                    (!disablePenaltyCalculatorBlockEditing || enablePenaltyCalculatorBlockRemoval) &&
                    (isActionsOpen || selected)
                "
                class="absolute top-2 right-2 z-20 rounded-md border border-neutral-300 shadow-sm dark:border-neutral-700 dark:bg-neutral-900/90"
            >
                <template v-if="!disablePenaltyCalculatorBlockEditing">
                    <UTooltip :text="$t('common.edit')">
                        <UButton
                            size="xs"
                            color="neutral"
                            variant="soft"
                            icon="i-mdi-pencil"
                            :label="$t('common.edit')"
                            @click="openEditDrawer"
                        />
                    </UTooltip>

                    <UTooltip :text="confirmClear ? $t('common.confirm') : $t('common.clear')">
                        <UButton
                            size="xs"
                            :color="confirmClear ? 'error' : 'neutral'"
                            :variant="confirmClear ? 'solid' : 'soft'"
                            :trailing-icon="confirmClear ? 'i-mdi-check' : 'i-mdi-clear'"
                            :label="$t('common.clear')"
                            @click="clearWithConfirm"
                        />
                    </UTooltip>
                </template>

                <UTooltip
                    v-if="enablePenaltyCalculatorBlockRemoval"
                    :text="confirmDelete ? $t('common.confirm') : $t('common.delete')"
                >
                    <UButton
                        size="xs"
                        color="error"
                        :variant="confirmDelete ? 'solid' : 'soft'"
                        :icon="confirmDelete ? 'i-mdi-check' : 'i-mdi-delete'"
                        @click="deleteNodeWithConfirm"
                    />
                </UTooltip>
            </UFieldGroup>

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
